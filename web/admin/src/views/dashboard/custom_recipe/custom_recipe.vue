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
            quickCreate = true;
            form.source_type = 'rss';
            showModal = true;
          }
        "
      >
        <template #icon><icon-plus /></template>
        {{ t('customRecipe.quickCreateRSS') }}
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
              time: dayjs(record.last_accessed_at).format(
                'YYYY-MM-DD HH:mm:ss',
              ),
            })
          "
        >
          <a-tag color="green" :default-checked="true">{{
            t('customRecipe.status.active')
          }}</a-tag>
        </a-tooltip>
        <a-tooltip v-else :content="t('customRecipe.status.inactiveTooltip')">
          <a-tag color="gray" :default-checked="true">{{
            t('customRecipe.status.inactive')
          }}</a-tag>
        </a-tooltip>
      </template>
      <template #source_config="{ record }">
        <a-space>
          <span
            style="
              max-width: 300px;
              display: inline-block;
              overflow: hidden;
              text-overflow: ellipsis;
              white-space: nowrap;
            "
            :title="humanReadableConfig(record.source_config)"
          >
            {{ humanReadableConfig(record.source_config) }}
          </span>
          <a-tooltip :content="t('customRecipe.viewConfig')">
            <a-button
              type="text"
              size="mini"
              :aria-label="t('customRecipe.viewConfig')"
              @click="viewConfig(record.source_config)"
            >
              <template #icon>
                <icon-eye />
              </template>
            </a-button>
          </a-tooltip>
        </a-space>
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
          <a-popconfirm
            :content="t('customRecipe.deleteConfirm')"
            @ok="deleteRecipe(record.id)"
          >
            <a-button status="danger">{{ t('customRecipe.delete') }}</a-button>
          </a-popconfirm>
          <a-link :href="`${baseUrl}/recipe/${record?.id}`">{{
            t('customRecipe.link')
          }}</a-link>
        </a-space>
      </template>
    </a-table>

    <!-- Create/Edit Modal -->
    <a-modal
      v-model:visible="showModal"
      :title="
        editing
          ? t('customRecipe.editModalTitle.edit')
          : quickCreate
            ? t('customRecipe.quickCreateRSS')
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
        <a-form-item
          :label="t('customRecipe.form.description')"
          field="description"
        >
          <a-input v-model="form.description" />
        </a-form-item>
        <a-form-item :label="t('customRecipe.form.craft')" field="craft">
          <CraftSelector
            v-model="craftList"
            mode="multiple"
            :placeholder="t('customRecipe.form.placeholder.craft')"
          />
        </a-form-item>

        <!-- Quick Create Fields -->
        <template v-if="quickCreate">
          <a-form-item
            :label="t('customRecipe.form.feedURL')"
            field="feed_url"
            :rules="[
              {
                required: true,
                message: t('customRecipe.form.rule.rssUrlRequired'),
              },
            ]"
          >
            <a-input
              v-model="rssUrl"
              :placeholder="t('customRecipe.form.placeholder.rssUrl')"
            />
          </a-form-item>
        </template>

        <!-- Advanced Fields -->
        <template v-else>
          <a-form-item
            :label="t('customRecipe.form.sourceType')"
            field="source_type"
          >
            <a-select v-model="form.source_type">
              <a-option value="rss">RSS</a-option>
              <a-option value="html">HTML</a-option>
              <a-option value="json">JSON</a-option>
            </a-select>
          </a-form-item>
          <a-form-item
            :label="t('customRecipe.form.sourceConfig')"
            field="source_config"
          >
            <a-textarea
              v-model="form.source_config"
              :auto-size="{ minRows: 3, maxRows: 10 }"
              :placeholder="t('customRecipe.form.placeholder.sourceConfig')"
            />
          </a-form-item>
        </template>
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
        <a-button type="primary" @click="saveRecipe">{{
          t('customRecipe.form.save')
        }}</a-button>
      </template>
    </a-modal>

    <!-- View Config Modal -->
    <a-modal
      v-model:visible="showConfigModal"
      :title="t('customRecipe.viewConfigModalTitle')"
      :footer="false"
    >
      <pre
        style="
          background-color: #f5f5f5;
          padding: 10px;
          border-radius: 4px;
          overflow: auto;
          max-height: 400px;
        "
        >{{ currentConfig }}</pre
      >
    </a-modal>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted, computed } from 'vue';
  import {
    createCustomRecipe,
    CustomRecipe,
    deleteCustomRecipe,
    getCustomRecipes,
    updateCustomRecipe,
  } from '@/api/custom_recipe';
  import XHeader from '@/components/header/x-header.vue';
  import { namingValidator } from '@/utils/validator';
  import { IconEye, IconPlus } from '@arco-design/web-vue/es/icon';
  import { Message } from '@arco-design/web-vue';
  import dayjs from 'dayjs';
  import { useI18n } from 'vue-i18n';
  import CraftSelector from '../craft_flow/CraftSelector.vue';

  const { t } = useI18n();

  const baseUrl = import.meta.env.VITE_API_BASE_URL ?? '';

  const recipes = ref<CustomRecipe[]>([]);
  const showModal = ref(false);
  const showConfigModal = ref(false);
  const currentConfig = ref('');
  const quickCreate = ref(false);
  const rssUrl = ref('');

  const form = ref<CustomRecipe>({
    id: '',
    description: '',
    craft: '',
    source_type: 'rss',
    source_config: '',
  });

  const craftList = computed({
    get: () =>
      form.value.craft ? form.value.craft.split(',').filter(Boolean) : [],
    set: (val: string[] | string) => {
      if (Array.isArray(val)) {
        form.value.craft = val.join(',');
      } else if (typeof val === 'string') {
        form.value.craft = val;
      }
    },
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
    { title: t('customRecipe.form.sourceType'), dataIndex: 'source_type' },
    {
      title: t('customRecipe.form.sourceConfig'),
      dataIndex: 'source_config',
      slotName: 'source_config',
    },
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
    source_type: [
      {
        required: true,
        message: t('customRecipe.form.rule.sourceTypeRequired'),
        trigger: 'change',
      },
    ],
    source_config: [
      {
        required: true,
        message: t('customRecipe.form.rule.sourceConfigRequired'),
        trigger: 'blur',
      },
    ],
  };

  const humanReadableConfig = (configStr: string) => {
    try {
      const config = JSON.parse(configStr);
      // Try to find the URL in common locations
      if (config.http_fetcher && config.http_fetcher.url) {
        return config.http_fetcher.url;
      }
      if (config.url) {
        return config.url;
      }
      return t('customRecipe.jsonConfigFallback');
    } catch (e) {
      return configStr;
    }
  };

  const viewConfig = (configStr: string) => {
    try {
      const obj = JSON.parse(configStr);
      currentConfig.value = JSON.stringify(obj, null, 2);
    } catch (e) {
      currentConfig.value = configStr;
    }
    showConfigModal.value = true;
  };

  const showEditModal = (recipe: CustomRecipe) => {
    editing.value = true;
    quickCreate.value = false; // Ensure we are not in quick create mode
    selectedRecipe.value = recipe;

    // Pretty print JSON for editing
    let prettyConfig = recipe.source_config;
    try {
      const obj = JSON.parse(recipe.source_config);
      prettyConfig = JSON.stringify(obj, null, 2);
    } catch (e) {
      // ignore error, keep original string
    }

    form.value = {
      id: recipe.id,
      description: recipe.description,
      craft: recipe.craft,
      source_type: recipe.source_type,
      source_config: prettyConfig,
    };
    showModal.value = true;
  };

  const saveRecipe = async () => {
    if (quickCreate.value) {
      // Construct JSON for Quick Create
      const config = {
        http_fetcher: {
          url: rssUrl.value,
        },
      };
      form.value.source_config = JSON.stringify(config);
      form.value.source_type = 'rss';
    } else {
      // Validate JSON before saving in Advanced mode
      try {
        JSON.parse(form.value.source_config);
      } catch (e) {
        Message.error(t('customRecipe.form.error.invalidJson'));
        return;
      }
    }

    if (editing.value) {
      if (selectedRecipe.value) {
        await updateCustomRecipe(form.value);
        selectedRecipe.value.description = form.value.description;
        selectedRecipe.value.craft = form.value.craft;
        selectedRecipe.value.source_type = form.value.source_type;
        selectedRecipe.value.source_config = form.value.source_config;
      }
    } else {
      await createCustomRecipe(form.value as CustomRecipe);
      await listCustomRecipes();
    }
    showModal.value = false;
    form.value = {
      id: '',
      description: '',
      craft: '',
      source_type: 'rss',
      source_config: '',
    };
    editing.value = false;
    isUpdating.value = false;
    selectedRecipe.value = null;
    quickCreate.value = false;
    rssUrl.value = '';
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
      source_type: 'rss',
      source_config: '',
    };
    quickCreate.value = false;
    rssUrl.value = '';
  }
</script>

<style scoped></style>
