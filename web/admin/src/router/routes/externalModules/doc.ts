import { DOCS_CENTER_ROUTE_NAME, DOCS_ORIGIN } from '@/utils/docsUrl';

export default {
  path: DOCS_ORIGIN,
  name: DOCS_CENTER_ROUTE_NAME,
  meta: {
    locale: 'menu.doc',
    icon: 'icon-book',
    requiresAuth: false,
    order: 8,
  },
};
