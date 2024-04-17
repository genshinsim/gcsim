FROM debian:stable-slim
WORKDIR /
ADD ui/packages/embed/dist /dist
COPY cmd/services/embedgenerator/embedgenerator /embedgenerator
RUN chmod +x /embedgenerator
RUN ls -la
ENTRYPOINT ["embedgenerator"]