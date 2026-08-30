export type PreviewMode = 'url' | 'recipe' | 'topic' | 'inbox' | 'uri';
export type PageMode = 'preview' | 'compare';

export type FeedViewerQuery = Record<string, unknown>;

export type FeedViewerRouteState = {
  pageMode: PageMode;
  previewMode: PreviewMode;
  feedUrl: string;
  advancedURI: string;
  selectedRecipeId: string;
  selectedTopicId: string;
  selectedInboxId: string;
  selectedCraft: string;
};

function firstQueryValue(value: unknown): string {
  if (Array.isArray(value)) return String(value[0] || '');
  return typeof value === 'string' ? value : '';
}

function emptyState(pageMode: PageMode): FeedViewerRouteState {
  return {
    pageMode,
    previewMode: 'url',
    feedUrl: '',
    advancedURI: '',
    selectedRecipeId: '',
    selectedTopicId: '',
    selectedInboxId: '',
    selectedCraft: '',
  };
}

function isComparePath(path: string): boolean {
  return /\/feed_compare\/?$/.test(path);
}

function selectInputURI(
  inputURI: string,
  state: FeedViewerRouteState
): FeedViewerRouteState {
  try {
    const parsed = new URL(inputURI);
    if (parsed.protocol === 'feedcraft:') {
      const id = parsed.pathname.replace(/^\/+/, '');
      if (parsed.hostname === 'recipe') {
        return { ...state, previewMode: 'recipe', selectedRecipeId: id };
      }
      if (parsed.hostname === 'topic') {
        return { ...state, previewMode: 'topic', selectedTopicId: id };
      }
      if (parsed.hostname === 'inbox') {
        return { ...state, previewMode: 'inbox', selectedInboxId: id };
      }
    }
    if (parsed.protocol === 'http:' || parsed.protocol === 'https:') {
      return { ...state, previewMode: 'url', feedUrl: inputURI };
    }
  } catch {
    // Fall through to advanced URI mode.
  }
  return { ...state, previewMode: 'uri', advancedURI: inputURI };
}

export function resolveFeedViewerRouteState(
  path: string,
  query: FeedViewerQuery
): FeedViewerRouteState {
  const pageMode: PageMode =
    isComparePath(path) || firstQueryValue(query.page_mode) === 'compare'
      ? 'compare'
      : 'preview';

  const state = {
    ...emptyState(pageMode),
    selectedCraft: firstQueryValue(query.craft || query.craft_name),
  };

  const target = firstQueryValue(query.target || query.mode);
  const id = firstQueryValue(query.id);
  if (target === 'recipe') {
    return {
      ...state,
      previewMode: 'recipe',
      selectedRecipeId: id || firstQueryValue(query.recipe_id),
    };
  }
  if (target === 'topic') {
    return {
      ...state,
      previewMode: 'topic',
      selectedTopicId: id || firstQueryValue(query.topic_id),
    };
  }
  if (target === 'inbox') {
    return {
      ...state,
      previewMode: 'inbox',
      selectedInboxId: id || firstQueryValue(query.inbox_id),
    };
  }

  const inputURI = firstQueryValue(query.input_uri || query.uri || query.url);
  if (inputURI) {
    return selectInputURI(inputURI, state);
  }
  return state;
}
