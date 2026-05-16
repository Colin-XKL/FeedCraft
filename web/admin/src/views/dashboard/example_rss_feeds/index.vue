<template>
  <div class="py-8 px-16">
    <x-header
      :title="t('menu.exampleRssFeeds')"
      :description="t('exampleRssFeeds.description')"
    />

    <a-card class="my-2">
      <a-alert type="info" show-icon class="mb-4">
        {{ t('exampleRssFeeds.windowNote') }}
      </a-alert>

      <a-spin :loading="loading" class="w-full">
        <a-list :bordered="false">
          <a-list-item v-for="feed in feeds" :key="feed.slug">
            <a-list-item-meta
              :title="feed.title"
              :description="feed.description"
            />
            <template #actions>
              <a-space wrap>
                <a-input
                  class="feed-url"
                  readonly
                  :model-value="buildFeedUrl(feed.path)"
                />
                <a-button @click="copyFeedUrl(feed.path)">
                  {{ t('exampleRssFeeds.copy') }}
                </a-button>
                <a-button @click="openFeed(feed.path)">
                  {{ t('exampleRssFeeds.open') }}
                </a-button>
                <a-button type="primary" @click="previewFeed(feed.path)">
                  {{ t('exampleRssFeeds.preview') }}
                </a-button>
              </a-space>
            </template>
          </a-list-item>
        </a-list>
      </a-spin>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { onMounted, ref } from 'vue';
  import { useI18n } from 'vue-i18n';
  import { useRouter } from 'vue-router';
  import { Message } from '@arco-design/web-vue';
  import XHeader from '@/components/header/x-header.vue';
  import buildPublicFeedUrl from '@/utils/publicFeedUrl';
  import {
    listExampleRssFeeds,
    type ExampleRssFeed,
  } from '@/api/example_rss_feeds';

  const { t } = useI18n();
  const router = useRouter();
  const loading = ref(false);
  const feeds = ref<ExampleRssFeed[]>([]);

  function buildFeedUrl(path: string) {
    return buildPublicFeedUrl(path);
  }

  async function loadFeeds() {
    loading.value = true;
    try {
      const response = await listExampleRssFeeds();
      feeds.value = response.data ?? [];
    } catch (error) {
      Message.error(t('exampleRssFeeds.loadError'));
    } finally {
      loading.value = false;
    }
  }

  async function copyFeedUrl(path: string) {
    try {
      await navigator.clipboard.writeText(buildFeedUrl(path));
      Message.success(t('exampleRssFeeds.copied'));
    } catch (error) {
      Message.error(t('exampleRssFeeds.copyError'));
    }
  }

  function openFeed(path: string) {
    window.open(buildFeedUrl(path), '_blank', 'noopener,noreferrer');
  }

  function previewFeed(path: string) {
    router.push({
      name: 'FeedViewer',
      query: {
        url: buildFeedUrl(path),
      },
    });
  }

  onMounted(loadFeeds);
</script>

<script lang="ts">
  export default {
    name: 'ExampleRssFeeds',
  };
</script>

<style scoped>
  .feed-url {
    width: min(52vw, 560px);
  }
</style>
