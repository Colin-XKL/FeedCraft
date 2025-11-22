import { mergeConfig, loadEnv, defineConfig } from 'vite';
import eslint from 'vite-plugin-eslint';
import { codeInspectorPlugin } from 'code-inspector-plugin';
import baseConfig from './vite.config.base';

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '');

  return mergeConfig(
    {
      mode: 'development',
      server: {
        open: true,
        fs: {
          strict: true,
        },
        proxy: {
          '/api': {
            target: env.VITE_API_PROXY_TARGET || 'http://localhost:8080',
            changeOrigin: true,
            // rewrite: (path) => path.replace(/^\/api/, ''),
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
    baseConfig
  );
});
