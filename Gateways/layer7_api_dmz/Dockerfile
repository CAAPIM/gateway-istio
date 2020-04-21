FROM caapim/gateway:latest

## Copying the deployment package
COPY build/gateway/layer7-api-dmz-1.0.0.gw7 /opt/docker/rc.d/otk5.gw7

## Make restman available
RUN touch /opt/SecureSpan/Gateway/node/default/etc/bootstrap/services/restman
