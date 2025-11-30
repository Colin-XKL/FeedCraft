/* eslint-disable */
import fs from 'fs';
import path from 'path';

// Since we cannot easily invoke the project's own TS environment (due to aliases and complexity),
// and we know the files are essentially static data structures (hopefully),
// we might be able to use a simpler approach.
// But the files use 'export default' and imports.
//
// The 'vite' build system can bundle these files for us!
// We can use 'vite build' in library mode or similar to produce a JS bundle for the locales.

// Let's create a temporary vite config to build the locales.
import { build } from 'vite';
import { fileURLToPath } from 'url';
import Module from 'module';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

async function generate() {
  const rootDir = path.resolve(__dirname, '..');
  const outputDir = path.resolve(rootDir, 'src/locales-generated');

  console.log('Building locales...');

  try {
      // Build en-US
      await build({
        configFile: false,
        root: rootDir,
        build: {
          lib: {
            entry: path.resolve(rootDir, 'src/locale/en-US.ts'),
            formats: ['cjs'],
            fileName: () => 'en-US.js',
            name: 'enUS'
          },
          outDir: path.resolve(rootDir, 'dist-temp-locales'),
          emptyOutDir: false, // Don't wipe on second pass
          rollupOptions: {
            external: ['vue', 'vue-i18n'],
          }
        },
        resolve: {
          alias: {
            '@': path.resolve(rootDir, 'src')
          }
        },
        logLevel: 'silent'
      });

      // Build zh-CN
      await build({
        configFile: false,
        root: rootDir,
        build: {
            lib: {
              entry: path.resolve(rootDir, 'src/locale/zh-CN.ts'),
              formats: ['cjs'],
              fileName: () => 'zh-CN.js',
              name: 'zhCN'
            },
            outDir: path.resolve(rootDir, 'dist-temp-locales'),
            emptyOutDir: false,
            rollupOptions: {
                external: ['vue', 'vue-i18n']
            }
          },
          resolve: {
            alias: {
              '@': path.resolve(rootDir, 'src')
            }
          },
          logLevel: 'silent'
      });

      console.log('Build complete. Generating JSON...');

      if (!fs.existsSync(outputDir)) {
          fs.mkdirSync(outputDir, { recursive: true });
      }

      // Use createRequire to load CJS modules in ESM context (if we are in ESM)
      // or just require if we are in CJS.
      // TypeScript compiler check might fail if we don't handle the type definition for Module.createRequire
      // but we can cast or use require directly if this is running in CJS.

      // @ts-ignore
      const createRequire = Module.createRequire || Module.default.createRequire;
      const require = createRequire(import.meta.url);

      const enPath = path.resolve(rootDir, 'dist-temp-locales/en-US.js');
      const cnPath = path.resolve(rootDir, 'dist-temp-locales/zh-CN.js');

      // Clear cache just in case
      delete require.cache[require.resolve(enPath)];
      delete require.cache[require.resolve(cnPath)];

      const en = require(enPath);
      const cn = require(cnPath);

      const enData = en.default || en;
      const cnData = cn.default || cn;

      fs.writeFileSync(path.join(outputDir, 'en-US.json'), JSON.stringify(enData, null, 2));
      fs.writeFileSync(path.join(outputDir, 'zh-CN.json'), JSON.stringify(cnData, null, 2));

      console.log('I18n JSON files generated successfully at ' + outputDir);

      // Cleanup
      fs.rmSync(path.resolve(rootDir, 'dist-temp-locales'), { recursive: true, force: true });

  } catch (e) {
      console.error('Error generating locales:', e);
      process.exit(1);
  }
}

generate();
