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
      row-key="row_key"
      :data="allCraftRows"
      :columns="columns"
      :loading="isLoading"
      :bordered="false"
      :pagination="{ pageSize: 10, showTotal: true }"
    >
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

  const { t } = useI18n();

  const isLoading = ref(false);
  const allCrafts = ref<CraftItem[]>([]);
  const allCraftRows = computed(() =>
    allCrafts.value.map((item) => ({
      ...item,
      row_key: `${item.type}:${item.name}`,
    }))
  );

  const columns = [
    { title: t('allCraftList.table.name'), dataIndex: 'name' },
    { title: t('allCraftList.table.type'), slotName: 'type', width: 160 },
    {
      title: t('allCraftList.table.templateOnly'),
      slotName: 'templateOnly',
      width: 140,
    },
    { title: t('allCraftList.table.description'), dataIndex: 'description' },
  ];

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
    } catch (error) {
      Message.error(t('allCraftList.message.fetchFailed'));
    } finally {
      isLoading.value = false;
    }
  };

  onBeforeMount(() => {
    fetchAllCrafts();
  });
</script>

<style scoped></style>
