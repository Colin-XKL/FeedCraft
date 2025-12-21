<template>
  <div class="container">
    <Breadcrumb :items="['menu.tools', 'menu.dependencyStatus']" />
    <a-card class="general-card" :title="$t('menu.dependencyStatus')">
      <template #extra>
        <a-button type="primary" :loading="loading" @click="handleCheck">
          <template #icon>
            <icon-refresh />
          </template>
          Check Health
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
              item.status
            }}</a-tag>
            <a-tag v-else-if="item.status === 'Unhealthy'" color="red">{{
              item.status
            }}</a-tag>
            <a-tag v-else-if="item.status === 'Configured'" color="blue">{{
              item.status
            }}</a-tag>
            <a-tag v-else color="gray">{{ item.status }}</a-tag>
            <span v-if="item.latency" class="latency">{{ item.latency }}</span>
          </template>
        </a-list-item>
      </a-list>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { ref, onMounted } from 'vue';
  import {
    fetchDependencyStatus,
    checkDependencyStatus,
    DependencyStatus,
  } from '@/api/monitor';

  const data = ref<DependencyStatus[]>([]);
  const loading = ref(false);

  const fetchConfig = async () => {
    try {
      const { data: res } = await fetchDependencyStatus();
      if (res.data) {
        data.value = res.data;
      }
    } catch (err) {
      console.error(err);
    }
  };

  const handleCheck = async () => {
    loading.value = true;
    try {
      const { data: res } = await checkDependencyStatus();
      if (res.data) {
        data.value = res.data;
      }
    } catch (err) {
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
