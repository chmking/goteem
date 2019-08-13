.PHONY: protobuf
protobuf:
	@protoc protobuf/public/public.proto --go_out=plugins=grpc:.
	@protoc protobuf/private/private.proto --go_out=plugins=grpc:.

.PHONY: test
test:
	ginkgo -r --randomizeAllSpecs --randomizeSuites --failFast --failOnPending --cover --trace --race --progress
