<template>
  <div class="container">
    <Breadcrumb :items="['menu.tools', 'menu.tools.dependencies']" />
    <a-card class="general-card" :title="$t('menu.tools.dependencies')">
      <a-table :data="data" :loading="loading" :pagination="false">
        <template #columns>
          <a-table-column :title="$t('monitor.dependencies.name')" data-index="name" />
          <a-table-column :title="$t('monitor.dependencies.configured')" data-index="configured">
            <template #cell="{ record }">
              <a-tag :color="record.configured ? 'green' : 'gray'">
                {{ record.configured ? $t('monitor.dependencies.yes') : $t('monitor.dependencies.no') }}
              </a-tag>
            </template>
          </a-table-column>
          <a-table-column :title="$t('monitor.dependencies.healthy')" data-index="healthy">
            <template #cell="{ record }">
              <a-tag v-if="record.configured" :color="record.healthy ? 'green' : 'red'">
                {{ record.healthy ? $t('monitor.dependencies.healthy') : $t('monitor.dependencies.unhealthy') }}
              </a-tag>
              <span v-else>-</span>
            </template>
          </a-table-column>
          <a-table-column :title="$t('monitor.dependencies.message')" data-index="message" />
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { ref, onMounted } from 'vue';
  import { getDependencyStatus, DependencyStatus } from '@/api/monitor';
  import { useI18n } from 'vue-i18n';

  const { t } = useI18n();
  const loading = ref(false);
  const data = ref<DependencyStatus[]>([]);

  const fetchData = async () => {
    loading.value = true;
    try {
      const res = await getDependencyStatus();
      data.value = res.data;
    } catch (err) {
      // console.error(err);
    } finally {
      loading.value = false;
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
</style>
