<template>
  <div class="py-8 px-16">
    <x-header title="RSS Feed Preview" description="RSS Feed Preview">
    </x-header>

    <a-space direction="horizontal" class="mb-6">
      <a-button type="primary" :loading="isLoading" @click="listAllCraftFlow">
        List
      </a-button>
      <a-button
        type="outline"
        @click="
          () => {
            showEditModal = true;
            isUpdating = false;
          }
        "
        >Create CraftFlow
      </a-button>
    </a-space>

    <a-table :data="craftFlows" :columns="columns" :loading="isLoading">
      <template #craft-flow-item-list="{ record }">
        <a-tag>START</a-tag>
        >
        <template
          v-for="(item, index) in record.craft_flow_config"
          :key="index"
        >
          <a-tag color="arcoblue">{{ item.craft_name }}</a-tag>
          >
        </template>
        <a-tag>END</a-tag>
      </template>
      <template #actions="{ record }">
        <a-space>
          <a-button type="outline" @click="editBtnHandler(record)"
            >Edit
          </a-button>
          <a-button status="danger" @click="deleteCraftFlowHandler(record.name)"
            >Delete
          </a-button>
        </a-space>
      </template>
    </a-table>

    <a-modal
      v-model:visible="showEditModal"
      :title="isUpdating ? 'Edit Craft Flow' : 'Create Craft Flow'"
    >
      <a-form
        :model="editedCraftFlow"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 18 }"
      >
        <a-form-item label="Name" name="name">
          <a-input v-model="editedCraftFlow.name" />
        </a-form-item>
        <a-form-item label="Description" name="description">
          <a-textarea v-model="editedCraftFlow.description" />
        </a-form-item>
        <a-form-item label="Flow" name="craftFlowConfig">
          <a-select
            v-model="editedCraftFlow.craftList"
            multiple
            allow-clear
            allow-create
            :options="optionList"
            option-group
          />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button
          @click="
            () => {
              showEditModal = false;
              isUpdating = false;
            }
          "
          >Cancel
        </a-button>
        <a-button type="primary" @click="saveCraftFlow">Save</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
  import XHeader from '@/components/header/x-header.vue';
  import { computed, onBeforeMount, ref } from 'vue';
  import {
    CraftFlow,
    createCraftFlow,
    deleteCraftFlow,
    listSysCraftAtoms,
    listCraftFlows,
    updateCraftFlow,
    listCraftAtoms,
  } from '@/api/craft_flow';

  const isLoading = ref(false);
  const craftFlows = ref<CraftFlow[]>([]);
  const editedCraftFlow = ref<CraftFlow>({
    name: '',
    description: '',
    craft_flow_config: [],
  });
  // const showCreateModal = ref(false);
  const showEditModal = ref(false);
  const isUpdating = ref(false);

  const columns = [
    { title: 'Name', dataIndex: 'name' },
    { title: 'Description', dataIndex: 'description' },
    { title: 'Craft Flow', slotName: 'craft-flow-item-list' },
    { title: 'Actions', slotName: 'actions' },
  ];
  const optionList = computed(() => {
    const mapper = (item: any) => {
      return {
        value: item.name,
        label: item.description?.length
          ? `${item.name} (${item.description})`
          : item.name,
      };
    };
    return [
      {
        label: 'System Craft Atoms',
        options: sysCraftAtomList.value.map(mapper),
      },
      {
        label: 'Craft Atoms',
        options: craftAtomList.value.map(mapper),
      },
      {
        label: 'Craft Flows',
        options: craftFlows.value.map(mapper),
      },
    ];
  });

  const editBtnHandler = (craftFlow: CraftFlow) => {
    editedCraftFlow.value = { ...craftFlow };
    showEditModal.value = true;
    isUpdating.value = true;
  };

  const deleteCraftFlowHandler = async (name: string) => {
    await deleteCraftFlow(name);
    await listAllCraftFlow();
  };

  // transform before sending requests
  function transformCraftForOption(craftFlowOrigin: any) {
    const { craftList, ...craftFlow } = craftFlowOrigin;
    // eslint-disable-next-line camelcase
    craftFlow.craft_flow_config =
      craftList?.map((item: string) => {
        return {
          craft_name: item,
          // todo implement custom option field
        };
      }) ?? [];
    return craftFlow;
  }

  async function listAllCraftFlow() {
    isLoading.value = true;
    craftFlows.value = (await listCraftFlows()).data.map((item) => {
      const ret = item;
      const craftFlowConfigList = item.craft_flow_config ?? [];
      ret.craftList =
        craftFlowConfigList.map(
          (craftConfigItem) => craftConfigItem.craft_name
        ) ?? [];
      return ret;
    });
    isLoading.value = false;
  }

  const saveCraftFlow = async () => {
    if (isUpdating.value) {
      await updateCraftFlow(
        editedCraftFlow.value.name,
        transformCraftForOption(editedCraftFlow.value)
      );
    } else {
      await createCraftFlow(transformCraftForOption(editedCraftFlow.value));
    }
    showEditModal.value = false;
    await listAllCraftFlow();
    isUpdating.value = false;
    editedCraftFlow.value = {
      name: '',
      description: '',
      craft_flow_config: [],
    };
  };
  const sysCraftAtomList = ref<any>([]);
  const craftAtomList = ref<any>([]);

  async function listAllSysCraftAtom() {
    sysCraftAtomList.value = (await listSysCraftAtoms()).data;
  }

  async function listAllCraftAtoms() {
    craftAtomList.value = (await listCraftAtoms()).data;
  }

  onBeforeMount(() => {
    listAllCraftFlow();
    listAllSysCraftAtom();
    listAllCraftAtoms();
  });
</script>

<script lang="ts">
  export default {
    name: 'CraftFlow',
  };
</script>
