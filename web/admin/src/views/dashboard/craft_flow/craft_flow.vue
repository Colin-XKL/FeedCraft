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
            editing = false;
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
          <a-button type="outline" @click="showEditModalHandler(record)"
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
      :title="editing ? 'Edit Craft Flow' : 'Create Craft Flow'"
    >
      <!--      @ok="updateCraftFlow"-->
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
          />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button
          @click="
            () => {
              showEditModal = false;
              isUpdating = false;
              editing = false;
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
  import { onBeforeMount, ref } from 'vue';
  import {
    CraftFlow,
    createCraftFlow,
    deleteCraftFlow,
    listCraftFlows,
    updateCraftFlow,
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
  const editing = ref(false);
  const isUpdating = ref(false);

  const columns = [
    { title: 'Name', dataIndex: 'name' },
    { title: 'Description', dataIndex: 'description' },
    { title: 'Craft Flow', slotName: 'craft-flow-item-list' },
    { title: 'Actions', slotName: 'actions' },
  ];

  const showEditModalHandler = (craftFlow: CraftFlow) => {
    editedCraftFlow.value = { ...craftFlow };
    showEditModal.value = true;
    editing.value = true;
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
          (craftConfigItem) => craftConfigItem.craftName
        ) ?? [];
      return ret;
    });
    isLoading.value = false;
  }

  const selectedCraftFlow = ref<CraftFlow | null>(null);
  const saveCraftFlow = async () => {
    if (editing.value) {
      if (selectedCraftFlow.value) {
        await updateCraftFlow(
          editedCraftFlow.value.name,
          transformCraftForOption(editedCraftFlow.value)
        );
        selectedCraftFlow.value.description = editedCraftFlow.value.description;
        selectedCraftFlow.value.name = editedCraftFlow.value.name;
        selectedCraftFlow.value.craft_flow_config =
          editedCraftFlow.value.craft_flow_config;
      }
    } else {
      await createCraftFlow(transformCraftForOption(editedCraftFlow.value));
      await listAllCraftFlow();
    }
    showEditModal.value = false;
    editedCraftFlow.value = {
      name: '',
      description: '',
      craft_flow_config: [],
    };
    editing.value = false;
    isUpdating.value = false;
    selectedCraftFlow.value = null;
  };

  onBeforeMount(() => {
    listAllCraftFlow();
  });
</script>

<script lang="ts">
  export default {
    name: 'CraftFlow',
  };
</script>
