import {createRouter, createWebHashHistory} from 'vue-router'
import requestManager from '@/utils/requestManager';
import {closeLoading} from '@/utils/loadingManager'; // 引入加载动画管理器
const routes = [
    {
        path: '/',
        redirect: '/login'
    },
    {
        path: '/init',
        name: 'Init',
        component: () => import('@/view/init/index.vue')
    },
    {
        path: '/licenseConfig',
        name: 'LicenseConfig',
        component: () => import('@/view/config/licenseConfig.vue'),
    },
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/view/login/index.vue')
    },
    {
        path: '/jump',
        name: 'Jump',
        component: () => import('@/view/jump/index.vue')
    },
    {
        path: '/:catchAll(.*)',
        meta: {
            closeTab: true,
        },
        component: () => import('@/view/error/index.vue')
    },
    {
        path: '/terminal',
        name: 'terminal',
        component: () => import('@/view/editor/cloudpods/Terminal.vue')
    },
    {
        path: '/virTerminal',
        name: 'virTerminal',
        component: () => import('@/view/editor/cloudpods/VirTerminal.vue')
    },
]

const router = createRouter({
    history: createWebHashHistory(),
    routes,
})
// 路由拦截
router.beforeEach((to, from, next) => {
    // 取消所有未完成的请求
    requestManager.cancelAll();

    // 确保加载动画被关闭
    closeLoading();
    next();
});

export default router
