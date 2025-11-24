<template>
  <div class="lobby-container">
    <div class="lobby-header">
      <div class="title-block">
        <h2>🤝 多人协作大厅</h2>
        <p>寻找队友，并肩作战！</p>
      </div>
      <a-button type="primary" size="large" @click="showConfigModal = true">
        <template #icon><icon-plus /></template> 创建房间
      </a-button>
    </div>

    <a-divider />

    <div class="room-grid">
      <a-card 
        v-for="room in rooms" 
        :key="room.id" 
        class="room-card" 
        hoverable
        :class="{ 'full': room.count >= 2 }"
      >
        <template #title>
          <span class="room-title">房间 #{{ room.id }}</span>
        </template>
        <template #extra>
          <a-tag :color="room.status === 'waiting' ? 'green' : 'blue'">
            {{ room.status === 'waiting' ? '等待中' : '进行中' }}
          </a-tag>
        </template>

        <div class="room-content">
          <div class="avatar-group">
             <a-avatar-group :size="32" :max-count="2">
               <a-avatar v-for="i in room.count" :key="i" :style="{ backgroundColor: '#165dff' }">P{{i}}</a-avatar>
             </a-avatar-group>
          </div>
          <div class="count-info">{{ room.count }} / 2 人</div>
        </div>

        <div class="actions">
          <a-button 
            type="primary" 
            status="success" 
            long 
            @click="handleJoinRoom(room.id)"
            :disabled="room.count >= 2 || room.status !== 'waiting'"
          >
            {{ room.count >= 2 ? '满员' : '加入房间' }}
          </a-button>
        </div>
      </a-card>
    </div>

    <div v-if="rooms.length === 0" class="empty-state">
      <a-empty description="暂无房间，快去创建一个吧！" />
    </div>

    <!-- 🚀 新增：创建房间配置弹窗 -->
    <a-modal v-model:visible="showConfigModal" title="创建房间设置" @ok="confirmCreate">
      <a-form :model="form" layout="vertical">
        <a-form-item label="单选题数量" field="single">
          <a-input-number v-model="form.single" :min="0" :max="50" />
        </a-form-item>
        <a-form-item label="多选题数量" field="multi">
          <a-input-number v-model="form.multi" :min="0" :max="50" />
        </a-form-item>
        <a-form-item label="判断题数量" field="judge">
          <a-input-number v-model="form.judge" :min="0" :max="50" />
        </a-form-item>
      </a-form>
      <div class="tip" style="color: #86909c; font-size: 12px;">
        * 设置为 0 则不包含该题型，创建后将无法修改。
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { IconPlus } from '@arco-design/web-vue/es/icon';
import { Message } from '@arco-design/web-vue';

const router = useRouter();
const rooms = ref<any[]>([]);
let refreshTimer: any = null;

// 弹窗控制
const showConfigModal = ref(false);
const form = reactive({
  single: 10,
  multi: 5,
  judge: 5
});

onMounted(() => {
  fetchRooms();
  refreshTimer = setInterval(fetchRooms, 3000);
});

onUnmounted(() => {
  if (refreshTimer) clearInterval(refreshTimer);
});

const fetchRooms = async () => {
  try {
    const res = await axios.get('/api/rooms');
    rooms.value = res.data.data || [];
  } catch (e) {
    console.error(e);
  }
};

// 确认创建
const confirmCreate = async () => {
  try {
    // 获取当前用户ID (可选)
    const userId = localStorage.getItem('user_id') || 0;
    
    const res = await axios.post(`/api/rooms?user_id=${userId}`, {
      single: form.single,
      multi: form.multi,
      judge: form.judge
    });
    
    const roomId = res.data.room_id;
    Message.success('房间创建成功！');
    showConfigModal.value = false;
    enterRoom(roomId);
  } catch (e) {
    console.error(e);
    Message.error('创建失败');
  }
};

const handleJoinRoom = (roomId: string) => {
  enterRoom(roomId);
};

const enterRoom = (roomId: string) => {
  router.push({ name: 'coop-room', params: { id: roomId } });
};
</script>

<style scoped>
.lobby-container { max-width: 1000px; margin: 40px auto; padding: 20px; }
.lobby-header { display: flex; justify-content: space-between; align-items: center; }
.title-block h2 { margin-bottom: 5px; }
.title-block p { color: #86909c; margin: 0; }

.room-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 24px; margin-top: 20px; }
.room-card { transition: transform 0.2s; }
.room-card:hover { transform: translateY(-4px); }
.room-card.full { opacity: 0.7; }

.room-content { display: flex; justify-content: space-between; align-items: center; margin: 20px 0; }
.count-info { color: #86909c; font-size: 14px; }
.empty-state { padding: 60px 0; }
</style>