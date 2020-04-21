# layer7 istio demo server a and b

This is a sub chart of layer7 (edge) and will automatically be installed when you run through the documentation on the main page

# Install the Demo Chart

## From this Repository

`helm install demo .`

## From the layer7 Chart Repo

`helm repo add layer7 https://charts.brcmlabs.com`
`helm repo update`
`helm install demo layer7/demo`

## Delete this Chart

`helm del <release name>`


## Configuration

The following table lists the configurable parameters of the Demo Chart.

| Parameter                        | Description                               | Default                                                      |
| -----------------------------    | -----------------------------------       | -----------------------------------------------------------  |
| `nameOverride`                   | Override the name        | `nil`                                                          |
| `fullnameOverride`             | Override the full name                      | `nil`                                              |
| `stsUrl`                | The STS URL, it's unlikely you will need to modify this   | `sts-gateway` |
| `stsNamespace`                      | The namespace you install the parent chart into                        | `default`                                                     |
| `stsPort`          | The STS Gateway Port for the Istio Auth service                   | `18888`                                                     |
| `stsUri`          | The STS Gateway Istio Auth service URI | `/auth/istio/v1/token`  |
| `servera.image`          | Server A image | `layer7api/server_a`  |
| `servera.tag`          | Server A Tag | `1.0`  |
| `servera.pullPolicy`          | Server A Pull Policy | `Always`  |
| `serverb.image`          |Server B image | `layer7api/server_b`  |
| `serverb.tag`    | Server B Tag | `1.0`  |
| `serverb.pullPolicy`     | Server B Pull Policy | `Always`  |