build: protoc
	mkdir -p bin
	CGO_ENABLED=0 go build -o bin ./...
.PHONY: build

install: protoc
	mkdir -p bin
	CGO_ENABLED=0 go install ./...
.PHONY: build

release: protoc
	mkdir -p bin
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-d -s" -o bin ./...
.PHONY: build

install: protoc
	go install ./...
.PHONY: install

images: protoc
	docker build -t foldrt ./docker/Dockerfile
.PHONY: docker

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

clean:
	rm -rf bin
.PHONY: install
