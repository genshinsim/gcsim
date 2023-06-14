FROM alpine:3.16.3 as backend
WORKDIR /
COPY binary/share /share
RUN ls -la
RUN chmod +x /share
ENTRYPOINT ["/share"]