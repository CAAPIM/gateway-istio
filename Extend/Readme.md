## Extending the EDGE Gateway
The EDGE Gateway is deployed with a volume mount to its bootstrap folder meaning you can deploy your own bundlefiles.

If you installed Reloader then these changes will automatically be detected and deployed with a RollingUpdate strategy

### Place .bundle files into the bundles directory

```
$ make

kubectl create configmap "bundle-configmap" --from-file=bundles/_0_env.req.bundle -o yaml --dry-run | kubectl replace -f -
configmap/bundle-configmap replaced
```
