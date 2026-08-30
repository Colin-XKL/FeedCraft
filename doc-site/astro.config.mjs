import { readFileSync } from "node:fs";
import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";
import starlightCatppuccin from "@catppuccin/starlight";
import { resolveDisplayVersion } from "./src/lib/displayVersion.js";

const pkg = JSON.parse(
  readFileSync(new URL("./package.json", import.meta.url), "utf-8")
);
const displayVersion = resolveDisplayVersion({
  explicitVersion: process.env.APP_VERSION,
  branch: process.env.VERCEL_GIT_COMMIT_REF,
  commitSha: process.env.VERCEL_GIT_COMMIT_SHA,
  packageVersion: pkg.version,
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
      pagefind: true,
      components: {
        Footer: "./src/components/Footer.astro",
        MarkdownContent: "./src/components/MarkdownContent.astro",
      },
      sidebar: [
        {
          label: "Getting Started",
          translations: {
            "zh-CN": "上手",
            "zh-TW": "上手",
          },
          autogenerate: { directory: "start" },
        },
        {
          label: "Guides",
          translations: {
            "zh-CN": "指南",
            "zh-TW": "指南",
          },
          autogenerate: { directory: "guides" },
        },
        {
          label: "Concepts",
          translations: {
            "zh-CN": "概念",
            "zh-TW": "概念",
          },
          autogenerate: { directory: "concepts" },
        },
        {
          label: "Reference",
          translations: {
            "zh-CN": "参考",
            "zh-TW": "參考",
          },
          autogenerate: { directory: "reference" },
        },
      ],
    }),
  ],
});
