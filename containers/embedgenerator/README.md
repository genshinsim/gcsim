# `embedgenerator`

Example docker-compose
```yaml
# version: "3.8"

services:
  # cloudflare to expose the embed generator to the internet
  cloudflared:
    container_name: cloudflared
    image: cloudflare/cloudflared:latest
    command: "tunnel --no-autoupdate run --token $CF_TUNNEL_TOKEN"
    volumes:
      - ./cloudflared:/etc/cloudflared
    environment:
      - CF_TUNNEL_TOKEN
    restart: unless-stopped
  # redis is used for ensure multiple request for same image resource concurrently will not spawn multiple
  # generations at the same time
  redis:
    # using redis-stack here for easy debug as it comes with web ui
    # however do not need this in production
    image: redis/redis-stack:latest 
    ports:
      - 6379:6379
      - 8001:8001
  # the go-rod container houses a custom manager and chrome instance for go-rod.
  # see more here: https://github.com/go-rod/rod
  rod:
    image: ghcr.io/go-rod/rod
    ports:
      - 7317:7317
  preview:
    container_name: preview
    depends_on:
      - redis
      - rod
    build:
      context: .
    environment:
      - HOST=localhost # <optional> host to listen up, this should be left out if listening to all host; defaults to blank
      - PORT=7777 # <optional> port to listen to, defaults to 3000
      - LAUNCHER_URL=ws://rod:7317 # address of the go-rod container for connecting to chrome instance
      - AUTH_KEY=<insert some secure key> # <optional> if set, incoming http request should have header X-CUSTOM-AUTH-KEY set with this key
      - PREVIEW_URL=http://preview:7777 # this is the callback url for the rod container to screenshot
      - STATIC_ASSETS=/dist # this is where the default static assets (i.e. index.html) is located; is /dist by default
      - ASSETS_PATH=/assets/assets
      - REDIS_URL=redis:6379
      
      - PROXY_TO
      - AUTH_KEY
      - ASSETS_PAT_TOKEN
    ports:
      - 7777:7777
    restart: unless-stopped
```