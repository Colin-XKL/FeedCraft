<template>
  <div class="py-8 px-16">
    <x-header
      :title="t('menu.craftAtom')"
      :description="t('craftAtom.description')"
    ></x-header>

    <a-space direction="horizontal" class="mb-6">
      <a-button type="primary" :loading="isLoading" @click="listAllCraftAtoms">
        {{ t('craftAtom.query') }}
      </a-button>
      <a-button
        type="outline"
        @click="
          () => {
            showEditModal = true;
            isUpdating = false;
          }
        "
        >{{ t('craftAtom.create') }}
      </a-button>
    </a-space>

    <a-table :data="craftAtoms" :columns="columns" :loading="isLoading">
      <template #actions="{ record }">
        <a-space>
          <a-button type="outline" @click="editBtnHandler(record)"
            >{{ t('craftAtom.edit') }}
          </a-button>
          <a-popconfirm
            :content="t('craftAtom.deleteConfirm')"
            @ok="deleteCraftAtomHandler(record.name)"
          >
            <a-button status="danger">{{ t('craftAtom.delete') }}</a-button>
          </a-popconfirm>
        </a-space>
      </template>
    </a-table>

    <a-modal
      v-model:visible="showEditModal"
      :title="
        isUpdating
          ? t('craftAtom.editModalTitle.edit')
          : t('craftAtom.editModalTitle.create')
      "
    >
      <a-form
        :model="editedCraftAtom"
        :rules="rules"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 18 }"
        layout="vertical"
      >
        <a-form-item :label="t('craftAtom.form.name')" field="name">
          <a-input v-model="editedCraftAtom.name" />
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
                    <a-textarea
                      v-model="param.value"
                      :placeholder="t('craftAtom.form.value')"
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
  </div>
</template>

<script setup lang="ts">
  import XHeader from '@/components/header/x-header.vue';
  import { onBeforeMount, ref } from 'vue';
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

  const { t } = useI18n();

  const isLoading = ref(false);
  const craftAtoms = ref<CraftAtom[]>([]);
  const editedCraftAtom = ref<CraftAtom>({
    name: '',
    description: '',
    template_name: '',
    params: {},
  });
  const formParams = ref<{ key: string; value: string }[]>([]);
  const showEditModal = ref(false);
  const isUpdating = ref(false);

  const columns = [
    { title: t('craftAtom.form.name'), dataIndex: 'name' },
    { title: t('craftAtom.form.description'), dataIndex: 'description' },
    { title: t('craftAtom.form.template'), dataIndex: 'template_name' },
    { title: t('craftAtom.form.params'), dataIndex: 'params' },
    { title: t('craftAtom.edit'), slotName: 'actions' },
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
      value: editedCraftAtom.value.params[param.key] || param.default,
    }));
  };

  onBeforeMount(() => {
    listAllCraftAtoms();
    fetchTemplates();
  });

  const editBtnHandler = (craftAtom: CraftAtom) => {
    editedCraftAtom.value = { ...craftAtom };
    formParams.value = Object.entries(editedCraftAtom.value.params).map(
      ([key, value]) => ({ key, value }),
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
    craftAtoms.value = (await listCraftAtoms()).data;
    isLoading.value = false;
  }

  const addParam = () => {
    formParams.value.push({ key: '', value: '' });
  };

  const removeParam = (index: number) => {
    formParams.value.splice(index, 1);
  };

  const saveCraftAtom = async () => {
    // Convert formParams to map
    const paramsMap: Record<string, string> = {};
    formParams.value.forEach((param) => {
      if (param.key && param.value) {
        paramsMap[param.key] = param.value;
      }
    });
    editedCraftAtom.value.params = paramsMap;

    if (isUpdating.value) {
      await updateCraftAtom(editedCraftAtom.value.name, editedCraftAtom.value);
    } else {
      await createCraftAtom(editedCraftAtom.value);
    }
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
