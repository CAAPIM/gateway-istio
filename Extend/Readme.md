# Extend the EDGE Gateway
The EDGE Gateway can be extended/modified with your own bundle files.

If you installed Reloader then these changes will automatically be detected and deployed with a RollingUpdate strategy

### Place .bundle files into the bundles directory
The below shows a deployment with '3' edge replicas. The update strategy will maintain '3' healthy replicas throughout the upgrade.
```
$ make

kubectl create configmap "bundle-configmap" --from-file=bundles/_0_env.req.bundle -o yaml --dry-run | kubectl replace -f -
configmap/bundle-configmap replaced

$ kubectl get pods

edge-gateway-57fff6f567-dm2vz       2/2     Running     0          21m
edge-gateway-57fff6f567-lcg5q       2/2     Running     0          10m
edge-gateway-57fff6f567-zf2qs       2/2     Running     0          10m
## edge-gateway-6f87867475-2rk8f       1/2     Running     0          12s
edge-hazelcast-0                    2/2     Running     0          30m
edge-hazelcast-1                    2/2     Running     0          29m
grafana-6bc9b9866d-gn86d            2/2     Running     0          30m
influxdb-6bdb755574-tq8jz           2/2     Running     0          30m
influxdb-job-4hfjs                  0/1     Completed   0          30m
otk-mysql-6c79bfff5d-bfmm9          2/2     Running     0          30m
reloader-reloader-5fdc7bcdd-8m7pq   2/2     Running     0          3h54m
server-a-dep-856d8fdc9c-qb9zg       2/2     Running     0          30m
server-b-app-84bc745dbd-vzfqw       2/2     Running     0          30m
sts-gateway-544fdf67f7-svcxf        2/2     Running     0          30m
sts-hazelcast-0                     2/2     Running     0          30m
sts-hazelcast-1                     2/2     Running     0          29m
```

