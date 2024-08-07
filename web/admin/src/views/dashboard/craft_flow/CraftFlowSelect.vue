<template>
  <a-select
    v-model="selectedCraftFlow"
    mode="multiple"
    placeholder="Select Craft Flow"
    allow-create
    allow-clear
  >
    <a-optgroup label="系统 Craft Atoms">
      <a-option
        v-for="item in sysCraftAtomList"
        :key="item.name"
        :value="item.name"
      >
        {{
          item.description?.length
            ? `${item.name} (${item.description})`
            : item.name
        }}
      </a-option>
    </a-optgroup>
    <a-optgroup label="用户 Craft Atoms">
      <a-option
        v-for="item in craftAtomList"
        :key="item.name"
        :value="item.name"
      >
        {{
          item.description?.length
            ? `${item.name} (${item.description})`
            : item.name
        }}
      </a-option>
    </a-optgroup>
    <a-optgroup label="Craft Flows">
      <a-option v-for="item in craftFlows" :key="item.name" :value="item.name">
        {{
          item.description?.length
            ? `${item.name} (${item.description})`
            : item.name
        }}
      </a-option>
    </a-optgroup>
  </a-select>
</template>

<script lang="ts">
  import { ref, onMounted, watch } from 'vue';
  import {
    CraftFlow,
    listCraftFlows,
    listSysCraftAtoms,
  } from '@/api/craft_flow';
  import { listCraftAtoms } from '@/api/craft_atom';

  export default {
    name: 'CraftFlowSelect',
    props: {
      modelValue: {
        type: Array as () => string[],
        default: () => [],
      },
    },
    emits: ['update:modelValue'],
    setup(props, { emit }) {
      const craftFlows = ref<CraftFlow[]>([]);
      const sysCraftAtomList = ref<CraftFlow[]>([]);
      const craftAtomList = ref<CraftFlow[]>([]);
      const selectedCraftFlow = ref<string[]>(props.modelValue);

      onMounted(async () => {
        const [craftFlowsResponse, sysCraftAtomsResponse, craftAtomsResponse] =
          await Promise.all([
            listCraftFlows(),
            listSysCraftAtoms(),
            listCraftAtoms(),
          ]);
        craftFlows.value = craftFlowsResponse.data as CraftFlow[];
        sysCraftAtomList.value = sysCraftAtomsResponse.data as {
          name: string;
          description?: string;
        }[];
        craftAtomList.value = craftAtomsResponse.data as {
          name: string;
          description?: string;
        }[];
      });

      watch(selectedCraftFlow, (newValue) => {
        emit('update:modelValue', newValue);
      });

      watch(
        () => props.modelValue,
        (newValue) => {
          selectedCraftFlow.value = newValue;
        }
      );

      return {
        craftFlows,
        sysCraftAtomList,
        craftAtomList,
        selectedCraftFlow,
      };
    },
  };
</script>
