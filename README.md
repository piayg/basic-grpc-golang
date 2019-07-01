# Get Start, Golang with gRPC

## Install Tools
- [Goalng](https://golang.org/dl/) :)
- Protocol Buffers, [protoc](https://github.com/protocolbuffers/protobuf/releases) ช่วยให้แปลงจาก `.proto` เป็นภาษาที่ต้องการ
- protoc-gen-go ช่วยให้แปลงจาก `.proto` เป็น `.pb.go` ทำให้ง่ายต่อการใช้งานในภาษา `golang`
  ```bash
  $ go install github.com/golang/protobuf/protoc-gen-go
  ```

## VSCode Extensions
  - [Clang-Format](https://marketplace.visualstudio.com/items?itemName=xaver.clang-format) จัด Format ของไฟล์ Protocol Buffers `.proto` ต้องมี
    - [LLVM](https://releases.llvm.org/download.html), Clang-Format binary เป็นส่วนที่ `Clang-Format Extensions` รันเมื่อทำการ Format
  - [vscode-proto3](https://marketplace.visualstudio.com/items?itemName=zxh404.vscode-proto3) Autocomplete ของ `proto3`

## Start
### Protocol Buffers
ข้อมูลเพิ่มเติมจาก https://developers.google.com/protocol-buffers/docs/proto3

สร้างไฟล์ `echo.proto` ในโฟรเดอร์ `pb`
```proto3
syntax = "proto3";
package pb;

// EchoRequest is the request for echo.
message EchoRequest { string message = 1; }

// EchoResponse is the response for echo.
message EchoResponse { string message = 1; }

// Echo is the echo service.
service Echo {
  // UnaryEcho is unary echo.
  rpc UnaryEcho(EchoRequest) returns (EchoResponse) {}
}

service StreamingEcho {
  // ServerStreamingEcho is server side streaming.
  rpc ServerStreamingEcho(EchoRequest) returns (stream EchoResponse) {}
  // ClientStreamingEcho is client side streaming.
  rpc ClientStreamingEcho(stream EchoRequest) returns (EchoResponse) {}
  // BidirectionalStreamingEcho is bidi streaming.
  rpc BidirectionalStreamingEcho(stream EchoRequest)
      returns (stream EchoResponse) {}
}
```

โค้ตจาก https://github.com/grpc/grpc-go/blob/master/examples/features/proto/echo/echo.proto

จากนั้นทำการแปลงไฟล์ `.proto` เพื่อใช้สำหรับ `golang`
```bash
$ protoc --go_out=plugins=grpc:. pb/*.proto
```

จะได้ไฟล์ `.pb.go` ออกมาสำหรับใช้ในภาษา `golang` ทั้ง Client/Server

### Client/Server
#### Echo Service
##### Echo Server
สร้างไฟล์ `main.go` ในโฟรเดอร์ `server` สำหรับเขียนโปรแกรมส่วน Server พร้อมกับโค้ตเริ่มต้น Server
```golang
package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

  // Some work

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
```

แล้วกำหนด type ขึ้นมาใหม่ เพื่อใช่สำหรับระบุการทำงานของ service โดยใน struct จะมีค่าอื่นหรือไม่มีก็ได้
```golang
type handleServer struct{}
```

ผูก type ที่สร้างมาใหม่และ gRPC server เข้ากับ EchoServer จะเป็นการกำหนดให้เมื่อ Client เรียกใช้การทำงานของ Echo service แล้วตัว gRPC server จะไปเรียกใช้ฟังก์ชั่นของ ``handleServer``
```golang
pb.RegisterEchoServer(s, handle)
```

RegisterEchoServer จะไม่สำเร็จถ้าไม่กำหนดฟังก์ชั่นของ ``handleServer`` ให้เป็นไปตามที่ ``interface EchoServer`` ต้องการ สามารถดูได้ใน ``pb/echo.pb.go`` ซึ่งก็คือ
```golang
func (hs handleServer) UnaryEcho(context.Context, *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{Message: "I'm Groot"}, nil
}
```
เท่านี้การตั้ง ``Server gRPC`` ในส่วนของ ``service Echo`` ก็เสร็จแล้ว

โค้ตทั้งหมด
```golang
package main

import (
	"context"
	"log"
	"net"

	"app/pb"

	"google.golang.org/grpc"
)

type handleServer struct{}

func (hs handleServer) UnaryEcho(context.Context, *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{Message: "I'm Groot"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	handle := &handleServer{}
	pb.RegisterEchoServer(s, handle)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
```

##### Echo Client
สร้างไฟล์ `main.go` ในโฟรเดอร์ `Client` สำหรับเขียนโปรแกรมส่วน Server พร้อมกับโค้ตเริ่มต้น Client

......> https://github.com/payidalasgo/basic-grpc-golang