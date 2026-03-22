const apiBaseUrl = import.meta.env.VITE_API_BASE_URL ?? '';

function normalizeBaseUrl(): string {
  if (!apiBaseUrl) {
    return window.location.origin;
  }

  if (/^https?:\/\//.test(apiBaseUrl)) {
    return apiBaseUrl.replace(/\/$/, '');
  }

  return `${window.location.origin}${apiBaseUrl}`.replace(/\/$/, '');
}

export default function buildPublicFeedUrl(path: string): string {
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;
  return `${normalizeBaseUrl()}${normalizedPath}`;
}
