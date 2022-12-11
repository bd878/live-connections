proto-compile:
	protoc ../proto/*proto \
		--go_out=./api/v1/ \
		--go-grpc_out=./api/v1/ \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=../proto/

proto-rm:
	rm ./api/v1/*go

rm:
	make proto-rm
	rm -rf ../files
