<template>
  <div class="py-8 px-16">
    <h1 class="text-3xl font-bold">LLM Debug Page</h1>
    <a-card class="my-2" title="输入链接">
      <p>输入要预览的RSS源地址 支持RSS/ATOM</p>
      <a-space>
        <a-input
          v-model="articleUrl"
          type="text"
          placeholder="Enter article URL"
        />
        <a-button @click="fetchArticle">Submit</a-button>
      </a-space>
    </a-card>
    <a-card title="结果预览" class="my-4">
      <div v-if="articleContent">
        <p>Article Content: {{ articleContent }}</p>
        <p>Is Advertorial: {{ isAdvertorial }}</p>
      </div>
      <a-empty v-else />
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { Message } from '@arco-design/web-vue';

const articleUrl = ref('');
const articleContent = ref('');
const isAdvertorial = ref(false);

import axios from 'axios';

async function fetchArticle() {
  try {
    const resp = await axios.post(`${baseUrl}/craft-debug/advertorial`, {
      url: articleUrl.value,
    });
    articleContent.value = resp.data.article_content;
    isAdvertorial.value = resp.data.is_advertorial;
  } catch (error) {
    Message.warning(error?.toString() ?? 'Unknown Error');
    console.error(error);
  }
}
</script>

<script lang="ts">
  export default {
    name: 'LLMDebug',
  };
</script>
