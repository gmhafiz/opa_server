.PHONY: proto check run install
proto:
	protoc api/v1/*.proto \
	--go_out=. \
	--go-grpc_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	--grpc-gateway_out . --grpc-gateway_opt paths=source_relative
	--proto_path=.

check:
	go mod vendor
	go fmt ./...
	go vet ./...
	golangci-lint run
	gosec -quiet ./...

dev:
	go mod tidy
	go run cmd/opa/main.go

install:
	go install \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
        google.golang.org/protobuf/cmd/protoc-gen-go \
        google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u golang.org/x/tools/...
	go install github.com/securego/gosec/v2/cmd/gosec@latest
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

docker_build:
	cp .env env.prod
	docker build -t opa/server -f Dockerfile .
	rm env.prod

docker_run:
	docker run -p 9091:9091 -p 8091:8091  --net=host --rm -it --name opa_container opa/server

benchmark:
	wrk -t2 -d1 -c200 'http://localhost:8091/api/v1/opa/check' -s wrk_post.lua
