
const fs = require('fs');
const path = require('path');
const { createJiti } = require('jiti');

// Use jiti to load TypeScript files with path alias support
const jiti = createJiti(__filename, {
  alias: {
    '@': path.resolve(__dirname, '../src')
  }
});

const localeDir = path.resolve(__dirname, '../src/locale');
const outputDir = path.resolve(__dirname, '../src/locale/json');

if (!fs.existsSync(outputDir)) {
  fs.mkdirSync(outputDir, { recursive: true });
}

const locales = ['en-US', 'zh-CN'];

async function generate() {
  for (const locale of locales) {
    try {
      const filePath = path.join(localeDir, `${locale}.ts`);
      console.log(`Processing ${filePath}...`);

      // Load the module using jiti
      const module = jiti(filePath);
      const messages = module.default;

      const jsonPath = path.join(outputDir, `${locale}.json`);
      fs.writeFileSync(jsonPath, JSON.stringify(messages, null, 2));
      console.log(`Generated ${jsonPath}`);
    } catch (e) {
      console.error(`Error processing ${locale}:`, e);
      process.exit(1);
    }
  }
}

generate();
