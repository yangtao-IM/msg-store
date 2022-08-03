#!/bin/bash

#生成普通grpc-pb文件，有service的client。生成_grpc.pb.go
protoc.exe -I.  --go-grpc_out=paths=source_relative:.  *.proto

#生成pb文件，直接结构体的go文件。生成.pb.go
protoc.exe -I.  --go_out=paths=source_relative:.  *.proto
go mod tidy
