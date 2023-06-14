FROM chromedp/headless-shell:latest 
WORKDIR /
COPY binary/preview /preview
RUN ls -la
RUN chmod +x /preview
ENTRYPOINT ["/preview"]