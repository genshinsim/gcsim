from golang:1.19.3-bullseye as build
ADD . /gcsim

WORKDIR /gcsim

RUN go mod download

WORKDIR /gcsim/backend/cmd/preview
RUN go build -o /preview

FROM chromedp/headless-shell:latest 
WORKDIR /
COPY --from=build /preview /preview
ENTRYPOINT ["/preview"] 