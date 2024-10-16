FROM node:22.4-slim

ENV NODE_ENV=development

WORKDIR /app

# Setup yarn package manager
RUN npm install -g corepack && \
    corepack enable

# Install dependencies
COPY .npmrc .yarnrc.yml package.json yarn.lock ./
RUN yarn install

# Copy non-src/static files
COPY .postcssrc.cjs .tailwindrc.cjs svelte.config.js tsconfig.json vite.config.ts wrangler.toml ./

# App
EXPOSE 3333

CMD [ "yarn", "dev" ]
