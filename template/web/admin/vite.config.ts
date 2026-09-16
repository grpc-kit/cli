import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from "path";

export default defineConfig({
  base: '/admin/',
  build: {
    outDir: '../../public/admin',
  },
  plugins: [react()],
  resolve: {
    alias: [
      { find: "@api/known", replacement: path.resolve(__dirname, "./node_modules/@grpc-kit/adm/src/api/known") },
      { find: "@api", replacement: path.resolve(__dirname, "./src/api") },
      { find: "@", replacement: path.resolve(__dirname, "./src") },
    ],
  },
  server: {
    proxy: {
      '/builtin': {
        target: 'http://127.0.0.1:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/builtin/, '/builtin'),
      },
      '/api': {
        target: 'http://127.0.0.1:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '/api'),
      },
    },
  },
})
