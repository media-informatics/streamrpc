package service

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative temperature.proto

// server location:
const Port = 50000
const Host = "localhost"
