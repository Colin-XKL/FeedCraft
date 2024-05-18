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
        title: '欢迎使用',
        requiresAuth: true,
      },
    },
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
      component: () => import('@/views/dashboard/workplace/index.vue'),
      meta: {
        locale: 'menu.dashboard.workplace',
        requiresAuth: true,
        roles: ['*'],
      },
    },
    {
      path: 'llm-debug',
      name: 'LLMDebug',
      component: () => import('@/views/dashboard/llm-debug/llm-debug.vue'),
      meta: {
        title: 'LLM Debug',
        requiresAuth: true,
      },
    },
  ],
};

export default DASHBOARD;
