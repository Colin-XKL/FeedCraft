<template>
  <div class="craft-flow-editor">
    <div
      v-if="modelValue.length === 0"
      class="text-center text-gray-400 py-4 border border-dashed border-gray-300 rounded mb-4"
    >
      {{ t('blueprint.editor.empty') }}
    </div>

    <div v-else class="space-y-3 mb-4">
      <div
        v-for="(processor, index) in modelValue"
        :key="index"
        class="flex items-center justify-between p-3 border border-gray-200 rounded bg-white shadow-sm hover:shadow-md transition-shadow"
      >
        <div class="flex items-center gap-3">
          <div
            class="w-6 h-6 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center text-xs font-bold"
          >
            {{ index + 1 }}
          </div>
          <span class="font-medium text-gray-800">{{ processor }}</span>
        </div>

        <div class="flex items-center gap-2">
          <a-tooltip :content="t('blueprint.editor.moveUp')">
            <a-button
              size="small"
              type="text"
              :aria-label="t('blueprint.editor.moveUp')"
              :disabled="index === 0"
              @click="moveUp(index)"
            >
              <template #icon><icon-arrow-up /></template>
            </a-button>
          </a-tooltip>
          <a-tooltip :content="t('blueprint.editor.moveDown')">
            <a-button
              size="small"
              type="text"
              :aria-label="t('blueprint.editor.moveDown')"
              :disabled="index === modelValue.length - 1"
              @click="moveDown(index)"
            >
              <template #icon><icon-arrow-down /></template>
            </a-button>
          </a-tooltip>
          <a-tooltip :content="t('blueprint.editor.remove')">
            <a-button
              size="small"
              type="text"
              status="danger"
              :aria-label="t('blueprint.editor.remove')"
              @click="removeProcessor(index)"
            >
              <template #icon><icon-delete /></template>
            </a-button>
          </a-tooltip>
        </div>
      </div>
    </div>

    <ProcessorSelector
      mode="single"
      :model-value="''"
      :placeholder="t('blueprint.editor.addPlaceholder')"
      @update:model-value="addProcessor"
    >
      <a-button type="dashed" long class="flex items-center justify-center">
        <template #icon><icon-plus /></template>
        {{ t('blueprint.editor.add') }}
      </a-button>
    </ProcessorSelector>
  </div>
</template>

<script setup lang="ts">
  import { useI18n } from 'vue-i18n';
  import ProcessorSelector from './ProcessorSelector.vue';

  const { t } = useI18n();

  const props = defineProps<{
    modelValue: string[];
  }>();

  const emit = defineEmits<{
    (e: 'update:modelValue', value: string[]): void;
  }>();

  const addProcessor = (craft: string | string[]) => {
    if (!craft || (Array.isArray(craft) && craft.length === 0)) return;
    const processorName = Array.isArray(craft) ? craft[0] : craft;
    const newList = [...props.modelValue, processorName];
    emit('update:modelValue', newList);
  };

  const removeProcessor = (index: number) => {
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
