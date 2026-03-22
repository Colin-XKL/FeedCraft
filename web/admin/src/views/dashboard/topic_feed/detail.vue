<template>
  <div class="py-8 px-16">
    <Breadcrumb :items="['menu.worktable', 'menu.topicFeed']" />
    <x-header
      :title="detail?.topic.title || String(route.params.id || '')"
      :description="t('topic.description')"
    />

    <a-spin :loading="loading" style="width: 100%">
      <a-space direction="vertical" fill size="large">
        <a-card
          v-if="detail"
          class="general-card"
          :title="t('topic.detail.overview')"
        >
          <a-row :gutter="16">
            <a-col :span="8">
              <a-statistic :title="t('topic.detail.currentStatus')">
                <template #value>
                  <a-tag :color="statusColor(detail.health.current_status)">
                    {{ formatStatus(detail.health.current_status) }}
                  </a-tag>
                </template>
              </a-statistic>
            </a-col>
            <a-col :span="8">
              <a-statistic
                :title="t('topic.inputCount')"
                :value="detail.topic.input_uris.length"
              />
            </a-col>
            <a-col :span="8">
              <a-statistic
                :title="t('topic.detail.executionCount')"
                :value="detail.recent_executions.length"
              />
            </a-col>
          </a-row>

          <a-descriptions :column="1" bordered style="margin-top: 16px">
            <a-descriptions-item :label="t('topic.id')">
              {{ detail.topic.id }}
            </a-descriptions-item>
            <a-descriptions-item :label="t('topic.title')">
              {{ detail.topic.title || '-' }}
            </a-descriptions-item>
            <a-descriptions-item :label="t('topic.descriptionLabel')">
              {{ detail.topic.description || '-' }}
            </a-descriptions-item>
            <a-descriptions-item :label="t('topic.publicUrl')">
              <a-space>
                <a-link :href="publicUrl" target="_blank">
                  {{ publicUrl }}
                </a-link>
                <a-button size="mini" @click="copyPublicUrl">
                  {{ t('topic.copyLink') }}
                </a-button>
              </a-space>
            </a-descriptions-item>
            <a-descriptions-item :label="t('topic.detail.lastSuccess')">
              {{ formatTime(detail.health.last_success_at) }}
            </a-descriptions-item>
            <a-descriptions-item :label="t('topic.detail.lastFailure')">
              {{ formatTime(detail.health.last_failure_at) }}
            </a-descriptions-item>
            <a-descriptions-item :label="t('topic.detail.latestError')">
              {{ detail.health.last_error_message || '-' }}
            </a-descriptions-item>
          </a-descriptions>
        </a-card>

        <a-card
          v-if="detail"
          class="general-card"
          :title="t('topic.detail.config')"
        >
          <a-row :gutter="16">
            <a-col :span="12">
              <div class="section-label">{{ t('topic.inputs') }}</div>
              <a-list bordered>
                <a-list-item
                  v-for="(uri, idx) in detail.topic.input_uris"
                  :key="`${uri}-${idx}`"
                >
                  {{ uri }}
                </a-list-item>
              </a-list>
            </a-col>
            <a-col :span="12">
              <div class="section-label">{{ t('topic.aggregatorConfig') }}</div>
              <a-list bordered>
                <a-list-item v-if="detail.topic.aggregator_config.length === 0">
                  {{ t('topic.noAggregator') }}
                </a-list-item>
                <a-list-item
                  v-for="(step, idx) in detail.topic.aggregator_config"
                  :key="`${step.type}-${idx}`"
                >
                  {{ formatAggregatorStep(step) }}
                </a-list-item>
              </a-list>
            </a-col>
          </a-row>
        </a-card>

        <a-card
          v-if="detail"
          class="general-card"
          :title="t('topic.detail.executions')"
        >
          <a-table
            :data="detail.recent_executions"
            :pagination="false"
            row-key="id"
          >
            <template #columns>
              <a-table-column :title="t('observability.time')">
                <template #cell="{ record }">
                  {{ formatTime(record.created_at) }}
                </template>
              </a-table-column>
              <a-table-column :title="t('observability.status')">
                <template #cell="{ record }">
                  <a-tag :color="statusColor(record.status)">
                    {{ formatStatus(record.status) }}
                  </a-tag>
                </template>
              </a-table-column>
              <a-table-column :title="t('observability.trigger')">
                <template #cell="{ record }">
                  {{ formatTrigger(record.trigger) }}
                </template>
              </a-table-column>
              <a-table-column :title="t('observability.errorType')">
                <template #cell="{ record }">
                  {{ formatErrorKind(record.error_kind) }}
                </template>
              </a-table-column>
              <a-table-column
                :title="t('observability.message')"
                data-index="message"
                :ellipsis="true"
              />
              <a-table-column :title="t('observability.actions')">
                <template #cell="{ record }">
                  <a-button
                    type="text"
                    size="small"
                    :disabled="!record.details && !record.details_json"
                    @click="openExecutionDetails(record)"
                  >
                    {{ t('topic.detail.viewExecutionDetails') }}
                  </a-button>
                </template>
              </a-table-column>
            </template>
          </a-table>
          <a-empty
            v-if="detail.recent_executions.length === 0"
            :description="t('topic.detail.emptyExecutions')"
          />
        </a-card>

        <a-card
          v-if="detail"
          class="general-card"
          :title="t('topic.detail.notifications')"
        >
          <a-list bordered>
            <a-list-item v-if="detail.related_notifications.length === 0">
              {{ t('topic.detail.emptyNotifications') }}
            </a-list-item>
            <a-list-item
              v-for="item in detail.related_notifications"
              :key="item.id"
            >
              <a-list-item-meta
                :title="item.title"
                :description="`${formatTime(item.created_at)} · ${formatStatus(
                  item.status_after
                )}`"
              />
              <div>{{ item.content }}</div>
            </a-list-item>
          </a-list>
        </a-card>
      </a-space>
    </a-spin>

    <a-modal
      v-model:visible="detailsModalVisible"
      :title="t('topic.detail.executionDetails')"
      :footer="false"
      width="720px"
    >
      <pre v-if="selectedExecutionDetails" class="details-json">{{
        selectedExecutionDetails
      }}</pre>
      <a-empty v-else :description="t('topic.detail.noExecutionDetails')" />
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
  import { computed, onMounted, ref } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import { useI18n } from 'vue-i18n';
  import { useRoute } from 'vue-router';
  import XHeader from '@/components/header/x-header.vue';
  import {
    formatObservabilityErrorKind,
    formatObservabilityStatus,
    formatObservabilityTrigger,
  } from '@/utils/observability';
  import buildPublicFeedUrl from '@/utils/publicFeedUrl';
  import { AggregatorStep, TopicDetail, getTopicFeedDetail } from '@/api/topic';

  const { t } = useI18n();
  const route = useRoute();
  const loading = ref(false);
  const detail = ref<TopicDetail | null>(null);
  const detailsModalVisible = ref(false);
  const selectedExecutionDetails = ref('');
  const publicUrl = computed(() =>
    detail.value ? buildPublicFeedUrl(detail.value.public_url) : ''
  );

  const formatTime = (value?: string) => {
    if (!value) return '-';
    return new Date(value).toLocaleString();
  };

  const formatStatus = (status?: string) => {
    return formatObservabilityStatus(t, status);
  };

  const formatTrigger = (trigger?: string) => {
    return formatObservabilityTrigger(t, trigger);
  };

  const formatErrorKind = (kind?: string) => {
    return formatObservabilityErrorKind(t, kind);
  };

  const statusColor = (status?: string) => {
    if (status === 'healthy' || status === 'success') return 'green';
    if (status === 'degraded' || status === 'partial_success') return 'orange';
    if (status === 'paused' || status === 'failure') return 'red';
    return 'gray';
  };

  const formatAggregatorStep = (step: AggregatorStep) => {
    if (step.type === 'deduplicate') {
      return `${t('topic.stepType.deduplicate')} · ${t(
        `topic.stepOption.strategy.${step.option?.strategy || 'by_link'}`
      )}`;
    }
    if (step.type === 'sort') {
      return `${t('topic.stepType.sort')} · ${t(
        `topic.stepOption.sort.${step.option?.by || 'date_desc'}`
      )}`;
    }
    if (step.type === 'limit') {
      return `${t('topic.stepType.limit')} · ${step.option?.max || '-'}`;
    }
    return step.type;
  };

  const copyPublicUrl = async () => {
    if (!detail.value) return;
    await navigator.clipboard.writeText(publicUrl.value);
    Message.success(t('topic.copyLink'));
  };

  const buildExecutionDetails = (
    record: TopicDetail['recent_executions'][number]
  ) => {
    if (record.details) {
      return JSON.stringify(record.details, null, 2);
    }
    if (record.details_json) {
      try {
        return JSON.stringify(JSON.parse(record.details_json), null, 2);
      } catch {
        return record.details_json;
      }
    }
    return '';
  };

  const openExecutionDetails = (
    record: TopicDetail['recent_executions'][number]
  ) => {
    selectedExecutionDetails.value = buildExecutionDetails(record);
    detailsModalVisible.value = true;
  };

  const fetchDetail = async () => {
    loading.value = true;
    try {
      const res = await getTopicFeedDetail(String(route.params.id));
      detail.value = res.data;
    } catch (err: any) {
      Message.error(err.message || t('topic.detail.loadFailed'));
    } finally {
      loading.value = false;
    }
  };

  onMounted(() => {
    fetchDetail();
  });
</script>

<script lang="ts">
  export default {
    name: 'TopicFeedDetail',
  };
</script>

<style scoped>
  .section-label {
    margin-bottom: 12px;
    font-weight: 600;
  }

  .details-json {
    max-height: 60vh;
    overflow: auto;
    padding: 12px;
    border-radius: 4px;
    background: var(--color-fill-2);
    white-space: pre-wrap;
    word-break: break-word;
  }
</style>
