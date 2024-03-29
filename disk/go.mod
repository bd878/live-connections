module github.com/bd878/live-connections/disk

go 1.19

require (
	github.com/bd878/live-connections/meta v0.0.0-20230415185003-1d802308ce63
	github.com/joho/godotenv v1.5.1
	google.golang.org/grpc v1.54.0
	google.golang.org/protobuf v1.30.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
)

replace github.com/bd878/live-connections/disk/pkg/fd => ./pkg/fd
