FROM alpine:3.16.3 as backend
WORKDIR /
COPY backend/artifacts/share /share
RUN ls -la
ENTRYPOINT ["/share"]