import { appRoutes, appExternalRoutes } from '@/router/routes';
import { applyMenuGrouping } from '@/router/menu-groups';

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
    };
  })
);

export default appClientMenus;
