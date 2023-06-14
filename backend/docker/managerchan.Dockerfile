FROM alpine:3.16.3 as backend
WORKDIR /
COPY backend/artifacts/managerchan /managerchan
RUN ls -la
ENTRYPOINT ["/managerchan"]