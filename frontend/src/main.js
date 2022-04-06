import {createApp} from 'vue'
import {createPinia} from 'pinia'
import App from './App.vue'
import router from "./router";

createApp(App).use(createPinia()).use(router).mount('#app');

new Vue({
    el: '#app',
    router,
    components: {App},
    template: '<App/>'
});

