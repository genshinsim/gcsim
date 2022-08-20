/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  swcMinify: true,
  reactStrictMode: false,
  async rewrites() {
    return [
      {
        source: "/api/:all*",
        destination: "http://localhost:3030/api/:all*",
      },
    ];
  },
};

module.exports = nextConfig;
