# You can update the istio version here. The `create` stanza is built from 
# the Istio installation instructions, and expects a particular directory 
# structure. This has been tested with istio 1.3.0 and 1.3.1. If you update
# past 1.3.1 and have a problem, you will need to see what's changed.
ISTIO_VERSION=1.4.5
ISTIO_DIR=istio-$(ISTIO_VERSION)
ISTIOCTL=$(ISTIO_DIR)/bin/istioctl

empty:
	@echo "use: make [install | uninstall]"
	@echo "Current Istio version is: $(ISTIO_VERSION)"

$(ISTIO_DIR):
	curl -L https://git.io/getLatestIstio | ISTIO_VERSION=$(ISTIO_VERSION) sh -

install: $(ISTIO_DIR)

	# NOTE: The following single ISTIOCTL apply statement doesn't work. The 3 statements that follow it
	# are required to get a properly functioning ingress & gateway. This seems to be a bug, because
	# the two stanzas appear to be functionally the same, other than a single invocation vs
	# three.  (The problem you see when it goes wrong is getting a 502 "Bad Gateway" error when 
	# curling server-a.localdomain).
	# I have reported this as a bug, and I am leaving this line here commented out in case someone
	# thinks to "optimize" the 3 lines into 1.
	#
	# yes | $(ISTIOCTL) manifest apply --set profile=demo --set values.global.mtls.enabled=true --set values.global.controlPlaneSecurityEnabled=true

	yes | $(ISTIOCTL) manifest apply && \
	$(ISTIOCTL) manifest apply --set values.global.mtls.enabled=true --set values.global.controlPlaneSecurityEnabled=true && \
	$(ISTIOCTL) manifest apply --set profile=demo

	kubectl label namespace default istio-injection=enabled
	@echo "-------"
	@echo "Remember to label any namespace that is to be istio controlled with this injection flag"

