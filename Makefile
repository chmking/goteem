.PHONY: protobuf
protobuf:
	@protoc horde.proto --go_out=plugins=grpc:.

.PHONY: test
test:
	ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --cover --trace --race --progress
