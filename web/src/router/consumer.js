// 用户前台路由

import ArticleListVue from '@/views/consumer/article/ArticleList.vue';
import ArticleMainVue from '@/views/consumer/article/ArticleMain.vue';



const ConsumerRouters = [
    {
        path: "/qa",
        component: () => import('@/views/consumer/qa/QaLayout.vue'),
        meta: {
            requiresAuth: true
        },
        children: [
            {
                path: ":classfily",
                component: ()=> import('@/views/consumer/qa/QaList.vue')
            },
            {
                path: "view/:id",
                name: "qaView",
                component: () => import('@/views/consumer/qa/QaView.vue')
            }
        ]
    },
    {
        path: "/article",
        component: ArticleMainVue,
        meta: {
            requiresAuth: true
        },
        children: [
            {
                path: ":classfily",
                component: ArticleListVue
            },
            {
                path: "view/:id",
                name: "articleView",
                component: () => import('@/views/consumer/article/ArticleView.vue')
            }
        ]
    },
    {
        path: "/user",
        component: () => import('@/views/consumer/user/UserLayout.vue'),
        meta: {
            requiresAuth: true
        },
        children: [
            {
                path: ":userId",
                name: "UserHome",
                component: () => import('@/views/consumer/user/UserHome.vue'),
                meta: {
                    hideSelf: true
                }
            },
            {
                path: "profile",
                name: "userProfile",
                component: () => import('@/views/consumer/user/UserProfile.vue')
            },
            {
                path: "article",
                name: "articleManager",
                component: () => import('@/views/consumer/article/ArticleManager.vue')
            },
            {
                path: "qa",
                name: "qAManager",
                component: () => import('@/views/consumer/qa/QaManager.vue')
            },
            {
                path: "comment",
                component: () => import('@/views/consumer/user/CommentManager.vue'),
            },
            {
                path: "message",
                component: () => import('@/views/consumer/user/UserMessage.vue')
            },
            {
                path: "subscribe",
                component: () => import('@/views/consumer/user/UserSubscribe.vue'),
            }
        ]
    }
]
export default ConsumerRouters