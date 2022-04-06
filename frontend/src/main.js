import {createApp} from 'vue'
import {createPinia} from 'pinia'
import axios from 'axios'
import VueAxios from "vue-axios";
import App from './App.vue'
import router from "./router";

createApp(App).use(createPinia().use(axios).use(VueAxios).use(router)).mount('#app')



