import { appRoutes, appExternalRoutes } from '@/router/routes';
import { applyMenuGrouping, type MenuRouteLike } from '@/router/menu-groups';

const mixinRoutes = [...appRoutes, ...appExternalRoutes];

const appClientMenus = applyMenuGrouping(
  mixinRoutes.map((el) => {
    const { name, path, meta, redirect, children } = el;
    return {
      name,
      path,
      meta,
      redirect,
      children,
    } as MenuRouteLike;
  })
);

export default appClientMenus;
