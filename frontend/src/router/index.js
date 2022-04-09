import {createRouter, createWebHistory} from 'vue-router'
import {ElMessage} from 'element-plus'
import 'element-plus/dist/index.css'

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
            component: () => import('../views/Manufacturer.vue'),
            beforeEnter: (to, from, next) => {
                let userInfo = JSON.parse(window.localStorage.getItem("user"));
                if (userInfo['orgType'] !== 'manufacturer') {
                    ElMessage({
                        message: '企业信息不匹配',
                        type: 'warning',
                    });
                    next('/home');
                }
                next()
            },
        },
        {
            path: '/carrier',
            component: () => import('../views/Carrier.vue'),
            beforeEnter: (to, from, next) => {
                let userInfo = JSON.parse(window.localStorage.getItem("user"));
                if (userInfo['orgType'] !== 'carrier') {
                    ElMessage({
                        message: '企业信息不匹配',
                        type: 'warning',
                    });
                    next('/home');
                }
                next()
            },
        },
        {
            path: '/supplier',
            component: () => import('../views/Supplier.vue'),
            beforeEnter: (to, from, next) => {
                let userInfo = JSON.parse(window.localStorage.getItem("user"));
                if (userInfo['orgType'] !== 'supplier') {
                    ElMessage({
                        message: '企业信息不匹配',
                        type: 'warning',
                    });
                    next('/home');
                }
                next()
            },
        },
        {
            path: '/middleman',
            component: () => import('../views/Middleman.vue'),
            beforeEnter: (to, from, next) => {
                let userInfo = JSON.parse(window.localStorage.getItem("user"));
                if (userInfo['orgType'] !== 'middleman') {
                    ElMessage({
                        message: '企业信息不匹配',
                        type: 'warning',
                    });
                    next('/home');
                }
                next()
            },
        }
    ]
});

router.beforeEach((to, from, next) => {
    if (window.localStorage.getItem("user") === null || window.localStorage.getItem("user") === undefined) {
        ElMessage({
            message: '请先登录',
            type: 'warning',
        });
        next('/home');
    }
    else{
        next();
    }
});
export default router