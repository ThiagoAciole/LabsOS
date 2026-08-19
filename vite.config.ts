import path from "path"
import tailwindcss from "@tailwindcss/vite"
import react from "@vitejs/plugin-react"
import { defineConfig } from "vite"

// https://vite.dev/config/
export default defineConfig({
  plugins: [react(), tailwindcss()],
  publicDir: path.resolve(import.meta.dirname, "./assets"),
  resolve: {
    alias: {
      "@": path.resolve(import.meta.dirname, "./web"),
    },
  },
  server: {
    proxy: {
      "/api": { target: "http://localhost:18080", changeOrigin: true },
      "/file-manager": { target: "http://localhost:18081", changeOrigin: true, rewrite: (path) => path.replace(/^\/file-manager/, "") },
    },
  },
})
