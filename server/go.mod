module github.com/bd878/live-connections/server

go 1.19

require (
	github.com/bd878/live-connections/disk v0.0.0-20230417105641-5eb252a44741
	github.com/bd878/live-connections/meta v0.0.0-20230415185003-1d802308ce63
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.5.0
	github.com/joho/godotenv v1.5.1
	google.golang.org/grpc v1.54.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace github.com/bd878/live-connections/disk/pkg/proto => ../disk/pkg/proto

replace github.com/bd878/live-connections/disk => ../disk
