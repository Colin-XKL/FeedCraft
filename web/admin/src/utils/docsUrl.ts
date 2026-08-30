export const DOCS_ORIGIN = 'https://feed-craft-doc.vercel.app';
export const DOCS_CENTER_ROUTE_NAME = 'doc-center';

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

export function buildDocsUrl(locale: string, path = ''): string {
  const lang = docsLangForLocale(locale);
  const normalizedPath = path.replace(/^\/+/, '').replace(/\/+$/, '');
  if (!normalizedPath) {
    return `${DOCS_ORIGIN}/${lang}/`;
  }
  return `${DOCS_ORIGIN}/${lang}/${normalizedPath}/`;
}

export function isDocsCenterLink(path: string, routeName?: unknown): boolean {
  return (
    routeName === DOCS_CENTER_ROUTE_NAME ||
    path === DOCS_ORIGIN ||
    path.startsWith(`${DOCS_ORIGIN}/`)
  );
}

export function resolveExternalMenuUrl(
  path: string,
  locale: string,
  routeName?: unknown
): string {
  if (isDocsCenterLink(path, routeName)) {
    return buildDocsUrl(locale, '');
  }
  return path;
}
