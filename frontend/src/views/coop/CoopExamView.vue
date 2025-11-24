<template>
  <div class="coop-exam-container">
    <!-- 左侧：答题区 (复用逻辑) -->
    <div class="exam-area">
      <div class="header-bar">
        <h3>协作考试 ({{ currentIndex + 1 }} / {{ questions.length }})</h3>
        <div class="timer">{{ formatTime(timeLeft) }}</div>
      </div>

      <div class="question-box" v-if="currentQuestion">
        <div class="q-title">
          <a-tag>{{ getTypeName(currentQuestion.question_type) }}</a-tag>
          {{ currentQuestion.question_text }}
        </div>
        
        <div class="options-list">
          <div 
            v-for="(optContent, optKey) in getDisplayOptions(currentQuestion)" 
            :key="optKey"
            class="option-item"
            :class="{ 'selected': userAnswers[currentQuestion.id]?.includes(String(optKey)) }"
            @click="selectOption(currentQuestion.id, String(optKey))"
          >
            <span class="key">{{ optKey }}.</span> {{ optContent }}
          </div>
        </div>

        <div class="nav-btns">
          <a-button @click="prev" :disabled="currentIndex===0">上一题</a-button>
          <a-button type="primary" @click="next">{{ isLast ? '交卷' : '下一题' }}</a-button>
        </div>
      </div>
    </div>

    <!-- 右侧：队友状态 -->
    <div class="sidebar-area">
      <a-card title="实时战况">
        <div class="player-status mine">
          <a-avatar :style="{ backgroundColor: '#165dff' }">我</a-avatar>
          <div class="info">
            <div class="name">我</div>
            <a-progress :percent="myProgress" status="success" />
            <div class="text">已答 {{ answeredCount }} / {{ questions.length }}</div>
          </div>
        </div>

        <a-divider />

        <div class="player-status opponent">
          <a-avatar :style="{ backgroundColor: '#ff7d00' }">友</a-avatar>
          <div class="info">
            <div class="name">队友</div>
            <!-- 队友进度条 -->
            <a-progress :percent="opponentProgress" status="warning" />
            <div class="text">已答 {{ opponentAnswered }} / {{ questions.length }}</div>
          </div>
        </div>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useSocketStore } from '../../stores/socket';
import { Message, Modal } from '@arco-design/web-vue';

const route = useRoute();
const router = useRouter();
const socketStore = useSocketStore();

const questions = ref<any[]>([]);
const currentIndex = ref(0);
const userAnswers = ref<Record<number, string[]>>({});
const timeLeft = ref(1800); // 30分钟

const waitingResult = ref(false);
const opponentFinished = ref(false);
let timer: any = null;

// 队友状态
const opponentAnswered = ref(0);
const myId = Number(localStorage.getItem('user_id') || 0);

const currentQuestion = computed(() => questions.value[currentIndex.value]);
const answeredCount = computed(() => Object.keys(userAnswers.value).length);
const isLast = computed(() => currentIndex.value === questions.value.length - 1);

const myProgress = computed(() => questions.value.length ? answeredCount.value / questions.value.length : 0);
const opponentProgress = computed(() => questions.value.length ? opponentAnswered.value / questions.value.length : 0);

onMounted(() => {
  // 1. 加载题目
  const data = localStorage.getItem('coop_questions');
  if (data) {
    questions.value = JSON.parse(data);
    startTimer();
  } else {
    Message.error('试卷丢失');
    router.push({ name: 'coop-lobby' });
  }

  // 2. 监听队友进度
  socketStore.on('progress_update', handleProgressUpdate);
    socketStore.on('user_finished', handleUserFinished);
  socketStore.on('game_over', handleGameOver);
});

onUnmounted(() => {
  if (timer) clearInterval(timer);
  socketStore.off('progress_update', handleProgressUpdate);
    socketStore.off('user_finished', handleUserFinished);
  socketStore.off('game_over', handleGameOver);
});

// 监听我的答题进度，发送给队友
watch(answeredCount, (newVal) => {
  socketStore.send('progress_update', {
    room_id: route.params.id,
    answered: newVal,
    total: questions.value.length
  });
});

const handleProgressUpdate = (payload: any) => {
  // 如果消息不是自己的，那就是队友的
  if (payload.user_id !== myId && payload.answered !== undefined) {
    opponentAnswered.value = payload.answered;
  }
};

const startTimer = () => {
  timer = setInterval(() => {
    if (timeLeft.value > 0) timeLeft.value--;
    else submitExam(); // 超时自动交卷
  }, 1000);
};

// 选项逻辑 (简化版)
const selectOption = (qid: number, key: string) => {
  // 假设全是单选，简化逻辑
  userAnswers.value = { ...userAnswers.value, [qid]: [key] };
};

const next = () => {
  if (isLast.value) {
    Modal.confirm({ title: '交卷确认', content: '确定要交卷吗？', onOk: submitExam });
  } else {
    currentIndex.value++;
  }
};
const prev = () => { if (currentIndex.value > 0) currentIndex.value--; };

const submitExam = () => {
  socketStore.send('submit_exam', {
    room_id: route.params.id,
    answers: userAnswers.value
  });
  
  waitingResult.value = true;
  // 显示一个 Loading 遮罩或提示
  Modal.info({
    title: '已交卷',
    content: '请等待队友完成考试...',
    hideCancel: true,
    okText: '等待中...',
    okButtonProps: { disabled: true }
  });
};

const handleUserFinished = (payload: any) => {
  if (payload.user_id !== myId) {
    opponentFinished.value = true;
    opponentAnswered.value = questions.value.length; // 进度条拉满
    Message.info('队友已交卷！');
  }
};

const handleGameOver = (payload: any) => {
  console.log("收到 Game Over:", payload); // 调试日志
  
  // 1. 强制关闭 Modal (不管是 Info 还是 Confirm)
  Modal.clear();
  
  // 2. 存数据
  localStorage.setItem('coop_report', JSON.stringify(payload.report));
  
  // 3. 强制跳转 (用 replace 防止回退)
  router.replace({ name: 'coop-result', params: { id: route.params.id } });
};
// helper
const getDisplayOptions = (q: any) => q?.question_type === 'judgment' ? { 'A': '对', 'B': '错' } : q?.options;
const getTypeName = (t: string) => ({ single_choice: '单选', multiple_choice: '多选', judgment: '判断' }[t] || '题');
const formatTime = (s: number) => `${Math.floor(s/60)}:${(s%60).toString().padStart(2,'0')}`;
</script>

<style scoped>
.coop-exam-container { display: flex; height: 100vh; padding: 20px; gap: 20px; background: #f7f8fa; }
.exam-area { flex: 1; background: #fff; border-radius: 8px; padding: 20px; display: flex; flex-direction: column; }
.sidebar-area { width: 300px; }

.header-bar { display: flex; justify-content: space-between; margin-bottom: 20px; font-size: 18px; font-weight: bold; }
.timer { color: #165dff; }

.question-box { flex: 1; }
.q-title { font-size: 18px; margin-bottom: 20px; }
.options-list { display: flex; flex-direction: column; gap: 12px; }
.option-item { padding: 15px; border: 1px solid #eee; border-radius: 6px; cursor: pointer; }
.option-item.selected { background: #e8f3ff; border-color: #165dff; color: #165dff; }

.nav-btns { margin-top: 30px; display: flex; justify-content: space-between; }

.player-status { display: flex; gap: 15px; margin-bottom: 20px; }
.info { flex: 1; }
.name { font-weight: bold; margin-bottom: 5px; }
.text { font-size: 12px; color: #86909c; margin-top: 5px; }
</style>