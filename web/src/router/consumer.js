// 用户前台路由

import ArticleListVue from '@/views/consumer/article/ArticleList.vue';
import ArticleMainVue from '@/views/consumer/article/ArticleMain.vue';
import Comments from "@/views/consumer/comment/Comment.vue";
import Message from "@/views/consumer/message/Message.vue";
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
        component: UserInfo,
        meta: {
            requiresAuth: true
        }
    },
    {
        path: "/comment",
        component: Comments,
        meta: {
            requiresAuth: true
        }
    },
    {
        path: "/message",
        component: Message
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