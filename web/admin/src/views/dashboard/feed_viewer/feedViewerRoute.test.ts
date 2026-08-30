import { describe, expect, it } from 'vitest';
import TOOLS from '@/router/routes/modules/tools';
import { resolveFeedViewerRouteState } from './feedViewerRoute';

describe('resolveFeedViewerRouteState', () => {
  it('opens compare mode for the legacy /tools/feed_compare path', () => {
    expect(
      resolveFeedViewerRouteState('/tools/feed_compare', {}).pageMode
    ).toBe('compare');
  });

  it('keeps preview mode on /tools/viewer without page_mode', () => {
    expect(resolveFeedViewerRouteState('/tools/viewer', {}).pageMode).toBe(
      'preview'
    );
  });

  it('opens compare mode when page_mode=compare is present', () => {
    expect(
      resolveFeedViewerRouteState('/tools/viewer', { page_mode: 'compare' })
        .pageMode
    ).toBe('compare');
  });

  it('still resolves recipe preview targets and can combine with compare mode', () => {
    const state = resolveFeedViewerRouteState('/tools/feed_compare', {
      target: 'recipe',
      id: 'demo-recipe',
    });
    expect(state.previewMode).toBe('recipe');
    expect(state.selectedRecipeId).toBe('demo-recipe');
    expect(state.pageMode).toBe('compare');
  });

  it('fills url and craft query params used by the old compare page', () => {
    const state = resolveFeedViewerRouteState('/tools/feed_compare', {
      url: 'https://hnrss.org/frontpage',
      craft: 'fulltext',
    });
    expect(state.previewMode).toBe('url');
    expect(state.feedUrl).toBe('https://hnrss.org/frontpage');
    expect(state.selectedCraft).toBe('fulltext');
  });
});

describe('tools routes', () => {
  it('keeps a hidden /tools/feed_compare route for compatibility', () => {
    const child = TOOLS.children?.find(
      (route) => route.path === 'feed_compare'
    );
    expect(child).toBeDefined();
    expect(child?.name).toBe('FeedCompare');
    expect(child?.meta?.hideInMenu).toBe(true);
  });
});
