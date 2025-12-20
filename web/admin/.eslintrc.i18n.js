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
      // Use the generated JSON files to avoid parsing errors
      localeDir: './src/locale/json/*.json',
      messageSyntaxVersion: '^9.0.0',
    },
  },
  rules: {
    ...baseConfig.rules,
    // Enable raw text detection, but it is turned off by user request
    '@intlify/vue-i18n/no-raw-text': 'off',
    // Enable missing keys check as requested
    '@intlify/vue-i18n/no-missing-keys': 'error',
    '@intlify/vue-i18n/no-v-html': 'off',
    // Disable other rules that might be too strict
    '@intlify/vue-i18n/no-unused-keys': 'off',
    '@intlify/vue-i18n/key-format-style': 'off',
    '@intlify/vue-i18n/no-dynamic-keys': 'off',
  },
};
