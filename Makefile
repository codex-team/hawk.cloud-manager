proto:
	protoc --go_out=plugins=grpc:./api/protobuf -I ./third_party/googleapis -I ./api/protobuf --go_opt=paths=source_relative api/protobuf/*.proto

build:
	go build root.go -o manager

init:
	git submodule update --init

prepare: init proto