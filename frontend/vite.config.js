import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [vue()],
    server: {
        port: 3333,
        proxy: {
            '/invoke': {
                target: 'localhost:8080/invoke',	//实际请求地址
                changeOrigin: true,
                pathRewrite: {
                    '^/invoke': ''   // 重写接口
                }
            },
            '/query': {
                target: 'localhost:8080/invoke',	//实际请求地址
                changeOrigin: true,
                pathRewrite: {
                    '^/invoke': ''   // 重写接口
                }
            }
        }
    }
})
