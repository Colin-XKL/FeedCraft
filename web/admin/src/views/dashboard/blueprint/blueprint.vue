<template>
  <div class="py-8 px-16">
    <x-header
      :title="t('menu.blueprint')"
      :description="t('blueprint.description')"
    ></x-header>

    <a-space direction="horizontal" class="mb-6">
      <a-button type="primary" :loading="isLoading" @click="listAllBlueprint">
        {{ t('blueprint.query') }}
      </a-button>
      <a-button
        type="outline"
        @click="
          () => {
            showEditModal = true;
            isUpdating = false;
            editedBlueprint = {
              name: '',
              description: '',
              processorList: [],
            };
          }
        "
        >{{ t('blueprint.create') }}
      </a-button>
    </a-space>

    <a-table :data="blueprints" :columns="columns" :loading="isLoading">
      <template #craft-flow-item-list="{ record }">
        <a-tag>开始</a-tag>
        >
        <template v-for="(item, index) in record.blueprint_config" :key="index">
          <a-tooltip :content="getCraftDescription(item.processor_name)">
            <a-tag color="arcoblue">{{ item.processor_name }}</a-tag>
          </a-tooltip>
          >
        </template>
        <a-tag>结束</a-tag>
      </template>
      <template #actions="{ record }">
        <a-space>
          <a-button type="outline" @click="editBtnHandler(record)"
            >{{ t('blueprint.edit') }}
          </a-button>
          <a-popconfirm
            :content="t('blueprint.deleteConfirm')"
            @ok="deleteBlueprintHandler(record.name)"
          >
            <a-button status="danger">{{ t('blueprint.delete') }}</a-button>
          </a-popconfirm>
        </a-space>
      </template>
    </a-table>

    <a-modal
      v-model:visible="showEditModal"
      :title="
        isUpdating
          ? t('blueprint.editModalTitle.edit')
          : t('blueprint.editModalTitle.create')
      "
    >
      <a-form
        ref="formRef"
        :model="editedBlueprint"
        :rules="rules"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 18 }"
      >
        <a-form-item :label="t('blueprint.form.name')" field="name">
          <a-input v-model="editedBlueprint.name" />
        </a-form-item>
        <a-form-item
          :label="t('blueprint.form.description')"
          field="description"
        >
          <a-textarea v-model="editedBlueprint.description" />
        </a-form-item>
        <a-form-item :label="t('blueprint.form.flow')" field="blueprintConfig">
          <BlueprintEditor v-model="editedBlueprint.processorList" />
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
          >{{ t('blueprint.form.cancel') }}
        </a-button>
        <a-button type="primary" :loading="saving" @click="saveBlueprint">{{
          t('blueprint.form.save')
        }}</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
  import XHeader from '@/components/header/x-header.vue';
  import { onBeforeMount, ref, computed } from 'vue';
  import {
    Blueprint,
    createBlueprint,
    deleteBlueprint,
    listSysTools,
    listBlueprints,
    updateBlueprint,
  } from '@/api/blueprint';
  import { listTools } from '@/api/tool';
  import { namingValidator } from '@/utils/validator';
  import BlueprintEditor from '@/views/dashboard/blueprint/BlueprintEditor.vue';
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
  const blueprints = ref<Blueprint[]>([]);
  const editedBlueprint = ref<any>({
    name: '',
    description: '',
    processorList: [], // processorList should be initialized
    blueprint_config: [],
  });
  // const showCreateModal = ref(false);
  const showEditModal = ref(false);
  const isUpdating = ref(false);

  const columns = [
    { title: t('blueprint.form.name'), dataIndex: 'name' },
    { title: t('blueprint.form.description'), dataIndex: 'description' },
    { title: t('blueprint.form.flow'), slotName: 'craft-flow-item-list' },
    { title: t('blueprint.edit'), slotName: 'actions' },
  ];

  const editBtnHandler = (blueprint: Blueprint) => {
    // Clone and ensure processorList exists
    const blueprintCopy = { ...blueprint } as any;
    if (!blueprintCopy.processorList && blueprintCopy.blueprint_config) {
      blueprintCopy.processorList = blueprintCopy.blueprint_config.map(
        (c: any) => c.processor_name,
      );
    } else if (!blueprintCopy.processorList) {
      blueprintCopy.processorList = [];
    }
    editedBlueprint.value = blueprintCopy;
    showEditModal.value = true;
    isUpdating.value = true;
  };

  const deleteBlueprintHandler = async (name: string) => {
    await deleteBlueprint(name);
    await listAllBlueprint();
  };

  // transform before sending requests
  function transformCraftForOption(blueprintOrigin: any) {
    const { processorList, ...blueprint } = blueprintOrigin;
    // eslint-disable-next-line camelcase
    blueprint.blueprint_config =
      processorList?.map((item: string) => {
        return {
          processor_name: item,
          // todo implement custom option field
        };
      }) ?? [];
    return blueprint;
  }

  async function listAllBlueprint() {
    isLoading.value = true;
    try {
      const res = await listBlueprints();
      blueprints.value = res.data.map((item) => {
        const ret = item as any;
        const blueprintConfigList = item.blueprint_config ?? [];
        ret.processorList =
          blueprintConfigList.map(
            (craftConfigItem) => craftConfigItem.processor_name,
          ) ?? [];
        return ret;
      });
    } finally {
      isLoading.value = false;
    }
  }

  const saveBlueprint = async () => {
    const res = await formRef.value?.validate();
    if (res) return;

    saving.value = true;
    try {
      if (isUpdating.value) {
        await updateBlueprint(
          editedBlueprint.value.name,
          transformCraftForOption(editedBlueprint.value),
        );
      } else {
        await createBlueprint(transformCraftForOption(editedBlueprint.value));
      }
      Message.success(t('blueprint.form.saveSuccess'));
      showEditModal.value = false;
      await listAllBlueprint();
      isUpdating.value = false;
      editedBlueprint.value = {
        name: '',
        description: '',
        processorList: [],
        blueprint_config: [],
      };
    } catch (err) {
      // Error handling is done by interceptor or default handling
    } finally {
      saving.value = false;
    }
  };
  const sysToolList = ref<any>([]);
  const toolList = ref<any>([]);

  async function listAllSysTool() {
    sysToolList.value = (await listSysTools()).data;
  }

  async function listAllTools() {
    toolList.value = (await listTools()).data;
  }

  const craftDescriptionMap = computed(() => {
    const map = new Map<string, string>();
    sysToolList.value.forEach((item: any) => {
      map.set(item.name, item.description);
    });
    toolList.value.forEach((item: any) => {
      map.set(item.name, item.description);
    });
    return map;
  });

  function getCraftDescription(name: string) {
    return craftDescriptionMap.value.get(name) || '';
  }

  onBeforeMount(() => {
    listAllBlueprint();
    listAllSysTool();
    listAllTools();
  });
</script>

<script lang="ts">
  export default {
    name: 'Blueprint',
  };
</script>
