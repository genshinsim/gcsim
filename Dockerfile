FROM golang:1.19.3-alpine3.16 AS build
# RUN mkdir /usr/src/app
ADD . /gcsim

WORKDIR /gcsim

RUN go mod tidy
RUN go mod download

WORKDIR /gcsim/backend/cmd/jadechamber
RUN go build -o /jadechamber

WORKDIR /gcsim/backend/cmd/preview
RUN go build -o /preview

WORKDIR /gcsim/backend/cmd/result
RUN go build -o /result

WORKDIR /gcsim/backend/cmd/db
RUN go build -o /db

RUN ls

# preview
FROM chromedp/headless-shell:latest as preview
WORKDIR /
COPY --from=build /preview /preview
RUN ls
ENTRYPOINT ["/preview"]

# the rest
FROM alpine:3.16.3 as backend
WORKDIR /
COPY --from=build /jadechamber /jadechamber
COPY --from=build /result /result
COPY --from=build /db /db
RUN ls