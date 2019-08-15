.PHONY: protobuf
protobuf:
	@protoc protobuf/public/public.proto --go_out=plugins=grpc:.
	@protoc protobuf/private/private.proto --go_out=plugins=grpc:.

.PHONY: build
build:
	go build -o build/cmd horde/main.go
	go build -o build/agent example/main.go
	go build -o build/manager manager/main.go

.PHONY: test
test:
	ginkgo -r --randomizeAllSpecs --randomizeSuites --failFast --failOnPending --cover --trace --race --progress
