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
    // {
    //   path: 'welcome',
    //   name: 'Welcome',
    //   component: () => import('@/views/dashboard/welcome/welcome.vue'),
    //   meta: {
    //     title: '欢迎使用',
    //     requiresAuth: true,
    //   },
    // },
    {
      path: 'viewer',
      name: 'FeedViewer',
      component: () => import('@/views/dashboard/feed-viewer/feed_viewer.vue'),
      meta: {
        title: 'FeedViewer',
        requiresAuth: false,
      },
    },
    {
      path: 'workplace',
      name: 'Workplace',
      component: () => import('@/views/dashboard/welcome/welcome.vue'),
      meta: {
        locale: 'menu.dashboard.workplace',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'ad-check-debug',
      name: 'AdCheckDebug',
      component: () => import('@/views/dashboard/llm-debug/ad-check-debug.vue'),
      meta: {
        title: 'Ad Check Debug',
        requiresAuth: true,
      },
    },
    {
      path: 'llm-debug',
      name: 'llm-debug',
      component: () => import('@/views/dashboard/llm-debug/llm-test.vue'),
      meta: {
        title: 'LLM Debug',
        requiresAuth: true,
      },
    },
    {
      path: 'custom_recipe',
      name: 'custom_recipe',
      component: () =>
        import('@/views/dashboard/custom_recipe/custom_recipe.vue'),
      meta: {
        requiresAuth: true,
        title: 'CustomRecipe',
      },
    },
    {
      path: 'craft_flow',
      name: 'craft_flow',
      component: () => import('@/views/dashboard/craft_flow/craft_flow.vue'),
      meta: {
        requiresAuth: true,
        title: 'Craft Flow',
      },
    },
  ],
};

export default DASHBOARD;
