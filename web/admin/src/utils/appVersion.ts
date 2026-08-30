export interface VersionSource {
  /** Version handed down by the release pipeline (e.g. v3.2.0 or dev-<hash>). */
  explicitVersion?: string;
  /** Git branch of the build, when the build platform exposes one. */
  branch?: string;
  /** Git commit of the build, when the build platform exposes one. */
  commitSha?: string;
  packageVersion: string;
}

export function resolveDisplayVersion({
  explicitVersion,
  branch,
  commitSha,
  packageVersion,
}: VersionSource): string {
  const explicit = explicitVersion?.trim();
  if (explicit) {
    return explicit;
  }
  const shortSha = commitSha?.trim().slice(0, 7);
  if (shortSha && branch?.trim() !== 'main') {
    return `dev-${shortSha}`;
  }
  return `v${packageVersion}`;
}

export function formatFeedCraftVersion(displayVersion: string): string {
  return `FeedCraft ${displayVersion}`;
}

export function getAdminFooterVersion(): string {
  // eslint-disable-next-line no-underscore-dangle
  return formatFeedCraftVersion(__APP_VERSION__);
}
