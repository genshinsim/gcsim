FROM golang:latest
# RUN mkdir /usr/src/app
WORKDIR /app

COPY go.mod go.sum ./
RUN ls
RUN go mod download
COPY pkg pkg
COPY embed embed 

FROM python:3
WORKDIR /app/embed/scripts

RUN pip install -r requirements.txt