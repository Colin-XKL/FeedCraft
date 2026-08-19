const apiBaseUrl = import.meta.env.VITE_API_BASE_URL ?? '';

function normalizeBaseUrl(): string {
  if (!apiBaseUrl) {
    return window.location.origin;
  }

  return new URL(apiBaseUrl, window.location.origin).href.replace(/\/$/, '');
}

export default function buildPublicFeedUrl(path: string): string {
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;
  return `${normalizeBaseUrl()}${normalizedPath}`;
}

export function isHttpUrl(value: string): boolean {
  try {
    const parsed = new URL(value.trim());
    return parsed.protocol === 'http:' || parsed.protocol === 'https:';
  } catch {
    return false;
  }
}

export function buildCraftFeedPath(
  craftName: string,
  inputUrl: string
): string {
  const craft = craftName.trim();
  const source = inputUrl.trim();
  return `/craft/${encodeURIComponent(craft)}?input_url=${encodeURIComponent(
    source
  )}`;
}

export function buildCraftFeedUrl(craftName: string, inputUrl: string): string {
  return buildPublicFeedUrl(buildCraftFeedPath(craftName, inputUrl));
}
