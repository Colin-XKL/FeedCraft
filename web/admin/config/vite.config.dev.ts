import { mergeConfig } from 'vite';
import eslint from 'vite-plugin-eslint';
import { codeInspectorPlugin } from 'code-inspector-plugin';
import baseConfig from './vite.config.base';

export default mergeConfig(
  {
    mode: 'development',
    server: {
      open: true,
      fs: {
        strict: true,
      },
      proxy: {
        '/api': {
          target: 'http://localhost:8080',
          changeOrigin: true,
        },
      },
    },
    plugins: [
      eslint({
        cache: false,
        include: ['src/**/*.ts', 'src/**/*.tsx', 'src/**/*.vue'],
        exclude: ['node_modules'],
      }),
      codeInspectorPlugin({
        bundler: 'vite',
      }),
    ],
  },
  baseConfig,
);
