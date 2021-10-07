FROM node:lts-alpine AS deps

WORKDIR /opt/app
COPY package.json package-lock.json ./
RUN npm install --frozen-lockfile

FROM node:lts-alpine AS builder

ENV NODE_ENV=production
WORKDIR /opt/app
COPY --from=deps /opt/app/node_modules ./node_modules
RUN npm run build

FROM node:lts-alpine AS runner

ARG X_TAG
WORKDIR /opt/app
ENV NODE_ENV=production
COPY --from=builder /opt/app/next.config.js ./
COPY --from=builder /opt/app/public ./public
COPY --from=builder /opt/app/.next ./.next
COPY --from=builder /opt/app/node_modules ./node_modules
COPY --from=builder /opt/app/data ./data
CMD ["node_modules/.bin/next", "start"]