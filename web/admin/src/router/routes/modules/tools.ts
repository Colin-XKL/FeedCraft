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
      path: 'html-to-rss',
      name: 'HtmlToRss',
      component: () =>
        import('@/views/dashboard/html_to_rss/html_to_rss.vue'),
      meta: {
        locale: 'menu.rssGenerator',
        requiresAuth: true,
      },
    },
    {
      path: 'curl-to-rss',
      name: 'CurlToRss',
      component: () =>
        import('@/views/dashboard/curl_to_rss/curl_to_rss.vue'),
      meta: {
        locale: 'menu.curlToRss',
        requiresAuth: true,
      },
    },
    {
      path: 'search-to-rss',
      name: 'SearchToRss',
      component: () =>
        import('@/views/dashboard/search_to_rss/index.vue'),
      meta: {
        locale: 'menu.searchToRss',
        requiresAuth: true,
      },
    },
  ],
};

export default TOOLS;
