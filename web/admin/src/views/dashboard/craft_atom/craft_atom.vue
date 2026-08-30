<template>
  <CraftManagePage
    :title="t('menu.craftAtom')"
    :description="t('craftAtom.description')"
  >
    <template #toolbar>
      <a-space wrap>
        <a-button :loading="isLoading" @click="listAllCraftAtoms">
          <template #icon>
            <icon-refresh />
          </template>
          {{ t('craftAtom.query') }}
        </a-button>
        <a-button type="primary" @click="handleAdd">
          <template #icon>
            <icon-plus />
          </template>
          {{ t('craftAtom.create') }}
        </a-button>
      </a-space>
    </template>

    <a-table
      v-if="isLoading || craftAtoms.length > 0"
      row-key="name"
      :data="craftAtoms"
      :columns="columns"
      :loading="isLoading"
      :bordered="false"
      :pagination="{ pageSize: 10, showTotal: true }"
    >
      <template #params="{ record }">
        <a-space v-if="Object.keys(record.params || {}).length" wrap>
          <a-tag v-for="(_, key) in record.params" :key="key" color="gray">
            {{ key }}
          </a-tag>
        </a-space>
        <span v-else class="text-gray-400">{{
          t('craftAtom.form.noParams')
        }}</span>
      </template>
      <template #actions="{ record }">
        <a-space wrap>
          <a-button type="text" size="small" @click="editBtnHandler(record)">
            {{ t('craftAtom.edit') }}
          </a-button>
          <a-popconfirm
            :content="t('craftAtom.deleteConfirm')"
            @ok="deleteCraftAtomHandler(record.name)"
          >
            <a-button type="text" status="danger" size="small">
              {{ t('craftAtom.delete') }}
            </a-button>
          </a-popconfirm>
        </a-space>
      </template>
    </a-table>

    <ListEmptyGuide
      v-else-if="!listFailed"
      :description="t('craftAtom.empty.description')"
      :hint="t('craftAtom.empty.hint')"
      :create-label="t('craftAtom.empty.createFirst')"
      :docs-label="t('craftAtom.empty.docs')"
      :docs-href="atomDocsHref"
      @create="handleAdd"
    />

    <a-modal
      v-model:visible="showEditModal"
      :title="
        isUpdating
          ? t('craftAtom.editModalTitle.edit')
          : t('craftAtom.editModalTitle.create')
      "
    >
      <a-form
        ref="formRef"
        :model="editedCraftAtom"
        :rules="rules"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 18 }"
        layout="vertical"
      >
        <a-form-item :label="t('craftAtom.form.name')" field="name">
          <a-input v-model="editedCraftAtom.name" :disabled="isUpdating" />
        </a-form-item>
        <a-form-item
          :label="t('craftAtom.form.description')"
          field="description"
        >
          <a-textarea v-model="editedCraftAtom.description" />
        </a-form-item>
        <a-form-item
          :label="t('craftAtom.form.template')"
          field="template_name"
        >
          <a-select
            v-model="editedCraftAtom.template_name"
            :options="templateOptions"
            :placeholder="t('craftAtom.form.selectTemplate')"
            @change="handleTemplateChange"
          />
        </a-form-item>
        <a-form-item :label="t('craftAtom.form.params')" field="params">
          <a-space direction="vertical" style="width: 100%">
            <a-list :split="false" size="small" :bordered="false">
              <div class="mb-2 text-gray-400">
                <div class="">{{ t('craftAtom.form.requiredParams') }}</div>
                <template
                  v-if="
                    paramTemplates[editedCraftAtom.template_name]?.length > 0
                  "
                >
                  <div
                    v-for="item in paramTemplates[
                      editedCraftAtom.template_name
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
                  <div>{{ t('craftAtom.form.noParams') }}</div>
                </template>
                <hr class="my-1" />
              </div>
              <div v-for="(param, index) in formParams" :key="index">
                <a-row :gutter="12">
                  <a-col :span="8">
                    <a-input
                      v-model="param.key"
                      :placeholder="t('craftAtom.form.key')"
                    />
                  </a-col>
                  <a-col :span="14">
                    <a-select
                      v-if="
                        isAIFilterExtraPayloadParam(
                          editedCraftAtom.template_name,
                          param.key
                        )
                      "
                      v-model="param.value"
                      multiple
                      allow-clear
                      :options="aiFilterExtraPayloadOptions"
                      :placeholder="t('craftAtom.form.value')"
                    />
                    <a-select
                      v-else-if="
                        isAIContentProcessPlacementParam(
                          editedCraftAtom.template_name,
                          param.key
                        )
                      "
                      v-model="param.value"
                      :options="aiContentProcessPlacementOptions"
                      :placeholder="t('craftAtom.form.value')"
                    />
                    <a-select
                      v-else-if="
                        isEmbeddingFilterModeParam(
                          editedCraftAtom.template_name,
                          param.key
                        )
                      "
                      v-model="param.value"
                      :options="embeddingFilterModeOptions"
                      :placeholder="t('craftAtom.form.value')"
                    />
                    <a-input-number
                      v-else-if="
                        isEmbeddingFilterThresholdParam(
                          editedCraftAtom.template_name,
                          param.key
                        )
                      "
                      v-model="param.value"
                      :min="0"
                      :max="1"
                      :step="0.05"
                      :precision="2"
                      :placeholder="t('craftAtom.form.value')"
                    />
                    <a-input-number
                      v-else-if="
                        isEmbeddingFilterMaxContentLengthParam(
                          editedCraftAtom.template_name,
                          param.key
                        )
                      "
                      v-model="param.value"
                      :min="1"
                      :step="100"
                      :placeholder="t('craftAtom.form.value')"
                    />
                    <a-textarea
                      v-else
                      :model-value="String(param.value)"
                      :placeholder="t('craftAtom.form.value')"
                      :auto-size="
                        isEmbeddingFilterAnchorsParam(
                          editedCraftAtom.template_name,
                          param.key
                        )
                          ? { minRows: 4, maxRows: 8 }
                          : { minRows: 2, maxRows: 4 }
                      "
                      @update:model-value="param.value = $event"
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
              t('craftAtom.form.addParam')
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
          >{{ t('craftAtom.form.cancel') }}
        </a-button>
        <a-button type="primary" @click="saveCraftAtom">{{
          t('craftAtom.form.save')
        }}</a-button>
      </template>
    </a-modal>
  </CraftManagePage>
</template>

<script setup lang="ts">
  import CraftManagePage from '@/components/craft/CraftManagePage.vue';
  import ListEmptyGuide from '@/components/list-empty-guide/index.vue';
  import { computed, onBeforeMount, ref } from 'vue';
  import { buildDocsUrl } from '@/utils/docsUrl';
  import {
    CraftAtom,
    createCraftAtom,
    deleteCraftAtom,
    listCraftAtoms,
    updateCraftAtom,
  } from '@/api/craft_atom';
  import { listCraftTemplates } from '@/api/craft_flow';
  import { namingValidator } from '@/utils/validator';
  import { useI18n } from 'vue-i18n';
  import { Message } from '@arco-design/web-vue';
  import {
    aiContentProcessPlacementOptions,
    aiFilterExtraPayloadOptions,
    CraftParamValue,
    embeddingFilterModeOptions,
    isAIContentProcessPlacementParam,
    isAIFilterExtraPayloadParam,
    isEmbeddingFilterAnchorsParam,
    isEmbeddingFilterMaxContentLengthParam,
    isEmbeddingFilterModeParam,
    isEmbeddingFilterThresholdParam,
    serializeCraftParamValue,
    toCraftParamFormValue,
  } from '@/views/dashboard/craft_atom/paramOptions';

  const { t, locale } = useI18n();
  const atomDocsHref = computed(() =>
    buildDocsUrl(locale.value, 'guides/customization')
  );

  const isLoading = ref(false);
  const listFailed = ref(false);
  const formRef = ref();
  const craftAtoms = ref<CraftAtom[]>([]);
  const editedCraftAtom = ref<CraftAtom>({
    name: '',
    description: '',
    template_name: '',
    params: {},
  });
  const formParams = ref<{ key: string; value: CraftParamValue }[]>([]);
  const showEditModal = ref(false);
  const isUpdating = ref(false);
  const originalName = ref('');

  const columns = [
    { title: t('craftAtom.form.name'), dataIndex: 'name' },
    { title: t('craftAtom.form.description'), dataIndex: 'description' },
    { title: t('craftAtom.form.template'), dataIndex: 'template_name' },
    { title: t('craftAtom.form.params'), slotName: 'params' },
    {
      title: t('craftAtom.edit'),
      slotName: 'actions',
      width: 140,
      align: 'right',
    },
  ];
  const rules = {
    template_name: [
      {
        required: true,
        message: t('craftAtom.form.rule.templateRequired'),
        trigger: 'blur',
      },
    ],
    name: [
      {
        required: true,
        message: t('craftAtom.form.rule.nameRequired'),
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
    const response = await listCraftTemplates();
    const sortedTemplates = [...response.data].sort((a, b) =>
      a.name.localeCompare(b.name, undefined, { sensitivity: 'base' })
    );
    templateOptions.value = sortedTemplates.map((template) => ({
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
      value: toCraftParamFormValue(
        templateName as string,
        param.key,
        editedCraftAtom.value.params[param.key] || param.default
      ),
    }));
  };

  const handleAdd = () => {
    editedCraftAtom.value = {
      name: '',
      description: '',
      template_name: '',
      params: {},
    };
    formParams.value = [];
    showEditModal.value = true;
    isUpdating.value = false;
  };

  onBeforeMount(() => {
    listAllCraftAtoms();
    fetchTemplates();
  });

  const editBtnHandler = (craftAtom: CraftAtom) => {
    editedCraftAtom.value = { ...craftAtom };
    originalName.value = craftAtom.name;
    formParams.value = Object.entries(editedCraftAtom.value.params).map(
      ([key, value]) => ({
        key,
        value: toCraftParamFormValue(craftAtom.template_name, key, value),
      })
    );
    showEditModal.value = true;
    isUpdating.value = true;
  };

  const deleteCraftAtomHandler = async (name: string) => {
    await deleteCraftAtom(name);
    await listAllCraftAtoms();
  };

  async function listAllCraftAtoms() {
    isLoading.value = true;
    try {
      craftAtoms.value = (await listCraftAtoms()).data ?? [];
      listFailed.value = false;
    } catch {
      listFailed.value = true;
    } finally {
      isLoading.value = false;
    }
  }

  const addParam = () => {
    formParams.value.push({ key: '', value: '' });
  };

  const removeParam = (index: number) => {
    formParams.value.splice(index, 1);
  };

  const validateEmbeddingFilterParams = (paramsMap: Record<string, string>) => {
    if (editedCraftAtom.value.template_name !== 'embedding-filter') {
      return true;
    }
    if (!paramsMap.anchors?.trim()) {
      Message.error('Embedding filter requires at least one anchor.');
      return false;
    }
    if (paramsMap.mode) {
      paramsMap.mode = paramsMap.mode.trim().toLowerCase();
    }
    if (paramsMap.mode && !['include', 'exclude'].includes(paramsMap.mode)) {
      Message.error('Embedding filter mode must be include or exclude.');
      return false;
    }
    return true;
  };

  const saveCraftAtom = async () => {
    const res = await formRef.value?.validate();
    if (res) return;

    // Convert formParams to map
    const paramsMap: Record<string, string> = {};
    formParams.value.forEach((param) => {
      const value = serializeCraftParamValue(param.value);
      if (param.key && value !== '') {
        paramsMap[param.key] = value;
      }
    });
    if (!validateEmbeddingFilterParams(paramsMap)) {
      return;
    }
    editedCraftAtom.value.params = paramsMap;

    if (isUpdating.value) {
      await updateCraftAtom(originalName.value, editedCraftAtom.value);
    } else {
      await createCraftAtom(editedCraftAtom.value);
    }
    Message.success(t('craftAtom.form.saveSuccess'));
    showEditModal.value = false;
    await listAllCraftAtoms();
    isUpdating.value = false;
    editedCraftAtom.value = {
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
    name: 'CraftAtomManage',
  };
</script>
