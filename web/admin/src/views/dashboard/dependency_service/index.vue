<template>
  <div class="container">
    <Breadcrumb :items="['menu.tools', 'menu.dependencyStatus']" />
    <a-card class="general-card" :title="$t('menu.dependencyStatus')">
      <a-list>
        <a-list-item v-for="item in data" :key="item.name">
          <a-list-item-meta
            :title="item.name"
            :description="item.error ? item.error : item.details"
          >
            <template #avatar>
              <a-avatar v-if="item.status === 'Healthy'" :style="{ backgroundColor: '#0fbf60' }">
                <icon-check />
              </a-avatar>
              <a-avatar v-else-if="item.status === 'Unhealthy'" :style="{ backgroundColor: '#f53f3f' }">
                <icon-close />
              </a-avatar>
              <a-avatar v-else :style="{ backgroundColor: '#c9cdd4' }">
                <icon-minus />
              </a-avatar>
            </template>
          </a-list-item-meta>
          <template #actions>
            <a-tag v-if="item.status === 'Healthy'" color="green">{{ item.status }}</a-tag>
            <a-tag v-else-if="item.status === 'Unhealthy'" color="red">{{ item.status }}</a-tag>
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
  import { fetchDependencyStatus, DependencyStatus } from '@/api/monitor';

  const data = ref<DependencyStatus[]>([]);

  const fetchData = async () => {
    try {
      const { data: res } = await fetchDependencyStatus();
      // res is APIResponse { code, msg, data }
      // We need res.data which is DependencyStatus[]
      if (res.data) {
        data.value = res.data;
      }
    } catch (err) {
      console.error(err);
    }
  };

  onMounted(() => {
    fetchData();
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
