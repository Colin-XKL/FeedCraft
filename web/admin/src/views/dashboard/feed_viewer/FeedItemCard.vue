<template>
  <a-card class="feed-item-card" :body-style="{ padding: '16px' }">
    <!-- Card header: title + meta -->
    <div class="flex items-start justify-between gap-3">
      <div class="flex-1 min-w-0">
        <a
          v-if="item.link"
          :href="item.link"
          target="_blank"
          rel="noopener noreferrer"
          class="no-underline hover:underline"
        >
          <h3 class="text-base font-semibold leading-snug m-0 text-inherit">
            {{ item.title || t('feedViewer.item.noContent') }}
          </h3>
        </a>
        <h3 v-else class="text-base font-semibold leading-snug m-0">
          {{ item.title || t('feedViewer.item.noContent') }}
        </h3>
        <div v-if="formattedDate" class="text-xs text-gray-400 mt-1">
          {{ formattedDate }}
        </div>
      </div>

      <!-- Actions -->
      <div class="flex items-center gap-1 flex-shrink-0">
        <a-tooltip :content="t('feedViewer.item.viewDetail')">
          <a-button
            size="mini"
            type="text"
            @click.stop="showDetailModal = true"
          >
            <template #icon>
              <icon-expand />
            </template>
          </a-button>
        </a-tooltip>
        <a-button size="mini" type="text" @click.stop="toggleExpand">
          <template #icon>
            <icon-up v-if="expanded" />
            <icon-down v-else />
          </template>
          {{
            expanded
              ? t('feedViewer.item.collapse')
              : t('feedViewer.item.expand')
          }}
        </a-button>
      </div>
    </div>

    <!-- Collapsed: snippet preview -->
    <div
      v-if="!expanded"
      class="mt-2 text-sm text-gray-500 leading-relaxed line-clamp-2 cursor-pointer"
      @click="toggleExpand"
    >
      {{ snippetText || t('feedViewer.item.noContent') }}
    </div>

    <!-- Expanded: full content -->
    <div v-else class="mt-3">
      <!-- View mode switcher -->
      <div class="mb-3 flex items-center gap-2">
        <a-radio-group v-model="localViewMode" type="button" size="small">
          <a-radio value="normal">{{ t('feedViewer.viewModeNormal') }}</a-radio>
          <a-radio value="rich">{{ t('feedViewer.viewModeRich') }}</a-radio>
          <a-radio value="html">{{ t('feedViewer.viewModeHtml') }}</a-radio>
        </a-radio-group>
      </div>

      <FeedItemContent :item="item" :view-mode="localViewMode" />
    </div>
  </a-card>

  <!-- Detail Modal -->
  <FeedItemDetailModal
    v-if="showDetailModal"
    :item="item"
    :initial-view-mode="localViewMode"
    @close="showDetailModal = false"
  />
</template>

<script lang="ts" setup>
  import { computed, ref, watch } from 'vue';
  import dayjs from 'dayjs';
  import { useI18n } from 'vue-i18n';
  import type { FeedViewerPreviewItem } from '@/api/feed_viewer';
  import FeedItemContent from './FeedItemContent.vue';
  import FeedItemDetailModal from './FeedItemDetailModal.vue';

  export type ViewMode = 'normal' | 'rich' | 'html';

  interface Props {
    item: FeedViewerPreviewItem;
    viewMode?: ViewMode;
  }

  const props = withDefaults(defineProps<Props>(), {
    viewMode: 'normal',
  });

  const { t } = useI18n();
  const expanded = ref(false);
  const showDetailModal = ref(false);
  const localViewMode = ref<ViewMode>(props.viewMode);

  watch(
    () => props.viewMode,
    (val) => {
      localViewMode.value = val;
    }
  );

  const formattedDate = computed(() => {
    const dateValue = props.item.isoDate || props.item.pubDate;
    if (!dateValue) return '';
    const parsed = dayjs(dateValue);
    return parsed.isValid() ? parsed.format('YYYY-MM-DD HH:mm') : dateValue;
  });

  const snippetText = computed(() => {
    return props.item.contentSnippet?.trim() || '';
  });

  function toggleExpand() {
    expanded.value = !expanded.value;
  }
</script>

<style scoped>
  .feed-item-card {
    transition: box-shadow 0.2s;
  }
  .feed-item-card:hover {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  }
  .line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
</style>
