FROM golang:1.17-buster AS src

WORKDIR /go/src/app/

# Copy dependencies first to take advantage of Docker caching
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

# Insert version using git tag and latest commit hash
# Build Go Binary
RUN set -ex; \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ./server ./cmd/opa/main.go;

# (Optional) Compress binary using upx https://upx.github.io/
#RUN apt update
#RUN apt install -y upx-ucl
#RUN upx /go/src/app/server/main

FROM scratch
LABEL com.gmhafiz.maintainers="User <author@example.com>"

WORKDIR /App

COPY --from=src /go/src/app/server /App/server

# Docker cannot copy hidden .env file. So in Makefile, we make a copy of it.
COPY --from=src /go/src/app/env.prod /App/.env

COPY --from=src /go/src/app/third_party/opa/rbac.rego /App/third_party/opa/rbac.rego


EXPOSE 9091
EXPOSE 8091

ENTRYPOINT ["/App/server/main"]
#ENTRYPOINT ["ls",  "-la", "/App/server"]
