FROM golang:1.19.3-alpine3.16 AS build 
ADD . /gcsim
WORKDIR /gcsim

RUN go mod download

WORKDIR /gcsim/backend/cmd/jadechamber
RUN go build -o /jadechamber

# the rest
FROM alpine:3.16.3 as backend
WORKDIR /
COPY --from=build /jadechamber /jadechamber
RUN ls
ENTRYPOINT ["/jadechamber"]