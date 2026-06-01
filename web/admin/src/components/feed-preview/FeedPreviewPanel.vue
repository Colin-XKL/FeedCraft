<template>
  <div class="feed-preview-panel">
    <a-spin :loading="loading" style="width: 100%">
      <a-alert v-if="errorMessage" type="error" class="mb-4" show-icon>
        {{ errorMessage }}
      </a-alert>
      <FeedViewContainer v-if="feedContent" :feed-data="feedContent" />
      <a-empty
        v-else-if="!loading && !errorMessage"
        :description="emptyDescription || t('feedViewer.resultPreview')"
      />
    </a-spin>
  </div>
</template>

<script lang="ts" setup>
  import { ref, watch } from 'vue';
  import { useI18n } from 'vue-i18n';
  import { previewFeed, type FeedViewerPreview } from '@/api/feed_viewer';
  import FeedViewContainer from '@/views/dashboard/feed_viewer/feed_view_container.vue';

  const props = withDefaults(
    defineProps<{
      inputUri?: string;
      autoLoad?: boolean;
      emptyDescription?: string;
    }>(),
    {
      inputUri: '',
      autoLoad: false,
      emptyDescription: '',
    }
  );

  const { t } = useI18n();
  const loading = ref(false);
  const errorMessage = ref('');
  const feedContent = ref<FeedViewerPreview | null>(null);
  let requestSeq = 0;

  const resetState = () => {
    requestSeq += 1;
    feedContent.value = null;
    errorMessage.value = '';
  };

  const loadPreview = async (inputUri = props.inputUri) => {
    const uri = inputUri?.trim();
    if (!uri) {
      resetState();
      return;
    }

    const currentSeq = requestSeq + 1;
    requestSeq = currentSeq;
    loading.value = true;
    errorMessage.value = '';

    try {
      const response = await previewFeed(uri);
      if (currentSeq !== requestSeq) return;
      feedContent.value = response.data ?? null;
      if (!feedContent.value) {
        errorMessage.value = response.msg || t('topic.inputPreview.failed');
      }
    } catch (err: any) {
      if (currentSeq !== requestSeq) return;
      feedContent.value = null;
      errorMessage.value = err.message || t('topic.inputPreview.failed');
    } finally {
      if (currentSeq === requestSeq) {
        loading.value = false;
      }
    }
  };

  watch(
    () => props.inputUri,
    (uri) => {
      if (!props.autoLoad) {
        resetState();
        return;
      }
      loadPreview(uri);
    }
  );

  defineExpose({
    loadPreview,
    resetState,
    loading,
    feedContent,
    errorMessage,
  });
</script>

<script lang="ts">
  export default {
    name: 'FeedPreviewPanel',
  };
</script>
