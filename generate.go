//go:generate protoc --go-grpc_out=. proto/backend/result.proto
//go:generate protoc --go_out=. proto/backend/result.proto
//go:generate sh scripts/preview/build.sh
package gcsim
