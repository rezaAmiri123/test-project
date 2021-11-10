.PHONY: proto
proto:
	protoc \
		--go_out=internal/common/genproto/users --go_opt=paths=source_relative \
		--go-grpc_out=internal/common/genproto/users --go-grpc_opt=paths=source_relative \
		--proto_path=api/protobuf api/protobuf/users.proto
