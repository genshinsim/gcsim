//go:generate sh scripts/generate/generate_protos.sh
//go:generate sh scripts/generate/build_preview.sh
//go:generate go run pipeline/cmd/pipeline/main.go
//go:generate go run scripts/generate/bsontags/main.go -dir backend/pkg/services/db -verbose
package gcsim
