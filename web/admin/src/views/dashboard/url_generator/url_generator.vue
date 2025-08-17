<template>
  <div class="container p-4 w-full">
    <div>
      <div class="mb-8 text-2xl">
        <p class="text-gray-700"
          >{{ t('urlGenerator.welcome') }}<br />
          <span
            class="text-4xl font-bold text-sky-700 underline decoration-sky-500 decoration-wavy hover:underline-offset-2 hover:decoration-4"
            >Feed Craft</span
          ></p
        >
      </div>
      <a-card class="rounded-lg w-full min-w-8xl">
        <h1 class="text-2xl mb-2">{{ t('urlGenerator.title') }}</h1>
        <p class="mb-8">
          {{ t('urlGenerator.description') }}
        </p>
        <div class="mb-4">
          <label for="siteSelector" class="mr-2">{{
            t('urlGenerator.selectCraft')
          }}</label>
          <CraftFlowSelect v-model="selectedSite" mode="single" />
        </div>
        <div class="mb-4">
          <label for="inputUrl" class="mr-2">{{
            t('urlGenerator.inputOriginalUrl')
          }}</label>
          <a-input
            id="inputUrl"
            v-model="inputUrl"
            type="text"
            :placeholder="t('urlGenerator.inputUrlPlaceholder')"
          />
        </div>
        <a-button @click="appendPrefix"
          >{{ t('urlGenerator.showCraftedUrl') }}
        </a-button>
        <div class="mt-8">
          <label for="resultUrl" class="mr-2">{{
            t('urlGenerator.resultUrl')
          }}</label>
          <span id="resultUrl">{{ resultUrl }}</span>
          <a-button
            id="copyButton"
            class="px-2 py-0.5 rounded ml-0.5"
            :style="{ display: resultUrl ? 'inline-block' : 'none' }"
            @click="copyUrl"
            >{{ copyButtonText }}
          </a-button>
        </div>
      </a-card>
    </div>
  </div>
</template>

<script setup>
  import { ref, watch } from 'vue';
  import CraftFlowSelect from '@/views/dashboard/craft_flow/CraftFlowSelect.vue';
  import { useI18n } from 'vue-i18n';

  const { t } = useI18n();

  const selectedSite = ref('proxy');
  const customCraft = ref('');
  const inputUrl = ref('');
  const resultUrl = ref('');
  const copyButtonText = ref(t('urlGenerator.copyUrl'));

  const copyUrl = () => {
    if (resultUrl.value) {
      navigator.clipboard
        .writeText(resultUrl.value)
        .then(() => {
          copyButtonText.value = t('urlGenerator.copied');
        })
        .catch((err) => {
          console.error('无法复制文本: ', err);
        });
    }
  };

  watch(customCraft, (newValue) => {
    if (newValue) {
      selectedSite.value = ''; // Clear selectedSite if customCraft is used
    }
  });

  watch(selectedSite, (newValue) => {
    if (newValue) {
      customCraft.value = ''; // Clear customCraft if selectedSite is used
    }
  });
</script>

<style scoped>
  .container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 90vh;
  }
</style>
