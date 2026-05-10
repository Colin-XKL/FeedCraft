<template>
  <div>
    <h2>{{ props.feedData.title }}</h2>
    <a-descriptions
      style="margin-top: 20px"
      :data="feedMetaList"
      title="Feed 信息"
      :column="1"
    />
    <div>总数: {{ feedData.items?.length }}</div>
    <div class="my-4">
      <a-radio-group v-model="viewMode" type="button">
        <a-radio value="normal">{{ t('feedViewer.viewModeNormal') }}</a-radio>
        <a-radio value="rich">{{ t('feedViewer.viewModeRich') }}</a-radio>
      </a-radio-group>
    </div>
    <ul>
      <li
        v-for="item in feedData.items?.slice(0, 10)"
        :key="item.guid || item.link"
      >
        <a-card class="my-2">
          <a-space>
            <a
              :href="item.link"
              target="_blank"
              class="hover:text-blue-600 no-underline"
            >
              <h3 class="font-bold cursor-pointer">{{ item.title }}</h3>
            </a>
            <p v-if="item.isoDate || item.pubDate">{{ formatDate(item) }}</p>
          </a-space>

          <!-- eslint-disable vue/no-v-html -->
          <div v-if="viewMode === 'rich'" class="rich-text-content">
            <div
              v-html="
                sanitizeContent(item.content || item.contentSnippet || '')
              "
            ></div>
          </div>
          <!-- eslint-enable vue/no-v-html -->
          <a-typography-paragraph
            v-else
            :ellipsis="{
              rows: 3,
              showTooltip: false,
              expandable: true,
            }"
          >
            {{ item.contentSnippet }}
          </a-typography-paragraph>
        </a-card>
      </li>
    </ul>
    <div v-if="feedData.items?.length > 10"
      >Feed 内容项过多, 只显示前10项。
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { computed, ref } from 'vue';
  import dayjs from 'dayjs';
  import DOMPurify from 'dompurify';
  import { useI18n } from 'vue-i18n';
  import type { FeedViewerPreview } from '@/api/feed_viewer';

  const { t } = useI18n();

  interface FeedViewerProp {
    feedData: FeedViewerPreview;
  }

  const props = defineProps<FeedViewerProp>();
  const viewMode = ref('normal');

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

  const sanitizeContent = (content: string) => {
    return DOMPurify.sanitize(content);
  };

  const formatDate = (item: FeedViewerPreview['items'][number]) => {
    const dateValue = item.isoDate || item.pubDate;
    if (!dateValue) return '';

    const parsed = dayjs(dateValue);
    return parsed.isValid() ? parsed.format('YYYY-MM-DD HH:mm:ss') : dateValue;
  };
</script>

<script lang="ts">
  export default {
    name: 'FeedViewContainer',
  };
</script>

<style scoped>
  .rich-text-content :deep(img) {
    max-width: 100%;
    height: auto;
  }
  .rich-text-content :deep(pre) {
    white-space: pre-wrap;
    word-wrap: break-word;
  }
</style>
