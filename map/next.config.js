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
  },
  serverRuntimeConfig: {
    path: {
      exportsMeasures: process.env.PATH_EXPORTS_MEASURES,
    },
  },
};
