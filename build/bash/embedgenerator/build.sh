cd ui
yarn workspace @gcsim/embed build
cd ..
cd cmd/services/embedgenerator
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build .
docker build -f build/docker/embedgenerator/Dockerfile --no-cache --progress=plain