<template>
  <div class="p-4">
    <a-card title="修改密码">
      <a-form :model="form" @submit="handleSubmit">
        <a-form-item label="当前用户" field="username">
          <a-input :model-value="currentUser" disabled></a-input>
        </a-form-item>
        <a-form-item field="currentPassword" label="当前密码">
          <a-input v-model="form.currentPassword" type="password" />
        </a-form-item>
        <a-form-item field="newPassword" label="新密码">
          <a-input v-model="form.newPassword" type="password" />
        </a-form-item>
        <a-form-item field="confirmPassword" label="确认新密码">
          <a-input v-model="form.confirmPassword" type="password" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit">修改密码</a-button>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { computed, ref } from 'vue';
  import { changePassword } from '@/api/user';
  import { Message } from '@arco-design/web-vue';
  import { useUserStore } from '@/store';

  const form = ref({
    currentPassword: '',
    newPassword: '',
    confirmPassword: '',
  });
  const userStore = useUserStore();
  const currentUser = computed(() => {
    return userStore.name || 'unknown user';
  });
  const handleSubmit = async (event: any) => {
    // event.preventDefault();
    if (form.value.newPassword.length < 6) {
      Message.error('新密码必须至少6个字符长');
      return;
    }
    if (/^\d+$/.test(form.value.newPassword)) {
      Message.error('新密码不能纯数字');
      return;
    }
    if (form.value.newPassword !== form.value.confirmPassword) {
      Message.error('新密码不匹配');
      return;
    }
    try {
      await changePassword({
        username: currentUser.value,
        currentPassword: form.value.currentPassword,
        newPassword: form.value.newPassword,
      });
      Message.success('密码修改成功');
    } catch (error) {
      Message.error('密码修改失败');
    }
  };
</script>

<script lang="ts">
  export default {
    name: 'ChangePassword',
  };
</script>
