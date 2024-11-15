cd ui
yarn workspace @gcsim/embed build
cd ..
ls -lh ./ui/packages/embed/dist
cd cmd/services/embedgenerator
GOOS=linux GOARCH=amd64 go build .