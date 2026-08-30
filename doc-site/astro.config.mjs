import { readFileSync } from "node:fs";
import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";
import starlightCatppuccin from "@catppuccin/starlight";
import { resolveDocSiteDisplayVersion } from "./src/lib/displayVersion.js";

const pkg = JSON.parse(
  readFileSync(new URL("./package.json", import.meta.url), "utf-8")
);
const branch = process.env.VERCEL_GIT_COMMIT_REF || "";
const sha = (process.env.VERCEL_GIT_COMMIT_SHA || "").slice(0, 7);
const displayVersion = resolveDocSiteDisplayVersion({
  pkgVersion: pkg.version,
  branch,
  sha,
});

// https://astro.build/config
export default defineConfig({
  site: "https://feed-craft-doc.vercel.app/",
  vite: {
    define: {
      "import.meta.env.PUBLIC_FEEDCRAFT_VERSION":
        JSON.stringify(displayVersion),
    },
  },
  integrations: [
    starlight({
      title: "FeedCraft",
      defaultLocale: "en",
      social: [
        {
          icon: "github",
          label: "GitHub",
          href: "https://github.com/Colin-XKL/FeedCraft",
        },
      ],
      plugins: [
        starlightCatppuccin({
          dark: { flavor: "mocha", accent: "sapphire" },
          light: { flavor: "latte", accent: "teal" },
        }),
      ],
      locales: {
        en: {
          label: "English",
          lang: "en",
        },
        zh: {
          label: "简体中文",
          lang: "zh-CN",
        },
        "zh-tw": {
          label: "繁體中文",
          lang: "zh-TW",
        },
      },
      components: {
        Footer: "./src/components/Footer.astro",
      },
      sidebar: [
        {
          label: "Quick Start",
          translations: {
            "zh-CN": "快速开始",
            "zh-TW": "快速開始",
          },
          autogenerate: { directory: "guides/start" },
        },
        {
          label: "Advanced Customization",
          translations: {
            "zh-CN": "进阶自定义",
            "zh-TW": "進階定制",
          },
          autogenerate: { directory: "guides/advanced" },
        },
      ],
    }),
  ],
});
