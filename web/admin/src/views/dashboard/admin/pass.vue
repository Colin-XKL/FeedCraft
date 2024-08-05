<template>
  <div class="p-4">
    <h1 class="text-2xl font-bold mb-4">Change Password</h1>
    <a-form :model="form" @submit="handleSubmit">
      <a-form-item field="currentPassword" label="Current Password">
        <a-input v-model="form.currentPassword" type="password" />
      </a-form-item>
      <a-form-item field="newPassword" label="New Password">
        <a-input v-model="form.newPassword" type="password" />
      </a-form-item>
      <a-form-item field="confirmPassword" label="Confirm New Password">
        <a-input v-model="form.confirmPassword" type="password" />
      </a-form-item>
      <a-form-item>
        <a-button type="primary" html-type="submit">Change Password</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  import { changePassword } from '@/api/user';
  import { Message } from '@arco-design/web-vue';

  const form = ref({
    currentPassword: '',
    newPassword: '',
    confirmPassword: '',
  });

  const handleSubmit = async (event: any) => {
    event.preventDefault();
    if (form.value.newPassword.length < 6) {
      Message.error('New password must be at least 6 characters long');
      return;
    }
    if (/^\d+$/.test(form.value.newPassword)) {
      Message.error('New password cannot be purely numeric');
      return;
    }
    if (form.value.newPassword !== form.value.confirmPassword) {
      Message.error('New passwords do not match');
      return;
    }
    try {
      await changePassword({
        username: 'currentUser', // 这里需要替换为当前登录用户的用户名
        currentPassword: form.value.currentPassword,
        newPassword: form.value.newPassword,
      });
      Message.success('Password changed successfully');
    } catch (error) {
      Message.error('Failed to change password');
    }
  };
</script>

<script lang="ts">
  export default {
    name: 'ChangePassword',
  };
</script>
