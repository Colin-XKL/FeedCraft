import { describe, expect, it } from 'vitest';
import {
  buildSearchFetchReq,
  buildSearchSourceConfig,
  searchModeOptions,
} from '@/views/dashboard/search_to_rss/searchMode';

describe('search to rss mode helpers', () => {
  it('maps keyword and semantic modes to the backend enhanced_mode contract', () => {
    expect(buildSearchFetchReq('AI news', 'keyword')).toEqual({
      query: 'AI news',
      enhanced_mode: false,
    });
    expect(buildSearchFetchReq('AI news', 'semantic')).toEqual({
      query: 'AI news',
      enhanced_mode: true,
    });
  });

  it('preserves the selected mode when building a search source config', () => {
    expect(
      buildSearchSourceConfig('SpaceX launch', 'semantic', {
        title: 'Search: SpaceX launch',
        description: 'Search results for SpaceX launch',
        link: 'https://google.com/search?q=SpaceX%20launch',
      }).search_fetcher
    ).toEqual({
      query: 'SpaceX launch',
      enhanced_mode: true,
    });
  });

  it('exposes keyword first so the simpler mode is the default choice', () => {
    expect(searchModeOptions.map((option) => option.value)).toEqual([
      'keyword',
      'semantic',
    ]);
  });
});
