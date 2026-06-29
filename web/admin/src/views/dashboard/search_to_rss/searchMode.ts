import type { SearchFetchReq } from '@/api/json_rss';

export type SearchMode = 'keyword' | 'semantic';

export interface SearchModeOption {
  value: SearchMode;
  titleKey: string;
  descriptionKey: string;
  badgeKey: string;
  placeholderKey: string;
  helpKey: string;
}

export interface SearchFeedMeta {
  title: string;
  description: string;
  link: string;
}

export const searchModeOptions: SearchModeOption[] = [
  {
    value: 'keyword',
    titleKey: 'searchToRss.mode.keyword.title',
    descriptionKey: 'searchToRss.mode.keyword.description',
    badgeKey: 'searchToRss.mode.keyword.badge',
    placeholderKey: 'searchToRss.mode.keyword.placeholder',
    helpKey: 'searchToRss.mode.keyword.help',
  },
  {
    value: 'semantic',
    titleKey: 'searchToRss.mode.semantic.title',
    descriptionKey: 'searchToRss.mode.semantic.description',
    badgeKey: 'searchToRss.mode.semantic.badge',
    placeholderKey: 'searchToRss.mode.semantic.placeholder',
    helpKey: 'searchToRss.mode.semantic.help',
  },
];

export function isSemanticSearchMode(mode: SearchMode) {
  return mode === 'semantic';
}

export function buildSearchFetchReq(
  query: string,
  mode: SearchMode
): SearchFetchReq {
  return {
    query,
    enhanced_mode: isSemanticSearchMode(mode),
  };
}

export function buildSearchSourceConfig(
  query: string,
  mode: SearchMode,
  feedMeta: SearchFeedMeta
) {
  return {
    type: 'search',
    search_fetcher: buildSearchFetchReq(query, mode),
    feed_meta: {
      title: feedMeta.title,
      description: feedMeta.description,
      link: feedMeta.link,
    },
  };
}
