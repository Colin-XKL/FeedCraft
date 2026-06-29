import { mergeConfig } from 'vite';
import baseConfig from './vite.config.base';
import configCompressPlugin from './plugin/compress';
import configVisualizerPlugin from './plugin/visualizer';
import configArcoResolverPlugin from './plugin/arcoResolver';
import configImageminPlugin from './plugin/imagemin';

export default mergeConfig(
  {
    mode: 'production',
    plugins: [
      // configCompressPlugin('gzip'),
      // configVisualizerPlugin(),
      configArcoResolverPlugin(),
      // configImageminPlugin(),
    ],
    build: {
      rollupOptions: {
        output: {
          manualChunks(id) {
            if (!id.includes('node_modules')) {
              return undefined;
            }
            if (id.includes('@arco-design/web-vue')) {
              return 'arco';
            }
            if (id.includes('echarts') || id.includes('vue-echarts')) {
              return 'chart';
            }
            if (
              id.includes('/vue/') ||
              id.includes('/vue-router/') ||
              id.includes('/pinia/') ||
              id.includes('/@vueuse/core/') ||
              id.includes('/vue-i18n/')
            ) {
              return 'vue';
            }
            return undefined;
          },
        },
      },
      chunkSizeWarningLimit: 2000,
    },
  },
  baseConfig
);
