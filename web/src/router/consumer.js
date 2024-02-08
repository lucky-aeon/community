// 用户前台路由

import ArticleListVue from '@/views/consumer/article/ArticleList.vue'
import ArticleMainVue from '@/views/consumer/article/ArticleMain.vue'
import UserInfo from "@/views/consumer/user/UserInfo.vue";
import Comments from "@/views/consumer/comment/Comment.vue";

/**
 * @type {import('vue-router').RouteRecordRaw}
 */
const ConsumerRouters = [
    {
        path: "/article",
        component: ArticleMainVue,
        children: [
            {
                path: ":classfily",
                component: ArticleListVue
            }   
        ]
    },
    {
        path: "/user",
        component: UserInfo
    },
    {
        path: "/comment",
        component: Comments
    }
]
export default ConsumerRouters