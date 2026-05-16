import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  base: '/app/',
  plugins: [react()],
  server: {
    proxy: {
      '/auth': 'http://localhost:9090',
      '/me': 'http://localhost:9090',
      '/webhook': 'http://localhost:9090',
      '/ws': { target: 'ws://localhost:9090', ws: true },
    },
  },
  build: {
    outDir: 'dist',
  },
})
