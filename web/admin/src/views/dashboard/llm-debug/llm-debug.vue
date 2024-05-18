<template>
  <div class="py-8 px-16">
    <h1 class="text-3xl font-bold">LLM Debug Page - Check If Advertorial</h1>
    <p>通过大模型能力,检测一个文章是不是广告或是推广的软文</p>
    <a-card class="my-2" title="输入链接">
      <p>输入要预览的文章链接</p>
      <a-space>
        <a-input
          v-model="articleUrl"
          type="text"
          placeholder="Enter article URL"
          class="w-full"
        />
        <a-button :loading="isLoading" @click="fetchArticle">Submit</a-button>
      </a-space>
    </a-card>
    <a-card title="结果预览" class="my-4" :loading="isLoading">
      <div v-if="articleContent?.length !== 0">
        <p class="mb-2"><b>Is Advertorial:</b>{{ isAdvertorial }}</p>
        <hr class="my-1" />
        <p class="my-1"><b>Article Content: </b></p>
        <p>{{ articleContent }}</p>
      </div>
      <a-empty v-else />
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  import { Message } from '@arco-design/web-vue';

  import axios from 'axios';

  const articleUrl = ref('');
  const articleContent = ref('');
  const isAdvertorial = ref(false);
  const isLoading = ref(false);

  async function fetchArticle() {
    isLoading.value = true;
    try {
      const resp = await axios.post(`/api/admin/craft-debug/advertorial`, {
        url: articleUrl.value,
      });
      articleContent.value = resp.data.article_content;
      isAdvertorial.value = resp.data.is_advertorial;
    } catch (error) {
      Message.warning(error?.toString() ?? 'Unknown Error');
      console.error(error);
    } finally {
      isLoading.value = false;
    }
  }
</script>

<script lang="ts">
  export default {
    name: 'LLMDebug',
  };
</script>
