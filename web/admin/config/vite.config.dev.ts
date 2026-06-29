import { mergeConfig } from 'vite';
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
        '/craft': {
          target: 'http://localhost:8080',
          changeOrigin: true,
        },
        '/recipe': {
          target: 'http://localhost:8080',
          changeOrigin: true,
        },
      },
    },
    plugins: [
      codeInspectorPlugin({
        bundler: 'vite',
      }),
    ],
  },
  baseConfig
);
