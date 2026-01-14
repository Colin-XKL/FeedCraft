import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const DASHBOARD: AppRouteRecordRaw = {
  path: '/dashboard',
  name: 'dashboard',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.dashboard',
    requiresAuth: true,
    icon: 'icon-dashboard',
    order: 0,
  },
  children: [
    {
      path: 'welcome',
      name: 'Welcome',
      component: () => import('@/views/dashboard/welcome/welcome.vue'),
      meta: {
        locale: 'menu.dashboard.welcome',
        requiresAuth: false,
        roles: ['*'],
      },
    },
    {
      path: 'quick_start',
      name: 'QuickStartFeedCraftUrlGenerator',
      component: () =>
        import('@/views/dashboard/url_generator/url_generator.vue'),
      meta: {
        requiresAuth: false,
        locale: 'menu.quickStart',
      },
    },
    {
      path: 'all_craft_list',
      name: 'AllCraftList',
      component: () =>
        import('@/views/dashboard/all_craft_list/all_craft_list.vue'),
      meta: {
        requiresAuth: true,
        locale: 'menu.allCraftList',
      },
    },
  ],
};

export default DASHBOARD;
