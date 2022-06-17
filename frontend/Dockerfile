FROM node:18-alpine3.15
RUN mkdir -p /app/node_modules && chown -R node:node /app/node_modules
WORKDIR /app
USER node

COPY --chown=node:node package.json package-lock.json ./
RUN npm ci --quiet
COPY --chown=node:node . .

