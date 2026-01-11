import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const SETTINGS: AppRouteRecordRaw = {
  path: '/settings',
  name: 'settings',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.settings',
    requiresAuth: true,
    icon: 'icon-settings',
    order: 3,
  },
  children: [
    {
      path: 'search_provider',
      name: 'SearchProvider',
      component: () => import('@/views/settings/search_provider/index.vue'),
      meta: {
        locale: 'menu.settings.searchProvider',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'dependencies',
      name: 'DependencyStatus',
      component: () => import('@/views/dashboard/dependency_service/index.vue'),
      meta: {
        locale: 'menu.dependencyStatus',
        requiresAuth: true,
      },
    },
    {
      path: 'change_password',
      name: 'ChangePassword',
      component: () => import('@/views/dashboard/admin/pass.vue'),
      meta: {
        locale: 'menu.changePassword',
        requiresAuth: true,
      },
    },
  ],
};

export default SETTINGS;
