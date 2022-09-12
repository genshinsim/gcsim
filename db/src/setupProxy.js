const { createProxyMiddleware } = require("http-proxy-middleware");

module.exports = function (app) {
  app.use(
    "/api",
    createProxyMiddleware({
      target: "http://localhost:3030",
      changeOrigin: true,
    })
  );
  app.use(
    "/assets",
    createProxyMiddleware({
      target: "http://localhost:3030",
      changeOrigin: true,
    })
  );
};
