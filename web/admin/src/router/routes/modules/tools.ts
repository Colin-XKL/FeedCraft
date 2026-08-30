import { DEFAULT_LAYOUT } from '@/router/routes/base';
import { AppRouteRecordRaw } from '@/router/routes/types';

const TOOLS: AppRouteRecordRaw = {
  path: '/tools',
  name: 'tools',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.tools',
    requiresAuth: true,
    icon: 'icon-tool',
    order: 2,
  },
  children: [
    {
      path: 'viewer',
      name: 'FeedViewer',
      component: () => import('@/views/dashboard/feed_viewer/feed_viewer.vue'),
      meta: {
        locale: 'menu.feedViewer',
        requiresAuth: false,
        menuGroup: 'feedSourcePreview',
      },
    },
    {
      path: 'example_rss_feeds',
      name: 'ExampleRssFeeds',
      component: () => import('@/views/dashboard/example_rss_feeds/index.vue'),
      meta: {
        locale: 'menu.exampleRssFeeds',
        requiresAuth: false,
        menuGroup: 'feedSourcePreview',
      },
    },
    {
      path: 'llm-debug',
      name: 'LlmDebug',
      component: () => import('@/views/dashboard/llm_debug/llm-test.vue'),
      meta: {
        locale: 'menu.llmDebug',
        requiresAuth: true,
        menuGroup: 'feedSourceDebug',
      },
    },
    {
      path: 'embedding-filter',
      name: 'EmbeddingFilterDebug',
      component: () =>
        import('@/views/dashboard/embedding_filter_debug/index.vue'),
      meta: {
        requiresAuth: true,
        locale: 'menu.embeddingFilterDebug',
        menuGroup: 'feedSourceDebug',
      },
    },
    {
      path: 'ad-check-debug',
      name: 'AdCheckDebug',
      component: () => import('@/views/dashboard/llm_debug/ad-check-debug.vue'),
      meta: {
        locale: 'menu.adCheckDebug',
        requiresAuth: true,
        menuGroup: 'feedSourceDebug',
      },
    },
    {
      path: 'observability',
      name: 'Observability',
      component: () => import('@/views/dashboard/observability/index.vue'),
      meta: {
        locale: 'menu.observability',
        requiresAuth: true,
        menuGroup: 'systemStatus',
      },
    },
    {
      path: 'health',
      name: 'SystemHealth',
      component: () => import('@/views/dashboard/health/index.vue'),
      meta: {
        locale: 'menu.systemHealth',
        requiresAuth: true,
        menuGroup: 'systemStatus',
      },
    },
  ],
};

export default TOOLS;
