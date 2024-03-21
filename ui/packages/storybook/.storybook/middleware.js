const proxy = require("http-proxy-middleware");

module.exports = function expressMiddleware(router) {
  router.use(
    "/api",
    proxy.createProxyMiddleware({
      target: "https://gcsim.app",
      changeOrigin: true,
    })
  );
};
