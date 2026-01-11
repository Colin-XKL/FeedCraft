<template>
  <div class="py-8 px-16">
    <a-card class="general-card" :title="$t('menu.dependencyStatus')">
      <template #extra>
        <a-button type="primary" :loading="loading" @click="handleCheck">
          <template #icon>
            <icon-refresh />
          </template>
          {{ $t('dependencyService.checkHealth') }}
        </a-button>
      </template>
      <a-list>
        <a-list-item v-for="item in data" :key="item.name">
          <a-list-item-meta
            :title="item.name"
            :description="item.error ? item.error : item.details"
          >
            <template #avatar>
              <a-avatar
                v-if="item.status === 'Healthy'"
                :style="{ backgroundColor: '#0fbf60' }"
              >
                <icon-check />
              </a-avatar>
              <a-avatar
                v-else-if="item.status === 'Unhealthy'"
                :style="{ backgroundColor: '#f53f3f' }"
              >
                <icon-close />
              </a-avatar>
              <a-avatar
                v-else-if="item.status === 'Configured'"
                :style="{ backgroundColor: '#165dff' }"
              >
                <icon-settings />
              </a-avatar>
              <a-avatar v-else :style="{ backgroundColor: '#c9cdd4' }">
                <icon-minus />
              </a-avatar>
            </template>
          </a-list-item-meta>
          <template #actions>
            <a-tag v-if="item.status === 'Healthy'" color="green">{{
              $t('dependencyService.status.healthy')
            }}</a-tag>
            <a-tag v-else-if="item.status === 'Unhealthy'" color="red">{{
              $t('dependencyService.status.unhealthy')
            }}</a-tag>
            <a-tag v-else-if="item.status === 'Configured'" color="blue">{{
              $t('dependencyService.status.configured')
            }}</a-tag>
            <a-tag v-else color="gray">{{
              $t('dependencyService.status.notConfigured')
            }}</a-tag>
            <span v-if="item.latency" class="latency">{{ item.latency }}</span>
          </template>
        </a-list-item>
      </a-list>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { ref, onMounted } from 'vue';
  import { useI18n } from 'vue-i18n';
  import { Message } from '@arco-design/web-vue';
  import {
    fetchDependencyStatus,
    checkDependencyStatus,
    DependencyStatus,
  } from '@/api/monitor';

  const { t } = useI18n();
  const data = ref<DependencyStatus[]>([]);
  const loading = ref(false);

  const fetchConfig = async () => {
    loading.value = true;
    try {
      const res = await fetchDependencyStatus();
      if (res.data) {
        data.value = res.data;
      }
    } catch (err) {
      Message.error(t('dependencyService.fetchError'));
      console.error(err);
    } finally {
      loading.value = false;
    }
  };

  const handleCheck = async () => {
    loading.value = true;
    try {
      const res = await checkDependencyStatus();
      if (res.data) {
        data.value = res.data;
        Message.success(t('dependencyService.checkSuccess'));
      }
    } catch (err) {
      Message.error(t('dependencyService.checkError'));
      console.error(err);
    } finally {
      loading.value = false;
    }
  };

  onMounted(() => {
    fetchConfig();
  });
</script>

<style scoped lang="less">
  .container {
    padding: 0 20px 20px 20px;
  }
  .latency {
    margin-left: 10px;
    color: var(--color-text-3);
  }
</style>
