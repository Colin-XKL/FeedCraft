<template>
  <div class="py-8 px-16">
    <x-header
      title="所有 Craft 列表"
      description="展示所有系统内置 Craft Atom、用户自定义 Craft Atom 和用户自定义 Craft Flow"
    ></x-header>

    <a-space direction="horizontal" class="mb-6">
      <a-button type="primary" :loading="isLoading" @click="listAllCrafts">
        查询
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

  interface CraftItem {
    name: string;
    description: string;
    type: string;
  }

  const isLoading = ref(false);
  const allCrafts = ref<CraftItem[]>([]);

  const columns = [
    { title: '名称', dataIndex: 'name' },
    { title: '类型', dataIndex: 'type' },
    { title: '描述', dataIndex: 'description' },
  ];

  const listAllCrafts = async () => {
    isLoading.value = true;
    try {
      const response = await axios.get('/api/list-all-craft');
      allCrafts.value = response.data;
    } catch (error) {
      Message.error('Failed to fetch all crafts');
    } finally {
      isLoading.value = false;
    }
  };

  onBeforeMount(() => {
    listAllCrafts();
  });
</script>

<style scoped></style>
