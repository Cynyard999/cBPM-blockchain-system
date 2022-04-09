import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [vue()],
    server: {
        port: 3333,
        proxy: {
            '/work': {
                target: "http://localhost:8080/work",
                changeOrigin: true,
                rewrite: (path) => path.replace(/^\/work/, '')
            },
            '/user': {
                target: "http://localhost:8080/user",
                changeOrigin: true,
                rewrite: (path) => path.replace(/^\/user/, '')
            }
        }
    }
})
