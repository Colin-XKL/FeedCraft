import assert from "node:assert/strict";
import { describe, it } from "node:test";
import {
  formatFeedCraftVersion,
  resolveDisplayVersion,
} from "./displayVersion.js";

describe("resolveDisplayVersion", () => {
  it("prefers the version handed down by the release pipeline", () => {
    assert.equal(
      resolveDisplayVersion({
        explicitVersion: "v3.2.0",
        branch: "main",
        commitSha: "becb6a35343b6f1dcac111fb105abddce70175c7",
        packageVersion: "3.1.0",
      }),
      "v3.2.0"
    );
  });

  it("uses the short commit sha for preview builds off main", () => {
    assert.equal(
      resolveDisplayVersion({
        branch: "cursor/col-36-release-please-90a2",
        commitSha: "becb6a35343b6f1dcac111fb105abddce70175c7",
        packageVersion: "3.1.0",
      }),
      "dev-becb6a3"
    );
  });

  it("uses the package version on main and for local builds", () => {
    assert.equal(
      resolveDisplayVersion({
        branch: "main",
        commitSha: "becb6a35343b6f1dcac111fb105abddce70175c7",
        packageVersion: "3.1.0",
      }),
      "v3.1.0"
    );
    assert.equal(resolveDisplayVersion({ packageVersion: "3.1.0" }), "v3.1.0");
  });
});

describe("formatFeedCraftVersion", () => {
  it("prefixes the product name", () => {
    assert.equal(
      formatFeedCraftVersion("dev-becb6a3"),
      "FeedCraft dev-becb6a3"
    );
  });
});
