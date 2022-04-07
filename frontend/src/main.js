import {createApp} from 'vue'
import {createPinia} from 'pinia'
import App from './App.vue'
import router from "./router";
import ElementPlus from 'element-plus'

createApp(App).use(createPinia()).use(router).use(ElementPlus).mount('#app');
