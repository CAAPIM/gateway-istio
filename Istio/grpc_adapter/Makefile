ROOT_FOLDER:=$(shell pwd)
GOPATH:=$(shell pwd)
ISTIO:=$(GOPATH)/src/istio.io
MIXER_REPO:=$(ISTIO)/istio/mixer
ARCHIVE=https://github.com/greenenergy/istio_custom_auth_adapter.git

default:
	@echo "Order should be: make init, make gen, make yamls, make docker, make install"
	@echo "Before running 'make gen', you will need to run 'make mount' in another shell, AND make sure "\
		"your kubectl default context is *not* minikube. You can switch it back after 'make gen'."
	@echo "Once you've gotten to the 'make docker' stage, you won't need to revisit the earlier stages."

mount:
	cd && minikube mount .:`pwd`

init:
	mkdir -p src/istio.io && \
	cd src/istio.io && \
	git clone -b release-1.4 --single-branch https://github.com/istio/istio
	mkdir -p $(MIXER_REPO)/adapter/grpc_adapter/config
	cp  $(ROOT_FOLDER)/config/config.proto $(MIXER_REPO)/adapter/grpc_adapter/config/config.proto
	mkdir -p $(MIXER_REPO)/adapter/grpc_adapter/cmd
	cp  $(ROOT_FOLDER)/*.go $(MIXER_REPO)/adapter/grpc_adapter	# These source copies look redundant, but are necessary for 'make gen'.
	cp  $(ROOT_FOLDER)/cmd/*.go $(MIXER_REPO)/adapter/grpc_adapter/cmd/


	# It might be a good idea to make sure the version of istio checked out here is the same as what's used
	# in the main installation. In that case use the next line.
	# git clone --single-branch https://github.com/istio/istio
	#git clone -b release-1.4 --single-branch https://github.com/istio/istio && \

gen:
	@echo "For this step, you must make sure your kubectl config is pointing at a non-minikube kubernetes instance."
	@echo "By the way - *have patience*. For some reason, this build process *seems* to hang after declaring 'Using local Kubernetes configuration...'"\
		"On my system it hangs there for almost 5 minutes, and then just proceeds."
	cd $(ISTIO)/istio && \
	BUILD_WITH_CONTAINER=1 make gen

yamls:
	mkdir -p $(MIXER_REPO)/adapter/grpc_adapter/testdata
	cp $(MIXER_REPO)/adapter/grpc_adapter/config/grpcadapter.yaml $(MIXER_REPO)/adapter/grpc_adapter/testdata
	cp $(MIXER_REPO)/testdata/config/attributes.yaml $(MIXER_REPO)/adapter/grpc_adapter/testdata
	cp $(MIXER_REPO)/template/authorization/template.yaml $(MIXER_REPO)/adapter/grpc_adapter/testdata
	cp  $(ROOT_FOLDER)/yaml/operator_cfg.yaml $(MIXER_REPO)/adapter/grpc_adapter/testdata

build:
	@echo "You can switch your kubectl context back to minikube"
	cd $(MIXER_REPO)/adapter/grpc_adapter && \
	go build ./...

docker:
	cp  $(ROOT_FOLDER)/*.go $(MIXER_REPO)/adapter/grpc_adapter
	cp  $(ROOT_FOLDER)/cmd/*.go $(MIXER_REPO)/adapter/grpc_adapter/cmd/
	(cd $(MIXER_REPO)/adapter/grpc_adapter && \
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -v -o bin/grpc_adapter cmd/main.go ) && \
	docker build . -t layer7api/oauth_adapter:1.0

#certs:
#	mkdir grpc-test-key-cert
#	cp -R $(MIXER_REPO)/pkg/protobuf/yaml/testdata/auth/. grpc-test-key-cert

install:
	kubectl apply -f yaml/cluster_service.yaml
	kubectl apply -f $(MIXER_REPO)/adapter/grpc_adapter/testdata/attributes.yaml -f $(MIXER_REPO)/adapter/grpc_adapter/testdata/template.yaml
	kubectl apply -f $(MIXER_REPO)/adapter/grpc_adapter/testdata/grpcadapter.yaml
	kubectl apply -f $(MIXER_REPO)/adapter/grpc_adapter/testdata/operator_cfg.yaml

uninstall:
	kubectl delete -f yaml/cluster_service.yaml
	kubectl delete -f $(MIXER_REPO)/adapter/grpc_adapter/testdata/attributes.yaml -f $(MIXER_REPO)/adapter/grpc_adapter/testdata/template.yaml
	kubectl delete -f $(MIXER_REPO)/adapter/grpc_adapter/testdata/grpcadapter.yaml
	kubectl delete -f $(MIXER_REPO)/adapter/grpc_adapter/testdata/operator_cfg.yaml

check:
	@echo "gopath: $(GOPATH)"
	@echo "istio: $(ISTIO)"
	@echo "root folder: $(ROOT_FOLDER)"
	@echo "mixer repo: $(MIXER_REPO)"

clean:
	# Need the sudo because the go mods get delete protected
	#sudo rm -rf src
	rm -rf src
