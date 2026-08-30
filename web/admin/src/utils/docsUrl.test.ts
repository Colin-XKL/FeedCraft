import { describe, expect, it } from 'vitest';
import { buildDocsUrl } from '@/utils/docsUrl';

describe('buildDocsUrl', () => {
  it('maps zh-CN to the simplified Chinese docs site', () => {
    expect(buildDocsUrl('zh-CN', 'start/quick-start')).toBe(
      'https://feed-craft-doc.vercel.app/zh/start/quick-start/'
    );
  });

  it('maps zh-TW to the traditional Chinese docs site', () => {
    expect(buildDocsUrl('zh-TW', 'start/quick-start')).toBe(
      'https://feed-craft-doc.vercel.app/zh-tw/start/quick-start/'
    );
  });

  it('maps en-US and unknown locales to the English docs site', () => {
    expect(buildDocsUrl('en-US', 'guides/customization')).toBe(
      'https://feed-craft-doc.vercel.app/en/guides/customization/'
    );
    expect(buildDocsUrl('fr-FR', 'start/quick-start')).toBe(
      'https://feed-craft-doc.vercel.app/en/start/quick-start/'
    );
  });

  it('keeps a trailing slash and strips extra leading slashes on the path', () => {
    expect(buildDocsUrl('zh-CN', '/start/quick-start/')).toBe(
      'https://feed-craft-doc.vercel.app/zh/start/quick-start/'
    );
  });
});
