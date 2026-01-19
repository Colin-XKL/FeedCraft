<template>
  <div class="py-8 px-16">
    <x-header
      :title="t('menu.systemHealth')"
      :description="t('health.description')"
    >
    </x-header>

    <a-card class="general-card" :title="t('menu.systemHealth')">
      <a-row style="margin-bottom: 16px">
        <a-col :span="24">
          <a-space>
            <a-button type="primary" :loading="loading" @click="fetchData">
              {{ t('health.analyze') }}
            </a-button>
          </a-space>
        </a-col>
      </a-row>

      <a-spin :loading="loading" style="width: 100%">
        <div v-if="treeData.length > 0">
           <a-alert v-if="missingCount > 0" type="error" style="margin-bottom: 16px">
             {{ t('health.issuesFound', { count: missingCount }) }}
           </a-alert>
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
                  <span style="font-weight: bold;">{{ node.name }}</span>
                  <a-tag v-if="node.type" :color="getTypeColor(node.type)" size="small">{{ node.type }}</a-tag>
                  <a-tag v-if="!node.exists" color="red" size="small">{{ t('health.missing') }}</a-tag>
                  <span v-if="node.details" style="color: #86909c; font-size: 12px">{{ node.details }}</span>
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
import { fetchDependencyHealth, DependencyNode } from '@/api/health';
import XHeader from '@/components/header/x-header.vue';

const { t } = useI18n();
const loading = ref(false);
const treeData = ref<DependencyNode[]>([]);
const missingCount = ref(0);

const getTypeColor = (type: string) => {
  switch (type) {
    case 'recipe': return 'arcoblue';
    case 'flow': return 'purple';
    case 'atom': return 'cyan';
    case 'built-in': return 'green';
    case 'missing': return 'red';
    case 'cycle': return 'orange';
    default: return 'gray';
  }
};

const countMissing = (nodes: DependencyNode[]) => {
  let count = 0;
  nodes.forEach(node => {
    if (!node.exists || node.type === 'missing') count++;
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
    treeData.value = res.data;
    missingCount.value = countMissing(res.data);
  } catch (err) {
    console.error(err);
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

<style scoped>
</style>
