<template>
  <div class="craft-list">
    <div class="mb-4">
      <a-input-search
        v-model="searchText"
        :placeholder="t('common.search', 'Search...')"
        allow-clear
      />
    </div>

    <a-tabs default-active-key="sys">
      <a-tab-pane v-for="tab in tabs" :key="tab.key" :title="tab.title">
        <div
          class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4 max-h-[500px] overflow-y-auto p-1"
        >
          <div
            v-for="item in tab.items"
            :key="item.name"
            class="group border border-[var(--color-neutral-3)] rounded p-3 cursor-pointer hover:shadow-md transition-all relative flex flex-col bg-[var(--color-bg-2)]"
            :class="{
              '!border-[rgb(var(--primary-6))] !bg-[var(--color-primary-light-1)]':
                isSelected(item.name),
            }"
            @click="handleSelect(item.name)"
          >
            <div
              class="font-bold text-base mb-1 break-all text-[var(--color-text-1)]"
            >
              {{ item.name }}
            </div>
            <div
              class="text-[var(--color-text-3)] text-xs line-clamp-3"
              :title="item.description"
            >
              {{ item.description || 'No description' }}
            </div>

            <!-- Selection Number for Multiple Mode -->
            <div
              v-if="isSelected(item.name) && multiple"
              class="absolute top-2 right-2 w-6 h-6 rounded-full bg-[rgb(var(--primary-6))] text-white flex items-center justify-center text-xs font-bold"
            >
              {{ getSelectionIndex(item.name) }}
            </div>
            <div
              v-else-if="isSelected(item.name)"
              class="absolute top-2 right-2 text-[rgb(var(--primary-6))]"
            >
              <icon-check-circle-fill size="20" />
            </div>
          </div>
          <div
            v-if="tab.items.length === 0"
            class="col-span-full text-center text-gray-400 py-8"
          >
            {{ t('common.noData', 'No results found') }}
          </div>
        </div>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted, computed } from 'vue';
  import {
    Blueprint,
    listBlueprints,
    listToolTemplates,
  } from '@/api/blueprint';
  import { listTools } from '@/api/tool';
  import { useI18n } from 'vue-i18n';
  import { IconCheckCircleFill } from '@arco-design/web-vue/es/icon';

  const { t } = useI18n();

  const props = defineProps({
    modelValue: {
      type: Array as () => string[],
      default: () => [],
    },
    multiple: {
      type: Boolean,
      default: false,
    },
  });

  const emit = defineEmits(['update:modelValue', 'change']);

  const searchText = ref('');
  const blueprints = ref<Blueprint[]>([]);
  const sysToolList = ref<{ name: string; description?: string }[]>([]);
  const toolList = ref<{ name: string; description?: string }[]>([]);

  // Fetch data
  onMounted(async () => {
    try {
      const [blueprintsResponse, sysToolsResponse, toolsResponse] =
        await Promise.all([listBlueprints(), listToolTemplates(), listTools()]);
      blueprints.value = blueprintsResponse.data || [];
      sysToolList.value = sysToolsResponse.data || [];
      toolList.value = toolsResponse.data || [];
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('Failed to fetch craft options', error);
    }
  });

  const filterAndSort = (items: { name: string; description?: string }[]) => {
    let result = items;
    if (searchText.value) {
      const lower = searchText.value.toLowerCase();
      result = items.filter(
        (item) =>
          item.name.toLowerCase().includes(lower) ||
          (item.description && item.description.toLowerCase().includes(lower)),
      );
    }
    return [...result].sort((a, b) => a.name.localeCompare(b.name));
  };

  const tabs = computed(() => [
    {
      key: 'sys',
      title: t('feedCompare.selectBlueprint.tabs.system'),
      items: filterAndSort(sysToolList.value),
    },
    {
      key: 'user',
      title: t('feedCompare.selectBlueprint.tabs.user'),
      items: filterAndSort(toolList.value),
    },
    {
      key: 'flow',
      title: t('feedCompare.selectBlueprint.tabs.flow'),
      items: filterAndSort(blueprints.value),
    },
  ]);

  const isSelected = (name: string) => {
    return props.modelValue.includes(name);
  };

  const getSelectionIndex = (name: string) => {
    return props.modelValue.indexOf(name) + 1;
  };

  const handleSelect = (name: string) => {
    const newValue = [...props.modelValue];

    if (props.multiple) {
      const idx = newValue.indexOf(name);
      if (idx > -1) {
        newValue.splice(idx, 1);
      } else {
        newValue.push(name);
      }
    } else {
      // For single selection, we replace the entire array with the new value
      // The parent/modal handles whether this causes an immediate close
      newValue.length = 0;
      newValue.push(name);
    }

    emit('update:modelValue', newValue);
    emit('change', newValue);
  };
</script>

<style scoped>
  /* Scoped styles */
</style>
