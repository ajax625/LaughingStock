import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

const backendURL = process.env.BACKEND_URL ?? 'http://localhost:9090'
const wsURL = backendURL.replace(/^http/, 'ws')

export default defineConfig({
  base: '/app/',
  plugins: [react()],
  server: {
    proxy: {
      '/auth': backendURL,
      '/me': backendURL,
      '/webhook': backendURL,
      '/ws': { target: wsURL, ws: true },
    },
  },
  build: {
    outDir: 'dist',
  },
})
