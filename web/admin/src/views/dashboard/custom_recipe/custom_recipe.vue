<template>
  <div class="py-8 px-16">
    <x-header title="Custom Recipe" description="自定义rss, 以及要使用的craft">
    </x-header>

    <a-space direction="horizontal" class="mb-4">
      <a-button type="primary" :loading="isLoading" @click="listCustomRecipes">
        List
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
        Create Custom Recipe
      </a-button>
    </a-space>

    <a-table
      :data="recipes"
      :columns="columns"
      :bordered="true"
      :loading="isLoading"
    >
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
            Edit
          </a-button>
          <a-button @click="deleteRecipe(record.id)">Delete</a-button>
          <a-link :href="`${baseUrl}/recipe/${record?.id}`">Link</a-link>
        </a-space>
      </template>
    </a-table>

    <a-modal
      v-model:visible="showModal"
      :title="editing ? 'Edit Custom Recipe' : 'Create Custom Recipe'"
    >
      <a-form
        :model="form"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 18 }"
      >
        <a-form-item label="RecipeName">
          <a-input v-model="form.id" :disabled="isUpdating" />
        </a-form-item>
        <a-form-item label="Description">
          <a-input v-model="form.description" />
        </a-form-item>
        <a-form-item label="Craft">
          <a-input v-model="form.craft" />
        </a-form-item>
        <a-form-item label="FeedURL">
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
          >Cancel
        </a-button>
        <a-button type="primary" @click="saveRecipe">Save</a-button>
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

  const baseUrl = import.meta.env.VITE_API_BASE_URL ?? '/';

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
    { title: 'Name', dataIndex: 'id' },
    { title: 'Description', dataIndex: 'description' },
    { title: 'Craft', dataIndex: 'craft' },
    { title: 'Feed URL', dataIndex: 'feed_url' },
    { title: 'Actions', slotName: 'actions' },
  ];

  async function listCustomRecipes() {
    isLoading.value = true;
    recipes.value = (await getCustomRecipes()).data;
    isLoading.value = false;
  }

  onMounted(() => {
    listCustomRecipes();
  });

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
