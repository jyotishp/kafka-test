SHELL=/bin/bash -o pipefail

install-proto:
	sudo apt update && sudo apt install -y protobuf-compiler
	go get -u -v github.com/golang/protobuf/protoc-gen-go

compile-proto: pre-build proto-dependencies
	protoc \
	-I pkg \
	-I ./vendor/github.com/grpc-ecosystem/grpc-gateway \
	--go_out=plugins=grpc:. \
	--grpc-gateway_out=logtostderr=true:. \
	--swagger_out=logtostderr=true,allow_merge=true,merge_file_name=app:swagger-ui \
	--proto_path pkg/proto tweets.proto

pre-build:
	mkdir -p build

proto-dependencies:
	go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

go-dependencies:
	go mod vendor

build: pre-build compile-proto
	go build -o server cmd/server/main.go

clean:
	rm -rf build
	rm -f swagger-ui/app.swagger.json

tests:
	go test ./... -v -coverprofile c.out
	go tool cover -html=c.out -o coverage.html

proto: compile-proto

run:
	go run cmd/server/main.go
