FROM golang:latest
# RUN mkdir /usr/src/app
WORKDIR /app

COPY go.mod go.sum ./
RUN ls
RUN go mod download
COPY pkg pkg
COPY embed embed

WORKDIR /app/embed
RUN go build -o /embed-server 

RUN ls

FROM python:3
WORKDIR /
COPY --from=0 /embed-server /embed-server

RUN ls

COPY embed/scripts /scripts

WORKDIR /scripts

RUN ls
RUN pip install -r ./requirements.txt

WORKDIR /

ENTRYPOINT ["/embed-server"]
