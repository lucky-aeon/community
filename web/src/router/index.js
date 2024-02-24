import { isLogin } from "@/utils/auth";
import LoginPageVue from "@/views/LoginPage.vue";
import ConsumerLayoutVue from "@/views/layout/ConsumerLayout.vue";
import { createRouter, createWebHistory } from "vue-router";
import ConsumerRouters from "./consumer";

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            component: ConsumerLayoutVue,
            children: ConsumerRouters,
            meta: {
                requiresAuth: true
            }
        },
        {
            path: '/auth',
            component: LoginPageVue,
            meta: {
                requiresAnonymous: true
            }
        }
    ]
})
router.beforeEach((to)=>{
    // 是否需要验证
    if(to.meta.requiresAuth) {
        // 是否登录
        if(!isLogin()) {
            return {path: '/auth'}
        }
    }else if(to.meta.requiresAnonymous && isLogin()) {
        return {path: '/user'}
    }
    return true
})
export default router