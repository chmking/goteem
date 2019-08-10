.PHONY: protobuf
protobuf:
	@protoc horde.proto --go_out=plugins=grpc:.
