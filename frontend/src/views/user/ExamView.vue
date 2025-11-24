<template>
    <div class="exam-container" v-if="questions.length > 0">
        <!-- 顶部进度条 -->
        <div class="progress-bar">
            <a-progress :percent="progressPercent" :color="{ '0%': 'rgb(22,93,255)', '100%': 'rgb(0,180,42)' }"
                :show-text="false" size="small" class="progress-line" />
            <div class="progress-text">{{ currentIndex + 1 }} / {{ questions.length }}</div>
        </div>

        <div class="layout-wrapper">
            <!-- 左侧：题目区域 -->
            <div class="main-column">
                <a-card class="question-card" :bordered="false">
                    <!-- 题目头部 -->
                    <div class="question-header">
                        <a-tag color="arcoblue" class="type-tag">
                            {{ getTypeName(currentQuestion?.question_type || '') }}
                        </a-tag>
                        <h3 class="question-text">{{ currentQuestion?.question_text }}</h3>
                    </div>

                    <!-- 选项区域 -->
                    <div class="options-list">
                        <div v-for="(optContent, optKey) in getDisplayOptions(currentQuestion)" :key="optKey"
                            class="option-item" :class="getOptionClass(String(optKey))"
                            @click="handleOptionClick(String(optKey))">
                            <div class="opt-key">{{ optKey }}</div>
                            <div class="opt-val">{{ optContent }}</div>
                            <!-- 结果图标 -->
                            <div class="opt-icon" v-if="hasCheckedResult">
                                <icon-check-circle-fill v-if="isOptionCorrect(String(optKey))" style="color:#00b42a" />
                                <icon-close-circle-fill v-else-if="isOptionSelectedAndWrong(String(optKey))"
                                    style="color:#f53f3f" />
                            </div>
                        </div>
                    </div>

                    <!-- 错题解析 (仅在判题后且答错时显示) -->
                    <transition name="fade">
                        <div v-if="hasCheckedResult && !isCurrentCorrect" class="analysis-box">
                            <a-alert type="error" title="回答错误">
                                <div class="analysis-content">
                                    正确答案是：<strong>{{ currentQuestion?.answer?.join('、') }}</strong>
                                </div>
                            </a-alert>
                        </div>
                    </transition>

                    <!-- 底部固定操作区 -->
                    <div class="footer-actions">
                        <a-button class="nav-btn" @click="prev" :disabled="currentIndex === 0">
                            <template #icon><icon-left /></template>
                            <span class="btn-text">上一题</span>
                        </a-button>

                        <!-- 🚀 核心按钮：点击进行判题或跳转 -->
                        <a-button type="primary" class="nav-btn next-btn" @click="handleNextAction"
                            :disabled="currentSelected.length === 0">
                            <span class="btn-text">{{ nextButtonText }}</span>
                            <template #icon><icon-right /></template>
                        </a-button>
                    </div>
                </a-card>
            </div>

            <!-- 右侧：答题卡 (移动端会变到底部) -->
            <div class="sidebar-column">
                <div class="sidebar-wrapper">
                    <AnswerCard 
                        :mode="route.params.mode === 'random' ? 'random' : 'sequential'" 
                        :statusList="statusList" 
                        :currentIndex="currentIndex" 
                        :showAccuracy="true"
                        @jump="handleJump" 
                    />
                </div>
            </div>
        </div>
    </div>

    <div v-else class="loading-container">
        <a-spin dot />
        <p>正在加载题库...</p>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import axios from 'axios';
import AnswerCard from '../../components/AnswerCard.vue';
import { 
  IconCheckCircleFill, IconCloseCircleFill, IconLeft, IconRight 
} from '@arco-design/web-vue/es/icon';

interface Question {
  id: number;
  question_type: string; 
  question_text: string;
  options: Record<string, string>;
  answer: string[];
}

const route = useRoute();
const questions = ref<Question[]>([]);
const currentIndex = ref(0);

// 状态管理
const userAnswers = ref<Record<number, string[]>>({}); 
const statusList = ref<string[]>([]); // 'unanswered'|'correct'|'wrong'

// 当前操作状态
const currentSelected = ref<string[]>([]);
const hasCheckedResult = ref(false); 
const isCurrentCorrect = ref(false); 

// 用户ID
const userId = Number(localStorage.getItem('user_id') || 0);

// 计算属性
const currentQuestion = computed(() => questions.value[currentIndex.value]);
const progressPercent = computed(() => (currentIndex.value + 1) / questions.value.length);

const nextButtonText = computed(() => {
  if (!hasCheckedResult.value) return '下一题';
  return '继续';
});

// 监听切题
watch(currentIndex, (newVal) => {
  if (userAnswers.value[newVal]) {
    currentSelected.value = [...userAnswers.value[newVal]];
    const status = statusList.value[newVal];
    if (status && status !== 'unanswered') {
      hasCheckedResult.value = true;
      isCurrentCorrect.value = (status === 'correct');
    } else {
      hasCheckedResult.value = false;
      isCurrentCorrect.value = false;
    }
  } else {
    currentSelected.value = [];
    hasCheckedResult.value = false;
    isCurrentCorrect.value = false;
  }
});

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
    goToNext();
  }
  // Previous question shortcuts
  else if (e.key === 'ArrowLeft' || e.key === 'PageUp' || e.key === 'p' || e.key === 'P') {
    e.preventDefault();
    prev();
  }
};

onMounted(async () => {
  const mode = route.params.mode; 
  let url = mode === 'sequential' ? '/api/practice/sequential' : '/api/practice/random';
  
  if (mode === 'random') {
    const { category, type, count } = route.query;
    const params = new URLSearchParams();
    if (category) params.append('category', String(category));
    if (type) params.append('type', String(type));
    if (count) params.append('count', String(count));
    url += `?${params.toString()}`;
  }
  
  try {
    const res = await axios.get(url);
    const list = Array.isArray(res.data.data) ? res.data.data : [];
    questions.value = list;
    statusList.value = Array.from({ length: list.length }, () => 'unanswered');
  } catch (e) {
    console.error(e);
  }

  // Register keyboard shortcuts
  window.addEventListener('keydown', onGlobalKeydown);
});

onUnmounted(() => {
  // Unregister keyboard shortcuts
  window.removeEventListener('keydown', onGlobalKeydown);
});

// 统一处理选项显示 (修复判断题)
const getDisplayOptions = (q: Question | undefined) => {
  if (!q) return {};
  if (q.question_type === 'judgment') {
    return { 'A': '对', 'B': '错' };
  }
  return q.options;
};

const handleOptionClick = (key: string) => {
  if (hasCheckedResult.value) return; 
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

const handleNextAction = async () => {
  if (!currentQuestion.value) return;

  if (!hasCheckedResult.value) {
    const q = currentQuestion.value;
    const correctArr = q.answer || [];
    let userArr = currentSelected.value;

    // 判断题映射：A->对, B->错
    if (q.question_type === 'judgment') {
        const map: Record<string, string> = { 'A': '对', 'B': '错' };
        userArr = userArr.map(k => map[k] || k);
    }

    const isRight = JSON.stringify(correctArr.sort()) === JSON.stringify(userArr.sort());
    
    isCurrentCorrect.value = isRight;
    hasCheckedResult.value = true;
    statusList.value[currentIndex.value] = isRight ? 'correct' : 'wrong';
    // 这里存原始选择(A/B)以便回显
    userAnswers.value[currentIndex.value] = [...currentSelected.value];

    // 错题录入
    if (!isRight) {
      try {
         await axios.post('/api/mistakes/record', {
           user_id: userId,
           question_id: q.id
         });
      } catch {
        // Silently ignore error - mistake recording is non-critical
      }
    }

    if (isRight) {
      setTimeout(() => { goToNext(); }, 300); 
    }
  } else {
    goToNext();
  }
};

const goToNext = () => {
  if (currentIndex.value < questions.value.length - 1) {
    currentIndex.value++;
  }
};

const prev = () => { if (currentIndex.value > 0) currentIndex.value--; };

const handleJump = (idx: number) => { currentIndex.value = idx; };

const isOptionCorrect = (key: string) => {
    let correctArr = currentQuestion.value?.answer || [];
    if (currentQuestion.value?.question_type === 'judgment') {
        correctArr = correctArr.map(ans => ans === '对' ? 'A' : (ans === '错' ? 'B' : ans));
    }
    return correctArr.includes(key);
};

const isOptionSelectedAndWrong = (key: string) => {
  const userArr = userAnswers.value[currentIndex.value] || currentSelected.value;
  return userArr.includes(key) && !isOptionCorrect(key);
};

const getOptionClass = (key: string) => {
  if (!hasCheckedResult.value) {
    return currentSelected.value.includes(key) ? 'opt-active' : '';
  }
  if (isOptionCorrect(key)) return 'opt-correct';
  if (isOptionSelectedAndWrong(key)) return 'opt-wrong';
  return 'opt-disabled';
};

const getTypeName = (type: string) => {
  const map: Record<string, string> = { 'single_choice': '单选', 'multiple_choice': '多选', 'judgment': '判断' };
  return map[type] || '题';
};
</script>

<style scoped>
/* --- 布局与容器 --- */
.exam-container {
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

/* min-width: 0 防止 flex 子项溢出 */
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
.question-card {
    height: 100%;
    display: flex;
    flex-direction: column;
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.question-header {
    display: flex;
    align-items: flex-start;
    margin-bottom: 20px;
}

.type-tag {
    margin-right: 10px;
    margin-top: 3px;
    flex-shrink: 0;
}

.question-text {
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

.opt-icon {
    margin-left: 10px;
    font-size: 20px;
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

.analysis-content {
    font-size: 15px;
    color: #1d2129;
}

/* --- 底部按钮 --- */
.footer-actions {
    margin-top: auto;
    /* 推到底部 */
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



/* --- 📱 移动端适配 (Max Width 768px) --- */
@media (max-width: 768px) {
    .exam-container {
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
        /* 放到下面 */
    }



    .question-text {
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

    /* 给底部留出按钮的空间 */
    .question-card {
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