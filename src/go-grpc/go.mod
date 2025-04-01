module github.com/guodongq/quickstart/go-grpc

go 1.24.1

require (
	github.com/guodongq/quickstart/proto v1.2.3
	google.golang.org/grpc v1.71.0
)

replace github.com/guodongq/quickstart/proto => ../../proto

require (
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)
