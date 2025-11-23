<template>
  <div class="login-container">
    <a-card class="login-card" title="NewBee System 登录">
      <a-form :model="form" @submit="handleSubmit">
        <a-form-item field="username" label="用户名">
          <a-input v-model="form.username" placeholder="请输入用户名" />
        </a-form-item>
        <a-form-item field="password" label="密码">
          <a-input-password v-model="form.password" placeholder="请输入密码" />
        </a-form-item>
        
        <div class="actions">
          <a-button type="primary" html-type="submit" long :loading="loading">
            {{ isRegister ? '注册并登录' : '登录' }}
          </a-button>
          <a-button type="text" long @click="toggleMode">
            {{ isRegister ? '已有账号？去登录' : '没有账号？去注册' }}
          </a-button>
        </div>
      </a-form>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { Message } from '@arco-design/web-vue';
import axios from 'axios';
import { useRouter } from 'vue-router';

const router = useRouter();
const isRegister = ref(false);
const loading = ref(false);

// 定义接口类型让 TS 更开心（可选，不写也行）
interface LoginForm {
  username: string;
  password: string;
}

const form = reactive<LoginForm>({
  username: '',
  password: ''
});

const toggleMode = () => {
  isRegister.value = !isRegister.value;
};

const handleSubmit = async () => {
  if (!form.username || !form.password) {
    Message.warning('请输入用户名和密码');
    return;
  }

  loading.value = true;
  const url = isRegister.value ? 'http://localhost:8080/register' : 'http://localhost:8080/login';

  try {
    const res = await axios.post(url, form);

    if (isRegister.value) {
      Message.success('注册成功，请登录');
      isRegister.value = false;
    } else {
      Message.success('登录成功！');
      localStorage.setItem('token', res.data.token);
      localStorage.setItem('role', res.data.role);
      // 👉 修改这里：跳转到新的路由名称或路径
      router.push({ name: 'dashboard' }); 
    }
  } catch (error: any) {
    Message.error(error.response?.data?.error || '操作失败');
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f2f3f5;
}
.login-card {
  width: 400px;
}
.actions {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
</style>