# The bin/$(1)/$(1) in the build commands looks odd, but it just makes it really easy
# to get a very small build context for docker.
define build-local
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/$(1)/$(1) ./cmd/$(1)
endef

define build-image
	docker build -t $(1) -f ./cmd/$(1)/images/Dockerfile ./bin/$(1)
endef

define build-release
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-d -s" -o bin/$(1)/$(1) ./cmd/$(1)
endef

# FOLD GATEWAY
foldgw: bin
	$(call build-local,foldgw)
.PHONY: foldgw

foldgw-release: bin
	$(call build-release,foldgw)
.PHONY: foldgw

foldgw-image: foldgw-release
	$(call build-image,foldgw)
.PHONY: foldgw-image

# FOLD RUNTIME
foldrt: bin
	$(call build-local,foldrt)
.PHONY: foldrt

foldrt-release: bin
	$(call build-release,foldrt)
.PHONY: foldrt-image

foldrt-image: foldrt-release
	$(call build-image,foldrt)
.PHONY: foldrt-image

# FOLD CTL
foldctl: bin
	$(call build-local,foldctl)
.PHONY: foldctl

foldctl-release: bin
	$(call build-release,foldctl)
.PHONY: foldctl-release

foldctl-image: foldctl-release
	$(call build-image,foldctl)
.PHONY: foldctl-image

protoc:
	protoc --proto_path=proto \
		--go_out=runtime/supervisor/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=runtime/supervisor/pb \
		--go-grpc_opt=paths=source_relative \
		proto/ingress.proto
	protoc --proto_path=proto \
		--go_out=manifest \
		--go_opt=paths=source_relative \
		proto/manifest.proto
.PHONY: protoc

genmocks:
	mockgen -source=ctl/container/docker_client.go \
		-destination=ctl/container/mock_docker_client.go \
		-package container

bin: ./bin
	mkdir -p bin

clean:
	rm -rf bin
.PHONY: install
