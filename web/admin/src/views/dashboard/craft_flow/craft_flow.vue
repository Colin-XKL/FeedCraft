<template>
  <div class="py-8 px-16">
    <x-header title="RSS 源预览" description="RSS 源预览">
    </x-header>

    <a-space direction="horizontal" class="mb-6">
      <a-button type="primary" :loading="isLoading" @click="listAllCraftFlow">
        列表
      </a-button>
      <a-button
        type="outline"
        @click="
          () => {
            showEditModal = true;
            isUpdating = false;
          }
        "
        >创建 CraftFlow
      </a-button>
    </a-space>

    <a-table :data="craftFlows" :columns="columns" :loading="isLoading">
      <template #craft-flow-item-list="{ record }">
        <a-tag>开始</a-tag>
        >
        <template
          v-for="(item, index) in record.craft_flow_config"
          :key="index"
        >
          <a-tag color="arcoblue">{{ item.craft_name }}</a-tag>
          >
        </template>
        <a-tag>结束</a-tag>
      </template>
      <template #actions="{ record }">
        <a-space>
          <a-button type="outline" @click="editBtnHandler(record)"
            >编辑
          </a-button>
          <a-button status="danger" @click="deleteCraftFlowHandler(record.name)"
            >删除
          </a-button>
        </a-space>
      </template>
    </a-table>

    <a-modal
      v-model:visible="showEditModal"
      :title="isUpdating ? '编辑 Craft Flow' : '创建 Craft Flow'"
    >
      <a-form
        :model="editedCraftFlow"
        :rules="rules"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 18 }"
      >
        <a-form-item label="名称" field="name">
          <a-input v-model="editedCraftFlow.name" />
        </a-form-item>
        <a-form-item label="描述" field="description">
          <a-textarea v-model="editedCraftFlow.description" />
        </a-form-item>
        <a-form-item label="流程" field="craftFlowConfig">
          <a-select
            v-model="editedCraftFlow.craftList"
            multiple
            allow-clear
            allow-create
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
              <a-option
                v-for="item in craftFlows"
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
          </a-select>
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
          >取消
        </a-button>
        <a-button type="primary" @click="saveCraftFlow">保存</a-button>
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
  } from '@/api/craft_flow';
  import { listCraftAtoms } from '@/api/craft_atom';
  import { namingValidator } from '@/utils/validator';

  const rules = {
    name: [
      {
        required: true,
        message: 'Name is required',
        trigger: 'blur',
      },
      namingValidator,
    ],
  };


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
