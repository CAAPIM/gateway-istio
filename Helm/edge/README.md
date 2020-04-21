# Layer7 API Gateway deployed as an Edge Service

This chart deploys the API Gateway with services to expose the Layer7 Istio demo

# Install EDGE Chart

## From this Repository

The API Gateway requires a license for its installation. 
`helm dep build`

`helm install layer7 --set-file "global.gateway.license.value=path/to/license.xml" --set "global.license.accept=true" .`

## From layer7 Chart Repo

`helm repo add layer7 https://charts.brcmlabs.com`

`helm repo update`

`helm install layer7 --set-file "global.gateway.license.value=path/to/license.xml" --set "global.license.accept=true" layer7/layer7`


## Delete this Chart
To delete STS installation

`helm del <release name>`

## Configuration

The following table lists the configurable parameters of the STS chart and their default values. See values.yaml for additional parameters and info

| Parameter                        | Description                               | Default                                                      |
| -----------------------------    | -----------------------------------       | -----------------------------------------------------------  |
| `replicas`                   | Number of Gateway service replicas        | `1`                                                          |
| `deploymentStrategy`             | Deployment Strategy                       | `rollingUpdate`                                              |
| `nameOverride`                | Name override   | `""` |
| `fullnameOverride`                      | Full name override                       | `sts-gateway`                                                     |
| `global.license.value`          | Gateway license file | `nil`  |
| `global.license.accept`          | Accept Gateway license EULA | `false`  |
| `service.type`    | Service Type               | `ClusterIP` |
| `service.ports`    | List of http external port mappings               | https: 8443 -> 8443, istio: 18888->18888 |
| `gateway.clusterHostname`          | Gateway Cluster Hostname  | `broadcom.localdomain`  |
| `gateway.influxdb.url`          | InfluxDb location | `influxdb`  |
| `gateway.influxdb.tags`          | InfluxDb Tags | `env=sts`  |
| `gateway.manager.enabled`          | Enable/Disable Policy Manager access | `true`  |
| `gateway.manager.username`          | Policy Manager Username | `admin`  |
| `gateway.manager.password`          | Policy Manager Password | `7layer`  |
| `gateway.perftestUrl`          | Built in EDGE gateway Performance test | `https://sts-gateway:8443/oauth/manager`  |
| `gateway.demoUrl`          | Location of the demo a service | `http://server-a-svc:8080`  |


### Logs & Audit Configuration

The API Gateway containers are configured to output logs and audits as JSON events, and to never write audits to the in-memory Derby database:

- System properties in the default template for the `gateway.javaArgs` value configure the log and audit behaviour:
  - Auditing to the database is disabled: `-Dcom.l7tech.server.audit.message.saveToInternal=false -Dcom.l7tech.server.audit.admin.saveToInternal=false -Dcom.l7tech.server.audit.system.saveToInternal=false`
  - JSON formatting is enabled: `-Dcom.l7tech.server.audit.log.format=json`
  - Default log output configuration is overridden by specifying an alternative configuration properties file: `-Djava.util.logging.config.file=/opt/SecureSpan/Gateway/node/default/etc/conf/log-override.properties`
- The alternative log configuration properties file `log-override.properties` is mounted on the container, via the `gateway-config` ConfigMap.
- System property to include well known Certificate Authorities Trust Anchors 
    - API Gateway does not implicitly trust certificates without importing it but If you want to avoid import step then configure Gateway to accept any certificate signed by well known CA's (Certificate Authorities)
      configure following property to true -
      Set '-Dcom.l7tech.server.pkix.useDefaultTrustAnchors=true' for well known Certificate Authorities be included as Trust Anchors (true/false)
- Allow wildcards when verifying hostnames (true/false)
    - Set '-Dcom.l7tech.security.ssl.hostAllowWildcard=true' to allow wildcards when verifying hostnames (true/false)

### Subcharts

This Chart uses the following SubCharts
*  Hazelcast ==> https://github.com/helm/charts/tree/master/stable/hazelcast
*  STS  ==> this repository
*  Demo   ==> this repository