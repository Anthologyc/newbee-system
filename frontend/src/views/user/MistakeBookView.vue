<template>
  <div class="mistake-container" v-if="questions.length > 0">
    <!-- 顶部进度条 -->
    <div class="progress-bar">
      <a-progress :percent="progressPercent" :color="{ '0%': 'rgb(22,93,255)', '100%': 'rgb(0,180,42)' }"
        :show-text="false" size="small" class="progress-line" />
      <div class="progress-text">{{ currentIndex + 1 }} / {{ questions.length }}</div>
    </div>

    <div class="layout-wrapper">
      <!-- 左侧：题目区域 -->
      <div class="main-column">
        <a-card class="q-card" :bordered="false">
        <!-- 题目区 -->
        <div class="q-header">
          <a-tag color="arcoblue" class="type-tag">
             {{ getTypeName(currentQuestion?.question_type || '') }}
          </a-tag>
          <h3 class="q-title">{{ currentQuestion?.question_text }}</h3>
        </div>

        <!-- 选项区 -->
        <div class="options-list">
          <div 
            v-for="(optContent, optKey) in getDisplayOptions(currentQuestion)" 
            :key="optKey"
            class="option-item"
            :class="getOptionClass(String(optKey))"
            @click="selectOption(String(optKey))"
          >
            <div class="opt-key">{{ optKey }}</div>
            <div class="opt-val">{{ optContent }}</div>
          </div>
        </div>

        <!-- 底部解析 -->
        <transition name="fade">
          <div v-if="showResult" class="analysis-box">
             <a-alert :type="isCorrect ? 'success' : 'error'">
               <template #icon>
                 <icon-check-circle-fill v-if="isCorrect" />
                 <icon-close-circle-fill v-else />
               </template>
               <span v-if="isCorrect">🎉 回答正确！已从错题本移除。</span>
               <span v-else>❌ 回答错误，正确答案是：{{ currentQuestion?.answer?.join('、') }}</span>
             </a-alert>
          </div>
        </transition>

        <!-- 底部导航 -->
        <div class="footer-actions">
           <a-button class="nav-btn" @click="prev" :disabled="currentIndex === 0">
             <template #icon><icon-left /></template>
             <span class="btn-text">上一题</span>
           </a-button>
           
           <a-button 
             v-if="!showResult" 
             type="primary" 
             class="nav-btn next-btn"
             @click="submitAnswer"
             :disabled="currentSelected.length === 0"
           >
             <span class="btn-text">提交验证</span>
           </a-button>

           <a-button 
             v-else 
             type="primary" 
             class="nav-btn next-btn"
             @click="next"
           >
             <span class="btn-text">{{ isCorrect ? '移除并下一题' : '下一题' }}</span>
             <template #icon><icon-right /></template>
           </a-button>
        </div>
      </a-card>
    </div>

    <!-- 右侧：答题卡 -->
    <div class="sidebar-column">
      <div class="sidebar-wrapper">
        <AnswerCard 
          mode="mistake" 
          variant="mistake"
          :statusList="statusList" 
          :currentIndex="currentIndex" 
          :showAccuracy="true"
          @jump="handleJump" 
        />
      </div>
    </div>
  </div>
  </div>

  <div v-else class="empty-state">
    <div class="empty-icon">🎉</div>
    <h3>太棒了！错题本空空如也</h3>
    <p>去刷几套卷子保持手感吧~</p>
    <a-button type="primary" @click="$router.push('/practice/mode')">去刷题</a-button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import axios from 'axios';
import { IconCheckCircleFill, IconCloseCircleFill, IconLeft, IconRight } from '@arco-design/web-vue/es/icon';
import AnswerCard from '../../components/AnswerCard.vue';

interface Question {
  id: number;
  question_type: string;
  question_text: string;
  options: Record<string, string>;
  answer: string[];
}

const questions = ref<Question[]>([]);
const currentIndex = ref(0);
const currentSelected = ref<string[]>([]);
const showResult = ref(false);
const isCorrect = ref(false);
const userId = Number(localStorage.getItem('user_id') || 0);

// 状态管理 - 用于答题卡显示
const statusList = ref<string[]>([]);

const currentQuestion = computed(() => questions.value[currentIndex.value]);

const progressPercent = computed(() => {
  if (questions.value.length === 0) return 0;
  return (currentIndex.value + 1) / questions.value.length;
});

const handleJump = (index: number) => {
  currentIndex.value = index;
  // 重置状态
  currentSelected.value = [];
  showResult.value = false;
  isCorrect.value = false;
};

const onGlobalKeydown = (e: KeyboardEvent) => {
  const target = e.target as HTMLElement;
  
  // Ignore if user is typing in input fields
  if (target.tagName === 'INPUT' || 
      target.tagName === 'TEXTAREA' || 
      target.tagName === 'SELECT' ||
      target.isContentEditable) {
    return;
  }

  // Next question shortcuts
  if (e.key === 'ArrowRight' || e.key === 'PageDown' || e.key === 'n' || e.key === 'N') {
    e.preventDefault();
    next();
  }
  // Previous question shortcuts
  else if (e.key === 'ArrowLeft' || e.key === 'PageUp' || e.key === 'p' || e.key === 'P') {
    e.preventDefault();
    prev();
  }
};

onMounted(async () => {
  await fetchMistakes();
  
  // Register keyboard shortcuts
  window.addEventListener('keydown', onGlobalKeydown);
});

onUnmounted(() => {
  // Unregister keyboard shortcuts
  window.removeEventListener('keydown', onGlobalKeydown);
});

const fetchMistakes = async () => {
  try {
    const res = await axios.get(`/api/mistakes?user_id=${userId}`);
    questions.value = res.data.data || [];
    // 初始化状态列表
    statusList.value = Array.from({ length: questions.value.length }, () => 'unanswered');
  } catch (e) { console.error(e); }
};

// 处理选项显示 (兼容判断题)
const getDisplayOptions = (q: Question | undefined) => {
  if (!q) return {};
  if (q.question_type === 'judgment') {
    return { 'A': '对', 'B': '错' };
  }
  return q.options;
};

const selectOption = (key: string) => {
  if (showResult.value) return;
  if (!currentQuestion.value) return;

  const type = currentQuestion.value.question_type;
  if (type === 'multiple_choice') {
    const idx = currentSelected.value.indexOf(key);
    if (idx > -1) currentSelected.value.splice(idx, 1);
    else { currentSelected.value.push(key); currentSelected.value.sort(); }
  } else {
    currentSelected.value = [key];
  }
};

const submitAnswer = async () => {
  if (!currentQuestion.value) return;
  
  const q = currentQuestion.value;
  const correctArr = q.answer || [];
  let userArr = currentSelected.value;

  // 判断题特殊处理：如果后端存的是 ["对"]，前端选的是 ["A"]，需要映射
  if (q.question_type === 'judgment') {
    const map: Record<string, string> = { 'A': '对', 'B': '错' };
    userArr = userArr.map(k => map[k] || k);
  }

  const isRight = JSON.stringify(correctArr.sort()) === JSON.stringify(userArr.sort());
  
  showResult.value = true;
  isCorrect.value = isRight;
  
  // 更新状态列表
  statusList.value[currentIndex.value] = isRight ? 'correct' : 'wrong';

  if (isRight) {
    try {
      await axios.post('/api/mistakes/remove', {
        user_id: userId,
        question_id: q.id
      });
    } catch {
      // Silently ignore error - removal will be retried on next correct answer
    }
  }
};

const next = () => {
  // 如果答对了，从列表中移除当前题
  if (isCorrect.value) {
    questions.value.splice(currentIndex.value, 1);
    if (currentIndex.value >= questions.value.length) {
      currentIndex.value = Math.max(0, questions.value.length - 1);
    }
  } else {
    if (currentIndex.value < questions.value.length - 1) currentIndex.value++;
  }
  
  // 重置状态
  currentSelected.value = [];
  showResult.value = false;
  isCorrect.value = false;
};

const prev = () => {
  if (currentIndex.value > 0) {
    currentIndex.value--;
    currentSelected.value = [];
    showResult.value = false;
    isCorrect.value = false;
  }
};

const getOptionClass = (key: string) => {
  if (!showResult.value) {
    return currentSelected.value.includes(key) ? 'opt-active' : '';
  }
  
  const q = currentQuestion.value;
  if (!q) return '';
  
  let correctArr = q.answer || [];
  if (q.question_type === 'judgment') {
      correctArr = correctArr.map(ans => ans === '对' ? 'A' : (ans === '错' ? 'B' : ans));
  }

  if (correctArr.includes(key)) return 'opt-correct';
  if (currentSelected.value.includes(key) && !correctArr.includes(key)) return 'opt-wrong';
  return 'opt-disabled';
};

const getTypeName = (type: string) => {
  const map: Record<string, string> = { 'single_choice': '单选', 'multiple_choice': '多选', 'judgment': '判断' };
  return map[type] || '题';
};
</script>

<style scoped>
/* --- 布局与容器 --- */
.mistake-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  min-height: 90vh;
  display: flex;
  flex-direction: column;
}

.layout-wrapper {
  display: flex;
  gap: 24px;
  flex: 1;
}

.main-column {
  flex: 1;
  min-width: 0;
}

.sidebar-column {
  width: 280px;
  flex-shrink: 0;
}

/* --- 进度条 --- */
.progress-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 15px;
}

.progress-line {
  flex: 1;
}

.progress-text {
  font-size: 12px;
  color: #86909c;
  font-variant-numeric: tabular-nums;
}

/* --- 题目卡片 --- */
.q-card {
  height: 100%;
  display: flex;
  flex-direction: column;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.q-header {
  display: flex;
  align-items: flex-start;
  margin-bottom: 20px;
}

.type-tag {
  margin-right: 10px;
  margin-top: 3px;
  flex-shrink: 0;
}

.q-title {
  margin: 0;
  font-size: 18px;
  line-height: 1.6;
  color: #1d2129;
  font-weight: 600;
}

/* --- 选项列表 --- */
.options-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 20px;
}

.option-item {
  display: flex;
  align-items: center;
  padding: 14px 16px;
  border: 2px solid #f2f3f5;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  background: #fff;
}

.option-item:active {
  transform: scale(0.99);
}

.opt-key {
  width: 30px;
  height: 30px;
  background: #f2f3f5;
  color: #4e5969;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  margin-right: 12px;
  font-size: 14px;
}

.opt-val {
  flex: 1;
  font-size: 15px;
  color: #1d2129;
  line-height: 1.4;
}

/* 状态样式 */
.opt-active {
  border-color: #165dff;
  background-color: #e8f3ff;
}

.opt-active .opt-key {
  background-color: #165dff;
  color: #fff;
}

.opt-correct {
  border-color: #00b42a;
  background-color: #e8ffea;
}

.opt-correct .opt-key {
  background-color: #00b42a;
  color: #fff;
}

.opt-wrong {
  border-color: #f53f3f;
  background-color: #ffece8;
}

.opt-wrong .opt-key {
  background-color: #f53f3f;
  color: #fff;
}

.opt-disabled {
  opacity: 0.5;
  cursor: default;
}

/* --- 解析盒 --- */
.analysis-box {
  margin-bottom: 20px;
}

/* --- 底部按钮 --- */
.footer-actions {
  margin-top: auto;
  padding-top: 20px;
  border-top: 1px solid #f2f3f5;
  display: flex;
  justify-content: space-between;
}

.nav-btn {
  height: 40px;
  padding: 0 25px;
  border-radius: 20px;
  font-size: 15px;
}

.next-btn {
  min-width: 120px;
}

/* --- 右侧侧边栏 --- */
.sidebar-wrapper {
  position: sticky;
  top: 20px;
}

/* --- 空状态 --- */
.empty-state {
  text-align: center;
  padding: 60px 0;
}

.empty-icon {
  font-size: 60px;
  margin-bottom: 20px;
}

/* --- 📱 移动端适配 (Max Width 768px) --- */
@media (max-width: 768px) {
  .mistake-container {
    padding: 10px;
  }

  .layout-wrapper {
    flex-direction: column;
    gap: 15px;
  }

  .main-column {
    width: 100%;
  }

  .sidebar-column {
    width: 100%;
    order: 2;
  }

  .q-title {
    font-size: 16px;
  }

  .option-item {
    padding: 12px;
  }

  .opt-val {
    font-size: 14px;
  }

  .footer-actions {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    background: #fff;
    padding: 15px 20px;
    box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.05);
    z-index: 100;
  }

  .q-card {
    padding-bottom: 60px;
    border: none;
    box-shadow: none;
  }

  .next-btn {
    flex: 1;
    margin-left: 15px;
  }
}

/* 动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>