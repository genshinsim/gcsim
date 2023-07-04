FROM alpine:3.16.3 as backend
WORKDIR /
COPY binary/managerchan /managerchan
RUN ls -la
RUN chmod +x /managerchan
ENTRYPOINT ["/managerchan"]