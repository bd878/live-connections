proto-compile:
	cp ../proto/disk.proto ./pkg/proto/disk.proto
	sed -r -i 's/^option go_package.*/option go_package = "github.com\/bd878\/live-connections\/disk\/pkg\/proto";/g' ./pkg/proto/disk.proto
	sed -r -i 's/^package .*/package disk;/g' ./pkg/proto/disk.proto

	protoc ./pkg/proto/*.proto \
		--go_out=./pkg/proto/ \
		--go-grpc_out=./pkg/proto/ \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=./pkg/proto/

proto-rm:
	rm -rf ./pkg/proto/*

rm:
	make proto-rm
	rm -rf ../files
