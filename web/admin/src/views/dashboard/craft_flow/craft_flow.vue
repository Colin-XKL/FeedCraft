<template>
  <div class="py-8 px-16">
    <x-header
      :title="t('menu.craftFlow')"
      :description="t('craftFlow.description')"
    ></x-header>

    <a-space direction="horizontal" class="mb-6">
      <a-button type="primary" :loading="isLoading" @click="listAllCraftFlow">
        {{ t('craftFlow.query') }}
      </a-button>
      <a-button
        type="outline"
        @click="
          () => {
            showEditModal = true;
            isUpdating = false;
            editedCraftFlow = {
              name: '',
              description: '',
              craftList: [],
            };
          }
        "
        >{{ t('craftFlow.create') }}
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
          <a-tooltip :content="getCraftDescription(item.craft_name)">
            <a-tag color="arcoblue">{{ item.craft_name }}</a-tag>
          </a-tooltip>
          >
        </template>
        <a-tag>结束</a-tag>
      </template>
      <template #actions="{ record }">
        <a-space>
          <a-button type="outline" @click="editBtnHandler(record)"
            >{{ t('craftFlow.edit') }}
          </a-button>
          <a-popconfirm
            :content="t('craftFlow.deleteConfirm')"
            @ok="deleteCraftFlowHandler(record.name)"
          >
            <a-button status="danger">{{ t('craftFlow.delete') }}</a-button>
          </a-popconfirm>
        </a-space>
      </template>
    </a-table>

    <a-modal
      v-model:visible="showEditModal"
      :title="
        isUpdating
          ? t('craftFlow.editModalTitle.edit')
          : t('craftFlow.editModalTitle.create')
      "
    >
      <a-form
        ref="formRef"
        :model="editedCraftFlow"
        :rules="rules"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 18 }"
      >
        <a-form-item :label="t('craftFlow.form.name')" field="name">
          <a-input v-model="editedCraftFlow.name" />
        </a-form-item>
        <a-form-item
          :label="t('craftFlow.form.description')"
          field="description"
        >
          <a-textarea v-model="editedCraftFlow.description" />
        </a-form-item>
        <a-form-item :label="t('craftFlow.form.flow')" field="craftFlowConfig">
          <CraftFlowEditor v-model="editedCraftFlow.craftList" />
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
          >{{ t('craftFlow.form.cancel') }}
        </a-button>
        <a-button type="primary" :loading="saving" @click="saveCraftFlow">{{
          t('craftFlow.form.save')
        }}</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
  import XHeader from '@/components/header/x-header.vue';
  import { onBeforeMount, ref, computed } from 'vue';
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
  import CraftFlowEditor from '@/views/dashboard/craft_flow/CraftFlowEditor.vue';
  import { useI18n } from 'vue-i18n';
  import { Message } from '@arco-design/web-vue';

  const { t } = useI18n();

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
  const saving = ref(false);
  const formRef = ref();
  const craftFlows = ref<CraftFlow[]>([]);
  const editedCraftFlow = ref<any>({
    name: '',
    description: '',
    craftList: [], // craftList should be initialized
    craft_flow_config: [],
  });
  // const showCreateModal = ref(false);
  const showEditModal = ref(false);
  const isUpdating = ref(false);

  const columns = [
    { title: t('craftFlow.form.name'), dataIndex: 'name' },
    { title: t('craftFlow.form.description'), dataIndex: 'description' },
    { title: t('craftFlow.form.flow'), slotName: 'craft-flow-item-list' },
    { title: t('craftFlow.edit'), slotName: 'actions' },
  ];

  const editBtnHandler = (craftFlow: CraftFlow) => {
    // Clone and ensure craftList exists
    const craftFlowCopy = { ...craftFlow } as any;
    if (!craftFlowCopy.craftList && craftFlowCopy.craft_flow_config) {
      craftFlowCopy.craftList = craftFlowCopy.craft_flow_config.map(
        (c: any) => c.craft_name
      );
    } else if (!craftFlowCopy.craftList) {
      craftFlowCopy.craftList = [];
    }
    editedCraftFlow.value = craftFlowCopy;
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
    try {
      const res = await listCraftFlows();
      craftFlows.value = res.data.map((item) => {
        const ret = item as any;
        const craftFlowConfigList = item.craft_flow_config ?? [];
        ret.craftList =
          craftFlowConfigList.map(
            (craftConfigItem) => craftConfigItem.craft_name
          ) ?? [];
        return ret;
      });
    } finally {
      isLoading.value = false;
    }
  }

  const saveCraftFlow = async () => {
    const res = await formRef.value?.validate();
    if (res) return;

    saving.value = true;
    try {
      if (isUpdating.value) {
        await updateCraftFlow(
          editedCraftFlow.value.name,
          transformCraftForOption(editedCraftFlow.value)
        );
      } else {
        await createCraftFlow(transformCraftForOption(editedCraftFlow.value));
      }
      Message.success(t('craftFlow.form.saveSuccess'));
      showEditModal.value = false;
      await listAllCraftFlow();
      isUpdating.value = false;
      editedCraftFlow.value = {
        name: '',
        description: '',
        craftList: [],
        craft_flow_config: [],
      };
    } catch (err) {
      // Error handling is done by interceptor or default handling
    } finally {
      saving.value = false;
    }
  };
  const sysCraftAtomList = ref<any>([]);
  const craftAtomList = ref<any>([]);

  async function listAllSysCraftAtom() {
    sysCraftAtomList.value = (await listSysCraftAtoms()).data;
  }

  async function listAllCraftAtoms() {
    craftAtomList.value = (await listCraftAtoms()).data;
  }

  const craftDescriptionMap = computed(() => {
    const map = new Map<string, string>();
    sysCraftAtomList.value.forEach((item: any) => {
      map.set(item.name, item.description);
    });
    craftAtomList.value.forEach((item: any) => {
      map.set(item.name, item.description);
    });
    return map;
  });

  function getCraftDescription(name: string) {
    return craftDescriptionMap.value.get(name) || '';
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
