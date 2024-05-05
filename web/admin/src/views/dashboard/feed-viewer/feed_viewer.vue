<template>
  <div class="py-8 px-16">
    <h1 class="text-3xl font-bold">RSS Feed Preview</h1>
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
        <h2>{{ feedContent.title }}</h2>
        <a-descriptions
          style="margin-top: 20px"
          :data="feedMetaList"
          title="Feed Info"
          :column="1"
        />

        <ul>
          <li v-for="item in feedContent.items?.slice(0, 5)" :key="item.guid">
            <a-card class="my-2">
              <a-space>
                <h3 class="font-bold">{{ item.title }}</h3>
                <p>{{ dayjs(item.isoDate).format('YYYY-MM-DD hh:mm:ss') }}</p>
              </a-space>
              <a-typography-paragraph
                :ellipsis="{
                  rows: 3,
                  showTooltip: false,
                  expandable: true,
                }"
              >
                {{ item.contentSnippet }}
              </a-typography-paragraph>
            </a-card>
          </li>
        </ul>
      </div>
      <a-empty v-else />
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { computed, ref } from 'vue';
  import Parser from 'rss-parser';
  import dayjs from 'dayjs';
  import { Message } from '@arco-design/web-vue';

  const feedMetaList = computed(() => {
    return Object.keys(feedContent.value).map((key) => {
      const item = feedContent.value[key];
      return {
        label: key,
        value: item,
      };
    });
  });
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
