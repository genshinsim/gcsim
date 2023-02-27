FROM golang:1.19.3-alpine3.16 AS build 

ARG COMMIT
LABEL commit=${COMMIT}
ENV COMMIT_SHA=${COMMIT}

RUN mkdir -p /gcsim
WORKDIR /gcsim

### download deps and caches
# COPY .git /gcsim/.git
COPY go.mod .
COPY go.sum .
RUN go mod download

## copy all files in
COPY backend /gcsim/backend
COPY pkg /gcsim/pkg

WORKDIR /gcsim/backend/cmd/db
RUN go build -ldflags "-X main.sha1ver=${COMMIT_SHA}"  -o /db
# RUN go build -o /db

# the rest
FROM alpine:3.16.3 as backend
WORKDIR /
COPY --from=build /db /db
RUN ls
ENTRYPOINT ["/db"]