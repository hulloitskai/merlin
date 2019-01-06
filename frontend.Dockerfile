##################################################
## BUILDER IMAGE
##################################################
FROM node:10-alpine AS builder

## Build frontend:
## Copy source, install dependencies.
WORKDIR /build/
COPY ./client/ .
RUN yarn install

## Build from source.
ENV NODE_ENV="production"
RUN yarn build


## Build Caddy entrypoint:
## Install dependencies.
RUN apk add bash upx

## Install Caddy (with plugins) from bash script.
ENV PLUGINS=""
RUN wget -qO- https://getcaddy.com | bash -s personal $PLUGINS

## Move Caddy to /build/.
RUN mv /usr/local/bin/caddy .

## Compress Caddy using UPX.
RUN upx -9 caddy


##################################################
## PRODUCTION IMAGE
##################################################
FROM alpine:3.8 AS production

ARG BUILD_VERSION="unset"

## Labels:
LABEL maintainer="Steven Xie <hello@stevenxie.me>"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name = "stevenxie/merlin-frontend"
LABEL org.label-schema.url="https://merlin.stevenxie.me/"
LABEL org.label-schema.description="Merlin Frontend Server"
LABEL org.label-schema.vcs-url="https://github.com/stevenxie/merlin"
LABEL org.label-schema.version="$BUILD_VERSION"

## Install external dependencies.
RUN apk add ca-certificates

## Copy Caddy entrypoint, frontend, and config files:
COPY --from=builder /build/caddy /usr/bin/caddy
COPY --from=builder /build/dist/ /srv/
COPY ./Caddyfile /etc/Caddyfile

## Setup Caddy environment and persistent volume:
ENV CADDYPATH=/etc/.caddy
VOLUME /etc/.caddy

## Serve from /srv using default index.html.
WORKDIR /srv

## Define healthcheck.
COPY ./scripts/healthcheck.sh /usr/bin/healthcheck.sh
ENV HEALTH_ENDPOINT=http://localhost:200
HEALTHCHECK --interval=30s --timeout=30s --start-period=45s --retries=2 \
  CMD [ "healthcheck.sh" ]

## Expose ports, define entrypoint:
EXPOSE 80
ENTRYPOINT ["caddy", "-conf", "/etc/Caddyfile", "-log", "stdout"]
