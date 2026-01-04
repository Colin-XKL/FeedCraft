<template>
  <div class="craft-selector">
    <!-- Trigger Input -->
    <div
      v-if="!$slots.default"
      class="trigger-input w-full border border-[var(--color-neutral-3)] rounded bg-[var(--color-bg-2)] px-3 py-1.5 min-h-[32px] cursor-pointer flex items-center justify-between hover:border-[rgb(var(--primary-6))] transition-colors focus:outline-none focus:ring-2 focus:ring-[rgba(var(--primary-6),0.2)] focus:border-[rgb(var(--primary-6))]"
      tabindex="0"
      role="button"
      aria-haspopup="dialog"
      :aria-label="placeholder || t('feedCompare.selectCraftFlow.placeholder')"
      @click="openModal"
      @keydown.enter.prevent="openModal"
      @keydown.space.prevent="openModal"
    >
      <div class="flex flex-wrap gap-1 flex-1 overflow-hidden">
        <span v-if="!hasSelection" class="text-[var(--color-text-3)]">{{
          placeholder || t('feedCompare.selectCraftFlow.placeholder')
        }}</span>
        <template v-else>
          <!-- Multiple mode: tags -->
          <template v-if="mode === 'multiple'">
            <a-tag
              v-for="val in localValue"
              :key="val"
              closable
              class="mr-1 my-0.5"
              @close.stop="removeItem(val)"
            >
              {{ val }}
            </a-tag>
          </template>
          <!-- Single mode: text -->
          <span v-else class="text-[var(--color-text-1)] py-0.5 pl-1">{{
            displaySingleValue
          }}</span>
        </template>
      </div>

      <div class="flex items-center">
        <!-- Clear Icon (only single mode, has selection, and allowClear is true) -->
        <a-tooltip
          v-if="allowClear && mode === 'single' && hasSelection"
          :content="t('feedCompare.selectCraftFlow.clear')"
        >
          <icon-close-circle
            class="text-[var(--color-text-3)] hover:text-[var(--color-text-2)] mr-2 z-10"
            role="button"
            :aria-label="t('feedCompare.selectCraftFlow.clear')"
            @click.stop="handleClear"
          />
        </a-tooltip>
        <icon-down class="text-[var(--color-text-3)]" />
      </div>
    </div>
    <div v-else @click="openModal">
      <slot></slot>
    </div>

    <!-- Selection Modal -->
    <a-modal
      v-model:visible="visible"
      :title="placeholder || t('feedCompare.selectCraftFlow.placeholder')"
      width="800px"
      :footer="mode === 'multiple'"
      @ok="confirmSelection"
      @cancel="cancelSelection"
    >
      <div class="mb-4">
        <a-input-search
          v-model="searchText"
          :placeholder="t('common.search') || 'Search...'"
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
                >{{ item.name }}</div
              >
              <div
                class="text-[var(--color-text-3)] text-xs line-clamp-3"
                :title="item.description"
              >
                {{ item.description || 'No description' }}
              </div>

              <!-- Selection Number for Multiple Mode -->
              <div
                v-if="isSelected(item.name) && mode === 'multiple'"
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
              {{ t('common.noData') || 'No results found' }}
            </div>
          </div>
        </a-tab-pane>
      </a-tabs>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted, computed, watch } from 'vue';
  import {
    CraftFlow,
    listCraftFlows,
    listCraftTemplates,
  } from '@/api/craft_flow';
  import { listCraftAtoms } from '@/api/craft_atom';
  import { useI18n } from 'vue-i18n';
  import {
    IconCloseCircle,
    IconDown,
    IconCheckCircleFill,
  } from '@arco-design/web-vue/es/icon';

  const { t } = useI18n();

  interface Props {
    modelValue?: string | string[];
    mode?: 'single' | 'multiple';
    placeholder?: string;
    allowClear?: boolean;
  }

  const props = withDefaults(defineProps<Props>(), {
    modelValue: () => [],
    mode: 'single',
    placeholder: '',
    allowClear: false,
  });

  const emit = defineEmits(['update:modelValue', 'change']);

  const visible = ref(false);
  const searchText = ref('');

  const craftFlows = ref<CraftFlow[]>([]);
  const sysCraftAtomList = ref<{ name: string; description?: string }[]>([]);
  const craftAtomList = ref<{ name: string; description?: string }[]>([]);

  // Internal value to track selection
  const localValue = ref<string[]>([]);

  const hasSelection = computed(
    () => localValue.value && localValue.value.length > 0,
  );

  const displaySingleValue = computed(() => {
    if (localValue.value.length > 0) return localValue.value[0];
    return '';
  });

  // Fetch data
  onMounted(async () => {
    try {
      const [craftFlowsResponse, sysCraftAtomsResponse, craftAtomsResponse] =
        await Promise.all([
          listCraftFlows(),
          listCraftTemplates(),
          listCraftAtoms(),
        ]);
      craftFlows.value = craftFlowsResponse.data || [];
      sysCraftAtomList.value = sysCraftAtomsResponse.data || [];
      craftAtomList.value = craftAtomsResponse.data || [];
    } catch (error) {
      // eslint-disable-next-line no-console
      console.error('Failed to fetch craft options', error);
    }
  });

  // Sync with prop
  watch(
    () => props.modelValue,
    (val) => {
      if (props.mode === 'single') {
        localValue.value = val ? [val as string] : [];
      } else {
        localValue.value = Array.isArray(val) ? [...val] : [];
      }
    },
    { immediate: true },
  );

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
      title: t('feedCompare.selectCraftFlow.tabs.system'),
      items: filterAndSort(sysCraftAtomList.value),
    },
    {
      key: 'user',
      title: t('feedCompare.selectCraftFlow.tabs.user'),
      items: filterAndSort(craftAtomList.value),
    },
    {
      key: 'flow',
      title: t('feedCompare.selectCraftFlow.tabs.flow'),
      items: filterAndSort(craftFlows.value),
    },
  ]);

  const isSelected = (name: string) => {
    return localValue.value.includes(name);
  };

  const getSelectionIndex = (name: string) => {
    return localValue.value.indexOf(name) + 1;
  };

  const emitUpdate = (val: string[]) => {
    let emitVal: string | string[] = val;
    if (props.mode === 'single') {
      emitVal = val.length > 0 ? val[0] : '';
    }
    emit('update:modelValue', emitVal);
    emit('change', emitVal);
  };

  const confirmSelection = () => {
    emitUpdate(localValue.value);
    visible.value = false;
  };

  const handleSelect = (name: string) => {
    if (props.mode === 'single') {
      localValue.value = [name];
      confirmSelection(); // Auto confirm for single selection
    } else {
      // Multiple
      const idx = localValue.value.indexOf(name);
      if (idx > -1) {
        localValue.value.splice(idx, 1);
      } else {
        localValue.value.push(name);
      }
    }
  };

  const removeItem = (name: string) => {
    const idx = localValue.value.indexOf(name);
    if (idx > -1) {
      const newValue = [...localValue.value];
      newValue.splice(idx, 1);
      localValue.value = newValue;
      emitUpdate(newValue);
    }
  };

  const handleClear = () => {
    localValue.value = [];
    emitUpdate([]);
  };

  const cancelSelection = () => {
    // Reset local value to prop value
    if (props.mode === 'single') {
      localValue.value = props.modelValue ? [props.modelValue as string] : [];
    } else {
      localValue.value = Array.isArray(props.modelValue)
        ? [...props.modelValue]
        : [];
    }
    visible.value = false;
  };

  const openModal = () => {
    visible.value = true;
    searchText.value = ''; // Reset search
    // Ensure localValue is synced with props when opening (in case prop changed externally)
    if (props.mode === 'single') {
      localValue.value = props.modelValue ? [props.modelValue as string] : [];
    } else {
      localValue.value = Array.isArray(props.modelValue)
        ? [...props.modelValue]
        : [];
    }
  };
</script>

<style scoped>
  /* Custom styles if needed, but Tailwind should handle most */
</style>
