APP?=ns_bridge

GOOS?=linux
GOARCH?=amd64
GO_VERSION?=1.12.1

PROJECT?=gitlab.zlibs.com/${APP}
VERSION?=$(shell git describe --tags --always)


default: build

.PHONY: build
build: clean
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags "-s -w -X ${PROJECT}/app.version=${VERSION} -X '${PROJECT}/app.buildDate=`date`' " -o ${APP}

.PHONY: build-local
build-local: clean agent govet swagger build

.PHONY: .govet
govet: 
	go vet . && go fmt ./... && \
	(if [[ "$(gofmt -d $(find . -type f -name '*.go' -not -path "./vendor/*" -not -path "./tests/*" -not -path "./assets/*"))" == "" ]]; then echo "Good format"; else echo "Bad format"; exit 33; fi);

.PHONY: .out
out: 
	protoc -I/usr/local/include -I. \
		-I ${GOPATH}/src \
		-I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:. \
		dns_pb/*.proto

.PHONY: .gateway
gateway: out
	protoc -I/usr/local/include -I. \
		-I  ${GOPATH}/src \
		-I  ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true:. \
		dns_pb/*.proto

.PHONY: .swagger
swagger: gateway
	protoc -I/usr/local/include -I. \
		-I ${GOPATH}/src \
		-I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--swagger_out=logtostderr=true:. \
		dns_pb/*.proto

.PHONY: agent
agent: 
	mkdir -p vendor/go-agent/blueware \
	&& cp -r ${GOPATH}/src/go-agent/blueware vendor/go-agent

.PHONY: .clean
clean: 
	rm -fr ${APP}
