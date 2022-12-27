FROM golang:1.19.3-alpine3.16 AS build 
ADD . /gcsim
WORKDIR /gcsim

RUN go mod download

WORKDIR /gcsim/backend/cmd/db
RUN go build -o /db

# the rest
FROM alpine:3.16.3 as backend
WORKDIR /
COPY --from=build /db /db
RUN ls
ENTRYPOINT ["/db"]