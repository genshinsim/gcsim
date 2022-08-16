FROM golang:latest
# RUN mkdir /usr/src/app
WORKDIR /app

COPY go.mod go.sum ./
RUN ls
RUN go mod download
COPY pkg pkg
COPY dbmanage dbmanage 