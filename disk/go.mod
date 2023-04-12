module github.com/bd878/live-connections/disk

go 1.19

require (
	github.com/bd878/live-connections/meta v0.0.0-20230408143603-08aac6a0bc9c
	github.com/joho/godotenv v1.4.0
	google.golang.org/grpc v1.51.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/text v0.4.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
)

replace github.com/bd878/live-connections/meta => ../meta
