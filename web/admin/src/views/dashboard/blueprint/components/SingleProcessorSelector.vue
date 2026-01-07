<template>
  <div class="single-craft-selector">
    <div
      class="trigger-input w-full border border-[var(--color-neutral-3)] rounded bg-[var(--color-bg-2)] px-3 py-1.5 min-h-[32px] cursor-pointer flex items-center justify-between hover:border-[rgb(var(--primary-6))] transition-colors focus:outline-none focus:ring-2 focus:ring-[rgba(var(--primary-6),0.2)] focus:border-[rgb(var(--primary-6))]"
      tabindex="0"
      role="button"
      aria-haspopup="dialog"
      :aria-label="placeholder || t('feedCompare.selectBlueprint.placeholder')"
      @click="openModal"
      @keydown.enter.prevent="openModal"
      @keydown.space.prevent="openModal"
    >
      <div class="flex flex-wrap gap-1 flex-1 overflow-hidden">
        <span v-if="!modelValue" class="text-[var(--color-text-3)]">{{
          placeholder || t('feedCompare.selectBlueprint.placeholder')
        }}</span>
        <span v-else class="text-[var(--color-text-1)] py-0.5 pl-1">{{
          modelValue
        }}</span>
      </div>

      <div class="flex items-center">
        <a-tooltip
          v-if="allowClear && modelValue"
          :content="t('feedCompare.selectBlueprint.clear')"
        >
          <icon-close-circle
            class="text-[var(--color-text-3)] hover:text-[var(--color-text-2)] mr-2 z-10"
            role="button"
            :aria-label="t('feedCompare.selectBlueprint.clear')"
            @click.stop="handleClear"
          />
        </a-tooltip>
        <icon-down class="text-[var(--color-text-3)]" />
      </div>
    </div>

    <ProcessorPickerModal
      v-model:visible="visible"
      :model-value="modelValue"
      mode="single"
      :title="placeholder"
      @update:model-value="handleUpdate"
    />
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue';
  import { useI18n } from 'vue-i18n';
  import { IconCloseCircle, IconDown } from '@arco-design/web-vue/es/icon';
  import ProcessorPickerModal from './ProcessorPickerModal.vue';

  const { t } = useI18n();

  const props = defineProps({
    modelValue: {
      type: String,
      default: '',
    },
    placeholder: {
      type: String,
      default: '',
    },
    allowClear: {
      type: Boolean,
      default: false,
    },
  });

  const emit = defineEmits(['update:modelValue', 'change']);

  const visible = ref(false);

  const openModal = () => {
    visible.value = true;
  };

  const handleUpdate = (val: string | string[]) => {
    const newVal = val as string;
    emit('update:modelValue', newVal);
    emit('change', newVal);
  };

  const handleClear = () => {
    emit('update:modelValue', '');
    emit('change', '');
  };
</script>
