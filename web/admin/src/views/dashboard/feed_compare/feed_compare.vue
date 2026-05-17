<template>
  <div class="feed-compare-page">
    <x-header
      :title="t('menu.feedCompare')"
      :description="t('feedCompare.description')"
    ></x-header>

    <a-card class="my-2" :title="t('feedCompare.inputLink')">
      <div class="feed-compare-controls">
        <a-input
          v-model="feedUrl"
          type="text"
          class="feed-compare-url"
          :placeholder="t('feedCompare.placeholder')"
          allow-clear
          @input="clearErrors"
          @keyup.enter="compareFeeds"
        />
        <CraftFlowSelect
          v-model="selectedCraft"
          mode="single"
          class="feed-compare-craft"
          @change="clearErrors"
        />
        <a-button
          class="feed-compare-action"
          :loading="isLoading"
          type="primary"
          @click="compareFeeds"
          >{{ t('feedCompare.compare') }}
        </a-button>
      </div>
    </a-card>

    <a-row :gutter="[24, 24]">
      <a-col :span="12" :xs="24" :lg="12">
        <a-card :title="t('feedCompare.originalFeed')" :loading="isLoading">
          <a-alert v-if="originalFeedError" type="error" class="mb-4" show-icon>
            {{ originalFeedError }}
          </a-alert>
          <div v-if="originalFeedContent">
            <FeedViewContainer :feed-data="originalFeedContent" />
          </div>
          <a-empty v-else-if="!originalFeedError" />
        </a-card>
      </a-col>
      <a-col :span="12" :xs="24" :lg="12">
        <a-card :title="t('feedCompare.craftAppliedFeed')" :loading="isLoading">
          <a-alert
            v-if="craftAppliedFeedError"
            type="error"
            class="mb-4"
            show-icon
          >
            {{ craftAppliedFeedError }}
          </a-alert>
          <div v-if="craftAppliedFeedContent">
            <FeedViewContainer :feed-data="craftAppliedFeedContent" />
          </div>
          <a-empty v-else-if="!craftAppliedFeedError" />
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import FeedViewContainer from '@/views/dashboard/feed_viewer/feed_view_container.vue';
  import XHeader from '@/components/header/x-header.vue';
  import CraftFlowSelect from '@/views/dashboard/craft_flow/CraftFlowSelect.vue';
  import { useI18n } from 'vue-i18n';
  import { previewFeed, type FeedViewerPreview } from '@/api/feed_viewer';

  const { t } = useI18n();

  const feedUrl = ref('');
  const selectedCraft = ref('');
  const originalFeedContent = ref<FeedViewerPreview | null>(null);
  const craftAppliedFeedContent = ref<FeedViewerPreview | null>(null);
  const originalFeedError = ref('');
  const craftAppliedFeedError = ref('');
  const isLoading = ref(false);

  function clearErrors() {
    originalFeedError.value = '';
    craftAppliedFeedError.value = '';
  }

  async function compareFeeds() {
    if (!feedUrl.value || !selectedCraft.value) {
      Message.warning(t('feedCompare.message.inputRequired'));
      return;
    }

    isLoading.value = true;
    clearErrors();
    originalFeedContent.value = null;
    craftAppliedFeedContent.value = null;

    const [originalResult, craftedResult] = await Promise.allSettled([
      previewFeed(feedUrl.value),
      previewFeed(feedUrl.value, { craftName: selectedCraft.value }),
    ]);

    if (originalResult.status === 'fulfilled') {
      originalFeedContent.value = originalResult.value.data;
    } else {
      originalFeedError.value =
        originalResult.reason instanceof Error
          ? originalResult.reason.message
          : t('feedCompare.message.unknownError');
    }

    if (craftedResult.status === 'fulfilled') {
      craftAppliedFeedContent.value = craftedResult.value.data;
    } else {
      craftAppliedFeedError.value =
        craftedResult.reason instanceof Error
          ? craftedResult.reason.message
          : t('feedCompare.message.unknownError');
    }

    isLoading.value = false;
  }
</script>

<script lang="ts">
  export default {
    name: 'FeedCompare',
  };
</script>

<style scoped>
  .feed-compare-page {
    padding: 32px clamp(20px, 4vw, 64px);
  }

  .feed-compare-controls {
    display: grid;
    grid-template-columns: minmax(20rem, 1fr) minmax(14rem, 20rem) auto;
    gap: 12px;
    align-items: center;
  }

  .feed-compare-url,
  .feed-compare-craft {
    width: 100%;
    min-width: 0;
  }

  @media (max-width: 768px) {
    .feed-compare-page {
      padding: 24px 16px;
    }

    .feed-compare-controls {
      grid-template-columns: 1fr;
    }

    .feed-compare-action {
      width: 100%;
    }
  }
</style>
