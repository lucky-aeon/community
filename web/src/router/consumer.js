// 用户前台路由

import File from "@/views/admin/file/File.vue";
import Member from "@/views/admin/member/Member.vue";
import ArticleListVue from '@/views/consumer/article/ArticleList.vue';
import ArticleMainVue from '@/views/consumer/article/ArticleMain.vue';


const ConsumerRouters = [
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
                path: "comment",
                component: () => import('@/views/consumer/user/CommentManager.vue'),
            },
            {
                path: "message",
                component: ()=> import('@/views/consumer/user/UserMessage.vue')
            },
            {
                path: "subscribe",
                component: () => import('@/views/consumer/user/UserSubscribe.vue'),
            }
        ]
    },
    {
        path: "/member",
        component: Member
    },
    {
        path: "/file",
        component: File
    }
]
export default ConsumerRouters