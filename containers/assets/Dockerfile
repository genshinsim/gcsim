FROM debian:stable-slim
RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates
RUN update-ca-certificates

WORKDIR /
COPY cmd/services/assets/assets /assets
RUN chmod +x /assets
RUN ls -la
ENTRYPOINT ["/assets"]