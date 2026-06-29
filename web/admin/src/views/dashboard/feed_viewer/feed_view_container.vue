<template>
  <div>
    <!-- Feed metadata -->
    <h2 class="text-lg font-semibold mb-1">{{ feedData.title }}</h2>
    <a-descriptions
      class="mt-3 mb-4"
      :data="feedMetaList"
      :title="t('feedViewer.feedInfo')"
      :column="1"
      size="small"
    />

    <!-- Toolbar: global view mode + item count -->
    <div class="flex items-center justify-between flex-wrap gap-2 mb-4">
      <a-radio-group v-model="viewMode" type="button" size="small">
        <a-radio value="normal">{{ t('feedViewer.viewModeNormal') }}</a-radio>
        <a-radio value="rich">{{ t('feedViewer.viewModeRich') }}</a-radio>
        <a-radio value="html">{{ t('feedViewer.viewModeHtml') }}</a-radio>
      </a-radio-group>
      <div class="text-sm text-gray-400">
        <span v-if="visibleItems.length < feedData.items?.length">
          {{ t('feedViewer.showingItems', { count: visibleItems.length }) }}
          &nbsp;/&nbsp;
          {{ t('feedViewer.totalItems', { count: feedData.items?.length }) }}
        </span>
        <span v-else>
          {{ t('feedViewer.totalItems', { count: feedData.items?.length }) }}
        </span>
      </div>
    </div>

    <!-- Article list -->
    <div class="flex flex-col gap-3">
      <FeedItemCard
        v-for="item in visibleItems"
        :key="item.guid || item.link"
        :item="item"
        :view-mode="viewMode"
      />
    </div>

    <div
      v-if="feedData.items?.length > MAX_VISIBLE_ITEMS"
      class="mt-4 text-sm text-gray-400 text-center"
    >
      {{ t('feedViewer.showingItems', { count: MAX_VISIBLE_ITEMS }) }}
      &nbsp;/&nbsp;
      {{ t('feedViewer.totalItems', { count: feedData.items?.length }) }}
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { computed, ref } from 'vue';
  import { useI18n } from 'vue-i18n';
  import type { FeedViewerPreview } from '@/api/feed_viewer';
  import FeedItemCard from '@/views/dashboard/feed_viewer/FeedItemCard.vue';
  import type { ViewMode } from '@/views/dashboard/feed_viewer/FeedItemCard.vue';

  const MAX_VISIBLE_ITEMS = 20;

  interface Props {
    feedData: FeedViewerPreview;
  }

  const props = defineProps<Props>();
  const { t } = useI18n();
  const viewMode = ref<ViewMode>('normal');

  const feedMetaList = computed(() => {
    const data = props.feedData;
    return [
      { label: 'description', value: data.description },
      { label: 'link', value: data.link },
      { label: 'feedUrl', value: data.feedUrl },
      { label: 'copyright', value: data.copyright },
      {
        label: 'image',
        value:
          data.image?.url || data.image?.title
            ? `${data.image?.title || ''} ${data.image?.url || ''}`.trim()
            : '',
      },
    ].filter((item) => item.value);
  });

  const visibleItems = computed(
    () => props.feedData.items?.slice(0, MAX_VISIBLE_ITEMS) ?? []
  );
</script>

<script lang="ts">
  export default {
    name: 'FeedViewContainer',
  };
</script>
