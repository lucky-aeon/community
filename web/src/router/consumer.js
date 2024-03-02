// 用户前台路由

import ArticleListVue from '@/views/consumer/article/ArticleList.vue';
import ArticleMainVue from '@/views/consumer/article/ArticleMain.vue';
import Comments from "@/views/consumer/comment/Comment.vue";
import Notice from "@/views/consumer/message/Notice.vue"
import At from "@/views/consumer/message/At.vue"
import Member from "@/views/admin/member/Member.vue";
import UserInfo from "@/views/consumer/user/UserInfo.vue";
import File from "@/views/admin/file/File.vue";


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

            }
        ]
    },
    {
        path: "/message",
        children: [
            {
                path: "/notice",
                component: Notice
            },
            {
                path: "/at",
               component: At
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