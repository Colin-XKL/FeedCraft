const DOCS_ORIGIN = 'https://feed-craft-doc.vercel.app';

const LOCALE_TO_DOCS_LANG: Record<string, string> = {
  'zh-CN': 'zh',
  'zh-TW': 'zh-tw',
  'en-US': 'en',
};

export function docsLangForLocale(locale?: string): string {
  if (!locale) {
    return 'en';
  }
  return LOCALE_TO_DOCS_LANG[locale] ?? 'en';
}

export function buildDocsUrl(locale: string, path: string): string {
  const lang = docsLangForLocale(locale);
  const normalizedPath = path.replace(/^\/+/, '').replace(/\/+$/, '');
  return `${DOCS_ORIGIN}/${lang}/${normalizedPath}/`;
}
