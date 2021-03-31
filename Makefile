GOOS=$(shell go env GOHOSTOS)
GOARCH=$(shell go env GOARCH)

install:
	go install ./...
.PHONY: install

version:
	@echo $(shell go run ./cmd/version)
.PHONY: version

foldctl: bin
	./scripts/build.sh --bin foldctl --os ${GOOS} --arch ${GOARCH}
.PHONY: foldctl

foldctl-release: bin
	./scripts/build.sh --bin foldctl --os "darwin linux" --arch "amd64" --tar
.PHONY: foldctl-release

foldrt: bin
	./scripts/build.sh --bin foldrt --os "${GOOS}" --arch "${GOARCH}" --images
.PHONY: foldrt

foldrt-release: bin
	./scripts/build.sh --bin foldrt --os "linux" --arch "amd64" --tar --images --latest-tag "alpine"
.PHONY: foldrt-release

foldgw: bin
	./scripts/build.sh --bin foldgw --os "${GOOS}" --arch "${GOARCH}" --images
.PHONY: foldgw

foldgw-release: bin
	./scripts/build.sh --bin foldgw --os "linux" --arch "amd64" --tar --images --latest-tag "scratch"
.PHONY: foldgw-release

local-release: foldctl-release foldrt-release foldgw-release
	@echo Finished building all binaries and image for release $(shell go run ./cmd/version)
.PHONY: release

publish: local-release
	./scripts/publish.sh
.PHONY: publish

protoc:
	protoc --proto_path=proto \
		--go_out=runtime/transport/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=runtime/transport/pb \
		--go-grpc_opt=paths=source_relative \
		proto/ingress.proto
	protoc --proto_path=proto \
		--go_out=manifest \
		--go_opt=paths=source_relative \
		proto/manifest.proto proto/http.proto
.PHONY: protoc

genmocks:
	mockgen -source=ctl/container/docker_client.go \
		-destination=ctl/container/mock_docker_client.go \
		-package container
	mockgen -source=ctl/project/container_api.go \
		-destination=ctl/project/mock_container_api_test.go \
		-package project_test
	mockgen -source=runtime/service.go \
		-destination=runtime/service_mocks_test.go \
		-package runtime_test

bin:
	mkdir -p bin

clean:
	rm -rf bin
.PHONY: install

install-gotools:
	mkdir -p .gotools
	cd .gotools && if [[ ! -f go.mod ]]; then \
		go mod init fold-tools ; \
	fi
	cd .gotools && go get -v github.com/golang/mock/mockgen@v1.5.0 github.com/mitchellh/gox github.com/mitchellh/gon
.PHONY: gotools
