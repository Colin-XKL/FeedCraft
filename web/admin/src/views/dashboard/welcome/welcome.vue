<template>
  <div class="p-8 w-full">
    <div class="h-screen flex items-center justify-center">
      <div>
        <img
          :src="logo"
          alt="FeedCraft Logo"
          style="height: 8rem; width: 8rem"
        />
        <h1 class="text-3xl my-4 font-light"
          >Welcome To
          <span class="font-bold"
            ><span style="color: #4d4f5a">Feed</span
            ><span style="color: #0d9488">Craft</span></span
          ></h1
        >
        <p class="text-xl"
          >你可以使用 FeedCraft 对你的 RSS 源进行任何你想要的处理</p
        >
        <a-row :gutter="20" class="my-8" style="min-width: 64rem">
          <a-col :span="8">
            <a-card title="项目主页" :bordered="false" hoverable>
              <a-link
                href="https://github.com/Colin-XKL/FeedCraft"
                target="_blank"
                rel="noopener noreferrer"
                icon
                >Colin-XKL/FeedCraft
              </a-link>
            </a-card>
          </a-col>
          <a-col :span="8">
            <a-card title="更多RSS" :bordered="false" hoverable>
              <a-link
                href="https://docs.rsshub.app"
                target="_blank"
                rel="noopener noreferrer"
                icon
              >
                RSSHub Doc
              </a-link>
            </a-card>
          </a-col>
          <a-col :span="8">
            <a-card title="了解RSS" :bordered="false" hoverable>
              <a-link
                href="https://sspai.com/post/56391"
                target="_blank"
                rel="noopener noreferrer"
                icon
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
              <div class="grid grid-cols-2 gap-2">
                <div>
                  <h3 class="font-bold">基础处理</h3>
                  <ul>
                    <li><b>proxy</b> - 仅代理模式</li>
                    <li><b>limit</b> - 限制单页条目数量</li>
                    <li><b>keyword</b> - 关键词过滤</li>
                    <li><b>guid-fix</b> - 修复RSS GUID</li>
                    <li><b>relative-link-fix</b> - 修复文章链接</li>
                  </ul>
                </div>
                <div>
                  <h3 class="font-bold">内容提取</h3>
                  <ul>
                    <li><b>fulltext</b> - 提取全文</li>
                    <li><b>fulltext-plus</b> - 模拟浏览器提取全文</li>
                    <li><b>cleanup</b> - 清理文章HTML内容</li>
                  </ul>
                </div>
                <div>
                  <h3 class="font-bold">AI处理</h3>
                  <ul>
                    <li><b>introduction</b> - AI生成导读</li>
                    <li><b>summary</b> - AI总结文章</li>
                    <li><b>ignore-advertorial</b> - 排除广告文章</li>
                    <li><b>translate-title</b> - 翻译标题</li>
                    <li><b>translate-content</b> - 翻译内容</li>
                  </ul>
                </div>
              </div>
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
                target="_blank"
                rel="noopener noreferrer"
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
