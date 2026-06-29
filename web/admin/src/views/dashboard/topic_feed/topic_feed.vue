<template>
  <div class="py-8 px-16">
    <Breadcrumb :items="['menu.worktable', 'menu.topicFeed']" />
    <x-header
      :title="$t('menu.topicFeed')"
      :description="t('topic.description')"
    />

    <a-card class="general-card" :title="$t('menu.topicFeed')">
      <template #extra>
        <a-space>
          <a-button :loading="loading" @click="fetchTopics">
            {{ t('topic.refresh') }}
          </a-button>
          <a-button type="primary" @click="handleAdd">
            <template #icon>
              <icon-plus />
            </template>
            {{ t('topic.create') }}
          </a-button>
        </a-space>
      </template>

      <a-table
        :data="topics"
        :loading="loading"
        :pagination="false"
        row-key="id"
      >
        <template #columns>
          <a-table-column :title="t('topic.id')" data-index="id" />
          <a-table-column :title="t('topic.title')" data-index="title">
            <template #cell="{ record }">
              {{ record.title || record.id }}
            </template>
          </a-table-column>
          <a-table-column
            :title="t('topic.descriptionLabel')"
            data-index="description"
            :ellipsis="true"
          />
          <a-table-column :title="t('topic.inputCount')">
            <template #cell="{ record }">
              <a-tag color="arcoblue">{{ inputCount(record) }}</a-tag>
            </template>
          </a-table-column>
          <a-table-column :title="t('topic.aggregator')">
            <template #cell="{ record }">
              <span>{{
                formatAggregatorSummary(record.aggregator_config, t)
              }}</span>
            </template>
          </a-table-column>
          <a-table-column :title="t('observability.actions')">
            <template #cell="{ record }">
              <a-space wrap>
                <a-button
                  type="text"
                  size="small"
                  @click="goToDetail(record.id)"
                >
                  {{ t('topic.viewDetails') }}
                </a-button>
                <a-button
                  type="text"
                  size="small"
                  @click="previewTopic(record.id)"
                >
                  {{ t('topic.preview') }}
                </a-button>
                <a-button
                  type="text"
                  size="small"
                  @click="handleEdit(record.id)"
                >
                  {{ t('topic.editAction') }}
                </a-button>
                <a-link :href="buildTopicFeedUrl(record.id)" target="_blank">
                  {{ t('topic.viewFeed') }}
                </a-link>
                <a-popconfirm
                  :content="t('topic.deleteConfirm')"
                  @ok="handleDelete(record.id)"
                >
                  <a-button type="text" status="danger" size="small">
                    {{ t('topic.deleteAction') }}
                  </a-button>
                </a-popconfirm>
              </a-space>
            </template>
          </a-table-column>
        </template>
      </a-table>

      <a-empty
        v-if="!loading && topics.length === 0"
        :description="t('topic.noTopics')"
      />
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { onMounted, ref } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import { useI18n } from 'vue-i18n';
  import { useRouter } from 'vue-router';
  import XHeader from '@/components/header/x-header.vue';
  import buildPublicFeedUrl from '@/utils/publicFeedUrl';
  import { TopicFeed, deleteTopicFeed, listTopicFeeds } from '@/api/topic';
  import { formatAggregatorSummary } from '@/views/dashboard/topic_feed/topicInputUtils';

  const { t } = useI18n();
  const router = useRouter();
  const topics = ref<TopicFeed[]>([]);
  const loading = ref(false);

  const inputCount = (record: TopicFeed) => {
    return record.inputs?.length || 0;
  };

  const fetchTopics = async () => {
    loading.value = true;
    try {
      const res = await listTopicFeeds();
      topics.value = res.data ?? [];
    } catch (err: any) {
      Message.error(err.message || t('topic.fetchFailed'));
    } finally {
      loading.value = false;
    }
  };

  const buildTopicFeedUrl = (id: string) => buildPublicFeedUrl(`/topic/${id}`);

  const previewTopic = (id: string) => {
    router.push({
      name: 'FeedViewer',
      query: { target: 'topic', id },
    });
  };

  const handleAdd = () => {
    router.push({ name: 'TopicFeedCreate' });
  };

  const handleEdit = (id: string) => {
    router.push({ name: 'TopicFeedEdit', params: { id } });
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteTopicFeed(id);
      Message.success(t('topic.deleteSuccess'));
      await fetchTopics();
    } catch (err: any) {
      Message.error(err.message || t('topic.deleteFailed'));
    }
  };

  const goToDetail = (id: string) => {
    router.push({ name: 'TopicFeedDetail', params: { id } });
  };

  onMounted(() => {
    fetchTopics();
  });
</script>

<script lang="ts">
  export default {
    name: 'TopicFeed',
  };
</script>
