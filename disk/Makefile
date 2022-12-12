proto-compile:
	cp ../proto/disk.proto ./api/v1/disk_server.proto
	sed -r -i 's/^option go_package.*/option go_package = "github.com\/teralion\/live-connections\/disk\/api\/v1";/g' ./api/v1/disk_server.proto
	sed -r -i 's/^package .*/package disk_server;/g' ./api/v1/disk_server.proto

	protoc ./api/v1/*.proto \
		--go_out=./api/v1/ \
		--go-grpc_out=./api/v1/ \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=./api/v1/

proto-rm:
	rm -rf ./api/v1/*

rm:
	make proto-rm
	rm -rf ../files
