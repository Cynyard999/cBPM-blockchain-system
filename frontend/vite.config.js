import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [vue()],
    server: {
        port: 3333,
        proxy: {
            '/query': {
                target: "http://localhost:8080/query",
                changeOrigin: true,
                rewrite: (path) => path.replace(/^\/query/, '')
            },
            '/invoke': {
                target: "http://localhost:8080/invoke",
                changeOrigin: true,
                rewrite: (path) => path.replace(/^\/invoke/, '')
            }
        }
    }
})
