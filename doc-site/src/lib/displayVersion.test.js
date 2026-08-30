import assert from "node:assert/strict";
import { describe, it } from "node:test";
import {
  formatDocSiteFooter,
  resolveDocSiteDisplayVersion,
} from "./displayVersion.js";

describe("resolveDocSiteDisplayVersion", () => {
  it("uses package version on main and when branch is unset (local build)", () => {
    assert.equal(
      resolveDocSiteDisplayVersion({ pkgVersion: "1.0.0", branch: "main" }),
      "v1.0.0"
    );
    assert.equal(
      resolveDocSiteDisplayVersion({ pkgVersion: "1.0.0", branch: "" }),
      "v1.0.0"
    );
  });

  it("uses short sha for non-main preview branches", () => {
    assert.equal(
      resolveDocSiteDisplayVersion({
        pkgVersion: "1.0.0",
        branch: "dev",
        sha: "abc1234",
      }),
      "dev-abc1234"
    );
  });
});

describe("formatDocSiteFooter", () => {
  it("prefixes FeedCraft", () => {
    assert.equal(formatDocSiteFooter("v3.2.0"), "FeedCraft v3.2.0");
  });
});
