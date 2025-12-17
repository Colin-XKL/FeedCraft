module.exports = {
  extends: [
    './.eslintrc.js',
    'plugin:@intlify/vue-i18n/recommended',
  ],
  settings: {
    'vue-i18n': {
      localeDir: './src/locale/**/*.{json,json5,yaml,yml,js,ts}',
      messageSyntaxVersion: '^9.0.0',
    },
  },
  rules: {
    '@intlify/vue-i18n/no-raw-text': ['error', {
      ignoreNodes: ['v-icon', 'x-icon', 'icon-font'],
      ignorePattern: '^[-#:()&]+$',
    }],
    '@intlify/vue-i18n/no-missing-keys': 'off',
    '@intlify/vue-i18n/no-v-html': 'off',
  },
};
