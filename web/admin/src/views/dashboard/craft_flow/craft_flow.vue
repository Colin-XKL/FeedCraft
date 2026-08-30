<template>
  <CraftManagePage
    :title="t('menu.craftFlow')"
    :description="t('craftFlow.description')"
  >
    <template #toolbar>
      <a-space wrap>
        <a-button :loading="isLoading" @click="listAllCraftFlow">
          <template #icon>
            <icon-refresh />
          </template>
          {{ t('craftFlow.query') }}
        </a-button>
        <a-button type="primary" @click="handleAdd">
          <template #icon>
            <icon-plus />
          </template>
          {{ t('craftFlow.create') }}
        </a-button>
      </a-space>
    </template>

    <a-table
      v-if="isLoading || craftFlows.length > 0"
      row-key="name"
      :data="craftFlows"
      :columns="columns"
      :loading="isLoading"
      :bordered="false"
      :pagination="{ pageSize: 10, showTotal: true }"
    >
      <template #craft-flow-item-list="{ record }">
        <div class="craft-flow-chain">
          <a-tag color="gray">{{ t('craftFlow.flow.start') }}</a-tag>
          <template
            v-for="(item, index) in record.craft_flow_config"
            :key="index"
          >
            <span class="craft-flow-chain__arrow">/</span>
            <a-tooltip :content="getCraftDescription(item.craft_name)">
              <a-tag color="arcoblue">{{ item.craft_name }}</a-tag>
            </a-tooltip>
          </template>
          <span class="craft-flow-chain__arrow">/</span>
          <a-tag color="gray">{{ t('craftFlow.flow.end') }}</a-tag>
        </div>
      </template>
      <template #actions="{ record }">
        <a-space wrap>
          <a-button type="text" size="small" @click="editBtnHandler(record)">
            {{ t('craftFlow.edit') }}
          </a-button>
          <a-popconfirm
            :content="t('craftFlow.deleteConfirm')"
            @ok="deleteCraftFlowHandler(record.name)"
          >
            <a-button type="text" status="danger" size="small">
              {{ t('craftFlow.delete') }}
            </a-button>
          </a-popconfirm>
        </a-space>
      </template>
    </a-table>

    <ListEmptyGuide
      v-else-if="!listFailed"
      :description="t('craftFlow.empty.description')"
      :hint="t('craftFlow.empty.hint')"
      :create-label="t('craftFlow.empty.createFirst')"
      :docs-label="t('craftFlow.empty.docs')"
      :docs-href="flowDocsHref"
      @create="handleAdd"
    />

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
        layout="vertical"
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
  </CraftManagePage>
</template>

<script setup lang="ts">
  import CraftManagePage from '@/components/craft/CraftManagePage.vue';
  import ListEmptyGuide from '@/components/list-empty-guide/index.vue';
  import { onBeforeMount, ref, computed } from 'vue';
  import { buildDocsUrl } from '@/utils/docsUrl';
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

  const { t, locale } = useI18n();
  const flowDocsHref = computed(() =>
    buildDocsUrl(locale.value, 'guides/advanced/customization')
  );

  const rules = {
    name: [
      {
        required: true,
        message: t('craftFlow.form.rule.nameRequired'),
        trigger: 'blur',
      },
      namingValidator,
    ],
  };

  const isLoading = ref(false);
  const listFailed = ref(false);
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
    {
      title: t('craftFlow.edit'),
      slotName: 'actions',
      width: 140,
      align: 'right',
    },
  ];

  const handleAdd = () => {
    editedCraftFlow.value = {
      name: '',
      description: '',
      craftList: [],
      craft_flow_config: [],
    };
    showEditModal.value = true;
    isUpdating.value = false;
  };

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
      craftFlows.value = (res.data ?? []).map((item) => {
        const ret = item as any;
        const craftFlowConfigList = item.craft_flow_config ?? [];
        ret.craftList =
          craftFlowConfigList.map(
            (craftConfigItem) => craftConfigItem.craft_name
          ) ?? [];
        return ret;
      });
      listFailed.value = false;
    } catch {
      listFailed.value = true;
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
    } catch {
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

<style scoped lang="less">
  .craft-flow-chain {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 6px;
  }

  .craft-flow-chain__arrow {
    color: var(--color-text-4);
  }
</style>
