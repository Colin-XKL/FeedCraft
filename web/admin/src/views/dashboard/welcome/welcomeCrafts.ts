export type WelcomeCraftGroup = 'basic' | 'extract' | 'ai';

export interface WelcomeCraftOption {
  group: WelcomeCraftGroup;
  value: string;
  labelKey: string;
}

export const DEFAULT_WELCOME_CRAFT = 'proxy';
export const EXAMPLE_RSS_URL = 'https://hnrss.org/frontpage';

export const WELCOME_CRAFT_GROUPS: {
  id: WelcomeCraftGroup;
  labelKey: string;
}[] = [
  { id: 'basic', labelKey: 'welcome.craftGroup.basic' },
  { id: 'extract', labelKey: 'welcome.craftGroup.extract' },
  { id: 'ai', labelKey: 'welcome.craftGroup.ai' },
];

export const WELCOME_CRAFT_OPTIONS: WelcomeCraftOption[] = [
  { group: 'basic', value: 'proxy', labelKey: 'welcome.craft.proxy' },
  { group: 'basic', value: 'limit', labelKey: 'welcome.craft.limit' },
  { group: 'basic', value: 'keyword', labelKey: 'welcome.craft.keyword' },
  { group: 'basic', value: 'guid-fix', labelKey: 'welcome.craft.guidFix' },
  {
    group: 'basic',
    value: 'relative-link-fix',
    labelKey: 'welcome.craft.relativeLinkFix',
  },
  { group: 'extract', value: 'fulltext', labelKey: 'welcome.craft.fulltext' },
  {
    group: 'extract',
    value: 'fulltext-plus',
    labelKey: 'welcome.craft.fulltextPlus',
  },
  { group: 'extract', value: 'cleanup', labelKey: 'welcome.craft.cleanup' },
  {
    group: 'ai',
    value: 'introduction',
    labelKey: 'welcome.craft.introduction',
  },
  { group: 'ai', value: 'summary', labelKey: 'welcome.craft.summary' },
  {
    group: 'ai',
    value: 'ignore-advertorial',
    labelKey: 'welcome.craft.ignoreAdvertorial',
  },
  {
    group: 'ai',
    value: 'translate-title',
    labelKey: 'welcome.craft.translateTitle',
  },
  {
    group: 'ai',
    value: 'translate-content',
    labelKey: 'welcome.craft.translateContent',
  },
];

export function craftsInGroup(group: WelcomeCraftGroup): WelcomeCraftOption[] {
  return WELCOME_CRAFT_OPTIONS.filter((craft) => craft.group === group);
}
