import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const DASHBOARD: AppRouteRecordRaw = {
  path: '/dashboard',
  name: 'dashboard',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.dashboard.workplace',
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
      path: 'custom_recipe',
      name: 'CustomRecipe',
      component: () =>
        import('@/views/dashboard/custom_recipe/custom_recipe.vue'),
      meta: {
        requiresAuth: true,
        locale: 'menu.customRecipe',
      },
    },
    {
      path: 'craft_flow',
      name: 'CraftFlow',
      component: () => import('@/views/dashboard/craft_flow/craft_flow.vue'),
      meta: {
        requiresAuth: true,
        locale: 'menu.craftFlow',
      },
    },
    {
      path: 'craft_atom',
      name: 'CraftAtom',
      component: () => import('@/views/dashboard/craft_atom/craft_atom.vue'),
      meta: {
        requiresAuth: true,
        locale: 'menu.craftAtom',
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
    {
      path: 'search_rss_generator',
      name: 'SearchRssGenerator',
      component: () =>
        import('@/views/dashboard/search_rss_generator/index.vue'),
      meta: {
        requiresAuth: true,
        locale: 'menu.searchRssGenerator',
      },
    },
  ],
};

export default DASHBOARD;
