# Build can take extremely long time on non-native Linux AMD64 platform
FROM --platform=linux/amd64 node:22.4-slim AS builder

ENV NODE_ENV=production

# Setup system
RUN apt update && \
    apt install --yes ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Setup yarn package manager
RUN npm install -g corepack && \
    corepack enable

# Install dependencies
COPY .npmrc .yarnrc.yml package.json yarn.lock ./
RUN yarn install --immutable > /dev/null

# Copy everything
COPY . ./

# Replace environmental values in wrangler config
RUN sed -i '/^name = / { s/ = "\(.*\)"/ = "\1-prod"/; s/ #:replace.*//; }' wrangler.toml

# Build frontend bundle with environment variables injected
RUN --mount=type=secret,id=env \
    env $(grep -v '^#' /run/secrets/env | awk '/=/ {print $1}') \
    yarn build

ENTRYPOINT [ "yarn" ]
