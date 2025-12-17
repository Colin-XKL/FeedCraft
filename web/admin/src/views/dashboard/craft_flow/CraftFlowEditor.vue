<template>
  <div class="craft-flow-editor">
    <div
      v-if="modelValue.length === 0"
      class="text-center text-gray-400 py-4 border border-dashed border-gray-300 rounded mb-4"
    >
      {{ t('craftFlow.editor.empty') }}
    </div>

    <div v-else class="space-y-3 mb-4">
      <div
        v-for="(craft, index) in modelValue"
        :key="index"
        class="flex items-center justify-between p-3 border border-gray-200 rounded bg-white shadow-sm hover:shadow-md transition-shadow"
      >
        <div class="flex items-center gap-3">
          <div
            class="w-6 h-6 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center text-xs font-bold"
          >
            {{ index + 1 }}
          </div>
          <span class="font-medium text-gray-800">{{ craft }}</span>
        </div>

        <div class="flex items-center gap-2">
          <a-button
            size="small"
            type="text"
            :disabled="index === 0"
            @click="moveUp(index)"
          >
            <template #icon><icon-arrow-up /></template>
          </a-button>
          <a-button
            size="small"
            type="text"
            :disabled="index === modelValue.length - 1"
            @click="moveDown(index)"
          >
            <template #icon><icon-arrow-down /></template>
          </a-button>
          <a-button
            size="small"
            type="text"
            status="danger"
            @click="removeCraft(index)"
          >
            <template #icon><icon-delete /></template>
          </a-button>
        </div>
      </div>
    </div>

    <CraftSelector
      mode="single"
      :model-value="''"
      :placeholder="t('craftFlow.editor.addPlaceholder')"
      @update:model-value="addCraft"
    >
      <a-button type="dashed" long class="flex items-center justify-center">
        <template #icon><icon-plus /></template>
        {{ t('craftFlow.editor.add') }}
      </a-button>
    </CraftSelector>
  </div>
</template>

<script setup lang="ts">
  import { useI18n } from 'vue-i18n';
  import CraftSelector from './CraftSelector.vue';

  const { t } = useI18n();

  const props = defineProps<{
    modelValue: string[];
  }>();

  const emit = defineEmits<{
    (e: 'update:modelValue', value: string[]): void;
  }>();

  const addCraft = (craft: string | string[]) => {
    if (!craft || (Array.isArray(craft) && craft.length === 0)) return;
    const craftName = Array.isArray(craft) ? craft[0] : craft;
    const newList = [...props.modelValue, craftName];
    emit('update:modelValue', newList);
  };

  const removeCraft = (index: number) => {
    const newList = [...props.modelValue];
    newList.splice(index, 1);
    emit('update:modelValue', newList);
  };

  const moveUp = (index: number) => {
    if (index <= 0) return;
    const newList = [...props.modelValue];
    const item = newList[index];
    newList.splice(index, 1);
    newList.splice(index - 1, 0, item);
    emit('update:modelValue', newList);
  };

  const moveDown = (index: number) => {
    if (index >= props.modelValue.length - 1) return;
    const newList = [...props.modelValue];
    const item = newList[index];
    newList.splice(index, 1);
    newList.splice(index + 1, 0, item);
    emit('update:modelValue', newList);
  };
</script>

<style scoped>
  .craft-flow-editor {
    width: 100%;
  }
</style>
