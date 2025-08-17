<template>
  <div class="container p-4 w-full">
    <div>
      <div class="mb-8 text-2xl">
        <p class="text-gray-700"
          >Welcome to<br />
          <span
            class="text-4xl font-bold text-sky-700 underline decoration-sky-500 decoration-wavy hover:underline-offset-2 hover:decoration-4"
            >Feed Craft</span
          ></p
        >
      </div>
      <div class="py-8 bg-gray-100 rounded-lg w-full h-72 mx-auto">
        <h1 class="text-2xl mb-4">快速生成 FeedCraft URL</h1>
        <p class="text-gray-700 mb-4">
          输入原 RSS URL，选择需要的 craft，即可生成最终URL。
        </p>
        <div class="mb-4">
          <label for="siteSelector" class="mr-2">选择一个 craft:</label>
          <select
            id="siteSelector"
            v-model="selectedSite"
            class="p-2 border border-gray-300 rounded"
          >
            <option value="proxy">proxy - 代理订阅源</option>
            <option value="limit">limit - 限制单页条目数量</option>
            <option value="relative-link-fix"
              >relative-link-fix - 修复文章链接</option
            >
            <option value="fulltext">fulltext - 提取全文</option>
            <option value="fulltext-plus"
              >fulltext-plus - 模拟浏览器提取全文</option
            >
            <option value="cleanup">cleanup - 清理文章HTML内容</option>
            <option value="introduction">introduction - AI生成导读</option>
            <option value="summary">summary - AI总结文章</option>
            <option value="ignore-advertorial"
              >ignore-advertorial - 排除广告文章</option
            >
            <option value="translate-title">translate-title - 翻译标题</option>
            <option value="translate-content"
              >translate-content - 翻译内容</option
            >
          </select>
          <label for="customCraft" class="ml-4 mr-2">或输入自定义 craft:</label>
          <input
            id="customCraft"
            v-model="customCraft"
            type="text"
            placeholder="输入自定义 craft"
            class="p-2 border border-gray-300 rounded"
          />
        </div>
        <div class="mb-4">
          <label for="inputUrl" class="mr-2">输入原 RSS URL:</label>
          <input
            id="inputUrl"
            v-model="inputUrl"
            type="text"
            placeholder="输入 URL"
            class="p-2 border border-gray-300 rounded"
          />
        </div>
        <button
          class="px-4 py-2 bg-blue-500 text-white rounded"
          @click="appendPrefix"
          >显示 Crafted Feed URL
        </button>
        <div class="mt-4">
          <label for="resultUrl" class="mr-2">结果 URL:</label>
          <span id="resultUrl">{{ resultUrl }}</span>
          <button
            id="copyButton"
            class="px-2 py-0.5 bg-gray-200 text-gray-700 rounded ml-0.5 hover:bg-teal-500"
            :style="{ display: resultUrl ? 'inline-block' : 'none' }"
            @click="copyUrl"
            >{{ copyButtonText }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { ref, watch } from 'vue';

  const selectedSite = ref('proxy');
  const customCraft = ref('');
  const inputUrl = ref('');
  const resultUrl = ref('');
  const copyButtonText = ref('复制 URL');

  const appendPrefix = () => {
    const currentSelectedSite = customCraft.value
      ? customCraft.value
      : selectedSite.value;
    const baseUrl = window.location.origin;
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
