import { AppRouteRecordRaw } from '@/router/routes/types';
import { DEFAULT_LAYOUT } from '@/router/routes/base';

const WORKTABLE: AppRouteRecordRaw = {
  path: '/worktable',
  name: 'WorkTableRoot',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.worktable',
    requiresAuth: true,
    icon: 'icon-common',
    order: 1,
  },
  children: [
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
      path: 'craft_atom',
      name: 'CraftAtom',
      component: () => import('@/views/dashboard/craft_atom/craft_atom.vue'),
      meta: {
        requiresAuth: true,
        locale: 'menu.craftAtom',
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
      path: 'html-to-rss',
      name: 'HtmlToRss',
      component: () => import('@/views/dashboard/html_to_rss/html_to_rss.vue'),
      meta: {
        locale: 'menu.rssGenerator',
        requiresAuth: true,
      },
    },
    {
      path: 'curl-to-rss',
      name: 'CurlToRss',
      component: () => import('@/views/dashboard/curl_to_rss/curl_to_rss.vue'),
      meta: {
        locale: 'menu.curlToRss',
        requiresAuth: true,
      },
    },
    {
      path: 'search-to-rss',
      name: 'SearchToRss',
      component: () => import('@/views/dashboard/search_to_rss/index.vue'),
      meta: {
        locale: 'menu.searchToRss',
        requiresAuth: true,
      },
    },
  ],
};

export default WORKTABLE;
