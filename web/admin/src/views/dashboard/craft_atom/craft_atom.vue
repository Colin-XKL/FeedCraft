<template>
  <div class="py-8 px-16">
    <x-header title="Craft Atom Management" description="Manage Craft Atoms">
    </x-header>

    <a-space direction="horizontal" class="mb-6">
      <a-button type="primary" :loading="isLoading" @click="listAllCraftAtoms">
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
        >Create CraftAtom
      </a-button>
    </a-space>

    <a-table :data="craftAtoms" :columns="columns" :loading="isLoading">
      <template #actions="{ record }">
        <a-space>
          <a-button type="outline" @click="editBtnHandler(record)"
            >Edit
          </a-button>
          <a-button status="danger" @click="deleteCraftAtomHandler(record.name)"
            >Delete
          </a-button>
        </a-space>
      </template>
    </a-table>

    <a-modal
      v-model:visible="showEditModal"
      :title="isUpdating ? 'Edit Craft Atom' : 'Create Craft Atom'"
    >
      <a-form
        :model="editedCraftAtom"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 18 }"
      >
        <a-form-item label="Name" name="name">
          <a-input v-model="editedCraftAtom.name" />
        </a-form-item>
        <a-form-item label="Description" name="description">
          <a-textarea v-model="editedCraftAtom.description" />
        </a-form-item>
        <a-form-item label="Template Name" name="template_name">
          <a-select
            v-model="editedCraftAtom.template_name"
            :options="templateOptions"
            placeholder="Select Template"
            @change="handleTemplateChange"
          />
        </a-form-item>
        <a-form-item label="Params" name="params">
          <a-space direction="vertical" style="width: 100%">
            <div v-for="(param, index) in formParams" :key="index">
              <a-row :gutter="16">
                <a-col :span="11">
                  <a-input v-model="param.key" placeholder="Key" disabled />
                </a-col>
                <a-col :span="11">
                  <a-input v-model="param.value" placeholder="Value" />
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
            <a-button type="dashed" @click="addParam">Add Param</a-button>
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
          >Cancel
        </a-button>
        <a-button type="primary" @click="saveCraftAtom">Save</a-button>
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
    { title: 'Name', dataIndex: 'name' },
    { title: 'Description', dataIndex: 'description' },
    { title: 'Template Name', dataIndex: 'template_name' },
    { title: 'Params', dataIndex: 'params' },
    { title: 'Actions', slotName: 'actions' },
  ];

  const templateOptions = ref<{ label: string; value: string }[]>([]);
  const paramTemplates = ref<{ [key: string]: { key: string; description: string; default: string }[] }>({});

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

  onBeforeMount(() => {
    listAllCraftAtoms();
    fetchTemplates();
  });

  const editBtnHandler = (craftAtom: CraftAtom) => {
    editedCraftAtom.value = { ...craftAtom };
    formParams.value = Object.entries(editedCraftAtom.value.params).map(
      ([key, value]) => ({ key, value })
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
