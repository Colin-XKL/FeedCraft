<template>
  <div>
    <x-header
      title="Feed Compare"
      description="指定feed源url与要应用的craft(craft flow 或 atom都可), 对比处理前和处理后的结果"
    ></x-header>
  </div>
</template>

<script lang="ts" setup>
  import XHeader from '@/components/header/x-header.vue';
</script>

<script lang="ts">
  export default {
    name: 'FeedCompare',
  };
</script>
