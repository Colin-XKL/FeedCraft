<template>
  <a-modal
    v-model:visible="visible"
    :title="item.title || t('feedViewer.item.detail')"
    :width="800"
    :footer="false"
    unmount-on-close
    @cancel="emit('close')"
  >
    <template #title>
      <div class="flex items-center justify-between w-full pr-6">
        <span class="font-semibold truncate flex-1 mr-4">
          {{ item.title || t('feedViewer.item.detail') }}
        </span>
        <a-radio-group v-model="viewMode" type="button" size="small">
          <a-radio value="normal">{{ t('feedViewer.viewModeNormal') }}</a-radio>
          <a-radio value="rich">{{ t('feedViewer.viewModeRich') }}</a-radio>
          <a-radio value="html">{{ t('feedViewer.viewModeHtml') }}</a-radio>
        </a-radio-group>
      </div>
    </template>

    <!-- Article meta -->
    <div class="mb-4 flex flex-wrap gap-3 text-xs text-gray-400">
      <span v-if="formattedDate">{{ formattedDate }}</span>
      <a
        v-if="item.link"
        :href="item.link"
        target="_blank"
        rel="noopener noreferrer"
        class="text-blue-500 hover:underline"
      >
        {{ t('feedViewer.item.viewDetail') }} ↗
      </a>
    </div>

    <FeedItemContent :item="item" :view-mode="viewMode" />
  </a-modal>
</template>

<script lang="ts" setup>
  import { computed, ref } from 'vue';
  import dayjs from 'dayjs';
  import { useI18n } from 'vue-i18n';
  import type { FeedViewerPreviewItem } from '@/api/feed_viewer';
  import FeedItemContent from './FeedItemContent.vue';
  import type { ViewMode } from './FeedItemCard.vue';

  interface Props {
    item: FeedViewerPreviewItem;
    initialViewMode?: ViewMode;
  }

  const props = withDefaults(defineProps<Props>(), {
    initialViewMode: 'normal',
  });

  const emit = defineEmits<{
    (e: 'close'): void;
  }>();

  const { t } = useI18n();
  const visible = ref(true);
  const viewMode = ref<ViewMode>(props.initialViewMode);

  const formattedDate = computed(() => {
    const dateValue = props.item.isoDate || props.item.pubDate;
    if (!dateValue) return '';
    const parsed = dayjs(dateValue);
    return parsed.isValid() ? parsed.format('YYYY-MM-DD HH:mm:ss') : dateValue;
  });
</script>
