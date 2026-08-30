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
      class="custom-recipe-table"
      :data="recipes"
      :columns="columns"
      :bordered="true"
      :loading="isLoading"
      :scroll="{ x: 1280 }"
    >
      <template #status="{ record }">
        <a-tooltip
          v-if="record.is_active"
          :content="
            t('customRecipe.status.activeTooltip', {
              time: dayjs(record.last_accessed_at).format(
                'YYYY-MM-DD HH:mm:ss'
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
      <template #craft="{ record }">
        <span class="custom-recipe-table__text" :title="record.craft">{{
          record.craft
        }}</span>
      </template>
      <template #actions="{ record }">
        <a-space>
          <a-button
            type="primary"
            size="small"
            @click="
              () => {
                isUpdating = true;
                showEditModal(record);
              }
            "
          >
            {{ t('customRecipe.edit') }}
          </a-button>
          <a-button type="text" size="small" @click="previewRecipe(record)">
            {{ t('customRecipe.preview') }}
          </a-button>
          <a-button type="text" size="small" @click="handleCopyLink(record.id)">
            {{ t('customRecipe.copyLink') }}
          </a-button>
          <a-popconfirm
            :content="t('customRecipe.deleteConfirm')"
            @ok="deleteRecipe(record.id)"
          >
            <a-button type="text" status="danger" size="small">
              {{ t('customRecipe.delete') }}
            </a-button>
          </a-popconfirm>
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
        ref="formRef"
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
          >
            <a-input
              v-model="form.feed_url"
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
              <a-option value="web_monitor">{{
                t('menu.webMonitor')
              }}</a-option>
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
              @blur="formatConfig"
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
        <a-button type="primary" :loading="saving" @click="saveRecipe">{{
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
      <div style="margin-bottom: 10px; text-align: right">
        <a-tooltip
          :content="
            copied ? t('customRecipe.copied') : t('customRecipe.copyConfig')
          "
        >
          <a-button size="mini" @click="handleCopyConfig">
            <template #icon><icon-copy /></template>
            {{ t('customRecipe.copyConfig') }}
          </a-button>
        </a-tooltip>
      </div>
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
  import { ref, onMounted, computed, nextTick, watch, type ComputedRef } from 'vue';
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
  import { useRouter } from 'vue-router';
  import { useClipboard } from '@vueuse/core';
  import buildPublicFeedUrl from '@/utils/publicFeedUrl';
  import CraftSelector from '@/views/dashboard/craft_flow/CraftSelector.vue';

  const { t } = useI18n();
  const router = useRouter();

  const recipes = ref<CustomRecipe[]>([]);
  const showModal = ref(false);
  const showConfigModal = ref(false);
  const currentConfig = ref('');
  const currentLink = ref('');
  type RecipeForm = CustomRecipe & { feed_url: string };

  const emptyRecipeForm = (): RecipeForm => ({
    id: '',
    description: '',
    craft: '',
    source_type: 'rss',
    source_config: '',
    feed_url: '',
  });

  const toRecipePayload = (value: RecipeForm): CustomRecipe => ({
    id: value.id,
    description: value.description,
    craft: value.craft,
    source_type: value.source_type,
    source_config: value.source_config,
  });

  const quickCreate = ref(false);

  const form = ref<RecipeForm>(emptyRecipeForm());
  const formRef = ref();

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

  const { copy: copyConfig, copied } = useClipboard({
    source: currentConfig,
    legacy: true,
    copiedDuring: 1500,
  });
  const { copy: copyLink } = useClipboard({
    source: currentLink,
    legacy: true,
    copiedDuring: 1500,
  });
  const buildRecipeFeedUrl = (id?: string) =>
    buildPublicFeedUrl(`/recipe/${id || ''}`);

  const handleCopyConfig = async () => {
    try {
      await copyConfig();
      Message.success(t('customRecipe.copied'));
    } catch (e: any) {
      Message.error(t('customRecipe.copyFailed', { msg: e.message || e }));
    }
  };

  const handleCopyLink = async (id: string) => {
    try {
      currentLink.value = buildRecipeFeedUrl(id);
      await copyLink();
      Message.success(t('customRecipe.copied'));
    } catch (e: any) {
      Message.error(t('customRecipe.copyFailed', { msg: e.message || e }));
    }
  };

  const formatConfig = () => {
    try {
      const obj = JSON.parse(form.value.source_config);
      form.value.source_config = JSON.stringify(obj, null, 2);
    } catch (e) {
      // invalid json, ignore
    }
  };

  const editing = ref(false);
  const selectedRecipe = ref<CustomRecipe | null>(null);
  const isLoading = ref(false);
  const isUpdating = ref(false);
  const saving = ref(false);

  const columns = [
    {
      title: t('customRecipe.form.name'),
      dataIndex: 'id',
      minWidth: 200,
      ellipsis: true,
      tooltip: true,
    },
    {
      title: t('customRecipe.form.description'),
      dataIndex: 'description',
      minWidth: 180,
      ellipsis: true,
      tooltip: true,
    },
    {
      title: t('customRecipe.form.craft'),
      dataIndex: 'craft',
      slotName: 'craft',
      minWidth: 200,
      ellipsis: true,
      tooltip: true,
    },
    {
      title: t('customRecipe.status.active'),
      slotName: 'status',
      width: 100,
    },
    {
      title: t('customRecipe.form.sourceType'),
      dataIndex: 'source_type',
      width: 110,
    },
    {
      title: t('customRecipe.form.sourceConfig'),
      dataIndex: 'source_config',
      slotName: 'source_config',
      minWidth: 240,
    },
    {
      title: t('customRecipe.actions'),
      slotName: 'actions',
      width: 280,
      fixed: 'right' as const,
    },
  ];

  async function listCustomRecipes() {
    isLoading.value = true;
    recipes.value = (await getCustomRecipes()).data;
    isLoading.value = false;
  }

  onMounted(() => {
    listCustomRecipes();
  });

  watch(showModal, (visible) => {
    if (visible) {
      nextTick(() => {
        formRef.value?.clearValidate();
      });
    }
  });
  const rules: ComputedRef<Record<string, unknown[]>> = computed(() => {
    const requiredNameAndCraft = {
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
          trigger: 'change',
        },
      ],
    };

    if (quickCreate.value) {
      return {
        ...requiredNameAndCraft,
        feed_url: [
          {
            required: true,
            message: t('customRecipe.form.rule.rssUrlRequired'),
            trigger: 'blur',
          },
        ],
      };
    }

    return {
      ...requiredNameAndCraft,
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
  });

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
      feed_url: '',
    };
    showModal.value = true;
  };

  const saveRecipe = async () => {
    const invalid = await formRef.value?.validate();
    if (invalid) {
      return;
    }

    if (quickCreate.value && !editing.value) {
      form.value.source_config = JSON.stringify({
        http_fetcher: {
          url: form.value.feed_url.trim(),
        },
      });
      form.value.source_type = 'rss';
    } else {
      try {
        JSON.parse(form.value.source_config);
      } catch (e) {
        Message.error(t('customRecipe.form.error.invalidJson'));
        return;
      }
    }

    saving.value = true;
    try {
      const payload = toRecipePayload(form.value);
      if (editing.value) {
        if (selectedRecipe.value) {
          await updateCustomRecipe(payload);
          selectedRecipe.value.description = payload.description;
          selectedRecipe.value.craft = payload.craft;
          selectedRecipe.value.source_type = payload.source_type;
          selectedRecipe.value.source_config = payload.source_config;
        }
      } else {
        await createCustomRecipe(payload);
        await listCustomRecipes();
      }
      showModal.value = false;
      form.value = emptyRecipeForm();
      editing.value = false;
      isUpdating.value = false;
      selectedRecipe.value = null;
      quickCreate.value = false;
    } catch {
      // Axios interceptor already displays the API error toast.
    } finally {
      saving.value = false;
    }
  };

  const deleteRecipe = async (id: string) => {
    await deleteCustomRecipe(id);
    await listCustomRecipes();
  };

  const previewRecipe = (record: CustomRecipe) => {
    router.push({
      name: 'FeedViewer',
      query: { target: 'recipe', id: record.id },
    });
  };

  function resetForm() {
    form.value = emptyRecipeForm();
    quickCreate.value = false;
  }
</script>

<style scoped>
  .custom-recipe-table :deep(.arco-table-th),
  .custom-recipe-table :deep(.arco-table-td) {
    white-space: nowrap;
    word-break: normal;
    overflow-wrap: normal;
  }

  .custom-recipe-table__text {
    display: inline-block;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    vertical-align: bottom;
  }
</style>
