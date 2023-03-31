proto-compile:
	mkdir ./pkg/proto
	cp ../proto/disk.proto ./pkg/proto/disk.proto
	sed -r -i 's/^option go_package.*/option go_package = "github.com\/teralion\/live-connections\/server\/rpc\/proto";/g' ./pkg/proto/disk.proto
	sed -r -i 's/^package .*/package rpc;/g' ./pkg/proto/disk.proto

	protoc ./pkg/proto/*.proto \
		--go_out=./pkg/proto/ \
		--go-grpc_out=./pkg/proto/ \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=./pkg/proto
