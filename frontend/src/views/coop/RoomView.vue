<template>
  <div class="room-waiting-container">
    <div class="room-header">
      <a-button @click="leaveRoom">
        <template #icon><icon-left /></template> 离开房间
      </a-button>
      <h2>房间 #{{ roomId }}</h2>
      <a-tag color="blue">等待开始</a-tag>
    </div>

    <div class="players-area">
      <!-- 自己 -->
      <div class="player-card mine">
        <a-avatar :size="80" :style="{ backgroundColor: '#165dff' }">我</a-avatar>
        <div class="player-name">我 ({{ myId }})</div>
        <div class="status">
           <a-tag :color="myReady ? 'green' : 'gray'">{{ myReady ? '已准备' : '未准备' }}</a-tag>
        </div>
      </div>

      <div class="vs-icon">VS</div>

      <!-- 队友 -->
      <div class="player-card opponent">
        <template v-if="opponentId">
           <a-avatar :size="80" :style="{ backgroundColor: '#ff7d00' }">友</a-avatar>
           <div class="player-name">队友 ({{ opponentId }})</div>
           <div class="status">
             <a-tag :color="opponentReady ? 'green' : 'gray'">{{ opponentReady ? '已准备' : '未准备' }}</a-tag>
           </div>
        </template>
        <template v-else>
           <div class="waiting-placeholder">
             <a-spin dot />
             <p>等待队友加入...</p>
           </div>
        </template>
      </div>
    </div>

    <div class="action-bar">
      <!-- 开始按钮 (只有全员准备且是房主或任意一人时才显示，这里简化为只要全员准备) -->
      <a-button 
        v-if="allReady"
        type="primary" 
        size="large" 
        status="danger"
        @click="startGame"
      >
        🚀 开始考试
      </a-button>

      <!-- 准备按钮 (如果还没全员准备) -->
      <a-button 
        v-else
        type="primary" 
        size="large" 
        :status="myReady ? 'success' : 'normal'"
        @click="toggleReady"
      >
        {{ myReady ? '取消准备' : '准备开始' }}
      </a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useSocketStore } from '../../stores/socket';
import { IconLeft } from '@arco-design/web-vue/es/icon';
import { Message } from '@arco-design/web-vue';

const route = useRoute();
const router = useRouter();
const socketStore = useSocketStore();

const roomId = route.params.id as string;
const myId = Number(localStorage.getItem('user_id') || 0);

const myReady = ref(false);
const opponentId = ref<number | null>(null);
const opponentReady = ref(false);
const allReady = ref(false);

onMounted(() => {
  // 确保连接
  if (!socketStore.isConnected) {
    socketStore.connect(myId);
  }

  // 注册监听
  socketStore.on('user_joined', handleUserJoined);
  socketStore.on('user_left', handleUserLeft);
  socketStore.on('ready_update', handleReadyUpdate);
  socketStore.on('game_start', handleGameStart);

  // 延迟发送加入，确保 socket 已连接
  setTimeout(() => {
    socketStore.send('join_room', { room_id: roomId });
  }, 500);
});

onUnmounted(() => {
  socketStore.off('user_joined', handleUserJoined);
  socketStore.off('user_left', handleUserLeft);
  socketStore.off('ready_update', handleReadyUpdate);
  socketStore.off('game_start', handleGameStart);
  
  socketStore.send('leave_room', {});
});

const leaveRoom = () => {
  router.push({ name: 'coop-lobby' });
};

const toggleReady = () => {
  // 发送状态，等后端广播 ready_update 再更新本地
  socketStore.send('ready', { room_id: roomId, status: !myReady.value });
};

const startGame = () => {
  socketStore.send('start_game', { room_id: roomId });
};

// --- Handlers ---

const handleUserJoined = (payload: any) => {
  console.log('User Joined Payload:', payload);
  // 如果进来的不是自己，那就是队友
  if (payload.user_id && payload.user_id !== myId) {
    opponentId.value = payload.user_id;
    Message.info(`队友 (${payload.user_id}) 加入了房间`);
  }
  // 如果 payload 包含 count >= 2，也可以辅助判断
};

const handleUserLeft = (payload: any) => {
  if (payload.user_id === opponentId.value) {
    opponentId.value = null;
    opponentReady.value = false;
    allReady.value = false;
    Message.warning('队友离开了房间');
  }
};

const handleReadyUpdate = (statusMap: any) => {
  console.log('Ready Update:', statusMap);
  // statusMap key 是 userId (string), value 是 boolean
  
  // 1. 更新自己
  if (statusMap[myId] !== undefined) {
    myReady.value = statusMap[myId];
  }

  // 2. 更新队友 & 确认队友ID (防止 user_joined 漏掉)
  let opId: string | null = null;
  for (const uidStr of Object.keys(statusMap)) {
    if (Number(uidStr) !== myId) {
      opId = uidStr;
      break;
    }
  }

  if (opId) {
    opponentId.value = Number(opId);
    opponentReady.value = statusMap[opId];
  }

  // 3. 判断是否全员准备 (需至少2人)
  const userIds = Object.keys(statusMap);
  const readyCount = Object.values(statusMap).filter(r => r).length;
  
  if (userIds.length >= 2 && readyCount === userIds.length) {
    allReady.value = true;
  } else {
    allReady.value = false;
  }
};

const handleGameStart = (payload: any) => {
  localStorage.setItem('coop_questions', JSON.stringify(payload.questions));
  Message.success('游戏开始！');
  router.push({ name: 'coop-exam', params: { id: roomId } });
};
</script>

<style scoped>
.room-waiting-container { max-width: 800px; margin: 40px auto; padding: 20px; text-align: center; }
.room-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 50px; }
.players-area { display: flex; justify-content: center; align-items: center; gap: 50px; margin-bottom: 50px; }
.player-card { width: 200px; padding: 30px; border: 1px solid #eee; border-radius: 12px; background: #fff; }
.player-name { font-size: 18px; font-weight: bold; margin: 15px 0; }
.vs-icon { font-size: 30px; font-weight: bold; color: #86909c; font-style: italic; }
.waiting-placeholder { color: #86909c; padding: 20px 0; }
</style>