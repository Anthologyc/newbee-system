<template>
  <div class="home-container">
    <div class="grid-layout">
      <!-- 左侧：欢迎卡片 -->
      <div class="welcome-section">
        <a-card class="welcome-card" :bordered="false">
          <div class="welcome-content">
            <h2>欢迎来到 NewBee 刷题中心 🚀</h2>
            <p>这里提供最全的题库，助你轻松通关！</p>
            <div class="actions">
              <a-button type="primary" size="large" shape="round" @click="$router.push({name: 'practice-mode'})">
                <template #icon><icon-play-circle /></template>
                开始练习
              </a-button>
            </div>
          </div>
          <img src="https://p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/cd7a1a23e382488d9c837d2f9d74ad3c.png~tplv-uwbnlip3yd-webp.webp" class="bg-img" />
        </a-card>
      </div>
      
      <!-- 右侧：公告 -->
      <div class="announcement-section">
        <a-card title="最新公告" :bordered="false" class="anno-card">
          <a-list :data="announcements" :bordered="false" :split="false">
            <template #item="{ item }">
              <div class="anno-item">
                <div class="anno-title">{{ item.title }}</div>
                <div class="anno-date">{{ new Date(item.created_at).toLocaleDateString() }}</div>
              </div>
            </template>
          </a-list>
        </a-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { IconPlayCircle } from '@arco-design/web-vue/es/icon';
import { useSocketStore } from '../../stores/socket';
const socketStore = useSocketStore();
// 假设 userId = 1
socketStore.connect(1);


const announcements = ref([]);

onMounted(async () => {
  const res = await axios.get('/api/announcements');
  announcements.value = res.data.data || [];
});
</script>

<style scoped>
.home-container { max-width: 1200px; margin: 0 auto; padding: 20px; }

.grid-layout {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 24px;
}

.welcome-card {
  height: 360px;
  background: linear-gradient(135deg, #e8f3ff 0%, #f2f9ff 100%);
  position: relative;
  overflow: hidden;
  border-radius: 12px;
}
.welcome-content { position: relative; z-index: 2; padding: 40px; }
h2 { font-size: 28px; color: #1d2129; margin-bottom: 10px; }
p { font-size: 16px; color: #4e5969; margin-bottom: 40px; }
.bg-img {
  position: absolute; right: 0; bottom: 0; height: 90%; opacity: 0.8; z-index: 1;
}

.anno-card { height: 100%; border-radius: 12px; }
.anno-item {
  display: flex; justify-content: space-between; padding: 12px 0;
  border-bottom: 1px solid #f2f3f5;
}
.anno-title { font-weight: 500; color: #1d2129; }
.anno-date { color: #86909c; font-size: 12px; }

/* 📱 移动端适配 */
@media (max-width: 768px) {
  .home-container { padding: 10px; }
  .grid-layout { grid-template-columns: 1fr; } /* 变单栏 */
  
  .welcome-card { height: auto; padding-bottom: 20px; }
  .welcome-content { padding: 20px; text-align: center; }
  .bg-img { display: none; } /* 手机上隐藏大背景图，太占地 */
  
  h2 { font-size: 22px; }
  .actions { display: flex; justify-content: center; }
}
</style>