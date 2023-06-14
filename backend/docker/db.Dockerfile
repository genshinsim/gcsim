FROM alpine:3.16.3 as backend
COPY binary/db /db
RUN ls -la
ENTRYPOINT ["/db"]