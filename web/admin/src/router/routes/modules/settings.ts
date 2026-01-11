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
    order: 2,
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
  ],
};

export default SETTINGS;
