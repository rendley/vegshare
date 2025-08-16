import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      // Строка-ключ: все запросы, начинающиеся с /api, будут проксироваться
      '/api': {
        // Цель: наш бэкенд-сервер
        target: 'http://localhost:8080',
        // Изменяем origin, чтобы бэкенд думал, что запрос пришел с того же хоста
        changeOrigin: true,
      },
    }
  }
})