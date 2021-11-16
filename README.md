# Introduction

[Open Policy Agent](https://www.openpolicyagent.org/) (OPA) microservice in gRPC with http gateway implementing Role Based Access Control (RBAC).

# Protobuf gRPC

Needs
1. Protobuf compiler
2. Language-specific plugin

Then compile with

    protoc api/v1/*.proto \
	--go_out=. \
	--go-grpc_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	--proto_path=.

Or simply use `Makefile` to compile by typing

    make proto


## 1. Compiler

    wget https://github.com/protocolbuffers/protobuf/releases/download/v3.18.0/protoc-3.18.0-linux-x86_64.zip
    unzip protoc-3.18.0-linux-x86_64.zip -d /usr/local/protobuf
    echo 'export PATH="$PATH:/usr/local/protobuf/bin"' >> ~/.bashrc


## 2. Plugin

     go get google.golang.org/protobuf/...@v1.27.1


## 3. Tools

    make install

## Database

    docker run --name postgres14 -p 5432:5432 -e POSTGRES_PASSWORD=password -e POSTGRES_USER=user -e POSTGRES_DB=opa -d postgres:14

    # or
    docker-compose -d postgres up 


# Run

    make dev
    # or
    go run cmd/opa/main.go
    # or
    docker-compose up -d server

# Benchmark

OPA policy is cached instead of hitting a database. As a result, evaluation is
extremely fast, capable of handling **24419** requests per second.

CPU: AMD 3600 3.6GHz

```shell
wrk -t2 -d60 -c200 -s wrk_post.lua 'http://localhost:8091/api/v1/opa/check'
Running 1m test @ http://localhost:8091/api/v1/opa/check
  2 threads and 200 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    14.03ms   17.34ms 191.75ms   86.19%
    Req/Sec    12.28k     1.13k   14.93k    76.42%
  1466847 requests in 1.00m, 173.46MB read
Requests/sec:  24419.12
Transfer/sec:      2.89MB
```

# TODO

 - [ ] Implement OPA Bundle API to update policy from external/internal database with schema of `users`, `roles` and `user_roles`.
 - [ ] Database to trigger API to update policy everytime any of the tables are modified.
