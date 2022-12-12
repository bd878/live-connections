proto-compile:
	cp ../proto/disk.proto ./proto/disk/disk_client.proto
	sed -r -i 's/^option go_package.*/option go_package = "github.com\/teralion\/live-connections\/server\/proto\/disk";/g' ./proto/disk/disk_client.proto
	sed -r -i 's/^package .*/package disk_client;/g' ./proto/disk/disk_client.proto

	protoc ./proto/disk/*.proto \
		--go_out=./proto/disk/ \
		--go-grpc_out=./proto/disk/ \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=./proto/disk/

proto-rm:
	rm -rf ./proto/disk/*
