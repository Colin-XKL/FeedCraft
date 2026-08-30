import { describe, expect, it } from 'vitest';
import {
  applyMenuGrouping,
  collectMenuLeafNames,
  collectVisibleGroupNames,
  MENU_GROUPS,
} from './menu-groups';

function leaf(name: string, path: string, group?: string, hideInMenu = false) {
  return {
    name,
    path,
    meta: {
      requiresAuth: true,
      hideInMenu,
      ...(group ? { menuGroup: group } : {}),
    },
  };
}

describe('applyMenuGrouping', () => {
  const routes = [
    {
      name: 'dashboard',
      path: '/dashboard',
      children: [
        leaf('Welcome', 'welcome'),
        leaf('QuickStartFeedCraftUrlGenerator', 'quick_start'),
        leaf('AllCraftList', 'all_craft_list'),
      ],
    },
    {
      name: 'WorkTableRoot',
      path: '/worktable',
      children: [
        leaf('TopicFeed', 'topic_feed', 'feedSourceManage'),
        leaf('TopicFeedCreate', 'topic_feed/create', undefined, true),
        leaf('CustomRecipe', 'custom_recipe', 'feedSourceManage'),
        leaf('CraftAtom', 'craft_atom', 'feedSourceProcess'),
        leaf('CraftFlow', 'craft_flow', 'feedSourceProcess'),
        leaf('HtmlToRss', 'html-to-rss', 'feedSourceGenerate'),
        leaf('WebMonitor', 'web-monitor', 'feedSourceGenerate'),
        leaf('JsonToRss', 'json-to-rss', 'feedSourceGenerate'),
        leaf('SearchToRss', 'search-to-rss', 'feedSourceGenerate'),
        leaf('InboxManager', 'inbox', 'feedSourceGenerate'),
      ],
    },
    {
      name: 'tools',
      path: '/tools',
      children: [
        leaf('FeedViewer', 'viewer', 'feedSourcePreview'),
        leaf('ExampleRssFeeds', 'example_rss_feeds', 'feedSourcePreview'),
        leaf('EmbeddingFilterDebug', 'embedding-filter', 'feedSourceDebug'),
        leaf('AdCheckDebug', 'ad-check-debug', 'feedSourceDebug'),
        leaf('LlmDebug', 'llm-debug', 'feedSourceDebug'),
        leaf('SystemHealth', 'health', 'systemStatus'),
        leaf('Observability', 'observability', 'systemStatus'),
      ],
    },
  ];

  it('nests worktable items into generate / process / manage groups in timeline order', () => {
    const grouped = applyMenuGrouping(routes);
    const worktable = grouped.find((item) => item.name === 'WorkTableRoot');
    const groupNames = (worktable?.children ?? [])
      .filter((child) => child.children?.length)
      .map((child) => child.name);

    expect(groupNames).toEqual([
      MENU_GROUPS.feedSourceGenerate.name,
      MENU_GROUPS.feedSourceProcess.name,
      MENU_GROUPS.feedSourceManage.name,
    ]);

    const generate = worktable?.children?.find(
      (child) => child.name === MENU_GROUPS.feedSourceGenerate.name
    );
    expect(generate?.children?.map((child) => child.name)).toEqual([
      'HtmlToRss',
      'JsonToRss',
      'SearchToRss',
      'WebMonitor',
      'InboxManager',
    ]);
    expect(generate?.children?.map((child) => child.path)).toEqual([
      'html-to-rss',
      'json-to-rss',
      'search-to-rss',
      'web-monitor',
      'inbox',
    ]);
  });

  it('moves all craft list into tools system status and hides it from dashboard', () => {
    const grouped = applyMenuGrouping(routes);
    const dashboard = grouped.find((item) => item.name === 'dashboard');
    const allCraftOnDashboard = dashboard?.children?.find(
      (child) => child.name === 'AllCraftList'
    );
    expect(allCraftOnDashboard?.meta?.hideInMenu).toBe(true);

    const tools = grouped.find((item) => item.name === 'tools');
    const systemStatus = tools?.children?.find(
      (child) => child.name === MENU_GROUPS.systemStatus.name
    );
    expect(systemStatus?.children?.map((child) => child.name)).toEqual([
      'Observability',
      'SystemHealth',
      'AllCraftList',
    ]);
  });

  it('keeps tools debug tools ordered from general to specific', () => {
    const grouped = applyMenuGrouping(routes);
    const tools = grouped.find((item) => item.name === 'tools');
    const debug = tools?.children?.find(
      (child) => child.name === MENU_GROUPS.feedSourceDebug.name
    );
    expect(debug?.children?.map((child) => child.name)).toEqual([
      'LlmDebug',
      'EmbeddingFilterDebug',
      'AdCheckDebug',
    ]);
  });

  it('does not expose hidden editor routes as menu leaves', () => {
    const grouped = applyMenuGrouping(routes);
    expect(collectMenuLeafNames(grouped)).not.toContain('TopicFeedCreate');
    expect(collectVisibleGroupNames(grouped)).toEqual(
      expect.arrayContaining([
        'WorkTableRoot',
        'FeedSourceGenerate',
        'tools',
        'SystemStatus',
      ])
    );
  });
});
