FROM golang:1.19.3-alpine3.16 AS build
# RUN mkdir /usr/src/app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY pkg pkg
COPY cmd cmd
RUN go mod tidy
RUN cat go.mod
RUN go mod download

WORKDIR /app/cmd/db
RUN go build -o /db

FROM alpine:3.16.3

WORKDIR /

COPY --from=build /db /db

ENTRYPOINT ["/db"]