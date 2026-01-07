<template>
  <div class="container flex flex-col items-center justify-center min-h-screen bg-gray-50 dark:bg-gray-900">
    <div class="text-center">
      <h1 class="text-5xl font-bold text-gray-900 dark:text-white mb-4">
        FeedCraft
      </h1>
      <p class="text-xl text-gray-600 dark:text-gray-300 mb-8 max-w-2xl mx-auto">
        A powerful RSS feed generator and management system. Turn any website into a structured feed with advanced customization options.
      </p>

      <div class="flex gap-4 justify-center">
        <a-button type="primary" size="large" @click="goToLogin">
          <template #icon>
            <icon-import />
          </template>
          Enter Admin Dashboard
        </a-button>

        <a-button size="large" @click="goToGithub">
          <template #icon>
            <icon-github />
          </template>
          GitHub
        </a-button>
      </div>
    </div>

    <div class="mt-16 grid grid-cols-1 md:grid-cols-3 gap-8 max-w-5xl px-4">
      <div class="p-6 bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-100 dark:border-gray-700">
        <h3 class="text-lg font-semibold mb-2 text-gray-800 dark:text-gray-200">RSS Generation</h3>
        <p class="text-gray-600 dark:text-gray-400">Convert standard web pages into RSS feeds with customizable selectors and rules.</p>
      </div>
      <div class="p-6 bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-100 dark:border-gray-700">
        <h3 class="text-lg font-semibold mb-2 text-gray-800 dark:text-gray-200">AI Powered</h3>
        <p class="text-gray-600 dark:text-gray-400">Leverage LLMs for content cleaning, summarization, and translation automatically.</p>
      </div>
      <div class="p-6 bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-100 dark:border-gray-700">
        <h3 class="text-lg font-semibold mb-2 text-gray-800 dark:text-gray-200">Flow Management</h3>
        <p class="text-gray-600 dark:text-gray-400">Design complex data processing flows with a visual editor and reusable atoms.</p>
      </div>
    </div>

    <footer class="absolute bottom-4 text-gray-400 text-sm">
      &copy; {{ new Date().getFullYear() }} FeedCraft. All rights reserved.
    </footer>
  </div>
</template>

<script lang="ts" setup>
  import { useRouter } from 'vue-router';
  import { IconGithub, IconImport } from '@arco-design/web-vue/es/icon';

  const router = useRouter();

  const goToLogin = () => {
    router.push({ name: 'login' });
  };

  const goToGithub = () => {
    window.open('https://github.com/colin-xkl/feed-craft', '_blank');
  };
</script>

<style scoped>
  .container {
    width: 100%;
    max-width: 100%;
  }
</style>
