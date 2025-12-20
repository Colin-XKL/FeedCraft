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
        locale: 'menu.feedViewer',
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
        locale: 'menu.feedCompare',
      },
    },
    {
      path: 'ad-check-debug',
      name: 'AdCheckDebug',
      component: () => import('@/views/dashboard/llm_debug/ad-check-debug.vue'),
      meta: {
        locale: 'menu.adCheckDebug',
        requiresAuth: true,
      },
    },
    {
      path: 'llm-debug',
      name: 'LlmDebug',
      component: () => import('@/views/dashboard/llm_debug/llm-test.vue'),
      meta: {
        locale: 'menu.llmDebug',
        requiresAuth: true,
      },
    },
    {
      path: 'rss-generator',
      name: 'RssGenerator',
      component: () =>
        import('@/views/dashboard/rss_generator/rss_generator.vue'),
      meta: {
        locale: 'menu.rssGenerator',
        requiresAuth: true,
      },
    },
    {
      path: 'json-rss-generator',
      name: 'JsonRssGenerator',
      component: () =>
        import('@/views/dashboard/json_rss_generator/json_rss_generator.vue'),
      meta: {
        locale: 'menu.jsonRssGenerator',
        requiresAuth: true,
      },
    },
    {
      path: 'dependencies',
      name: 'DependencyStatus',
      component: () =>
        import('@/views/dashboard/dependency_service/index.vue'),
      meta: {
        locale: 'menu.dependencyStatus',
        requiresAuth: true,
      },
    },
  ],
};

export default TOOLS;
