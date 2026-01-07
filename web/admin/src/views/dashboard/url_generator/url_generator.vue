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
            t('urlGenerator.selectProcessor')
          }}</label>
          <BlueprintSelect v-model="selectedProcessor" mode="single" />
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
        <a-button type="primary" @click="generateUrl"
          >{{ t('urlGenerator.showProcessedUrl') }}
        </a-button>
        <div class="mt-8">
          <label for="resultUrl" class="mr-2">{{
            t('urlGenerator.resultUrl')
          }}</label>
          <span id="resultUrl">{{ resultUrl }}</span>
          <a-button
            id="copyButton"
            class="px-2 py-0.5 rounded ml-0.5"
            @click="copyUrl"
            >{{ copyButtonText }}
          </a-button>
        </div>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { Ref, ref } from 'vue';
  import BlueprintSelect from '@/views/dashboard/blueprint/BlueprintSelect.vue';
  import { useI18n } from 'vue-i18n';

  const { t } = useI18n();

  const selectedProcessor = ref<string[]>([]);
  const customProcessor = ref('');
  const inputUrl = ref('');
  const resultUrl = ref('');
  const copyButtonText = ref(t('urlGenerator.copyUrl'));

  const generateUrl = () => {
    const currentSelectedCraft = customProcessor.value
      ? customProcessor.value
      : selectedProcessor.value;
    const baseUrl = import.meta.env.VITE_API_BASE_URL ?? window.location.origin;
    resultUrl.value = `${baseUrl}/craft/${currentSelectedCraft}?input_url=${encodeURIComponent(
      inputUrl.value,
    )}`;
    copyButtonText.value = t('urlGenerator.copyUrl');
  };

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
</script>

<style scoped>
  .container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 90vh;
  }
</style>
