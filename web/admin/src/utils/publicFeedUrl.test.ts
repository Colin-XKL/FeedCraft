import { describe, expect, it } from 'vitest';
import { buildCraftFeedPath, isHttpUrl } from '@/utils/publicFeedUrl';

describe('public feed URL helpers', () => {
  it('accepts only http(s) URLs', () => {
    expect(isHttpUrl('https://hnrss.org/frontpage')).toBe(true);
    expect(isHttpUrl('http://example.com/feed.xml')).toBe(true);
    expect(isHttpUrl('  https://example.com/rss  ')).toBe(true);
    expect(isHttpUrl('/craft/proxy?input_url=https://example.com')).toBe(false);
    expect(isHttpUrl('ftp://example.com/feed.xml')).toBe(false);
    expect(isHttpUrl('not a url')).toBe(false);
    expect(isHttpUrl('')).toBe(false);
  });

  it('builds a craft path with encoded input_url', () => {
    expect(buildCraftFeedPath('proxy', 'https://hnrss.org/frontpage')).toBe(
      '/craft/proxy?input_url=https%3A%2F%2Fhnrss.org%2Ffrontpage'
    );
  });

  it('trims craft name and source URL before encoding', () => {
    expect(
      buildCraftFeedPath('  fulltext-plus  ', '  https://example.com/a b.xml  ')
    ).toBe(
      '/craft/fulltext-plus?input_url=https%3A%2F%2Fexample.com%2Fa%20b.xml'
    );
  });
});
