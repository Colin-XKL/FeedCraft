<template>
  <div class="py-8 px-16">
    <x-header
      title="Feed Compare"
      description="指定feed源url与要应用的craft(craft flow 或 atom都可), 对比处理前和处理后的结果"
    ></x-header>

    <a-card class="my-2" title="输入链接">
      <a-space>
        <a-input
          v-model="feedUrl"
          type="text"
          class="min-w-48"
          placeholder="Enter RSS feed URL"
        />
        <CraftFlowSelect v-model="selectedCraft" mode="single" />
        <a-button :loading="isLoading" type="primary" @click="compareFeeds"
          >Compare
        </a-button>
      </a-space>
    </a-card>

    <a-row :gutter="24">
      <a-col :span="12">
        <a-card title="Original Feed" :loading="isLoading">
          <div v-if="originalFeedContent">
            <FeedViewContainer :feed-data="originalFeedContent" />
          </div>
          <a-empty v-else />
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="Craft Applied Feed" :loading="isLoading">
          <div v-if="craftAppliedFeedContent">
            <FeedViewContainer :feed-data="craftAppliedFeedContent" />
          </div>
          <a-empty v-else />
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  import Parser from 'rss-parser';
  import { Message } from '@arco-design/web-vue';
  import FeedViewContainer from '@/views/dashboard/feed_viewer/feed_view_container.vue';
  import XHeader from '@/components/header/x-header.vue';
  import CraftFlowSelect from '@/views/dashboard/craft_flow/CraftFlowSelect.vue';

  const feedUrl = ref('');
  const selectedCraft = ref<string[]>([]);
  // const crafts = ref(['craft1', 'craft2', 'craft3']); // 这里需要从后端获取craft列表
  const originalFeedContent = ref<any>(null);
  const craftAppliedFeedContent = ref<any>(null);
  const isLoading = ref(false);
  const baseUrl = import.meta.env.VITE_API_BASE_URL ?? '';

  async function fetchFeed(url: string) {
    const parser = new Parser();
    const resp = await fetch(url);
    return parser.parseString(await resp?.text());
  }

  async function compareFeeds() {
    if (!feedUrl.value || !selectedCraft.value) {
      Message.warning('Please enter a feed URL and select a craft');
      return;
    }

    isLoading.value = true;
    try {
      originalFeedContent.value = await fetchFeed(
        `${baseUrl}/craft/proxy?input_url=${encodeURIComponent(feedUrl.value)}`
      );
      craftAppliedFeedContent.value = await fetchFeed(
        `${baseUrl}/craft/${selectedCraft.value}?input_url=${encodeURIComponent(
          feedUrl.value
        )}`
      );
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
    name: 'FeedCompare',
  };
</script>
