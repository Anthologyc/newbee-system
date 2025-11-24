<template>
  <div class="user-layout">
    <!-- 顶部导航栏 -->
    <div class="nav-header">
      <div class="logo" @click="$router.push('/practice')">NewBee 刷题</div>
      
      <div class="nav-items">
        <router-link to="/practice" class="nav-link">首页</router-link>
        <router-link v-if="isAdmin" to="/admin/dashboard" class="nav-link admin-link">后台</router-link>
        <a-button type="text" size="small" @click="logout" class="logout-btn">退出</a-button>
      </div>
    </div>
    
    <!-- 内容区域 -->
    <div class="content">
      <router-view />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const isAdmin = computed(() => localStorage.getItem('role') === 'admin');

const logout = () => {
  localStorage.clear();
  router.push('/login');
};
</script>

<style scoped>
.nav-header {
  height: 56px;
  background: #fff;
  border-bottom: 1px solid #eee;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  position: sticky;
  top: 0;
  z-index: 999;
}
.logo {
  font-size: 18px;
  font-weight: bold;
  color: #165dff;
  cursor: pointer;
}
.nav-items { display: flex; align-items: center; gap: 15px; }

.nav-link { color: #4e5969; text-decoration: none; font-size: 14px; }
.nav-link.router-link-active { color: #165dff; font-weight: 500; }

.content {
  padding: 20px;
  background: #f7f8fa;
  min-height: calc(100vh - 56px);
}

/* 📱 移动端适配 */
@media (max-width: 768px) {
  .nav-header { padding: 0 15px; }
  .content { padding: 10px; } /* 减少内边距 */
  .admin-link { display: none; } /* 手机上一般不进后台，隐藏保持简洁 */
}
</style>