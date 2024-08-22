import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          lodash: ['diff2html', 'showdown'],
        },
      },
    },
  },
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: 'vitest.setup.ts',
  },
  server: {
    port: 3000,
    hmr: {
      host: 'local.the-zula.app'
    },
    watch: {
      usePolling: true,
    },
    host: true,
  },
})
