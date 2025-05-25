<template>
  <div class="py-8 px-16">
    <x-header title="自定义配方" description="自定义RSS，以及要使用的工艺">
    </x-header>

    <a-space direction="horizontal" class="mb-4">
      <a-button type="primary" :loading="isLoading" @click="listCustomRecipes">
        查询
      </a-button>
      <a-button
        type="outline"
        @click="
          () => {
            resetForm();
            showModal = true;
          }
        "
      >
        创建自定义Recipe
      </a-button>
    </a-space>

    <a-table
      :data="recipes"
      :columns="columns"
      :bordered="true"
      :loading="isLoading"
    >
      <template #status="{ record }">
        <a-tooltip
          v-if="record.is_active"
          :content="`该 recipe 处于活跃状态, 会自动定期获取最新内容. 最近一次用户请求时间 ${dayjs(
            record.last_accessed_at
          ).format('YYYY-MM-DD HH:mm:ss')}`"
        >
          <a-tag color="green" :default-checked="true">活跃</a-tag>
        </a-tooltip>
        <a-tooltip
          v-else
          content="该 recipe 近3天没有请求,已经进入休眠状态,不会自动定期获取最新内容"
        >
          <a-tag color="gray" :default-checked="true">不活跃</a-tag>
        </a-tooltip>
      </template>
      <template #actions="{ record }">
        <a-space direction="horizontal">
          <a-button
            type="outline"
            @click="
              () => {
                isUpdating = true;
                showEditModal(record);
              }
            "
          >
            编辑
          </a-button>
          <a-button @click="deleteRecipe(record.id)">删除</a-button>
          <a-link :href="`${baseUrl}/recipe/${record?.id}`">链接</a-link>
        </a-space>
      </template>
    </a-table>

    <a-modal
      v-model:visible="showModal"
      :title="editing ? '编辑自定义配方' : '创建自定义配方'"
    >
      <a-form
        :model="form"
        :label-col="{ span: 6 }"
        :rules="rules"
        :wrapper-col="{ span: 18 }"
      >
        <a-form-item label="名称" field="id">
          <a-input v-model="form.id" :disabled="isUpdating" />
        </a-form-item>
        <a-form-item label="描述" field="description">
          <a-input v-model="form.description" />
        </a-form-item>
        <a-form-item label="工艺" field="craft">
          <a-input v-model="form.craft" />
        </a-form-item>
        <a-form-item label="FeedURL" field="feed_url">
          <a-input v-model="form.feed_url" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-button
          @click="
            () => {
              showModal = false;
              isUpdating = false;
            }
          "
          >取消
        </a-button>
        <a-button type="primary" @click="saveRecipe">保存</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue';
  import {
    createCustomRecipe,
    CustomRecipe,
    deleteCustomRecipe,
    getCustomRecipes,
    updateCustomRecipe,
  } from '@/api/custom_recipe';
  import XHeader from '@/components/header/x-header.vue';
  import { namingValidator } from '@/utils/validator';
  import dayjs from 'dayjs';

  const baseUrl = import.meta.env.VITE_API_BASE_URL ?? '';

  const recipes = ref<CustomRecipe[]>([]);
  const showModal = ref(false);
  const form = ref<CustomRecipe>({
    id: '',
    description: '',
    craft: '',
    feed_url: '',
  });
  const editing = ref(false);
  const selectedRecipe = ref<CustomRecipe | null>(null);
  const isLoading = ref(false);
  const isUpdating = ref(false);

  const columns = [
    { title: '名称', dataIndex: 'id' },
    { title: '描述', dataIndex: 'description' },
    { title: '使用的Craft', dataIndex: 'craft' },
    { title: '状态', slotName: 'status' },
    { title: 'URL', dataIndex: 'feed_url' },
    { title: '操作', slotName: 'actions' },
  ];

  async function listCustomRecipes() {
    isLoading.value = true;
    recipes.value = (await getCustomRecipes()).data;
    isLoading.value = false;
  }

  onMounted(() => {
    listCustomRecipes();
  });
  const rules = {
    id: [
      {
        required: true,
        message: 'Name is required',
        trigger: 'blur',
      },
      namingValidator,
    ],
    craft: [
      {
        required: true,
        message: 'Name is required',
        trigger: 'blur',
      },
    ],
    feed_url: [
      {
        required: true,
        message: 'Name is required',
        trigger: 'blur',
      },
    ],
  };
  const showEditModal = (recipe: CustomRecipe) => {
    editing.value = true;
    selectedRecipe.value = recipe;
    form.value = {
      id: recipe.id,
      description: recipe.description,
      craft: recipe.craft,
      feed_url: recipe.feed_url,
    };
    showModal.value = true;
  };

  const saveRecipe = async () => {
    if (editing.value) {
      if (selectedRecipe.value) {
        await updateCustomRecipe(form.value);
        selectedRecipe.value.description = form.value.description;
        selectedRecipe.value.craft = form.value.craft;
        selectedRecipe.value.feed_url = form.value.feed_url;
      }
    } else {
      await createCustomRecipe(form.value as CustomRecipe);
      await listCustomRecipes();
    }
    showModal.value = false;
    form.value = { id: '', description: '', craft: '', feed_url: '' };
    editing.value = false;
    isUpdating.value = false;
    selectedRecipe.value = null;
  };

  const deleteRecipe = async (id: string) => {
    await deleteCustomRecipe(id);
    await listCustomRecipes();
  };

  function resetForm() {
    form.value = {
      id: '',
      description: '',
      craft: '',
      feed_url: '',
    };
  }
</script>

<style scoped></style>
