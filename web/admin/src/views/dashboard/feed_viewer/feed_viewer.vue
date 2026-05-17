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
          @input="errorMessage = ''"
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
      <a-alert v-if="errorMessage" type="error" class="mb-4" show-icon>
        {{ errorMessage }}
      </a-alert>
      <div v-if="feedContent">
        <FeedViewContainer :feed-data="feedContent" />
      </div>
      <a-empty v-else-if="!errorMessage" />
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { ref, onMounted } from 'vue';
  import FeedViewContainer from '@/views/dashboard/feed_viewer/feed_view_container.vue';
  import XHeader from '@/components/header/x-header.vue';
  import { useI18n } from 'vue-i18n';
  import { useRoute } from 'vue-router';
  import { previewFeed, type FeedViewerPreview } from '@/api/feed_viewer';

  const { t } = useI18n();
  const route = useRoute();

  const feedUrl = ref('');
  const feedContent = ref<FeedViewerPreview | null>(null);
  const errorMessage = ref('');
  const isLoading = ref(false);

  async function fetchFeed() {
    if (!feedUrl.value) return;
    isLoading.value = true;
    errorMessage.value = '';
    try {
      const response = await previewFeed(feedUrl.value);
      feedContent.value = response.data;
    } catch (error) {
      feedContent.value = null;
      errorMessage.value =
        error instanceof Error
          ? error.message
          : t('feedViewer.message.unknownError');
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
