<template>
  <div class="container p-4 w-full">
    <div>
      <div class="mb-8 text-2xl">
        <p
          >Welcome to<br />
          <span
            class="text-4xl font-bold text-sky-700 underline decoration-sky-500 decoration-wavy hover:underline-offset-2 hover:decoration-4"
            >Feed Craft</span
          ></p
        >
      </div>
      <a-card class="rounded-lg w-full min-w-8xl">
        <h1 class="text-2xl mb-2">快速生成 FeedCraft URL</h1>
        <p class="mb-8">
          输入原 RSS URL，选择需要的 craft，即可生成最终URL。
        </p>
        <div class="mb-4">
          <label for="siteSelector" class="mr-2">选择一个 craft:</label>
          <CraftFlowSelect v-model="selectedSite" mode="single" />
        </div>
        <div class="mb-4">
          <label for="inputUrl" class="mr-2">输入原 RSS URL:</label>
          <a-input
            id="inputUrl"
            v-model="inputUrl"
            type="text"
            placeholder="输入 URL"
          />
        </div>
        <a-button @click="appendPrefix">显示 Crafted Feed URL </a-button>
        <div class="mt-8">
          <label for="resultUrl" class="mr-2">结果 URL:</label>
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

  const selectedSite = ref('proxy');
  const customCraft = ref('');
  const inputUrl = ref('');
  const resultUrl = ref('');
  const copyButtonText = ref('复制 URL');

  const appendPrefix = () => {
    const currentSelectedSite = customCraft.value
      ? customCraft.value
      : selectedSite.value;
    const baseUrl = import.meta.env.VITE_API_BASE_URL ?? window.location.origin;
    resultUrl.value = `${baseUrl}/craft/${currentSelectedSite}?input_url=${encodeURIComponent(
      inputUrl.value
    )}`;
    copyButtonText.value = '复制 URL';
  };

  const copyUrl = () => {
    if (resultUrl.value) {
      navigator.clipboard
        .writeText(resultUrl.value)
        .then(() => {
          copyButtonText.value = '已复制!';
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
    height: 100vh;
  }
</style>
