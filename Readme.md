Broadcom Layer7 Istio Mesh Networking
====================================
Public Alpha May 2020
----------------------
This project supplies a Kubernetes based mesh networking system to use the Broadcom Layer7 API Management Suite with mesh as an integrated and easy to use part of enterprise deployment of APIs.

Here's a [Discussion and demo of the project](https://apiacademy.co/2020/01/tech-talk-building-api-aware-service-meshes-with-layer7-istio/)
  
This is provided as source configuration, and is supplied as-is, with no implied warranty as to fitness to purpose. Please feel free to use this as source material for your own project. If we get something wrong, report an issue and we'll try to get it fixed, or send us a pull request with your fix. 

The Mixer code especially is not performance tested, and serves as a sample. When Istio's authorization framework is complete, this portion of the repository will be obsolete. At that point, we will augment this to include configuration for the auth framework. 

The focus of the project is to bring the world of API Management and the world of Mesh networking together. Our view is that this is a natural progression of the world of API Managment, and an obvious enhancement to Mesh & Istio.

Table of Contents
=================
<!--ts-->
   * [Broadcom Layer7 Istio Mesh Networking](#broadcom-layer7-istio-mesh-networking)
      * [Public Alpha March 2020](#public-alpha-march-2020)
   * [Table of Contents](#table-of-contents)
   * [Prerequisites](#prerequisites)
      * [Enterprise, Hosted or Cloud-Based Kubernetes Deployment prerequisites](#enterprise-hosted-or-cloud-based-kubernetes-deployment-prerequisites)
   * [Setup](#setup)
      * [Getting Started](#getting-started)
      * [We'll also need Helm](#well-also-need-helm)
      * [Start using toolchain](#start-using-toolchain)
      * [Install Reloader](#install-reloader)
      * [Check status](#check-status)
   * [Project Setup](#project-setup)
   * [Install Demo components](#install-demo-components)
   * [Validate your installation](#validate-your-installation)
   * [Configure hosts](#configure-hosts)
   * [Test that the auth adapter is working correctly](#test-that-the-auth-adapter-is-working-correctly)
   * [Evaluating the toolchain](#evaluating-the-toolchain)
      * [Built in sample projects](#built-in-sample-projects)
      * [Dashboard](#dashboard)
   * [Possible error Conditions](#possible-error-conditions)
      * [No APIM License](#no-apim-license)

<!-- Added by: jaythorne, at: Fri 24 Apr 2020 11:52:13 PDT -->

<!--te-->

# Prerequisites

The project is based on a Kubernetes environment. There are optional components and this project is specifically focused on Istio based mesh networking. Other projects in future will address other mesh networking systems.


* ## Enterprise, Hosted or Cloud-Based Kubernetes Deployment prerequisites
  * You will need to have a kubernetes system built and installed and have credentials for the kubernetes control plane. This implies that kubernetes is running on a Docker hosting system, perhaps more than one.  
  
  * Docker and Kubernetes setup is external to this project. We recommend a 'clean' Kubernetes for your initial project. We used Google's [GKE](https://cloud.google.com/kubernetes-engine) during our development effort. They have a free tier that this project is compatible with. 

  * As this project uses and focuses on the Broadcom Layer7 API Management suite, a commercial product, you will need access to those components and a license file. If you want to evaluate this prior to becoming a customer, our sales organization can help set that up for you. [Sales](mailto:sales@broadcom.com)

  * You will need "kubectl" available as your command interface into your environment: [Instructions for installing that on the Kubernetes site.](https://kubernetes.io/docs/tasks/tools/install-kubectl). Kubectl needs access to your kubernetes cluster [Setting that up is also explained on the kubernetes site](https://kubernetes.io/docs/tasks/access-application-cluster/access-cluster/) 

  * The project supplies ingress artifacts for "bare" kubernetes as well as Istio. Those can be used independently or in concert with the mesh interfaces. 
  
  * While Istio is technically an optional component, it's a strong requirement for some of the capabilities included. The default flow pulls Istio automatically. Other options are noted later in this document. 

# Setup

## Getting Started
There are a few pre-requisites that we'll need to cover before we get Started

* `PreReq Step 1:` Kubernetes

  *   If you're using GKE: [Use these instructions](GKE-Readme.md)

  *   Local install using Kubernetes [Have a look at this](K8S-Readme.md)

## We'll also need Helm Version 3
* `PreReq Step 2:` Helm v3

  - Helm Installation [Documentation](https://helm.sh/docs/intro/install/)

  - Verify your Helm version. V3 is absolutely required
  ```
  $ helm version
  version.BuildInfo{Version:"v3.1.2", GitCommit:"d878d4d45863e42fd5cff6743294a11d28a9abce", GitTreeState:"clean", GoVersion:"go1.13.8"}

   ```
   MacOS Example using HomeBrew (against a GKE Cluster)
   $ brew install helm
   $ kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin --user <youremailaddress.com>


   ```
  
* `PreReq Step 3:` Add Helm Repositories

   ```
   $ helm repo add stable https://kubernetes-charts.storage.googleapis.com/
  "stable" has been added to your repositories
  ```

  - Add the layer7 alpha helm repository

   ```
  $ helm repo add layer7 https://charts.brcmlabs.com
  "layer7" has been added to your repositories
   ```

## Start using toolchain
* At this point we've got most of a working base, so lets test out our tooling configuration

* `Step 1:` Deploy nginx from stable as a quick test, plus we'll need it later for grafana

  ```
  $ helm install nginx-ingress stable/nginx-ingress
  NAME: nginx-ingress
  LAST DEPLOYED: Mon Mar 23 11:55:26 2020
  NAMESPACE: default
  STATUS: deployed
  REVISION: 1
  TEST SUITE: None
  NOTES:
  The nginx-ingress controller has been installed.
  It may take a few minutes for the LoadBalancer IP to be available.
  ... etc ..
  ```
## Install Reloader
* `Step 1a:` This is optional, Reloader automatically watches the Layer7 Gateway configMaps.
  ```
  $ helm install reloader stable/reloader --set reloader.watchGlobally=false --namespace default 
  NAME: reloader
  LAST DEPLOYED: Tue Mar 24 22:15:20 2020
  NAMESPACE: default
  STATUS: deployed
  REVISION: 1
  TEST SUITE: None
  NOTES:
  ... etc ... 
  ```

## Check status
  * `Step 2:` check status both in kubernetes and in helm
    ```
    $ kubectl get pods
    get pods
    NAME                                            READY   STATUS    RESTARTS   AGE
    nginx-ingress-controller-6b85b64f49-kppsp       1/1     Running   0          20h
    nginx-ingress-default-backend-7db6cc5bf-rv78n   1/1     Running   0          20h
    reloader-reloader-7596457dd9-vtxqd              1/1     Running   0          20h
    
    $ helm list
    NAME         	NAMESPACE	REVISION	UPDATED                             	STATUS  	CHART               	APP VERSION
    nginx-ingress	default  	1       	2020-03-25 16:23:17.933944 -0700 PDT	deployed	nginx-ingress-1.34.2	0.30.0     
    reloader     	default  	1       	2020-03-25 16:23:48.813671 -0700 PDT	deployed	reloader-1.2.0      	v0.0.41    
    ```
   
# Project Setup

  * `Step 3:` At this point we want to start installing our project. Let's Check out the code:
    ```
    $ mkdir <my project>
    $ cd <my project>
    $ git clone https://github.com/CAAPIM/gateway-istio.git
    $ cd gateway-istio
    ```

# Install Demo components
 * `Step 4a:` Install Istio

    ```
    Continuing on from the previous step
    $ cd Istio
    $ make install
    .. lots of stuff ...
    ends with
    kubectl label namespace default istio-injection=enabled
    namespace/default labeled
    -------
    Remember to label any namespace that is to be istio controlled with this injection flag
    ```

  * `Step 4b:` Deploy the Layer7 Helm Chart
    ```
    Note:
    You can seamlessly switch Gateway versions using 'upgrade' in place of 'install' and the relevant v9.x or v10.x license
    
    Gateway v9.4
    $ helm install layer7 --set-file "global.license.value=/<file path>/license_94.xml" --set "global.license.accept=true,global.version=9.4" layer7/layer7
    
    Gateway v10.0
    $ helm install layer7 --set-file "global.license.value=/<file path>/license_10.xml" --set "global.license.accept=true,global.version=10.0" layer7/layer7
    
    NAME: layer7
    LAST DEPLOYED: Tue Mar 31 16:40:13 2020
    NAMESPACE: default
    STATUS: deployed
    REVISION: 3
    NOTES:
    1. Get the application URL by running these commands:
    ```
    
  * `Step 4c:` Install the Layer7 Custom Auth Adapter
    ```
    $ cd <my project>/Istio/custom_auth_adapter
    $ kubectl apply -f .
    attributemanifest.config.istio.io/istio-proxy created
    attributemanifest.config.istio.io/kubernetes configured
    service/layer7authadapterservice created
    deployment.apps/layer7authadapter created
    adapter.config.istio.io/layer7authadapter created
    handler.config.istio.io/h1 created
    instance.config.istio.io/icheck created
    rule.config.istio.io/r1 created
    template.config.istio.io/authorization created
    ```

# Validate your installation

  * `Step 5:` At this point the demo should be running in your cluster. Check out the Troubleshooting section if you have any errors:

  * Show the Helm Charts we installed
    ```
    $ helm list
    NAME         	NAMESPACE	REVISION	UPDATED                             	STATUS  	CHART               	APP VERSION
    layer7       	default  	1       	2020-03-26 13:09:09.985029 -0700 PDT	deployed	layer7-1.0.0           	9.4        
    nginx-ingress	default  	1       	2020-03-25 16:23:17.933944 -0700 PDT	deployed	nginx-ingress-1.34.2	0.30.0     
    reloader     	default  	1       	2020-03-25 16:23:48.813671 -0700 PDT	deployed	reloader-1.2.0      	v0.0.41  
    ```

  * Verify Istio and the custom auth adapater (layer7authadapter-xxx-xxx) are installed - (Optional)
  
    ```
    $ kubectl get pods --namespace istio-system
    
    NAME                                      READY   STATUS    RESTARTS   AGE
    grafana-5f798469fd-l48dj                  1/1     Running   0          10d
    istio-citadel-58bb67f9b8-wscmc            1/1     Running   0          10d
    istio-egressgateway-6fd57475b5-tsbg6      1/1     Running   0          10d
    istio-galley-7d4b9874c8-qr4jw             1/1     Running   0          10d
    istio-ingressgateway-7d65bf7fdf-n8fz7     1/1     Running   0          10d
    istio-pilot-65f8557545-qt5ln              1/1     Running   0          10d
    istio-policy-6c6449c56f-7mbgt             1/1     Running   1          10d
    istio-sidecar-injector-774969d686-rh4c5   1/1     Running   0          10d
    istio-telemetry-585cc965f7-rgcqm          1/1     Running   1          10d
    istio-tracing-cd67ddf8-qk7nv              1/1     Running   0          10d
    kiali-7964898d8c-twkd9                    1/1     Running   0          10d
    layer7authadapter-97f6bdb86-n9r5x         1/1     Running   0          10d
    prometheus-586d4445c7-xp775               1/1     Running   0          10d 
    ```

# Configure hosts

`Step 6:` Now we need to able to get requests to our system
* `6a:` Gateway External name "broadcom.localdomain"
    ```
    $ kubectl get svc --namespace istio-system | grep ingressgateway
    
    NAME                       TYPE           CLUSTER-IP      EXTERNAL-IP    PORT(S)  
    istio-ingressgateway       LoadBalancer   10.15.255.101    1.2.3.4   15020:30832/TCP,80:30473/TCP,443:32584/TCP,15029:31754/TCP,15030:30651/TCP,15031:32536/TCP,15032:32716/TCP,15443:30051/TCP   67m
    ```

* Second IP address is external facing address for gateway

* `6b:` Grafana External name "grafana.localdomain"
    ```
    $ kubectl get svc | grep ingress-controller
    NAME                       TYPE           CLUSTER-IP      EXTERNAL-IP    PORT(S)  
    nginx-ingress-controller   LoadBalancer   10.15.254.100   1.2.3.5   80:32283/TCP,443:31979/TCP   25h
    ```
* Second IP address is the external ip for the grafana dashboard

* `6c:` Add those to /etc/hosts: Add these in as names
    ```
      $ sudo vi /etc/hosts
      Append: 
      # my layer7 project
      1.2.3.4 broadcom.localdomain
      1.2.3.5 grafana.localdomain
    ````
* `6d:` Open up a modern browser (i.e. Chrome) with two tabs
    ```
      tab 1: https://broadcom.localdomain
        - username: admin
        - password: 7layer
      tab 2: https://grafana.localdomain
        - username: admin
        - password: password
      Append: 
      # my layer7 project
      1.2.3.4 broadcom.localdomain
      1.2.3.5 grafana.localdomain
    ````

# Test that the auth adapter is working correctly
    
      $ kubectl get pods -n istio-system | grep layer7
      
      layer7authadapter-xxx-xxx         1/1     Running   0          74m
      
      $ kubectl logs -f layer7authadapter-97f6bdb86-fqnsz -n istio-system
      
      In your browser (https://broadcom.localdomain) click 'login' in the top right corner

      username: admin
      password: 7layer

      GET ACCESS TOKEN

      Type a name in the Add Topic field and click Create

      You should now see logs from the Auth Adapter
        - the first request to retrieve the access token is ignored
        - subsequent requests are checked where the access token is validated
      
      You will also see a call to istio-mixer in Grafana!

      listening on "[::]:44225"
      2020-04-03T15:02:04.916877Z     info    transport: loopyWriter.run returning. connection error: desc = "transport is closing"
      2020-04-03T16:10:31.168177Z     info    received request {&InstanceMsg{Subject:&SubjectMsg{User:,Groups:,Properties:map[string]*v1beta1.Value      {auth_token_header: &Value{Value:&Value_StringValue{StringValue:,},},custom_token_header: &Value{Value:&Value_StringValue{StringValue:abc,},},      fingerprint_header: &Value{Value:&Value_StringValue{StringValue:,},},},},Action:nil,Name:icheck.instance.istio-system,} &Any{TypeUrl:type.      googleapis.com/adapter.grpc_adapter.config.Params,Value:[10 3 97 98 99],XXX_unrecognized:[],} 13466988563306109762}
      
      2020-04-03T16:10:31.168225Z     info    k: fingerprint_header, v: 
      
      2020-04-03T16:10:31.168234Z     info    k: auth_token_header, v: 
      
      2020-04-03T16:10:31.168240Z     info    found auth_token_header, set tokenStr to: ""
      2020-04-03T16:10:31.168243Z     info    k: custom_token_header, v: abc
      
      2020-04-03T16:10:31.168246Z     info    found custom_token_header, NOT doing the auth check
      2020-04-03T16:10:31.922759Z     info    received request {&InstanceMsg{Subject:&SubjectMsg{User:,Groups:,Properties:map[string]*v1beta1.Value      {auth_token_header: &Value{Value:&Value_StringValue{StringValue:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.      eyJyb2xlcyI6bnVsbCwiYXV0aF90b2tlbiI6IjUyMDE1ZmY0LTJjYjYtNDFmMC1hYTIxLWE1NjRlNjcyYjI3NCIsImV4cCI6MTU4NjAxNjYzMSwiaXNzIjoic2VydmVyX2EifQ.      94NP_kBNXep1imvOSlp5eUFDXUt0Gof8wQXkcptpcYs,},},custom_token_header: &Value{Value:&Value_StringValue{StringValue:,},},fingerprint_header: &Value      {Value:&Value_StringValue{StringValue:,},},},},Action:nil,Name:icheck.instance.istio-system,} &Any{TypeUrl:type.googleapis.com/adapter.      grpc_adapter.config.Params,Value:[10 3 97 98 99],XXX_unrecognized:[],} 13466988563306109763}
      
      2020-04-03T16:10:31.922821Z     info    k: custom_token_header, v: 
      
      2020-04-03T16:10:31.922826Z     info    found custom_token_header but its empty, so doing the auth check
      2020-04-03T16:10:31.922831Z     info    k: fingerprint_header, v: 
      
      2020-04-03T16:10:31.922836Z     info    k: auth_token_header, v: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.      eyJyb2xlcyI6bnVsbCwiYXV0aF90b2tlbiI6IjUyMDE1ZmY0LTJjYjYtNDFmMC1hYTIxLWE1NjRlNjcyYjI3NCIsImV4cCI6MTU4NjAxNjYzMSwiaXNzIjoic2VydmVyX2EifQ.      94NP_kBNXep1imvOSlp5eUFDXUt0Gof8wQXkcptpcYs
      
      2020-04-03T16:10:31.922850Z     info    found auth_token_header, set tokenStr to: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.      eyJyb2xlcyI6bnVsbCwiYXV0aF90b2tlbiI6IjUyMDE1ZmY0LTJjYjYtNDFmMC1hYTIxLWE1NjRlNjcyYjI3NCIsImV4cCI6MTU4NjAxNjYzMSwiaXNzIjoic2VydmVyX2EifQ.      94NP_kBNXep1imvOSlp5eUFDXUt0Gof8wQXkcptpcYs"
      2020-04-03T16:10:31.922863Z     info    Token to parse: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.      eyJyb2xlcyI6bnVsbCwiYXV0aF90b2tlbiI6IjUyMDE1ZmY0LTJjYjYtNDFmMC1hYTIxLWE1NjRlNjcyYjI3NCIsImV4cCI6MTU4NjAxNjYzMSwiaXNzIjoic2VydmVyX2EifQ.      94NP_kBNXep1imvOSlp5eUFDXUt0Gof8wQXkcptpcYs"
      
      2020-04-03T16:10:31.922902Z     info    parsing: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.      eyJyb2xlcyI6bnVsbCwiYXV0aF90b2tlbiI6IjUyMDE1ZmY0LTJjYjYtNDFmMC1hYTIxLWE1NjRlNjcyYjI3NCIsImV4cCI6MTU4NjAxNjYzMSwiaXNzIjoic2VydmVyX2EifQ.      94NP_kBNXep1imvOSlp5eUFDXUt0Gof8wQXkcptpcYs"
      
      2020-04-03T16:10:31.923086Z     info    header to be checked: "Bearer 52015ff4-2cb6-41f0-aa21-a564e672b274"
      2020-04-03T16:10:31.923160Z     info    ---------CHECK TOKEN-----------
      2020-04-03T16:10:31.923182Z     info    Headers: map[Authorization:[Bearer 52015ff4-2cb6-41f0-aa21-a564e672b274]]
      2020-04-03T16:10:32.043206Z     info    ------->Check Token<--------
    

# Evaluating the toolchain

You should start by getting a look at publishing an API, and then using it. We provide tools here that automate a "my first published API" project.

## Built in sample projects

* The "Server_A" and "Server_B" are our provided East/West sample project. They provide an extremely crude topic system, that composes a server (Server A) with a back end API (Server B). These are connected via a ruleset that uses a mixer based authorization, backed by the API Gateway in OTK Secure Token Service mode (sts)
** Access that via the https://broadcom.localdomain/ URL assuming you have used the suggested domains as noted above. 

* There's a built in ingress test chain, based on an oversimplified policy. That's accessed via https://broadcom.localdomain/perftest and is controlled via http headers. Output is small and is meant only to demonstrate latency of the ingress path.  
* Performance testing via apache benchmark: 

   ```ab -c 1 -n 10 -H "foo: Foo" -H "route: no" https://broadcom.localdomain/perftest```
   
   Provides a simple header based flow: "foo: Foo" is required or the policy will emit an error, this specifically to show errors in statistics to indicate policy based errors. 

  ```ab -c 1 -n 10 -H "foo: Foo" -H "route: Yes" https://broadcom.localdomain/perftest```

  Changing the route header to Yes will show back end latency - it routes to the STS login page to provide both latency and response message size. 

  Testing via tools like gatling, jmeter or other performance tests are certainly possible. Note the headers and URL as above. By default this uses a self signed certificate. 

## Dashboard
* Grafana with Influx backing is included with a pre-built dashboard accessible via https://grafana.localdomain/. This is a default grafana security setup, so username / password is the default. See https://docs.gitlab.com/ee/administration/monitoring/performance/grafana_configuration.html for more details.  




# Possible error Conditions
## No APIM License

This is marked by a restart as the availability check won't find a useful connection.

```
$kubectl get pods
NAME                               READY   STATUS    RESTARTS   AGE
gatewaydev-84df67fbd5-xs9hj        0/2     Running   12          20m

```
