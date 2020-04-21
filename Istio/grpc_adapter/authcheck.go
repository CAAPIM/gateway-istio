// nolint:lll
// Generates the grpc_adapter's resource yaml. It contains the adapter's configuration, name, supported template
// names (auth in this case), and whether it is session or no-session based.
//go:generate $REPO_ROOT/bin/mixer_codegen.sh -a mixer/adapter/grpc_adapter/config/config.proto -x "-s=false -n grpcadapter -t authorization"

// The name can't have an underscore in it, which is why the above generator reads '-n grpcadapter'. These names must follow DNS rules becuase they become service names.

package grpc_adapter

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"

	"github.com/pkg/errors"
	"istio.io/api/mixer/adapter/model/v1beta1"
	policy "istio.io/api/policy/v1beta1"
	"istio.io/istio/mixer/adapter/grpc_adapter/config"
	"istio.io/istio/mixer/pkg/status"
	"istio.io/istio/mixer/template/authorization"
	"istio.io/pkg/log"
)

const (
	mySigningKey = "youshallnotpass!"
)

// The ServerAClaims is used as the claims portion of the JWT.
type ServerAClaims struct {
	Roles     []string `json:"roles"`
	AuthToken string   `json:"auth_token"`
	jwt.StandardClaims
}

type (
	// Server is basic server interface
	Server interface {
		Addr() string
		Close() error
		Run(shutdown chan error)
	}

	// AuthCheckAdapter supports metric template.
	AuthCheckAdapter struct {
		listener net.Listener
		server   *grpc.Server
	}
)

var _ authorization.HandleAuthorizationServiceServer = &AuthCheckAdapter{}

// HandleAuthorization - Perform the authorization check.
func (s *AuthCheckAdapter) HandleAuthorization(ctx context.Context, r *authorization.HandleAuthorizationRequest) (*v1beta1.CheckResult, error) {

	log.Infof("received request %v\n", *r)

	cfg := &config.Params{}

	if r.AdapterConfig != nil {
		if err := cfg.Unmarshal(r.AdapterConfig.Value); err != nil {
			log.Errorf("error unmarshalling adapter config: %v", err)
			return nil, err
		}
	}

	decodeValue := func(in interface{}) interface{} {
		switch t := in.(type) {
		case *policy.Value_StringValue:
			return t.StringValue
		case *policy.Value_Int64Value:
			return t.Int64Value
		case *policy.Value_DoubleValue:
			return t.DoubleValue
		default:
			return fmt.Sprintf("%v", in)
		}
	}

	decodeValueMap := func(in map[string]*policy.Value) map[string]interface{} {
		out := make(map[string]interface{}, len(in))
		for k, v := range in {
			out[k] = decodeValue(v.GetValue())
		}
		return out
	}

	props := decodeValueMap(r.Instance.Subject.Properties)
	//log.Infof("%v", props)

	/*
		There are some service calls that should be allowed to take place without having an
		authentication header - specifically login & register. The correct way to do this would
		involve a rule that looked for either of those two endpoints. If either endpoint was
		being accessed, then the authentication step would be bypassed.

		I don't know if you can make a rule that applies to a HTTP/REST endpoint specifically,
		and I'm running out of time, so I'm putting a hack here. If the `x-custom-token: abc` is
		present, that will be the indicator that says to skip the actual authentication check. This
		will be used for the register & login scripts, but the others will not have it and therefore
		they will need the token.

		Obviously this is purely for demonstration purposes. The real implementation would
		require the URL the client is calling to be used to determine whether to check
		for an authentic token.
	*/

	check := true
	tokenStr := ""
	for k, v := range props {
		log.Infof("k: %s, v: %s\n", k, v)

		if k == "auth_token_header" {
			tokenStr = v.(string)
			log.Infof("found auth_token_header, set tokenStr to: %q", tokenStr)
		}

		if k == "custom_token_header" {
			cust := v.(string)
			if cust != "" {
				check = false
				log.Infof("found custom_token_header, NOT doing the auth check")
			} else {
				log.Infof("found custom_token_header but its empty, so doing the auth check")
			}
		}

		/*
			if (k == "custom_token_header") && v == cfg.AuthKey {
				log.Infof("success!!")
				return &v1beta1.CheckResult{
					Status: status.OK,
				}, nil
			}
		*/
	}

	if !check {
		return &v1beta1.CheckResult{
			Status: status.OK,
		}, nil
	}

	log.Infof("Token to parse: %q\n", tokenStr)

	var server_a_claims ServerAClaims
	if tokenStr != "" {
		parts := strings.Split(tokenStr, " ")
		if len(parts) < 2 {
			return nil, errors.New("malformed auth header")
		}

		log.Infof("parsing: %q\n", parts[1])

		token, err := jwt.ParseWithClaims(parts[1], &server_a_claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		})
		if err != nil {
			fmt.Printf("*** problem parsing jwt: %s\n", err.Error())
			return nil, err
		}
		//log.Infof("parsed token:", token)
		_ = token
	}

	t := fmt.Sprintf("Bearer %s", server_a_claims.AuthToken)

	log.Infof("header to be checked: %q", t)

	if ok, err := CheckToken(t); err == nil {
		if ok {
			return &v1beta1.CheckResult{
				Status: status.OK,
			}, nil
		}
	} else {
		log.Errorf("problem checking token: %q", err.Error())
	}

	log.Infof("failure; header not provided")
	return &v1beta1.CheckResult{
		Status: status.WithPermissionDenied("Unauthorized..."),
	}, nil
}

/* Failure
{
  "result": "invalid",
   "message": {
  "error":"invalid_request",
  "error_description":"Validation error"
	}
}
*/

/* Success
{
   "result": "valid",
   "expiration": 1582010254
}

*/

type AuthResult struct {
	Message struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	} `json:"message,omitempty"`
	Result     string `json:"result"`
	Expiration int64  `json:"expiration"`
}

func CheckToken(tokenStr string) (bool, error) {
	stsUrl := os.Getenv("stsUrl")
	stsNamespace := os.Getenv("stsNamespace")
	stsPort := os.Getenv("stsPort")
	stsUri := os.Getenv("stsUri")
	finalUrl := "http://" + stsUrl + "." + stsNamespace + ":" + stsPort + stsUri
	req, err := http.NewRequest("GET", finalUrl, nil)
	header := req.Header
	header.Set("Authorization", tokenStr)

	log.Infof("---------CHECK TOKEN-----------")
	defer log.Infof("------->Check Token<--------")

	log.Infof("Headers: %v", header)

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true}, // disable cert checking
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return false, errors.Wrap(err, "problem calling authenticator")
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return false, errors.Wrap(err, "problem reading body")
	}
	var res AuthResult
	err = json.Unmarshal(body, &res)
	if err != nil {
		return false, errors.Wrap(err, "problem unmarshalling body")
	}

	return res.Result == "valid", nil
}

// Addr returns the listening address of the server
func (s *AuthCheckAdapter) Addr() string {
	return s.listener.Addr().String()
}

// Run starts the server run
func (s *AuthCheckAdapter) Run(shutdown chan error) {
	shutdown <- s.server.Serve(s.listener)
}

// Close gracefully shuts down the server; used for testing
func (s *AuthCheckAdapter) Close() error {
	if s.server != nil {
		s.server.GracefulStop()
	}

	if s.listener != nil {
		_ = s.listener.Close()
	}

	return nil
}

// NewAuthCheckAdapter creates a new IBP adapter that listens at provided port.
func NewAuthCheckAdapter(addr string) (Server, error) {
	if addr == "" {
		addr = "0"
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", addr))
	if err != nil {
		return nil, fmt.Errorf("unable to listen on socket: %v", err)
	}
	s := &AuthCheckAdapter{
		listener: listener,
	}
	fmt.Printf("listening on \"%v\"\n", s.Addr())
	s.server = grpc.NewServer()
	authorization.RegisterHandleAuthorizationServiceServer(s.server, s)
	return s, nil
}
