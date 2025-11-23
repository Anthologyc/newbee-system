<template>
  <a-layout class="layout-demo">
    <!-- 左侧侧边栏 -->
    <a-layout-sider collapsible breakpoint="xl">
      <div class="logo">
        NewBee Admin
      </div>
      <a-menu
        :default-open-keys="['1']"
        :default-selected-keys="currentPath"
        @menu-item-click="onClickMenuItem"
      >
        <a-menu-item key="dashboard">
          <template #icon><icon-dashboard /></template>
          仪表盘
        </a-menu-item>
        
        <a-sub-menu key="questions">
          <template #icon><icon-book /></template>
          <template #title>题目管理</template>
          <a-menu-item key="question-list">题目列表</a-menu-item>
          <a-menu-item key="question-create">录入题目</a-menu-item>
        </a-sub-menu>

        <a-menu-item key="announcements">
          <template #icon><icon-notification /></template>
          公告发布
        </a-menu-item>

        <a-menu-item key="users">
          <template #icon><icon-user-group /></template>
          用户管理
        </a-menu-item>
      </a-menu>
    </a-layout-sider>

    <a-layout>
      <!-- 顶部 Header -->
      <a-layout-header style="padding-left: 20px; background: #fff; border-bottom: 1px solid #eee;">
        <div class="header-content">
          <h3>后台管理系统</h3>
          <a-dropdown @select="handleCommand">
            <a-avatar style="background-color: #3370ff; cursor: pointer;">
              <IconUser />
            </a-avatar>
            <template #content>
              <a-doption value="logout">退出登录</a-doption>
            </template>
          </a-dropdown>
        </div>
      </a-layout-header>

      <!-- 右侧内容区 -->
      <a-layout-content style="padding: 24px;">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { 
  IconDashboard, 
  IconBook, 
  IconNotification, 
  IconUserGroup, 
  IconUser 
} from '@arco-design/web-vue/es/icon';

const router = useRouter();
const route = useRoute();

// 计算当前选中的菜单项，用于高亮
const currentPath = computed(() => [route.name as string]);

// 菜单点击跳转
const onClickMenuItem = (key: string) => {
  router.push({ name: key });
};

// 顶部下拉菜单处理
const handleCommand = (value: any) => {
  if (value === 'logout') {
    localStorage.clear();
    router.push('/');
  }
};
</script>

<style scoped>
.layout-demo {
  height: 100vh;
}
.logo {
  height: 32px;
  margin: 16px;
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
  text-align: center;
  line-height: 32px;
  font-weight: bold;
  font-size: 16px;
}
.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
  padding-right: 20px;
}
</style>