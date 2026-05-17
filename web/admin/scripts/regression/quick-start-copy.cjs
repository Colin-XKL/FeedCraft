const fs = require('fs');
const path = require('path');

const projectRoot = path.resolve(__dirname, '..', '..');
const componentPath = path.join(
  projectRoot,
  'src/views/dashboard/url_generator/url_generator.vue'
);
const zhLocalePath = path.join(projectRoot, 'src/locale/zh-CN/urlGenerator.ts');
const enLocalePath = path.join(projectRoot, 'src/locale/en-US/urlGenerator.ts');

const component = fs.readFileSync(componentPath, 'utf8');
const zhLocale = fs.readFileSync(zhLocalePath, 'utf8');
const enLocale = fs.readFileSync(enLocalePath, 'utf8');

const failures = [];

if (!component.includes("import { useClipboard } from '@vueuse/core';")) {
  failures.push('quick_start should use VueUse clipboard composable');
}

if (!component.includes('legacy: true')) {
  failures.push('quick_start copy should enable legacy clipboard fallback');
}

if (component.includes('navigator.clipboard')) {
  failures.push(
    'quick_start copy should not directly depend on navigator.clipboard'
  );
}

if (!component.includes("t('urlGenerator.copyError')")) {
  failures.push(
    'quick_start copy failures should use copy-specific error text'
  );
}

if (!zhLocale.includes("'urlGenerator.copyError'")) {
  failures.push('zh-CN locale should define urlGenerator.copyError');
}

if (!enLocale.includes("'urlGenerator.copyError'")) {
  failures.push('en-US locale should define urlGenerator.copyError');
}

if (failures.length > 0) {
  console.error('quick_start copy regression check failed:');
  failures.forEach((failure) => console.error(`- ${failure}`));
  process.exit(1);
}

console.log('quick_start copy regression check passed');
