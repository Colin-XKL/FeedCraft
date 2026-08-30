<template>
  <div class="craft-manage-page">
    <x-header :title="title" :description="description" />

    <a-card class="craft-manage-page__card" :bordered="false">
      <div v-if="$slots.toolbar" class="craft-manage-page__toolbar">
        <slot name="toolbar" />
      </div>
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

  .craft-manage-page__toolbar {
    display: flex;
    justify-content: flex-start;
    margin-bottom: 16px;
  }

  :deep(.arco-card-body) {
    padding: 20px;
  }

  :deep(.arco-table-th) {
    background: var(--color-fill-1);
  }

  /* Override Arco's default word-break: break-all so English identifiers
     wrap at word/hyphen boundaries instead of mid-token. */
  :deep(.arco-table-td) {
    word-break: normal;
    overflow-wrap: break-word;
  }

  @media (max-width: 768px) {
    .craft-manage-page {
      padding: 20px 16px;
    }
  }
</style>
