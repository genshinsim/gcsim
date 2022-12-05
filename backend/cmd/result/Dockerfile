FROM golang:1.19.3-alpine3.15
# RUN mkdir /usr/src/app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY pkg pkg
COPY cmd cmd
RUN go mod tidy
RUN cat go.mod
RUN go mod download

WORKDIR /app/cmd/result
RUN go build -o /result

ENTRYPOINT ["/result"]