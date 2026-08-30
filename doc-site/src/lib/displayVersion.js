/**
 * @param {{ pkgVersion: string, branch?: string, sha?: string }} opts
 * @returns {string}
 */
export function resolveDocSiteDisplayVersion({
  pkgVersion,
  branch = "",
  sha = "",
}) {
  return branch === "main" || !branch ? `v${pkgVersion}` : `dev-${sha}`;
}

/**
 * @param {string} displayVersion
 * @returns {string}
 */
export function formatDocSiteFooter(displayVersion) {
  return `FeedCraft ${displayVersion}`;
}
