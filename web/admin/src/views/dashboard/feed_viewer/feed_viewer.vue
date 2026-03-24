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
          @keyup.enter="fetchFeed"
        />
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
  import { ref, onMounted } from 'vue';
  import Parser from 'rss-parser';
  import { Message } from '@arco-design/web-vue';
  import FeedViewContainer from '@/views/dashboard/feed_viewer/feed_view_container.vue';
  import XHeader from '@/components/header/x-header.vue';
  import { useI18n } from 'vue-i18n';
  import { useRoute } from 'vue-router';
  import { normalizeBaseUrl } from '@/utils/publicFeedUrl';

  const { t } = useI18n();
  const route = useRoute();

  const feedUrl = ref('');
  const feedContent = ref<any>(null);
  const isLoading = ref(false);

  async function fetchFeed() {
    if (!feedUrl.value) return;
    isLoading.value = true;

    let requestUrl = feedUrl.value;

    try {
      let isInternal = false;
      if (feedUrl.value.startsWith('/')) {
        isInternal = true;
      } else {
        const urlObj = new URL(feedUrl.value, window.location.origin);
        if (urlObj.origin === window.location.origin) {
          isInternal = true;
        } else {
          // If the API base URL is fully qualified and matches the input URL origin, consider it internal
          const apiBase = normalizeBaseUrl();
          if (apiBase) {
            const apiBaseObj = new URL(apiBase, window.location.origin);
            if (urlObj.origin === apiBaseObj.origin) {
              isInternal = true;
            }
          }
        }
      }

      if (!isInternal) {
        // If it's an external URL, wrap it in our proxy
        const prefix = normalizeBaseUrl();
        requestUrl = `${prefix}/craft/proxy?input_url=${encodeURIComponent(
          feedUrl.value
        )}`;
      }
    } catch (e) {
      // Invalid URL format, default to using proxy
      const prefix = normalizeBaseUrl();
      requestUrl = `${prefix}/craft/proxy?input_url=${encodeURIComponent(
        feedUrl.value
      )}`;
    }

    try {
      const parser = new Parser();
      const resp = await fetch(requestUrl);
      const feed = await parser.parseString(await resp?.text()).then((resp) => {
        return resp;
      });
      feedContent.value = feed;
    } catch (error) {
      Message.warning(
        error?.toString() ?? t('feedViewer.message.unknownError')
      );
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
