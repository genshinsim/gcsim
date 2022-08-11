FROM golang:latest
# RUN mkdir /usr/src/app
WORKDIR /app

COPY go.mod go.sum ./
RUN ls
RUN go mod download
COPY dbmanage dbmanage 
RUN ls
RUN ls dbmanage/

WORKDIR /app/dbmanage

RUN go build .