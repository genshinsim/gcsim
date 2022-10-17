FROM golang:latest
# RUN mkdir /usr/src/app
WORKDIR /app

COPY go.mod go.sum ./
RUN ls
RUN go mod download
COPY assets assets

WORKDIR /app/assets
RUN go build -o /assets-server 

ENTRYPOINT ["/assets-server"]
