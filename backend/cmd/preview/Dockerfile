from golang:1.19.3-bullseye as build

RUN mkdir -p /gcsim
WORKDIR /gcsim

### download deps and caches
COPY go.mod .
COPY go.sum .
RUN go mod download

## copy all files in
COPY backend /gcsim/backend
COPY pkg /gcsim/pkg

WORKDIR /gcsim/backend/cmd/preview
RUN go build -o /preview

FROM chromedp/headless-shell:latest 
WORKDIR /
COPY --from=build /preview /preview
ENTRYPOINT ["/preview"] 