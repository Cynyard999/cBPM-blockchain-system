import {createApp} from 'vue'
import {createPinia} from 'pinia'
import App from './App.vue'
import router from "./router";
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as Icons from '@element-plus/icons-vue'

// 创建全局的vue app
const app = createApp(App);
// 使用elementPlus以及router，router的配置在./router/index.js中
app.use(createPinia()).use(router).use(ElementPlus).mount('#app');
Object.keys(Icons).forEach(key => {
    app.component(key, Icons[key])
});
