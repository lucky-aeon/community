// 用户前台路由

import ArticleViewVue from '@/views/consumer/article/ArticleView.vue'

/**
 * @type {import('vue-router').RouteRecordRaw}
 */
const ConsumerRouters = [
    {
        path: "/article",
        component: ArticleViewVue
    }
]
export default ConsumerRouters