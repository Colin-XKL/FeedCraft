<template>
  <div class="mb-4">
    <h1 class="text-3xl font-bold">{{ props.title }}</h1>
    <p class="text-gray-500">{{ props.description }}</p>
  </div>
</template>

<script lang="ts" setup>
  interface HeaderProps {
    title: string;
    description?: string;
  }

  const props = defineProps<HeaderProps>();
</script>

<script lang="ts">
  export default {
    name: 'XHeader',
  };
</script>

<style scoped lang="less"></style>
