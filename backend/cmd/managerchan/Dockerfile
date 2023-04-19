FROM golang:1.19.3-alpine3.16 AS build 

RUN mkdir -p /gcsim
WORKDIR /gcsim

### download deps and caches
COPY go.mod .
COPY go.sum .
RUN go mod download

## copy all files in
COPY backend /gcsim/backend
COPY pkg /gcsim/pkg

WORKDIR /gcsim/backend/cmd/managerchan
RUN go build -o /managerchan

FROM alpine:3.16.3 as backend
WORKDIR /
COPY --from=build /managerchan /managerchan
RUN ls
ENTRYPOINT ["/managerchan"]