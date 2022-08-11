FROM golang:latest
# RUN mkdir /usr/src/app
WORKDIR /app

COPY go.mod go.sum ./
RUN ls
RUN go mod download
COPY embed embed 
RUN ls
RUN ls embed/

WORKDIR /app/embed

RUN go build .