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
    <ul>
      <li v-for="item in feedData.items?.slice(0, 10)" :key="item.guid">
        <a-card class="my-2">
          <a-space>
            <h3 class="font-bold">{{ item.title }}</h3>
            <p>{{ dayjs(item.isoDate).format('YYYY-MM-DD hh:mm:ss') }}</p>
          </a-space>
          <a-typography-paragraph
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
  import { computed } from 'vue';
  import dayjs from 'dayjs';

  interface FeedViewerProp {
    feedData: Parser.Output<any>;
  }

  const props = defineProps<FeedViewerProp>();
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
</script>

<script lang="ts">
  export default {
    name: 'FeedViewContainer',
  };
</script>
