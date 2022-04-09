import {createApp} from 'vue'
import {createPinia} from 'pinia'
import App from './App.vue'
import router from "./router";
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as Icons from '@element-plus/icons-vue'

const app = createApp(App);
app.use(createPinia()).use(router).use(ElementPlus).mount('#app');
Object.keys(Icons).forEach(key => {
    app.component(key, Icons[key])
});
