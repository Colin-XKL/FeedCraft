export function formatFeedCraftVersion(
  injectedVersion: string | undefined,
  packageVersion: string
): string {
  const injected = injectedVersion?.trim();
  if (injected) {
    return `FeedCraft ${injected}`;
  }
  return `FeedCraft v${packageVersion}`;
}

export function getAdminFooterVersion(): string {
  const fromEnv = import.meta.env.VITE_APP_VERSION;
  // Vite compile-time define; conventional double-underscore name.
  // eslint-disable-next-line no-underscore-dangle
  const compiledVersion = __APP_VERSION__;
  const packageVersion =
    typeof compiledVersion === 'string' ? compiledVersion : '0.0.0';
  return formatFeedCraftVersion(
    typeof fromEnv === 'string' ? fromEnv : undefined,
    packageVersion
  );
}
