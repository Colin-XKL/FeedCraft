<template>
  <div class="craft-manage-page">
    <x-header :title="title" :description="description" />

    <a-card class="craft-manage-page__card" :bordered="false">
      <template v-if="$slots.toolbar" #extra>
        <slot name="toolbar" />
      </template>
      <slot />
    </a-card>
  </div>
</template>

<script setup lang="ts">
  import XHeader from '@/components/header/x-header.vue';

  defineProps<{
    title: string;
    description?: string;
  }>();
</script>

<script lang="ts">
  export default {
    name: 'CraftManagePage',
  };
</script>

<style scoped lang="less">
  .craft-manage-page {
    padding: 32px 64px;
  }

  .craft-manage-page__card {
    border: 1px solid var(--color-border-2);
    border-radius: 14px;
    box-shadow: 0 8px 24px rgba(15, 23, 42, 0.04);
  }

  :deep(.arco-card-header) {
    min-height: 64px;
  }

  :deep(.arco-card-body) {
    padding: 20px;
  }

  :deep(.arco-table-th) {
    background: var(--color-fill-1);
  }

  @media (max-width: 768px) {
    .craft-manage-page {
      padding: 20px 16px;
    }
  }
</style>
