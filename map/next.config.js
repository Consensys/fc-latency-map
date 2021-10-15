/** @type {import('next').NextConfig} */
module.exports = {
  reactStrictMode: true,
  publicRuntimeConfig: {
    app: {
      name: process.env.SERVICE_NAME,
    },
    path: {
      exportsMeasures: process.env.PATH_EXPORTS_MEASURES,
    },
    latency: {
      low: process.env.LATENCY_LOW_LIMIT,
      medium: process.env.LATENCY_MEDIUM_LIMIT,
    },
  },
  serverRuntimeConfig: {
    path: {
      exportsMeasures: process.env.PATH_EXPORTS_MEASURES,
    },
  },
};
