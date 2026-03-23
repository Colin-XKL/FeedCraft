const apiBaseUrl = import.meta.env.VITE_API_BASE_URL ?? '';

function normalizeBaseUrl(): string {
  if (!apiBaseUrl) {
    return window.location.origin;
  }

  if (/^https?:\/\//.test(apiBaseUrl)) {
    return apiBaseUrl.replace(/\/$/, '');
  }

  const normalizedApiBaseUrl = apiBaseUrl.startsWith('/')
    ? apiBaseUrl
    : `/${apiBaseUrl}`;
  return `${window.location.origin}${normalizedApiBaseUrl}`.replace(/\/$/, '');
}

export default function buildPublicFeedUrl(path: string): string {
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;
  return `${normalizeBaseUrl()}${normalizedPath}`;
}
