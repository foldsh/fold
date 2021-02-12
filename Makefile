VERSION := $(shell go run ./cmd/version)
# The bin/$(1)/$(1) in the build commands looks odd, but it just makes it really easy
# to get a very small build context for docker.
define build-local
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/$(1)/$(1) ./cmd/$(1)
endef

define build-image
	docker build -t $(1) -f ./cmd/$(1)/$(2)/Dockerfile ./bin/$(1)
endef

define tag-image
	docker tag $(1):latest foldsh/$(1):$(2)
endef

define build-release-linux
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-d -s" -o bin/$(1)/$(1) ./cmd/$(1)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-d -s" -o bin/release/$(1)/$(1)-$(VERSION)-linux-amd64 ./cmd/$(1)
	tar czvf bin/release/$(1)/$(1)-$(VERSION)-linux-amd64.tar.gz bin/release/$(1)/$(1)-$(VERSION)-linux-amd64
endef

define build-release-darwin
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s" -o bin/release/$(1)/$(1)-$(VERSION)-darwin-amd64 ./cmd/$(1)
	tar czvf bin/release/$(1)/$(1)-$(VERSION)-darwin-amd64.tar.gz bin/release/$(1)/$(1)-$(VERSION)-darwin-amd64
endef

install:
	go install ./...
.PHONY: install

version:
	@echo $(VERSION)
.PHONY: version

# FOLD GATEWAY
foldgw: bin
	$(call build-local,foldgw)
.PHONY: foldgw

foldgw-release: bin
	$(call build-release-linux,foldgw)
.PHONY: foldgw

foldgw-image: foldgw-release
	$(call build-image,foldgw,images)
	$(call tag-image,foldgw,$(VERSION))
.PHONY: foldgw-image

# FOLD RUNTIME
foldrt: bin
	$(call build-local,foldrt)
.PHONY: foldrt

foldrt-release: bin
	$(call build-release-linux,foldrt)
	$(call build-release-darwin,foldrt)
.PHONY: foldrt-image

foldrt-image: foldrt-release
	$(call build-image,foldrt,images/alpine)
	$(call tag-image,foldrt,$(VERSION))
	$(call tag-image,foldrt,$(VERSION)-alpine)

	$(call build-image,foldrt,images/node)
	$(call tag-image,foldrt,$(VERSION)-node)
.PHONY: foldrt-image

# FOLD CTL
foldctl: bin
	$(call build-local,foldctl)
.PHONY: foldctl

foldctl-release: bin
	$(call build-release-linux,foldctl)
	$(call build-release-darwin,foldctl)
.PHONY: foldctl-release

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
	mockgen -source=ctl/project/container_api.go \
		-destination=ctl/project/mock_container_api_test.go \
		-package project_test

bin:
	mkdir -p bin

clean:
	rm -rf bin
.PHONY: install
