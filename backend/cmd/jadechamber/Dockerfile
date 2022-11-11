FROM golang:latest
# RUN mkdir /usr/src/app
WORKDIR /app

COPY go.mod go.sum ./
RUN ls
RUN go mod download
COPY pkg pkg
COPY cmd/jadechamber cmd/jadechamber

WORKDIR /app/cmd/jadechamber
RUN go build -o /jadechamber

ENTRYPOINT ["/jadechamber"]