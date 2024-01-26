import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const DASHBOARD: AppRouteRecordRaw = {
  path: '/user',
  name: 'Consumer',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.consumer',
    requiresAuth: true,
    icon: 'icon-user',
    order: 0,
  },
  children: [
    {
      path: 'workplace',
      name: 'ConsumerWorkplace',
      component: () => import('@/views/consumer/workplace/index.vue'),
      meta: {
        locale: 'menu.dashboard.workplace',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'info',
      name: 'info',
      component: () => import('@/views/consumer/info/index.vue'),
      meta: {
        locale: 'menu.consumer.info',
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
};

export default DASHBOARD;
