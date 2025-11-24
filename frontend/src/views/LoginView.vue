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
            {{ isRegister ? '立即注册' : '登录' }}
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

const form = reactive({
  username: '',
  password: ''
});

const toggleMode = () => {
  isRegister.value = !isRegister.value;
  form.password = ''; // 切换时清空密码
};

const handleSubmit = async () => {
  if (!form.username || !form.password) {
    Message.warning('请输入用户名和密码');
    return;
  }

  loading.value = true;
  const url = isRegister.value ? '/register' : '/login';

  try {
    const res = await axios.post(url, form);

    if (isRegister.value) {
      // 注册逻辑
      Message.success('注册成功，正在自动登录...');

      // 注册成功后尝试自动登录（可选，或者直接切回登录页让用户再输一遍）
      // 这里为了稳妥，先切回登录页，让用户确认
      isRegister.value = false;
      form.password = ''; // 稍微清一下，防误触

    } else {
      // 登录逻辑
      Message.success('登录成功！');
      const { token, role } = res.data;

      localStorage.setItem('token', token);
      localStorage.setItem('role', role || 'user'); // 默认 user

      // 根据角色跳转
      if (role === 'admin') {
        router.push({ name: 'dashboard' });
      } else {
        // 🚀 加上这一行：如果后端返回了 id 或 user_id
        // 如果后端 login 接口目前只返回了 token, role, username
        // 你可能需要去后端 handlers/auth.go 的 Login 函数里把 ID 也返回出来
        // 比如： c.JSON(http.StatusOK, gin.H{"token": tokenString, "role": user.Role, "user_id": user.ID})

        localStorage.setItem('user_id', res.data.user_id);

        router.push({ name: 'practice-home' });
      }
    }
  } catch (error: any) {
    console.error(error);
    const msg = error.response?.data?.error || '请求失败';
    Message.error(msg);
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