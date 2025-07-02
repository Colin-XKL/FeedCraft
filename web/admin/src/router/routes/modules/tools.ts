import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const TOOLS: AppRouteRecordRaw = {
  path: '/tools',
  name: 'tools',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.tools',
    requiresAuth: true,
    icon: 'icon-tool',
    order: 1,
  },
  children: [
    {
      path: 'viewer',
      name: 'FeedViewer',
      component: () => import('@/views/dashboard/feed_viewer/feed_viewer.vue'),
      meta: {
        title: 'RSS 预览',
        requiresAuth: false,
      },
    },
    {
      path: 'feed_compare',
      name: 'FeedCompare',
      component: () =>
        import('@/views/dashboard/feed_compare/feed_compare.vue'),
      meta: {
        requiresAuth: false,
        title: 'RSS 对比',
      },
    },
    {
      path: 'ad-check-debug',
      name: 'AdCheckDebug',
      component: () => import('@/views/dashboard/llm_debug/ad-check-debug.vue'),
      meta: {
        title: '广告软文检测',
        requiresAuth: true,
      },
    },
    {
      path: 'llm-debug',
      name: 'llm-debug',
      component: () => import('@/views/dashboard/llm_debug/llm-test.vue'),
      meta: {
        title: 'LLM API 调试',
        requiresAuth: true,
      },
    },
  ],
};

export default TOOLS;
