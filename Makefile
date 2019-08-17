PACKAGE_PATH="github.com/chmking/horde"

.PHONY: protobuf
protobuf:
	@TMPDIR=$(shell mktemp -d)
	@protoc protobuf/public/public.proto -I protobuf --go_out=plugins=grpc:$(TMPDIR) && \
		cp $(TMPDIR)$(PACKAGE_PATH)/protobuf/public/* protobuf/public/
	@protoc protobuf/private/private.proto -I protobuf --go_out=plugins=grpc:$(TMPDIR) && \
		cp $(TMPDIR)$(PACKAGE_PATH)/protobuf/private/* protobuf/private/

.PHONY: build
build:
	go build -o build/cmd horde/main.go
	go build -o build/agent example/agent/main.go
	go build -o build/manager example/manager/main.go

.PHONY: test
test:
	ginkgo -r --randomizeAllSpecs --randomizeSuites --failFast --failOnPending --cover --trace --race --progress
	go test -v -covermode=count -coverprofile=coverage.out ./...
