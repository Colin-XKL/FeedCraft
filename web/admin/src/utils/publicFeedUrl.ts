const apiBaseUrl = import.meta.env.VITE_API_BASE_URL ?? '';

export function normalizeBaseUrl(): string {
  if (!apiBaseUrl) {
    return window.location.origin;
  }

  let base = apiBaseUrl;
  if (base.endsWith('/api')) {
    base = base.substring(0, base.length - 4);
  } else if (base.endsWith('/api/')) {
    base = base.substring(0, base.length - 5);
  }

  if (!base) {
    base = window.location.origin;
  }

  return new URL(base, window.location.origin).href.replace(/\/$/, '');
}

export default function buildPublicFeedUrl(path: string): string {
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;
  return `${normalizeBaseUrl()}${normalizedPath}`;
}
