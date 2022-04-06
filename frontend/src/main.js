import {createApp} from 'vue'
import {createPinia} from 'pinia'
import axios from 'axios'
import VueAxios from "vue-axios";
import App from './App.vue'

createApp(App).use(createPinia()).mount('#app')

const app = createApp(App).use(VueAxios, axios);
app.mount("#app");

