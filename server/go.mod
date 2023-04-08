module github.com/teralion/live-connections/server

go 1.19

require (
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.5.0
	github.com/joho/godotenv v1.4.0
	github.com/teralion/live-connections/meta v0.0.0-20230322104303-2e53c52f017d
	google.golang.org/grpc v1.51.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/teralion/live-connections/disk v0.0.0-20230408141444-9c53e6156fa3 // indirect
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/text v0.4.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
)

replace github.com/teralion/live-connections/meta => ../meta
