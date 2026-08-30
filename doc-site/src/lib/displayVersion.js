/**
 * Mirrors web/admin/src/utils/appVersion.ts so both frontends label a build the
 * same way.
 *
 * @param {{ packageVersion: string, explicitVersion?: string, branch?: string, commitSha?: string }} source
 * @returns {string}
 */
export function resolveDisplayVersion({
  explicitVersion,
  branch,
  commitSha,
  packageVersion,
}) {
  const explicit = explicitVersion?.trim();
  if (explicit) {
    return explicit;
  }
  const shortSha = commitSha?.trim().slice(0, 7);
  if (shortSha && branch?.trim() !== "main") {
    return `dev-${shortSha}`;
  }
  return `v${packageVersion}`;
}

/**
 * @param {string} displayVersion
 * @returns {string}
 */
export function formatFeedCraftVersion(displayVersion) {
  return `FeedCraft ${displayVersion}`;
}
