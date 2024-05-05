<template>
  <div class="p-8 w-full">
    <div class="h-screen flex items-center justify-center">
      <div>
        <img :src="logo" style="height: 8rem; width: 8rem" />
        <h1 class="text-3xl my-4"
          >Welcome To
          <span class="font-bold"
            >Feed<span style="color: #79f1bf">Craft</span></span
          ></h1
        >
        <p class="text-xl"
          >你可以使用 FeedCraft 对你的 RSS 源进行任何你想要的处理</p
        >
        <a-row :gutter="20" class="my-8" style="min-width: 64rem">
          <a-col :span="8">
            <a-card title="项目主页" :bordered="false" hoverable>
              <a-link href="https://github.com/Colin-XKL/FeedCraft" icon
                >Colin-XKL/FeedCraft
              </a-link>
            </a-card>
          </a-col>
          <a-col :span="8">
            <a-card title="更多RSS" :bordered="false" hoverable>
              <a-link href="https://docs.rsshub.app" icon>RSSHub Doc</a-link>
            </a-card>
          </a-col>
          <a-col :span="8">
            <a-card title="了解RSS" :bordered="false" hoverable>
              <a-link href="https://sspai.com/post/56391" icon
                >RSS 入门指南 - sspai
              </a-link>
            </a-card>
          </a-col>
        </a-row>
        <a-row :gutter="20">
          <a-col :span="16">
            <a-card title="快速开始" :bordered="false" hoverable>
              <p class="py-2">
                地址:
                <span class="bg-slate-200 rounded-md border-0 px-2 py-0.5"
                  >/craft/{choose_craft_option}?input_url={input_rss_url}</span
                >
              </p>
              可使用的Craft Option:
              <li><b>proxy</b>仅代理模式</li>
              <li><b>fulltext</b>提取全文</li>
              <li><b>fulltext-plus</b>提取全文,使用模拟浏览器渲染网页内容</li>
              <li><b>introduction</b>使用AI生成文章简介,附加在原文开头</li>
            </a-card>
          </a-col>
          <a-col :span="8">
            <a-card
              hoverable
              title="问题反馈与建议"
              :bordered="false"
              :style="{ width: '100%' }"
            >
              <a-link
                href="https://github.com/Colin-XKL/FeedCraft/discussions"
                icon
                >Github Discussion
              </a-link>
            </a-card>
          </a-col>
        </a-row>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import logo from '@/assets/logo.png';
</script>

<script lang="ts">
  export default {
    name: 'WelcomePage',
  };
</script>
