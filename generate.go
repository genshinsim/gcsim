//go:generate sh -c "protoc --go_out=module=github.com/genshinsim/gcsim:. --go-grpc_out=module=github.com/genshinsim/gcsim:. protos/**/*.proto"
//go:generate go run scripts/generate/bsontags/main.go -dir backend/pkg/services/db -verbose
//go:generate go run scripts/generate/bsontags/main.go -dir pkg/model -verbose
//go:generate sh scripts/generate/build_preview.sh
package gcsim
