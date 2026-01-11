<template>
  <div class="py-8 px-16">
    <x-header
      :title="t('menu.allCraftList')"
      :description="t('allCraftList.description')"
    ></x-header>

    <a-space direction="horizontal" class="mb-6">
      <a-button type="primary" :loading="isLoading" @click="listAllCrafts">
        {{ t('allCraftList.query') }}
      </a-button>
    </a-space>

    <a-table
      :pagination="true"
      :data="allCrafts"
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

  interface CraftItem {
    name: string;
    description: string;
    type: string;
  }

  const isLoading = ref(false);
  const allCrafts = ref<CraftItem[]>([]);

  const columns = [
    { title: t('allCraftList.table.name'), dataIndex: 'name' },
    { title: t('allCraftList.table.type'), dataIndex: 'type' },
    { title: t('allCraftList.table.description'), dataIndex: 'description' },
  ];

  const listAllCrafts = async () => {
    isLoading.value = true;
    try {
      const response = await axios.get('/api/list-all-craft');
      allCrafts.value = response.data.data;
    } catch (error) {
      Message.error(t('allCraftList.message.fetchFailed'));
    } finally {
      isLoading.value = false;
    }
  };

  onBeforeMount(() => {
    listAllCrafts();
  });
</script>

<style scoped></style>
