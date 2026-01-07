<template>
  <div class="py-8 px-16">
    <x-header
      :title="t('menu.allProcessorList')"
      :description="t('allProcessorList.description')"
    ></x-header>

    <a-space direction="horizontal" class="mb-6">
      <a-button type="primary" :loading="isLoading" @click="listAllProcessors">
        {{ t('allProcessorList.query') }}
      </a-button>
    </a-space>

    <a-table
      :pagination="true"
      :data="allProcessors"
      :columns="columns"
      :loading="isLoading"
    ></a-table>
  </div>
</template>

<script setup lang="ts">
  import XHeader from '@/components/header/x-header.vue';
  import { onBeforeMount, ref } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import axios from 'axios';
  import { useI18n } from 'vue-i18n';

  const { t } = useI18n();

  interface ProcessorItem {
    name: string;
    description: string;
    type: string;
  }

  const isLoading = ref(false);
  const allProcessors = ref<ProcessorItem[]>([]);

  const columns = [
    { title: t('allProcessorList.table.name'), dataIndex: 'name' },
    { title: t('allProcessorList.table.type'), dataIndex: 'type' },
    { title: t('allProcessorList.table.description'), dataIndex: 'description' },
  ];

  const listAllProcessors = async () => {
    isLoading.value = true;
    try {
      const response = await axios.get('/api/list-all-craft');
      allProcessors.value = response.data.data;
    } catch (error) {
      Message.error(t('allProcessorList.message.fetchFailed'));
    } finally {
      isLoading.value = false;
    }
  };

  onBeforeMount(() => {
    listAllProcessors();
  });
</script>

<style scoped></style>
