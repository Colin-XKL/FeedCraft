// eslint-disable-next-line @typescript-eslint/no-var-requires
const baseConfig = require('./.eslintrc.js');

module.exports = {
  ...baseConfig,
  extends: [
    ...baseConfig.extends,
    'plugin:@intlify/vue-i18n/recommended',
  ],
  settings: {
    ...baseConfig.settings,
    'vue-i18n': {
      // Point to the dummy locale to avoid parsing errors with complex TS files
      localeDir: './src/locale/dummy/*.json',
      messageSyntaxVersion: '^9.0.0',
    },
  },
  rules: {
    ...baseConfig.rules,
    // Enable raw text detection
    '@intlify/vue-i18n/no-raw-text': ['error', {
      ignoreNodes: ['v-icon', 'x-icon', 'icon-font'],
      ignorePattern: '^[-#:()&]+$',
    }],
    // Disable missing keys check as it's not reliable with the dummy locale
    '@intlify/vue-i18n/no-missing-keys': 'off',
    '@intlify/vue-i18n/no-v-html': 'off',
    // Disable other rules that rely on accurate locale files
    '@intlify/vue-i18n/no-unused-keys': 'off',
    '@intlify/vue-i18n/key-format-style': 'off',
    '@intlify/vue-i18n/no-dynamic-keys': 'off',
  },
};
