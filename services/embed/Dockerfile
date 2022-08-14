FROM golang:latest
# RUN mkdir /usr/src/app
WORKDIR /app

COPY go.mod go.sum ./
RUN ls
RUN go mod download
COPY pkg pkg
COPY embed embed

WORKDIR /app/embed
RUN go build .

FROM python:3
COPY --from=0 /app /app
WORKDIR /app/embed/scripts

RUN ls
RUN pip install -r ./requirements.txt