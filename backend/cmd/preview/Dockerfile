from golang:1.19.3-bullseye as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY pkg pkg
COPY cmd cmd
RUN go mod tidy
RUN cat go.mod
RUN go mod download

WORKDIR /app/cmd/preview
RUN go build -o /preview


from chromedp/headless-shell:latest

WORKDIR /
COPY --from=build /preview /preview

RUN ls

ENTRYPOINT ["/preview"]