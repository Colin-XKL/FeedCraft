<template>
  <div class="py-8 px-16">
    <x-header
      :title="t('menu.systemHealth')"
      :description="t('health.description')"
    >
    </x-header>

    <a-card class="general-card" :title="t('menu.systemHealth')">
      <template #extra>
        <a-button type="primary" :loading="loading" @click="fetchData">
          <template #icon>
            <icon-refresh />
          </template>
          {{ t('health.analyze') }}
        </a-button>
      </template>

      <a-spin :loading="loading" style="width: 100%">
        <div v-if="treeData.length > 0">
          <div v-if="missingCount > 0" class="mb-4">
            <a-alert type="error" style="margin-bottom: 16px">
              {{ t('health.issuesFound', { count: missingCount }) }}
            </a-alert>
            <a-card
              class="mb-4 border-red-200"
              style="background-color: var(--color-danger-light-1);"
            >
              <template #title>
                <span class="font-medium" style="color: rgb(var(--danger-6))">{{
                  t('health.missingCrafts')
                }}</span>
              </template>
              <div class="flex flex-wrap gap-2">
                <a-tag
                  v-for="node in missingNodes"
                  :key="node.key"
                  color="red"
                  size="large"
                  class="font-medium px-3 py-1"
                >
                  <template #icon>
                    <icon-exclamation-circle-fill />
                  </template>
                  {{ node.name }}
                  <span v-if="node.details" class="ml-2 text-xs opacity-80"
                    >({{ node.details }})</span
                  >
                </a-tag>
              </div>
            </a-card>
          </div>
          <a-alert v-else type="success" style="margin-bottom: 16px">
            {{ t('health.allHealthy') }}
          </a-alert>

          <a-tree
            :data="treeData"
            :show-line="true"
            block-node
            default-expand-all
          >
            <template #title="node">
              <a-space>
                <span
                  :class="{
                    'font-bold': true,
                    'text-red-500': !node.exists || node.type === 'missing',
                  }"
                  >{{ node.name }}</span
                >
                <a-tag
                  v-if="node.type"
                  :color="getTypeColor(node.type)"
                  size="small"
                  >{{ node.type }}</a-tag
                >
                <a-tag
                  v-if="!node.exists || node.type === 'missing'"
                  color="red"
                  size="small"
                  >{{ t('health.missing') }}</a-tag
                >
                <span
                  v-if="node.details"
                  :class="{
                    'text-xs': true,
                    'text-gray-400': node.exists && node.type !== 'missing',
                    'text-red-400': !node.exists || node.type === 'missing',
                  }"
                  >{{ node.details }}</span
                >
              </a-space>
            </template>
          </a-tree>
        </div>
        <a-empty v-else :description="t('health.noData')" />
      </a-spin>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  import { useI18n } from 'vue-i18n';
  import { Message } from '@arco-design/web-vue';
  import { fetchDependencyHealth, DependencyNode } from '@/api/health';
  import XHeader from '@/components/header/x-header.vue';
  import {
    IconRefresh,
    IconExclamationCircleFill,
  } from '@arco-design/web-vue/es/icon';

  const { t } = useI18n();
  const loading = ref(false);
  const treeData = ref<DependencyNode[]>([]);
  const missingCount = ref(0);
  const missingNodes = ref<DependencyNode[]>([]);

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'recipe':
        return 'arcoblue';
      case 'flow':
        return 'purple';
      case 'atom':
        return 'cyan';
      case 'built-in':
        return 'green';
      case 'missing':
        return 'red';
      case 'cycle':
        return 'orange';
      default:
        return 'gray';
    }
  };

  const collectMissingNodes = (
    nodes: DependencyNode[],
    missingList: DependencyNode[]
  ) => {
    nodes.forEach((node) => {
      if (!node.exists || node.type === 'missing') {
        // Prevent duplicates
        if (!missingList.some((n) => n.name === node.name)) {
          missingList.push(node);
        }
      }
      if (node.children) {
        collectMissingNodes(node.children, missingList);
      }
    });
  };

  const countMissing = (nodes: DependencyNode[]) => {
    let count = 0;
    nodes.forEach((node) => {
      if (!node.exists || node.type === 'missing') count += 1;
      if (node.children) {
        count += countMissing(node.children);
      }
    });
    return count;
  };

  const fetchData = async () => {
    loading.value = true;
    try {
      const res = await fetchDependencyHealth();
      const data = res.data ?? [];
      treeData.value = data;
      missingCount.value = countMissing(data);

      const missingList: DependencyNode[] = [];
      collectMissingNodes(data, missingList);
      missingNodes.value = missingList;
    } catch (err: any) {
      treeData.value = [];
      missingCount.value = 0;
      missingNodes.value = [];
      Message.error(err.message || t('health.fetchError'));
    } finally {
      loading.value = false;
    }
  };
</script>

<script lang="ts">
  export default {
    name: 'SystemHealth',
  };
</script>

<style scoped></style>
