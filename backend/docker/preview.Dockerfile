FROM alpine:3.16.3 as backend
WORKDIR /
COPY binary/preview /preview
RUN ls -la
RUN chmod +x /preview
ENTRYPOINT ["/preview"]