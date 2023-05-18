module github.com/bd878/live-connections/server

go 1.19

require (
	github.com/bd878/live-connections/disk v0.0.0-20230517103142-1cbdb362a3a9
	github.com/bd878/live-connections/meta v0.0.0-20230517103142-1cbdb362a3a9
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.5.0
	github.com/joho/godotenv v1.5.1
	google.golang.org/grpc v1.55.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace github.com/bd878/live-connections/disk/pkg/proto => ../disk/pkg/proto

replace github.com/bd878/live-connections/disk => ../disk

replace github.com/bd878/live-connections/server/pkg/mock => ./pkg/mock

replace github.com/bd878/live-connections/server/pkg/rpc => ./pkg/rpc
