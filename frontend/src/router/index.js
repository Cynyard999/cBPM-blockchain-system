import {createRouter, createWebHistory} from 'vue-router'

const routerHistory = createWebHistory();

const router = createRouter({
    history: routerHistory,
    routes: [
        {
            path: '/',
            redirect: '/home'
        },
        {
            path: '/home',
            component: () => import('../views/Home.vue')
        },
        {
            path: '/manufacturer',
            component: () => import('../views/Manufacturer.vue')
        },
        {
            path: '/carrier',
            component: () => import('../views/Carrier.vue')
        },
        {
            path: '/supplier',
            component: () => import('../views/Supplier.vue')
        },
        {
            path: '/middleman',
            component: () => import('../views/Middleman.vue')
        }
    ]
});

export default router