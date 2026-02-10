import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";
import starlightCatppuccin from '@catppuccin/starlight'

// https://astro.build/config
export default defineConfig({
  integrations: [
    starlight({
      title: "FeedCraft",
      defaultLocale: "en",
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
