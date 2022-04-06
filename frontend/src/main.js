import {createApp} from 'vue'
import {createPinia} from 'pinia'
import axios from 'axios'
import VueAxios from "vue-axios";
import App from './App.vue'

createApp(App).use(createPinia()).mount('#app')

new Vue({
    el: '#app',
    components: {App},
    template: '<App/>'
})
app.mount("#app");

