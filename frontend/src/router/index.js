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
            children: [
                {
                    path: "assets",
                    name: "manufacturer-assets",
                    components: {
                        manufacturerSubpage: () => import("../components/mamichannel/AssetsPanel.vue")
                    }
                },
                {
                    path: "orders",
                    name: "manufacturer-orders",
                    components: {
                        manufacturerSubpage: () => import("../components/mamichannel/OrdersPanel.vue")
                    }
                },
                {
                    path: "delivery-details",
                    name: "manufacturer-delivery-details",
                    components: {
                        manufacturerSubpage: () => import("../components/cmachannel/DeliveryDetailsPanel.vue")
                    }
                },
            ],
            beforeEnter: (to, from, next) => {
                let userInfo = JSON.parse(window.localStorage.getItem("user"));
                if (userInfo['orgType'] !== 'manufacturer') {
                    ElMessage({
                        message: '企业信息不匹配',
                        type: 'warning',
                    });
                    next('/home');
                } else {
                    next()
                }
            },
        },
        {
            path: '/carrier',
            component: () => import('../views/Carrier.vue'),
            children: [
                {
                    path: "delivery-arrangements",
                    name: "carrier-delivery-arrangements",
                    components: {
                        carrierSubpage: () => import("../components/micchannel/DeliveryArrangementsPanel.vue")
                    }
                },
                {
                    path: "delivery-details",
                    name: "carrier-delivery-details",
                    components: {
                        carrierSubpage: () => import("../components/cmachannel/DeliveryDetailsPanel.vue")
                    }
                },
                {
                    path: "delivery-orders",
                    name: "carrier-delivery-orders",
                    components: {
                        carrierSubpage: () => import("../components/scchannel/DeliveryOrdersPanel.vue")
                    }
                },
            ],
            beforeEnter: (to, from, next) => {
                let userInfo = JSON.parse(window.localStorage.getItem("user"));
                if (userInfo['orgType'] !== 'carrier') {
                    ElMessage({
                        message: '企业信息不匹配',
                        type: 'warning',
                    });
                    next('/home');
                } else {
                    next()
                }
            },
        },
        {
            path: '/supplier',
            component: () => import('../views/Supplier.vue'),
            children: [
                {
                    path: "assets",
                    name: "supplier-orders",
                    components: {
                        supplierSubpage: () => import("../components/mischannel/AssetsPanel.vue")
                    }
                },
                {
                    path: "supply-orders",
                    name: "supplier-supply-orders1",
                    components: {
                        supplierSubpage: () => import("../components/mischannel/SupplyOrdersPanel.vue")
                    }
                },
                {
                    path: "delivery-orders",
                    name: "supplier-delivery-orders",
                    components: {
                        supplierSubpage: () => import("../components/scchannel/DeliveryOrdersPanel.vue")
                    }
                },
            ],
            beforeEnter: (to, from, next) => {
                let userInfo = JSON.parse(window.localStorage.getItem("user"));
                if (userInfo['orgType'] !== 'supplier') {
                    ElMessage({
                        message: '企业信息不匹配',
                        type: 'warning',
                    });
                    next('/home');
                } else {
                    next()
                }
            },
        },
        {
            path: '/middleman',
            component: () => import('../views/Middleman.vue'),
            children: [
                {
                    path: "assets",
                    name: "middleman-orders",
                    components: {
                        middlemanSubpage: () => import("../components/mischannel/AssetsPanel.vue")
                    }
                },
                {
                    path: "supply-orders",
                    name: "middleman-supply-orders",
                    components: {
                        middlemanSubpage: () => import("../components/mischannel/SupplyOrdersPanel.vue")
                    }
                },
                {
                    path: "goods",
                    name: "middleman-goods",
                    components: {
                        middlemanSubpage: () => import("../components/mamichannel/AssetsPanel.vue")
                    }
                },
                {
                    path: "orders",
                    name: "middleman-orders",
                    components: {
                        middlemanSubpage: () => import("../components/mamichannel/OrdersPanel.vue")
                    }
                },
                {
                    path: "delivery-arrangements",
                    name: "middleman-delivery-arrangements",
                    components: {
                        middlemanSubpage: () => import("../components/micchannel/DeliveryArrangementsPanel.vue")
                    }
                },
            ],
            beforeEnter: (to, from, next) => {
                let userInfo = JSON.parse(window.localStorage.getItem("user"));
                if (userInfo['orgType'] !== 'middleman') {
                    ElMessage({
                        message: '企业信息不匹配',
                        type: 'warning',
                    });
                    next('/home');
                } else {
                    next()
                }
            },
        }
    ]
});

router.beforeEach((to, from, next) => {
    if ((to.path !== "/home" && to.path !== "/") && (window.localStorage.getItem("user") === null || window.localStorage.getItem("user") === undefined)) {
        ElMessage({
            message: '请先登录',
            type: 'warning',
        });
        next('/home');
    } else {
        next();
    }
});
export default router