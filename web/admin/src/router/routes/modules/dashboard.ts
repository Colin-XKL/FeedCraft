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
      path: 'workplace',
      name: 'Workplace',
      component: () => import('@/views/dashboard/welcome/welcome.vue'),
      meta: {
        title: '欢迎页',
        requiresAuth: false,
        roles: ['*'],
      },
    },
    {
      path: 'custom_recipe',
      name: 'custom_recipe',
      component: () =>
        import('@/views/dashboard/custom_recipe/custom_recipe.vue'),
      meta: {
        requiresAuth: true,
        title: 'Craft Recipe',
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
    {
      path: 'craft_atom',
      name: 'craft_atom',
      component: () => import('@/views/dashboard/craft_atom/craft_atom.vue'),
      meta: {
        requiresAuth: true,
        title: 'Craft Atom',
      },
    },
  ],
};

export default DASHBOARD;
