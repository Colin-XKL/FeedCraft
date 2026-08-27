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
      class="all-craft-list-table"
      :data="allCrafts"
      :columns="columns"
      :loading="isLoading"
      :scroll="{ x: ALL_CRAFT_LIST_SCROLL_X }"
    ></a-table>
  </div>
</template>

<script setup lang="ts">
  import XHeader from '@/components/header/x-header.vue';
  import { onBeforeMount, ref } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import axios from 'axios';
  import { useI18n } from 'vue-i18n';
  import {
    ALL_CRAFT_LIST_SCROLL_X,
    getAllCraftListColumns,
  } from './all_craft_list.columns';

  const { t } = useI18n();

  interface CraftItem {
    name: string;
    description: string;
    type: string;
    template_only?: boolean;
  }

  const isLoading = ref(false);
  const allCrafts = ref<CraftItem[]>([]);
  const columns = getAllCraftListColumns(t);

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

<style scoped>
  /* Arco Table defaults to word-break: break-all, which splits English
     identifiers mid-token when a column is squeezed. Prefer wrapping at
     word/hyphen boundaries; identifier columns also use nowrap+ellipsis. */
  .all-craft-list-table :deep(.arco-table-td) {
    word-break: normal;
    overflow-wrap: break-word;
  }
</style>
