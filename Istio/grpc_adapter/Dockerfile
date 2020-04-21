FROM alpine:3.8
RUN apk --no-cache add ca-certificates
WORKDIR /bin/
COPY src/istio.io/istio/mixer/adapter/grpc_adapter/bin/grpc_adapter .
ENTRYPOINT [ "/bin/grpc_adapter" ]
CMD [ "44225" ]
EXPOSE 44225
