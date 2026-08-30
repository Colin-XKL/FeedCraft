import { describe, expect, it } from 'vitest';
import { buildDocsUrl } from '@/utils/docsUrl';

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
});
