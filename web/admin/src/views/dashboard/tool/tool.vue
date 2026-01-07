<template>
  <div class="py-8 px-16">
    <x-header
      :title="t('menu.tool')"
      :description="t('tool.description')"
    ></x-header>

    <a-space direction="horizontal" class="mb-6">
      <a-button type="primary" :loading="isLoading" @click="listAllTools">
        {{ t('tool.query') }}
      </a-button>
      <a-button
        type="outline"
        @click="
          () => {
            showEditModal = true;
            isUpdating = false;
          }
        "
        >{{ t('tool.create') }}
      </a-button>
    </a-space>

    <a-table :data="tools" :columns="columns" :loading="isLoading">
      <template #actions="{ record }">
        <a-space>
          <a-button type="outline" @click="editBtnHandler(record)"
            >{{ t('tool.edit') }}
          </a-button>
          <a-popconfirm
            :content="t('tool.deleteConfirm')"
            @ok="deleteToolHandler(record.name)"
          >
            <a-button status="danger">{{ t('tool.delete') }}</a-button>
          </a-popconfirm>
        </a-space>
      </template>
    </a-table>

    <a-modal
      v-model:visible="showEditModal"
      :title="
        isUpdating
          ? t('tool.editModalTitle.edit')
          : t('tool.editModalTitle.create')
      "
    >
      <a-form
        :model="editedTool"
        :rules="rules"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 18 }"
        layout="vertical"
      >
        <a-form-item :label="t('tool.form.name')" field="name">
          <a-input v-model="editedTool.name" />
        </a-form-item>
        <a-form-item
          :label="t('tool.form.description')"
          field="description"
        >
          <a-textarea v-model="editedTool.description" />
        </a-form-item>
        <a-form-item
          :label="t('tool.form.template')"
          field="template_name"
        >
          <a-select
            v-model="editedTool.template_name"
            :options="templateOptions"
            :placeholder="t('tool.form.selectTemplate')"
            @change="handleTemplateChange"
          />
        </a-form-item>
        <a-form-item :label="t('tool.form.params')" field="params">
          <a-space direction="vertical" style="width: 100%">
            <a-list :split="false" size="small" :bordered="false">
              <div class="mb-2 text-gray-400">
                <div class="">{{ t('tool.form.requiredParams') }}</div>
                <template
                  v-if="
                    paramTemplates[editedTool.template_name]?.length > 0
                  "
                >
                  <div
                    v-for="item in paramTemplates[
                      editedTool.template_name
                    ]"
                    :key="item.key"
                  >
                    <p class="text-sm"
                      ><span
                        class="font-bold px-1 py-0.5 bg-gray-200 rounded"
                        >{{ item.key }}</span
                      >: {{ item.description }}</p
                    >
                  </div>
                </template>
                <template v-else>
                  <div>{{ t('tool.form.noParams') }}</div>
                </template>
                <hr class="my-1" />
              </div>
              <div v-for="(param, index) in formParams" :key="index">
                <a-row :gutter="12">
                  <a-col :span="8">
                    <a-input
                      v-model="param.key"
                      :placeholder="t('tool.form.key')"
                    />
                  </a-col>
                  <a-col :span="14">
                    <a-textarea
                      v-model="param.value"
                      :placeholder="t('tool.form.value')"
                    />
                  </a-col>
                  <a-col :span="2">
                    <a-button type="text" @click="removeParam(index)">
                      <template #icon>
                        <icon-delete />
                      </template>
                    </a-button>
                  </a-col>
                </a-row>
              </div>
            </a-list>

            <a-button type="dashed" @click="addParam">{{
              t('tool.form.addParam')
            }}</a-button>
          </a-space>
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
          >{{ t('tool.form.cancel') }}
        </a-button>
        <a-button type="primary" @click="saveTool">{{
          t('tool.form.save')
        }}</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
  import XHeader from '@/components/header/x-header.vue';
  import { onBeforeMount, ref } from 'vue';
  import {
    Tool,
    createTool,
    deleteTool,
    listTools,
    updateTool,
  } from '@/api/tool';
  import { listToolTemplates } from '@/api/blueprint';
  import { namingValidator } from '@/utils/validator';
  import { useI18n } from 'vue-i18n';

  const { t } = useI18n();

  const isLoading = ref(false);
  const tools = ref<Tool[]>([]);
  const editedTool = ref<Tool>({
    name: '',
    description: '',
    template_name: '',
    params: {},
  });
  const formParams = ref<{ key: string; value: string }[]>([]);
  const showEditModal = ref(false);
  const isUpdating = ref(false);

  const columns = [
    { title: t('tool.form.name'), dataIndex: 'name' },
    { title: t('tool.form.description'), dataIndex: 'description' },
    { title: t('tool.form.template'), dataIndex: 'template_name' },
    { title: t('tool.form.params'), dataIndex: 'params' },
    { title: t('tool.edit'), slotName: 'actions' },
  ];
  const rules = {
    template_name: [
      {
        required: true,
        message: t('tool.form.rule.templateRequired'),
        trigger: 'blur',
      },
    ],
    name: [
      {
        required: true,
        message: t('tool.form.rule.nameRequired'),
        trigger: 'blur',
      },
      namingValidator,
    ],
  };
  const templateOptions = ref<{ label: string; value: string }[]>([]);
  const paramTemplates = ref<{
    [key: string]: { key: string; description: string; default: string }[];
  }>({});

  const fetchTemplates = async () => {
    const response = await listToolTemplates();
    templateOptions.value = response.data.map((template) => ({
      label: template.name,
      value: template.name,
    }));
    response.data.forEach((template) => {
      paramTemplates.value[template.name] = template.param_template_define;
    });
  };

  const handleTemplateChange = (templateName: any) => {
    const params = paramTemplates.value[templateName as string] || [];
    formParams.value = params.map((param) => ({
      key: param.key,
      value: editedTool.value.params[param.key] || param.default,
    }));
  };

  onBeforeMount(() => {
    listAllTools();
    fetchTemplates();
  });

  const editBtnHandler = (tool: Tool) => {
    editedTool.value = { ...tool };
    formParams.value = Object.entries(editedTool.value.params).map(
      ([key, value]) => ({ key, value }),
    );
    showEditModal.value = true;
    isUpdating.value = true;
  };

  const deleteToolHandler = async (name: string) => {
    await deleteTool(name);
    await listAllTools();
  };

  async function listAllTools() {
    isLoading.value = true;
    tools.value = (await listTools()).data;
    isLoading.value = false;
  }

  const addParam = () => {
    formParams.value.push({ key: '', value: '' });
  };

  const removeParam = (index: number) => {
    formParams.value.splice(index, 1);
  };

  const saveTool = async () => {
    // Convert formParams to map
    const paramsMap: Record<string, string> = {};
    formParams.value.forEach((param) => {
      if (param.key && param.value) {
        paramsMap[param.key] = param.value;
      }
    });
    editedTool.value.params = paramsMap;

    if (isUpdating.value) {
      await updateTool(editedTool.value.name, editedTool.value);
    } else {
      await createTool(editedTool.value);
    }
    showEditModal.value = false;
    await listAllTools();
    isUpdating.value = false;
    editedTool.value = {
      name: '',
      description: '',
      template_name: '',
      params: {},
    };
    formParams.value = [];
  };
</script>

<script lang="ts">
  export default {
    name: 'ToolManage',
  };
</script>
