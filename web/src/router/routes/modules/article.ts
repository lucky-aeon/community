import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const DASHBOARD: AppRouteRecordRaw = {
  path: '/article',
  name: 'Article',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.article',
    requiresAuth: true,
    icon: 'icon-user',
    order: 1,
  },
  children: [
    {
      path: 'self',
      name: 'ArticleSelf',
      component: () => import('@/views/article/self/index.vue'),
      meta: {
        locale: 'menu.article.self',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'public',
      name: 'ArticlePublic',
      component: () => import('@/views/article/public/index.vue'),
      meta: {
        locale: 'menu.article.public',
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
};

export default DASHBOARD;
