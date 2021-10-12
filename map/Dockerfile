FROM node:lts-alpine AS deps

WORKDIR /opt/app
COPY package.json package.json
COPY package-lock.json package-lock.json
RUN npm ci

FROM node:lts-alpine AS builder

ENV NODE_ENV=production
WORKDIR /opt/app
COPY --from=deps /opt/app/package.json ./package.json
COPY --from=deps /opt/app/package-lock.json ./package-lock.json
COPY --from=deps /opt/app/node_modules ./node_modules
COPY public public
COPY src src
COPY next-env.d.ts next-env.d.ts
COPY next.config.js next.config.js
COPY tsconfig.json tsconfig.json
RUN npm run build

FROM node:lts-alpine AS runner

ARG X_TAG
WORKDIR /opt/app
ENV NODE_ENV=production
COPY --from=builder /opt/app/next.config.js ./
COPY --from=builder /opt/app/public ./public
COPY --from=builder /opt/app/.next ./.next
COPY --from=builder /opt/app/node_modules ./node_modules
CMD ["node_modules/.bin/next", "start"]
EXPOSE 3000
