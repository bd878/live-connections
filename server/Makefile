proto-compile:
	protoc ../proto/*proto \
		--go_out=./proto/ \
		--go-grpc_out=./proto/ \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=../proto/

proto-rm:
	rm ./proto/*go
