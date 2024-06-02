<template>
  <div class="py-8 px-16">
    <x-header title="RSS Feed Preview" description="RSS Feed Preview">
    </x-header>

    <a-card class="my-2" title="输入链接">
      <p>输入要预览的RSS源地址 支持RSS/ATOM</p>
      <a-space>
        <a-input
          v-model="feedUrl"
          type="text"
          placeholder="Enter RSS feed URL"
        />
        <a-button :loading="isLoading" @click="fetchFeed">Preview</a-button>
      </a-space>
    </a-card>
    <a-card title="结果预览" class="my-4" :loading="isLoading">
      <div v-if="feedContent">
        <FeedViewContainer :feed-data="feedContent" />
      </div>
      <a-empty v-else />
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { computed, ref } from 'vue';
  import Parser from 'rss-parser';
  import { Message } from '@arco-design/web-vue';
  import FeedViewContainer from '@/views/dashboard/feed-viewer/feed_view_container.vue';
  import XHeader from '@/components/header/x-header.vue';

  const feedUrl = ref('');
  const feedContent = ref<any>(null);
  const isLoading = ref(false);
  const baseUrl = import.meta.env.VITE_API_BASE_URL;

  async function fetchFeed() {
    isLoading.value = true;
    const requestUrl = `${baseUrl}/craft/proxy?input_url=${encodeURIComponent(
      feedUrl.value
    )}`;
    try {
      const parser = new Parser();
      const resp = await fetch(requestUrl);
      const feed = await parser.parseString(await resp?.text()).then((resp) => {
        return resp;
      });
      feedContent.value = feed;
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
    name: 'FeedViewer',
  };
</script>
