FROM golang:1.19.3-alpine3.16 AS build 
ADD . /gcsim
WORKDIR /gcsim

RUN go mod download

WORKDIR /gcsim/backend/cmd/result
RUN go build -o /result

# the rest
FROM alpine:3.16.3 as backend
WORKDIR /
COPY --from=build /result /result
RUN ls
ENTRYPOINT ["/result"]