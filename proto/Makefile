compile:
	protoc ./v1/*proto \
		--go_out=./v1/ \
		--go-grpc_out=./v1/ \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=./v1/