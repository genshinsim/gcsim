FROM alpine:3.16.3 as backend
COPY binary/db /db
RUN ls -la
RUN chmod +x /db
ENTRYPOINT ["/db"]