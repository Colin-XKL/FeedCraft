<template>
  <CraftManagePage
    :title="t('menu.allCraftList')"
    :description="t('allCraftList.description')"
  >
    <template #toolbar>
      <a-button :loading="isLoading" @click="fetchAllCrafts">
        <template #icon>
          <icon-refresh />
        </template>
        {{ t('allCraftList.query') }}
      </a-button>
    </template>

    <a-table
      class="all-craft-list-table"
      row-key="row_key"
      :data="allCraftRows"
      :columns="columns"
      :loading="isLoading"
      :bordered="false"
      :pagination="{ pageSize: 10, showTotal: true }"
      :scroll="{ x: ALL_CRAFT_LIST_SCROLL_X }"
    >
      <template #name="{ record }">
        <span class="all-craft-name" :title="record.name">{{
          record.name
        }}</span>
      </template>
      <template #type="{ record }">
        <a-tag :color="getCraftTypeColor(record.type)">
          {{ record.type }}
        </a-tag>
      </template>
      <template #templateOnly="{ record }">
        <a-tag :color="record.template_only ? 'orange' : 'green'">
          {{
            record.template_only
              ? t('allCraftList.table.templateOnlyYes')
              : t('allCraftList.table.templateOnlyNo')
          }}
        </a-tag>
      </template>
    </a-table>
  </CraftManagePage>
</template>

<script setup lang="ts">
  import CraftManagePage from '@/components/craft/CraftManagePage.vue';
  import { computed, onBeforeMount, ref } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import { useI18n } from 'vue-i18n';
  import { CraftItem, listAllCrafts } from '@/api/craft_flow';
  import {
    ALL_CRAFT_LIST_SCROLL_X,
    getAllCraftListColumns,
  } from './all_craft_list.columns';

  const { t } = useI18n();

  const isLoading = ref(false);
  const allCrafts = ref<CraftItem[]>([]);
  const allCraftRows = computed(() =>
    allCrafts.value.map((item) => ({
      ...item,
      row_key: `${item.type}:${item.name}`,
    }))
  );

  const columns = getAllCraftListColumns(t);

  const getCraftTypeColor = (type: string) => {
    if (type?.toLowerCase().includes('flow')) return 'purple';
    if (type?.toLowerCase().includes('atom')) return 'arcoblue';
    return 'gray';
  };

  const fetchAllCrafts = async () => {
    isLoading.value = true;
    try {
      const response = await listAllCrafts();
      allCrafts.value = response.data;
    } catch {
      Message.error(t('allCraftList.message.fetchFailed'));
    } finally {
      isLoading.value = false;
    }
  };

  onBeforeMount(() => {
    fetchAllCrafts();
  });
</script>

<style scoped>
  /* Arco Table defaults to word-break: break-all, which splits English
     identifiers mid-token (ignore-advertorial → advert / orial). */
  .all-craft-list-table :deep(.arco-table-td) {
    word-break: normal;
  }

  .all-craft-list-table :deep(.all-craft-id-cell),
  .all-craft-list-table :deep(.all-craft-id-cell .arco-table-td-content),
  .all-craft-name {
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
    word-break: keep-all;
    overflow-wrap: normal;
  }

  .all-craft-name {
    display: inline-block;
    max-width: 100%;
    vertical-align: bottom;
  }

  .all-craft-list-table :deep(.all-craft-desc-cell) {
    overflow-wrap: break-word;
  }
</style>
