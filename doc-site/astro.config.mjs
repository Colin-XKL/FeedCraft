import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
  integrations: [
    starlight({
      title: 'FeedCraft',
      defaultLocale: 'root',
      locales: {
        root: {
          label: 'English',
          lang: 'en',
        },
        zh: {
          label: '简体中文',
          lang: 'zh-CN',
        },
      },
      sidebar: [
        {
          label: 'Quick Start',
          translations: {
            'zh-CN': '快速开始',
          },
          autogenerate: { directory: 'guides/start' },
        },
        {
          label: 'Advanced Customization',
          translations: {
            'zh-CN': '进阶自定义',
          },
          autogenerate: { directory: 'guides/advanced' },
        },
      ],
    }),
  ],
});
