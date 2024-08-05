<template>
  <a-select
    v-model="selectedCraftFlow"
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
  import { listCraftFlows, listSysCraftAtoms } from '@/api/craft_flow';
  import { listCraftAtoms } from '@/api/craft_atom';

  export default {
    name: 'CraftFlowSelect',
    props: {
      modelValue: {
        type: String,
        default: '',
      },
    },
    emits: ['update:modelValue'],
    setup(props, { emit }) {
      const craftFlows = ref([]);
      const sysCraftAtomList = ref([]);
      const craftAtomList = ref([]);
      const selectedCraftFlow = ref(props.modelValue);

      onMounted(async () => {
        const [craftFlowsResponse, sysCraftAtomsResponse, craftAtomsResponse] =
          await Promise.all([
            listCraftFlows(),
            listSysCraftAtoms(),
            listCraftAtoms(),
          ]);
        craftFlows.value = craftFlowsResponse.data;
        sysCraftAtomList.value = sysCraftAtomsResponse.data;
        craftAtomList.value = craftAtomsResponse.data;
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
