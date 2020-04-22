# This Makefile replaces the EDGE Gateway Configmap
# if you installed reloader it will automatically redeploy the EDGE Gateway with your updates

CONFIGMAPNAME="bundle-configmap"
DIRS := $(patsubst %,--from-file=%,$(wildcard bundles/*))
CHECKGWCONFEXISTS="$(shell kubectl get configmap $(CONFIGMAPNAME) 2>&1)"

.phony: all
all:
ifneq (,$(findstring (NotFound),$(CHECKGWCONFEXISTS)))
	kubectl create  configmap ${CONFIGMAPNAME} $(DIRS)
else
	kubectl create configmap ${CONFIGMAPNAME} $(DIRS) -o yaml --dry-run | kubectl replace -f -
endif
