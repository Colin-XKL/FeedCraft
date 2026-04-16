<template>
  <CraftSelector
    :model-value="modelValue"
    :mode="mode"
    :placeholder="t('feedCompare.selectCraftFlow.placeholder')"
    allow-clear
    @update:model-value="onUpdate"
  />
</template>

<script setup lang="ts">
  import { useI18n } from 'vue-i18n';
  import CraftSelector from '@/views/dashboard/craft_flow/CraftSelector.vue';

  const { t } = useI18n();
  defineProps<{
    modelValue: string[] | string;
    mode: 'single' | 'multiple';
  }>();

  const emit = defineEmits<{
    (event: 'update:modelValue', value: string | string[]): void;
  }>();

  const onUpdate = (val: string | string[]) => {
    emit('update:modelValue', val);
  };
</script>
