import { describe, expect, it } from 'vitest';
import enMenu from '@/locale/en-US/menu';
import zhCNMenu from '@/locale/zh-CN/menu';
import zhTWMenu from '@/locale/zh-TW/menu';
import SETTINGS from '@/router/routes/modules/settings';

const REQUIRED_MENU_KEYS = [
  'menu.feedSourceGenerate',
  'menu.feedSourceProcess',
  'menu.feedSourceManage',
  'menu.feedSourcePreview',
  'menu.feedSourceDebug',
  'menu.systemStatus',
  'menu.rssGenerator',
  'menu.jsonToRss',
  'menu.searchToRss',
  'menu.inbox',
  'menu.craftAtom',
  'menu.craftFlow',
  'menu.customRecipe',
  'menu.topicFeed',
  'menu.feedViewer',
  'menu.exampleRssFeeds',
  'menu.allCraftList',
  'menu.observability',
  'menu.systemHealth',
];

describe('admin menu locales', () => {
  it.each([
    ['zh-CN', zhCNMenu],
    ['en-US', enMenu],
    ['zh-TW', zhTWMenu],
  ] as const)('%s includes the restructured menu keys', (_locale, messages) => {
    REQUIRED_MENU_KEYS.forEach((key) => {
      expect(messages[key as keyof typeof messages]).toBeTruthy();
    });
  });

  it('uses product-facing names instead of HTML/JSON jargon in zh-CN', () => {
    expect(zhCNMenu['menu.rssGenerator']).toBe('网页转 RSS');
    expect(zhCNMenu['menu.jsonToRss']).toBe('接口转 RSS');
    expect(zhCNMenu['menu.customRecipe']).toBe('配方管理');
    expect(zhCNMenu['menu.topicFeed']).toBe('主题聚合');
    expect(zhCNMenu['menu.inbox']).toBe('推送收件箱');
    expect(zhCNMenu['menu.craftAtom']).toBe('原子工艺');
  });
});

describe('settings menu order', () => {
  it('orders config, then status, then account', () => {
    expect(SETTINGS.children?.map((child) => child.name)).toEqual([
      'SearchProvider',
      'FaviconProvider',
      'SystemAuthTokenManager',
      'DependencyStatus',
      'ChangePassword',
    ]);
  });
});
