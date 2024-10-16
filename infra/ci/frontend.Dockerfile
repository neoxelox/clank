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
RUN sed -i '/^name = / { s/ = "\(.*\)"/ = "\1-ci"/; s/ #:replace.*//; }' wrangler.toml

# Build frontend bundle with environment variables injected
RUN --mount=type=secret,id=env \
    env $(grep -v '^#' /run/secrets/env | awk '/=/ {print $1}') \
    yarn build

FROM --platform=linux/amd64 node:22.4-slim AS app

ENV NODE_ENV=production

WORKDIR /app

# Setup wrangler manager
RUN npm install -g @cloudflare/wrangler

# Copy frontend bundle
COPY --from=builder /app/.svelte-kit ./.svelte-kit

# Copy wrangler config
COPY wrangler.toml ./

# Replace environmental values in wrangler config
RUN sed -i '/^name = / { s/ = "\(.*\)"/ = "\1-ci"/; s/ #:replace.*//; }' wrangler.toml

# Copy dependencies
COPY --from=builder /app/node_modules ./node_modules

# Run container as non-root
RUN groupadd -g 1001 app && \
    useradd -u 1001 -g app -d /app -s /bin/sh -m app && \
    chown -R app:app /app

# App
EXPOSE 3333

USER app

CMD [ "wrangler", "pages", "dev", "--port", "3333" ]
