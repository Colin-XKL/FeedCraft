<template>
  <div class="py-8 px-16">
    <x-header
      :title="t('menu.feedViewer')"
      :description="t('feedViewer.description')"
    >
    </x-header>

    <a-card class="my-2" :title="t('feedViewer.inputLink')">
      <p>{{ t('feedViewer.inputTip') }}</p>
      <a-space>
        <a-input
          v-model="feedUrl"
          type="text"
          :placeholder="t('feedViewer.placeholder')"
          allow-clear
          @press-enter="fetchFeed"
        >
          <template #append>
            <a-button @click="handlePaste">
              <template #icon><icon-paste /></template>
            </a-button>
          </template>
        </a-input>
        <a-button
          :loading="isLoading"
          :disabled="!feedUrl"
          @click="fetchFeed"
          >{{ t('feedViewer.preview') }}</a-button
        >
      </a-space>
    </a-card>
    <a-card
      :title="t('feedViewer.resultPreview')"
      class="my-4"
      :loading="isLoading"
    >
      <div v-if="feedContent">
        <FeedViewContainer :feed-data="feedContent" />
      </div>
      <a-empty v-else />
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { computed, ref, onMounted } from 'vue';
  import Parser from 'rss-parser';
  import { Message } from '@arco-design/web-vue';
  import { IconPaste } from '@arco-design/web-vue/es/icon';
  import FeedViewContainer from '@/views/dashboard/feed_viewer/feed_view_container.vue';
  import XHeader from '@/components/header/x-header.vue';
  import { useI18n } from 'vue-i18n';
  import { useRoute } from 'vue-router';

  const { t } = useI18n();
  const route = useRoute();

  const feedUrl = ref('');
  const feedContent = ref<any>(null);
  const isLoading = ref(false);
  const baseUrl = import.meta.env.VITE_API_BASE_URL ?? '';

  async function handlePaste() {
    try {
      const text = await navigator.clipboard.readText();
      if (text) {
        feedUrl.value = text;
      }
    } catch (error) {
      Message.warning('Failed to read clipboard');
    }
  }

  async function fetchFeed() {
    if (!feedUrl.value) return;
    isLoading.value = true;
    const requestUrl = `${baseUrl}/craft/proxy?input_url=${encodeURIComponent(
      feedUrl.value,
    )}`;
    try {
      const parser = new Parser();
      const resp = await fetch(requestUrl);
      const feed = await parser.parseString(await resp?.text()).then((resp) => {
        return resp;
      });
      feedContent.value = feed;
    } catch (error) {
      Message.warning(
        error?.toString() ?? t('feedViewer.message.unknownError'),
      );
      console.error(error);
    } finally {
      isLoading.value = false;
    }
  }

  onMounted(() => {
    if (route.query.url) {
      feedUrl.value = route.query.url as string;
      fetchFeed();
    }
  });
</script>

<script lang="ts">
  export default {
    name: 'FeedViewer',
  };
</script>
