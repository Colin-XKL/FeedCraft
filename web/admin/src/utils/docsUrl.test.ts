import { describe, expect, it } from 'vitest';
import {
  buildDocsUrl,
  DOCS_ORIGIN,
  resolveExternalMenuUrl,
} from '@/utils/docsUrl';

describe('buildDocsUrl', () => {
  it('maps zh-CN to the simplified Chinese docs site', () => {
    expect(buildDocsUrl('zh-CN', 'guides/start/quick-start')).toBe(
      'https://feed-craft-doc.vercel.app/zh/guides/start/quick-start/'
    );
  });

  it('maps zh-TW to the traditional Chinese docs site', () => {
    expect(buildDocsUrl('zh-TW', 'guides/start/quick-start')).toBe(
      'https://feed-craft-doc.vercel.app/zh-tw/guides/start/quick-start/'
    );
  });

  it('maps en-US and unknown locales to the English docs site', () => {
    expect(buildDocsUrl('en-US', 'guides/advanced/customization')).toBe(
      'https://feed-craft-doc.vercel.app/en/guides/advanced/customization/'
    );
    expect(buildDocsUrl('fr-FR', 'guides/start/quick-start')).toBe(
      'https://feed-craft-doc.vercel.app/en/guides/start/quick-start/'
    );
  });

  it('keeps a trailing slash and strips extra leading slashes on the path', () => {
    expect(buildDocsUrl('zh-CN', '/guides/start/quick-start/')).toBe(
      'https://feed-craft-doc.vercel.app/zh/guides/start/quick-start/'
    );
  });

  it('returns the locale homepage without a double slash when path is empty', () => {
    expect(buildDocsUrl('zh-CN', '')).toBe(
      'https://feed-craft-doc.vercel.app/zh/'
    );
    expect(buildDocsUrl('zh-TW')).toBe(
      'https://feed-craft-doc.vercel.app/zh-tw/'
    );
    expect(buildDocsUrl('en-US', '')).toBe(
      'https://feed-craft-doc.vercel.app/en/'
    );
    expect(buildDocsUrl('fr-FR', '/')).toBe(
      'https://feed-craft-doc.vercel.app/en/'
    );
  });
});

describe('resolveExternalMenuUrl', () => {
  it('rewrites the docs center menu to the current UI language homepage', () => {
    expect(resolveExternalMenuUrl(DOCS_ORIGIN, 'zh-CN', 'doc-center')).toBe(
      'https://feed-craft-doc.vercel.app/zh/'
    );
    expect(
      resolveExternalMenuUrl(`${DOCS_ORIGIN}/en`, 'zh-TW', 'doc-center')
    ).toBe('https://feed-craft-doc.vercel.app/zh-tw/');
    expect(
      resolveExternalMenuUrl(`${DOCS_ORIGIN}/en`, 'en-US', 'doc-center')
    ).toBe('https://feed-craft-doc.vercel.app/en/');
  });

  it('leaves unrelated external links unchanged', () => {
    expect(
      resolveExternalMenuUrl('https://github.com/Colin-XKL/FeedCraft', 'zh-CN')
    ).toBe('https://github.com/Colin-XKL/FeedCraft');
  });
});
