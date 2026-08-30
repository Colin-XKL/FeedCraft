export type MenuGroupKey =
  | 'feedSourceGenerate'
  | 'feedSourceProcess'
  | 'feedSourceManage'
  | 'feedSourcePreview'
  | 'feedSourceDebug'
  | 'systemStatus';

export interface MenuGroupDef {
  name: string;
  locale: string;
  icon: string;
  order: number;
}

export const MENU_GROUPS: Record<MenuGroupKey, MenuGroupDef> = {
  feedSourceGenerate: {
    name: 'FeedSourceGenerate',
    locale: 'menu.feedSourceGenerate',
    icon: 'icon-edit',
    order: 0,
  },
  feedSourceProcess: {
    name: 'FeedSourceProcess',
    locale: 'menu.feedSourceProcess',
    icon: 'icon-layers',
    order: 1,
  },
  feedSourceManage: {
    name: 'FeedSourceManage',
    locale: 'menu.feedSourceManage',
    icon: 'icon-book',
    order: 2,
  },
  feedSourcePreview: {
    name: 'FeedSourcePreview',
    locale: 'menu.feedSourcePreview',
    icon: 'icon-eye',
    order: 0,
  },
  feedSourceDebug: {
    name: 'FeedSourceDebug',
    locale: 'menu.feedSourceDebug',
    icon: 'icon-tool',
    order: 1,
  },
  systemStatus: {
    name: 'SystemStatus',
    locale: 'menu.systemStatus',
    icon: 'icon-check-circle',
    order: 2,
  },
};

export const DEFAULT_OPEN_MENU_KEYS = ['WorkTableRoot'];

const LEAF_ORDER: Record<string, number> = {
  HtmlToRss: 0,
  JsonToRss: 1,
  SearchToRss: 2,
  WebMonitor: 3,
  InboxManager: 4,
  CraftAtom: 0,
  CraftFlow: 1,
  CustomRecipe: 0,
  TopicFeed: 1,
  FeedViewer: 0,
  ExampleRssFeeds: 1,
  LlmDebug: 0,
  EmbeddingFilterDebug: 1,
  AdCheckDebug: 2,
  Observability: 0,
  SystemHealth: 1,
  AllCraftList: 2,
};

export interface MenuRouteLike {
  name?: string | symbol;
  path: string;
  meta?: Record<string, any>;
  children?: MenuRouteLike[];
  redirect?: unknown;
}

function routeName(route: MenuRouteLike): string {
  return String(route.name ?? '');
}

function nestByGroup(children: MenuRouteLike[]): MenuRouteLike[] {
  const hidden: MenuRouteLike[] = [];
  const grouped = new Map<MenuGroupKey, MenuRouteLike[]>();
  const ungrouped: MenuRouteLike[] = [];

  children.forEach((child) => {
    if (child.meta?.hideInMenu === true) {
      hidden.push(child);
      return;
    }
    const groupKey = child.meta?.menuGroup as MenuGroupKey | undefined;
    if (groupKey && MENU_GROUPS[groupKey]) {
      const bucket = grouped.get(groupKey) ?? [];
      bucket.push(child);
      grouped.set(groupKey, bucket);
      return;
    }
    ungrouped.push(child);
  });

  const groupNodes = [...grouped.entries()]
    .sort(
      ([a], [b]) =>
        (MENU_GROUPS[a].order ?? 0) - (MENU_GROUPS[b].order ?? 0) ||
        MENU_GROUPS[a].name.localeCompare(MENU_GROUPS[b].name)
    )
    .map(([key, items]) => {
      const def = MENU_GROUPS[key];
      const sorted = [...items].sort(
        (a, b) =>
          (LEAF_ORDER[routeName(a)] ?? 99) - (LEAF_ORDER[routeName(b)] ?? 99)
      );
      return {
        path: '',
        name: def.name,
        meta: {
          locale: def.locale,
          icon: def.icon,
          requiresAuth: true,
        },
        children: sorted,
      };
    });

  return [...groupNodes, ...ungrouped, ...hidden];
}

function findNamed(
  routes: MenuRouteLike[],
  name: string
): MenuRouteLike | undefined {
  return routes.reduce<MenuRouteLike | undefined>((found, route) => {
    if (found) {
      return found;
    }
    if (routeName(route) === name) {
      return route;
    }
    if (route.children?.length) {
      return findNamed(route.children, name);
    }
    return undefined;
  }, undefined);
}

export function applyMenuGrouping(routes: MenuRouteLike[]): MenuRouteLike[] {
  const allCraft = findNamed(routes, 'AllCraftList');

  return routes.map((route) => {
    const name = routeName(route);
    if (name === 'WorkTableRoot') {
      return {
        ...route,
        children: nestByGroup(route.children ?? []),
      };
    }
    if (name === 'tools') {
      const extra = allCraft
        ? [
            {
              ...allCraft,
              meta: {
                ...allCraft.meta,
                hideInMenu: false,
                menuGroup: 'systemStatus' as MenuGroupKey,
              },
            },
          ]
        : [];
      return {
        ...route,
        children: nestByGroup([...(route.children ?? []), ...extra]),
      };
    }
    if (name === 'dashboard') {
      return {
        ...route,
        children: (route.children ?? []).map((child) =>
          routeName(child) === 'AllCraftList'
            ? { ...child, meta: { ...child.meta, hideInMenu: true } }
            : child
        ),
      };
    }
    return route;
  });
}

export function collectMenuLeafNames(routes: MenuRouteLike[]): string[] {
  const names: string[] = [];
  const walk = (nodes: MenuRouteLike[]) => {
    nodes.forEach((node) => {
      if (node.meta?.hideInMenu === true) {
        return;
      }
      if (node.children?.length) {
        walk(node.children);
        return;
      }
      names.push(routeName(node));
    });
  };
  walk(routes);
  return names;
}

export function collectVisibleGroupNames(routes: MenuRouteLike[]): string[] {
  const names: string[] = [];
  const walk = (nodes: MenuRouteLike[]) => {
    nodes.forEach((node) => {
      if (node.meta?.hideInMenu === true) {
        return;
      }
      if (node.children?.length) {
        names.push(routeName(node));
        walk(node.children);
      }
    });
  };
  walk(routes);
  return names;
}
