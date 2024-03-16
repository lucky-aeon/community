import usePermission from "@/stores/PermissionStore";
import { isLogin } from "@/utils/auth";
import LoginPageVue from "@/views/LoginPage.vue";
import AdminLayout from "@/views/layout/AdminLayout.vue";
import ConsumerLayoutVue from "@/views/layout/ConsumerLayout.vue";
import { Message } from "@arco-design/web-vue";
import { createRouter, createWebHistory } from "vue-router";
import AdminRouters from "./admin";
import ConsumerRouters from "./consumer";

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            component: ConsumerLayoutVue,
            children: ConsumerRouters,
            redirect: "/article",
            meta: {
                requiresAuth: true,
                roles: ["*"]
            }
        },
        {
            path: '/auth',
            component: LoginPageVue,
            meta: {
                requiresAnonymous: true,
                roles: ["*"]
            }
        },
        {
            path: '/admin',
            component: AdminLayout,
            children: AdminRouters,
            redirect: "/admin/code",
            meta: {
                requiresAuth: true,
                requiresAdmin: true,
                roles: ["admin"]
            }
        }
    ]
})
router.beforeEach((to, _, next) => {
    async function crossroads() {
        const Permission = usePermission();
        if (Permission.accessRouter(to)) next();
        else {
            Message.error("您的等级不够，没有权限访问该页面!")
            next({ path: '/auth' });
        }
    }
    // 是否需要验证 
    if (to.meta.requiresAuth) {
        // 是否登录
        if (!isLogin()) {
            next({ path: '/auth' })
            return
        } else {
            crossroads()
        }
    } else if (to.meta.requiresAnonymous && isLogin()) {
        next({ path: '/user' })
        return
    }
})
export default router