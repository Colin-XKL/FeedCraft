<template>
  <div class="py-8 px-16">
    <x-header
      :title="t('menu.customRecipe')"
      :description="t('customRecipe.description')"
    >
    </x-header>

    <a-space direction="horizontal" class="mb-4">
      <a-button type="primary" :loading="isLoading" @click="listCustomRecipes">
        {{ t('customRecipe.query') }}
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
        {{ t('customRecipe.create') }}
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
          :content="
            t('customRecipe.status.activeTooltip', {
              time: dayjs(record.last_accessed_at).format('YYYY-MM-DD HH:mm:ss'),
            })
          "
        >
          <a-tag color="green" :default-checked="true">{{ t('customRecipe.status.active') }}</a-tag>
        </a-tooltip>
        <a-tooltip
          v-else
          :content="t('customRecipe.status.inactiveTooltip')"
        >
          <a-tag color="gray" :default-checked="true">{{ t('customRecipe.status.inactive') }}</a-tag>
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
            {{ t('customRecipe.edit') }}
          </a-button>
          <a-button @click="deleteRecipe(record.id)">{{ t('customRecipe.delete') }}</a-button>
          <a-link :href="`${baseUrl}/recipe/${record?.id}`">{{ t('customRecipe.link') }}</a-link>
        </a-space>
      </template>
    </a-table>

    <a-modal
      v-model:visible="showModal"
      :title="
        editing
          ? t('customRecipe.editModalTitle.edit')
          : t('customRecipe.editModalTitle.create')
      "
    >
      <a-form
        :model="form"
        :label-col="{ span: 6 }"
        :rules="rules"
        :wrapper-col="{ span: 18 }"
      >
        <a-form-item :label="t('customRecipe.form.name')" field="id">
          <a-input v-model="form.id" :disabled="isUpdating" />
        </a-form-item>
        <a-form-item :label="t('customRecipe.form.description')" field="description">
          <a-input v-model="form.description" />
        </a-form-item>
        <a-form-item :label="t('customRecipe.form.craft')" field="craft">
          <a-input v-model="form.craft" />
        </a-form-item>
        <a-form-item :label="t('customRecipe.form.feedURL')" field="feed_url">
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
          >{{ t('customRecipe.form.cancel') }}
        </a-button>
        <a-button type="primary" @click="saveRecipe">{{ t('customRecipe.form.save') }}</a-button>
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
  import { useI18n } from 'vue-i18n';

  const { t } = useI18n();

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
    { title: t('customRecipe.form.name'), dataIndex: 'id' },
    { title: t('customRecipe.form.description'), dataIndex: 'description' },
    { title: t('customRecipe.form.craft'), dataIndex: 'craft' },
    { title: t('customRecipe.status.active'), slotName: 'status' },
    { title: t('customRecipe.form.feedURL'), dataIndex: 'feed_url' },
    { title: t('customRecipe.edit'), slotName: 'actions' },
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
        message: t('customRecipe.form.rule.nameRequired'),
        trigger: 'blur',
      },
      namingValidator,
    ],
    craft: [
      {
        required: true,
        message: t('customRecipe.form.rule.craftRequired'),
        trigger: 'blur',
      },
    ],
    feed_url: [
      {
        required: true,
        message: t('customRecipe.form.rule.feedURLRequired'),
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
