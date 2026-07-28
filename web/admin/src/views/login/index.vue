<template>
  <div class="container">
    <div class="logo">
      <div class="logo-text">FeedCraft Admin Panel</div>
    </div>
    <div class="content">
      <div class="content-inner">
        <LoginForm />
      </div>
      <div class="footer">
        <Footer />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import Footer from '@/components/footer/index.vue';
  import LoginForm from '@/views/login/components/login-form.vue';
</script>

<style lang="less" scoped>
  .container {
    display: flex;
    height: 100vh;

    .content {
      position: relative;
      display: flex;
      flex: 1;
      align-items: center;
      justify-content: center;
      padding-bottom: 40px;
    }

    .footer {
      position: absolute;
      right: 0;
      bottom: 0;
      width: 100%;
    }
  }

  .logo {
    position: fixed;
    top: 24px;
    left: 22px;
    z-index: 1;
    display: inline-flex;
    align-items: center;

    &-text {
      margin-right: 4px;
      margin-left: 4px;
      color: var(--color-text-1);
      font-size: 20px;
      font-weight: 600;
    }
  }
</style>
