import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import tsconfigPaths from 'vite-tsconfig-paths';


export default defineConfig({
  plugins: [react(), tsconfigPaths()],
  server: {
    proxy: {
      "/api": {
        target: "https://gcsim.app",
        changeOrigin: true
      },
      "/hastebin/post": {
        target: "https://hastebin.com/documents",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/hastebin\/post/, '')
      },
      "/hastebin/get": {
        target: "https://hastebin.com/raw/",
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/hastebin\/get/, '')
      }
    }
  }
});