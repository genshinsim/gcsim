FROM alpine:3.16.3 as backend
COPY backend/artifacts/db /db
RUN ls -la
ENTRYPOINT ["/db"]