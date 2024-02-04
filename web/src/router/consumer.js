// 用户前台路由

import ArticleListVue from '@/views/consumer/article/ArticleList.vue'
import ArticleMainVue from '@/views/consumer/article/ArticleMain.vue'

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
    }
]
export default ConsumerRouters