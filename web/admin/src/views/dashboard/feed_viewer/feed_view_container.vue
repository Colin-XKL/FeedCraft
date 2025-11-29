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
      <li v-for="item in feedData.items?.slice(0, 10)" :key="item.guid">
        <a-card class="my-2">
          <a-space>
            <h3 class="font-bold">{{ item.title }}</h3>
            <p>{{ dayjs(item.isoDate).format('YYYY-MM-DD hh:mm:ss') }}</p>
          </a-space>

          <div v-if="viewMode === 'rich'" class="rich-text-content">
            <!-- eslint-disable-next-line vue/no-v-html -->
            <div
              v-html="
                sanitizeContent(item.content || item.contentSnippet || '')
              "
            ></div>
          </div>
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
  import Parser from 'rss-parser';
  import { computed, ref } from 'vue';
  import dayjs from 'dayjs';
  import DOMPurify from 'dompurify';
  import { useI18n } from 'vue-i18n';

  const { t } = useI18n();

  interface FeedViewerProp {
    feedData: Parser.Output<any>;
  }

  const props = defineProps<FeedViewerProp>();
  const viewMode = ref('normal');

  const feedMetaList = computed(() => {
    return Object.keys(props.feedData)
      .filter((key) => key !== 'items')
      .map((key) => {
        const feed = props.feedData as any;

        const item = feed[key];
        return {
          label: key,
          value: item,
        };
      });
  });

  const sanitizeContent = (content: string) => {
    return DOMPurify.sanitize(content);
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
