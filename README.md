Simple gRPC demo with bi-directional stream and simple RPC

### Start Server

```
go run cmd/server/main.go
```

### Run Client

```
go run cmd/client/main.go
```

### Generate Proto

```
protoc -I proto/ --go_out=plugins=grpc:lib/proto proto/api.proto
```
