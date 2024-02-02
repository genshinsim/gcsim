FROM alpine:3.16.3 as backend
WORKDIR /
COPY binary/notification /notification
RUN ls -la
RUN chmod +x /notification
ENTRYPOINT ["/notification"]