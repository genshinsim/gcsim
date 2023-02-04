FROM golang:1.19.3-alpine3.16 AS build 

RUN mkdir -p /gcsim
WORKDIR /gcsim

### download deps and caches
COPY go.mod .
COPY go.sum .
RUN go mod download

RUN ls /gcsim

## copy all files in
COPY backend /gcsim/backend
COPY pkg /gcsim/pkg

RUN ls /gcsim

WORKDIR /gcsim/backend/cmd/jadechamber
RUN go build -o /jadechamber

# the rest
FROM alpine:3.16.3 as backend
WORKDIR /
COPY --from=build /jadechamber /jadechamber
RUN ls
ENTRYPOINT ["/jadechamber"]