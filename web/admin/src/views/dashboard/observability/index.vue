<template>
  <div class="py-8 px-16">
    <x-header
      :title="t('menu.observability')"
      :description="t('observability.description')"
    />

    <a-tabs default-active-key="resources">
      <a-tab-pane key="resources" :title="t('observability.resources')">
        <a-card class="general-card">
          <template #extra>
            <a-space>
              <a-select
                v-model="resourceType"
                style="width: 160px"
                :placeholder="t('observability.resourceType')"
              >
                <a-option value="">{{ t('observability.all') }}</a-option>
                <a-option value="recipe">{{
                  t('observability.resourceType.recipe')
                }}</a-option>
                <a-option value="topic">{{
                  t('observability.resourceType.topic')
                }}</a-option>
              </a-select>
              <a-button
                type="primary"
                :loading="resourceLoading"
                @click="loadResources"
              >
                {{ t('observability.refresh') }}
              </a-button>
            </a-space>
          </template>

          <a-table
            :data="resources"
            :loading="resourceLoading"
            :pagination="false"
            row-key="resource_id"
          >
            <template #columns>
              <a-table-column
                :title="t('observability.type')"
                data-index="resource_type"
              >
                <template #cell="{ record }">
                  {{ formatResourceType(record.resource_type) }}
                </template>
              </a-table-column>
              <a-table-column
                :title="t('observability.resource')"
                data-index="resource_name"
              />
              <a-table-column
                :title="t('observability.status')"
                data-index="current_status"
              >
                <template #cell="{ record }">
                  <a-tag :color="statusColor(record.current_status)">
                    {{ formatStatus(record.current_status) }}
                  </a-tag>
                </template>
              </a-table-column>
              <a-table-column
                :title="t('observability.failures')"
                data-index="consecutive_failures"
              />
              <a-table-column
                :title="t('observability.lastSuccess')"
                data-index="last_success_at"
              >
                <template #cell="{ record }">
                  {{ formatTime(record.last_success_at) }}
                </template>
              </a-table-column>
              <a-table-column
                :title="t('observability.lastFailure')"
                data-index="last_failure_at"
              >
                <template #cell="{ record }">
                  {{ formatTime(record.last_failure_at) }}
                </template>
              </a-table-column>
              <a-table-column :title="t('observability.actions')">
                <template #cell="{ record }">
                  <a-space>
                    <a-link
                      :href="
                        buildFeedUrl(record.resource_type, record.resource_id)
                      "
                      target="_blank"
                    >
                      {{ t('observability.link') }}
                    </a-link>
                    <!-- TopicFeed 功能当前仍在开发完善中，先隐藏详情入口；待功能 ready 后再重新开放。 -->
                    <a-button
                      v-if="false && record.resource_type === 'topic'"
                      type="text"
                      size="small"
                      @click="goToTopicDetail(record.resource_id)"
                    >
                      {{ t('topic.viewDetails') }}
                    </a-button>
                    <a-button
                      v-if="record.current_status === 'paused'"
                      type="outline"
                      size="small"
                      @click="handleResume(record)"
                    >
                      {{ t('observability.resume') }}
                    </a-button>
                  </a-space>
                </template>
              </a-table-column>
            </template>
          </a-table>
        </a-card>
      </a-tab-pane>

      <a-tab-pane key="executions" :title="t('observability.executions')">
        <a-card class="general-card">
          <template #extra>
            <a-space>
              <a-select
                v-model="logStatus"
                style="width: 180px"
                :placeholder="t('observability.status')"
              >
                <a-option value="">{{ t('observability.all') }}</a-option>
                <a-option value="success">{{
                  formatStatus('success')
                }}</a-option>
                <a-option value="partial_success">{{
                  formatStatus('partial_success')
                }}</a-option>
                <a-option value="failure">{{
                  formatStatus('failure')
                }}</a-option>
                <a-option value="paused_skip">{{
                  formatStatus('paused_skip')
                }}</a-option>
              </a-select>
              <a-button type="primary" :loading="logLoading" @click="loadLogs">
                {{ t('observability.refresh') }}
              </a-button>
            </a-space>
          </template>

          <a-table
            :data="logs"
            :loading="logLoading"
            :pagination="false"
            row-key="id"
          >
            <template #columns>
              <a-table-column
                :title="t('observability.time')"
                data-index="created_at"
              >
                <template #cell="{ record }">
                  {{ formatTime(record.created_at) }}
                </template>
              </a-table-column>
              <a-table-column
                :title="t('observability.type')"
                data-index="resource_type"
              >
                <template #cell="{ record }">
                  {{ formatResourceType(record.resource_type) }}
                </template>
              </a-table-column>
              <a-table-column
                :title="t('observability.id')"
                data-index="resource_id"
              />
              <a-table-column
                :title="t('observability.trigger')"
                data-index="trigger"
              >
                <template #cell="{ record }">
                  {{ formatTrigger(record.trigger) }}
                </template>
              </a-table-column>
              <a-table-column
                :title="t('observability.status')"
                data-index="status"
              >
                <template #cell="{ record }">
                  <a-tag :color="statusColor(record.status)">{{
                    formatStatus(record.status)
                  }}</a-tag>
                </template>
              </a-table-column>
              <a-table-column
                :title="t('observability.errorType')"
                data-index="error_kind"
              >
                <template #cell="{ record }">
                  {{ formatErrorKind(record.error_kind) }}
                </template>
              </a-table-column>
              <a-table-column
                :title="t('observability.message')"
                data-index="message"
                :ellipsis="true"
              />
            </template>
          </a-table>
        </a-card>
      </a-tab-pane>

      <a-tab-pane key="notifications" :title="t('observability.notifications')">
        <a-card class="general-card">
          <template #extra>
            <a-space>
              <a-link :href="systemNotificationUrl" target="_blank">
                {{ systemNotificationUrl }}
              </a-link>
              <a-button
                type="primary"
                :loading="notificationLoading"
                @click="loadNotifications"
              >
                {{ t('observability.refresh') }}
              </a-button>
            </a-space>
          </template>

          <a-list :bordered="false" :loading="notificationLoading">
            <a-list-item v-for="item in notifications" :key="item.id">
              <a-list-item-meta
                :title="item.title"
                :description="`${formatTime(
                  item.created_at
                )} · ${formatResourceType(item.resource_type)} / ${
                  item.resource_id
                }`"
              />
              <template #actions>
                <a-tag :color="statusColor(item.status_after)">
                  {{ formatStatus(item.status_after) }}
                </a-tag>
              </template>
              <div class="notification-content">{{ item.content }}</div>
            </a-list-item>
          </a-list>
        </a-card>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script lang="ts" setup>
  import { onMounted, ref, watch } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import { useI18n } from 'vue-i18n';
  import { useRouter } from 'vue-router';
  import XHeader from '@/components/header/x-header.vue';
  import {
    formatObservabilityErrorKind,
    formatObservabilityResourceType,
    formatObservabilityStatus,
    formatObservabilityTrigger,
  } from '@/utils/observability';
  import buildPublicFeedUrl from '@/utils/publicFeedUrl';
  import {
    ExecutionLog,
    ObservableResource,
    SystemNotification,
    fetchExecutionLogs,
    fetchObservableResources,
    fetchSystemNotifications,
    resumeObservableResource,
  } from '@/api/observability';

  const { t } = useI18n();
  const router = useRouter();
  const resourceLoading = ref(false);
  const logLoading = ref(false);
  const notificationLoading = ref(false);

  const resourceType = ref('');
  const logStatus = ref('');

  const resources = ref<ObservableResource[]>([]);
  const logs = ref<ExecutionLog[]>([]);
  const notifications = ref<SystemNotification[]>([]);
  const systemNotificationUrl = buildPublicFeedUrl('/system/notifications/rss');

  const formatTime = (value?: string) => {
    if (!value) return '-';
    return new Date(value).toLocaleString();
  };

  const formatStatus = (status?: string) => {
    return formatObservabilityStatus(t, status);
  };

  const formatResourceType = (type?: string) => {
    return formatObservabilityResourceType(t, type);
  };

  const formatTrigger = (trigger?: string) => {
    return formatObservabilityTrigger(t, trigger);
  };

  const formatErrorKind = (kind?: string) => {
    return formatObservabilityErrorKind(t, kind);
  };

  const statusColor = (status?: string) => {
    switch (status) {
      case 'healthy':
      case 'success':
        return 'green';
      case 'degraded':
      case 'partial_success':
        return 'orange';
      case 'paused':
      case 'failure':
      case 'paused_skip':
        return 'red';
      default:
        return 'gray';
    }
  };

  const loadResources = async () => {
    resourceLoading.value = true;
    try {
      const res = await fetchObservableResources({
        resource_type: resourceType.value || undefined,
      });
      resources.value = res.data ?? [];
    } catch (err: any) {
      Message.error(err.message || t('observability.loadFailed'));
    } finally {
      resourceLoading.value = false;
    }
  };

  const loadLogs = async () => {
    logLoading.value = true;
    try {
      const res = await fetchExecutionLogs({
        status: logStatus.value || undefined,
      });
      logs.value = res.data ?? [];
    } catch (err: any) {
      Message.error(err.message || t('observability.loadFailed'));
    } finally {
      logLoading.value = false;
    }
  };

  const loadNotifications = async () => {
    notificationLoading.value = true;
    try {
      const res = await fetchSystemNotifications();
      notifications.value = res.data ?? [];
    } catch (err: any) {
      Message.error(err.message || t('observability.loadFailed'));
    } finally {
      notificationLoading.value = false;
    }
  };

  const handleResume = async (record: ObservableResource) => {
    try {
      await resumeObservableResource(record.resource_type, record.resource_id);
      Message.success(t('observability.resumeSuccess'));
      loadResources();
      loadNotifications();
    } catch (err: any) {
      Message.error(err.message || t('observability.resumeFailed'));
    }
  };

  // TopicFeed 功能当前仍在开发完善中，先隐藏详情跳转；待功能 ready 后再重新开放。
  const goToTopicDetail = (id: string) => {
    router.push({ name: 'TopicFeedDetail', params: { id } });
  };

  const buildFeedUrl = (resourceType: string, resourceId: string) => {
    if (resourceType === 'topic') {
      return buildPublicFeedUrl(`/topic/${resourceId}`);
    }
    return buildPublicFeedUrl(`/recipe/${resourceId}`);
  };

  watch(resourceType, () => loadResources());
  watch(logStatus, () => loadLogs());

  onMounted(() => {
    loadResources();
    loadLogs();
    loadNotifications();
  });
</script>

<script lang="ts">
  export default {
    name: 'Observability',
  };
</script>

<style scoped lang="less">
  .notification-content {
    color: var(--color-text-2);
    white-space: pre-wrap;
  }
</style>
