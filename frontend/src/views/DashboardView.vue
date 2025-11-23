<template>
  <div class="dashboard-container">
    <a-card title="工作台数据">
      <a-grid :cols="3" :colGap="20">
        <a-grid-item>
          <a-statistic 
            title="总题目数" 
            :value="stats.question_count" 
            show-group-separator 
            animation
          >
            <template #prefix><icon-book /></template>
          </a-statistic>
        </a-grid-item>
        <a-grid-item>
          <a-statistic 
            title="注册用户" 
            :value="stats.user_count" 
            show-group-separator 
            animation
          >
            <template #prefix><icon-user-group /></template>
          </a-statistic>
        </a-grid-item>
        <a-grid-item>
          <!-- 🚀 修改这里：不使用 :value，改用 value 插槽 -->
          <a-statistic title="系统版本">
            <template #value>
              v1.0.0
            </template>
          </a-statistic>
        </a-grid-item>
      </a-grid>
    </a-card>
    
    <a-card title="最新公告" style="margin-top: 20px;">
      <div v-if="stats.announcement && stats.announcement.id">
        <h3>{{ stats.announcement.title }}</h3>
        <div style="color: #86909c; font-size: 12px; margin-bottom: 10px;">
          发布于：{{ new Date(stats.announcement.created_at).toLocaleString() }}
        </div>
        <p style="line-height: 1.6; white-space: pre-wrap;">
          {{ stats.announcement.content }}
        </p>
      </div>
      <div v-else>
        <a-empty description="暂无公告" />
      </div>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { IconBook, IconUserGroup } from '@arco-design/web-vue/es/icon';

const stats = ref({
  question_count: 0,
  user_count: 0,
  announcement: {
    id: 0,
    title: '',
    content: '',
    created_at: ''
  }
});

onMounted(async () => {
  try {
    const res = await axios.get('/api/dashboard/stats');
    stats.value = res.data;
  } catch (e) {
    console.error('获取仪表盘数据失败', e);
  }
});
</script>

<style scoped>
.dashboard-container {
  padding: 0 10px;
}
</style>