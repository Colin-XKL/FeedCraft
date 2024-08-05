import { AppRouteRecordRaw } from '@/router/routes/types';
import { DEFAULT_LAYOUT } from '@/router/routes/base';

const ADMIN_PAGE: AppRouteRecordRaw = {
  path: '/admin',
  name: 'adminRoot',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.user',
    requiresAuth: true,
    icon: 'icon-user',
    order: 10,
  },
  children: [
    {
      path: 'change_password',
      name: 'ChangePassword',
      component: () => import('@/views/dashboard/admin/pass.vue'),
      meta: {
        title: '修改密码',
        requiresAuth: true,
      },
    },
  ],
};

export default ADMIN_PAGE;
