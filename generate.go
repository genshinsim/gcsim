//go:generate protoc --go-grpc_out=. proto/backend/result.proto
//go:generate protoc --go_out=. proto/backend/result.proto
//go:generate protoc --go-grpc_out=. proto/backend/preview.proto
//go:generate protoc --go_out=. proto/backend/preview.proto
package gcsim
