FROM alpine:3.16.3 as backend
WORKDIR /
COPY backend/artifacts/jadechamber /jadechamber
RUN ls -la
ENTRYPOINT ["/jadechamber"]