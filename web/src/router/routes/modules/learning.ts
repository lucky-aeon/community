import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const DASHBOARD: AppRouteRecordRaw = {
  path: '/learning',
  name: 'Learning',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.learning',
    requiresAuth: true,
    icon: 'icon-user',
    order: 1,
  },
  children: [
    {
      path: 'qa',
      name: 'learningQa',
      component: () => import('@/views/learning/qa/index.vue'),
      meta: {
        locale: 'menu.learning.qa',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'project',
      name: 'learningProejct',
      component: () => import('@/views/learning/project/index.vue'),
      meta: {
        locale: 'menu.learning.project',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'share',
      name: 'learningShare',
      component: () => import('@/views/learning/share/index.vue'),
      meta: {
        locale: 'menu.learning.share',
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
};

export default DASHBOARD;
