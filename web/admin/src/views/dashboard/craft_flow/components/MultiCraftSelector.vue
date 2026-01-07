<template>
  <div class="multi-craft-selector">
    <div
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
          <a-tag
            v-for="val in modelValue"
            :key="val"
            closable
            class="mr-1 my-0.5"
            @close.stop="removeItem(val)"
          >
            {{ val }}
          </a-tag>
        </template>
      </div>

      <div class="flex items-center">
        <icon-down class="text-[var(--color-text-3)]" />
      </div>
    </div>

    <CraftPickerModal
      v-model:visible="visible"
      :model-value="modelValue"
      mode="multiple"
      :title="placeholder"
      @update:model-value="handleUpdate"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { IconDown } from '@arco-design/web-vue/es/icon';
import CraftPickerModal from './CraftPickerModal.vue';

const { t } = useI18n();

const props = defineProps({
  modelValue: {
    type: Array as () => string[],
    default: () => [],
  },
  placeholder: {
    type: String,
    default: '',
  },
});

const emit = defineEmits(['update:modelValue', 'change']);

const visible = ref(false);

const hasSelection = computed(
  () => props.modelValue && props.modelValue.length > 0,
);

const openModal = () => {
  visible.value = true;
};

const handleUpdate = (val: string | string[]) => {
  const newVal = val as string[];
  emit('update:modelValue', newVal);
  emit('change', newVal);
};

const removeItem = (val: string) => {
  const newVal = [...props.modelValue];
  const idx = newVal.indexOf(val);
  if (idx > -1) {
    newVal.splice(idx, 1);
    emit('update:modelValue', newVal);
    emit('change', newVal);
  }
};
</script>
