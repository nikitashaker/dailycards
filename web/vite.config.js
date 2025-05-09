import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
    tailwindcss(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },

  // ────────────────────────────────────────────────
  // Добавлено: прокси‑сервер для запросов /api → :8080
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080', // ваш Go‑бэкенд
        changeOrigin: true,
        secure: false,
      },
    },
  },
})
