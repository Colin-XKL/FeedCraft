<template>
  <div>
    <h2>{{ props.feedData.title }}</h2>
    <a-descriptions
      style="margin-top: 20px"
      :data="feedMetaList"
      title="Feed Info"
      :column="1"
    />
    <ul>
      <li v-for="item in feedData.items?.slice(0, 5)" :key="item.guid">
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
