<template>
  <div class="py-8 px-16">
    <x-header
      :title="t('menu.feedViewer')"
      :description="t('feedViewer.description')"
    > </x-header>

    <a-card class="my-2" :title="t('feedViewer.inputLink')">
      <p>{{ t('feedViewer.inputTip') }}</p>
      <a-space>
        <a-input
          v-model="feedUrl"
          type="text"
          :placeholder="t('feedViewer.placeholder')"
        />
        <a-button :loading="isLoading" @click="fetchFeed">{{ t('feedViewer.preview') }}</a-button>
      </a-space>
    </a-card>
    <a-card :title="t('feedViewer.resultPreview')" class="my-4" :loading="isLoading">
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
  import FeedViewContainer from '@/views/dashboard/feed_viewer/feed_view_container.vue';
  import XHeader from '@/components/header/x-header.vue';
  import { useI18n } from 'vue-i18n';

  const { t } = useI18n();

  const feedUrl = ref('');
  const feedContent = ref<any>(null);
  const isLoading = ref(false);
  const baseUrl = import.meta.env.VITE_API_BASE_URL ?? '';

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
      Message.warning(error?.toString() ?? t('feedViewer.message.unknownError'));
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
