<template>
  <div class="mistake-container" v-if="questions.length > 0">
    <div class="progress-bar">
       <div class="progress-text">当前错题: {{ currentIndex + 1 }} / {{ questions.length }}</div>
       <a-tag color="orange" style="margin-left:auto">错题本模式</a-tag>
    </div>

    <a-card class="q-card">
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

      <!-- 底部导航 -->
      <div class="nav-buttons">
         <a-button @click="prev" :disabled="currentIndex === 0">上一题</a-button>
         
         <a-button 
           v-if="!showResult" 
           type="primary" 
           @click="submitAnswer"
           :disabled="currentSelected.length === 0"
         >
           提交验证
         </a-button>

         <a-button 
           v-else 
           type="primary" 
           @click="next"
         >
           {{ isCorrect ? '移除并下一题' : '下一题' }}
         </a-button>
      </div>
    </a-card>
  </div>

  <div v-else class="empty-state">
    <div class="empty-icon">🎉</div>
    <h3>太棒了！错题本空空如也</h3>
    <p>去刷几套卷子保持手感吧~</p>
    <a-button type="primary" @click="$router.push('/practice/mode')">去刷题</a-button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import axios from 'axios';
import { IconCheckCircleFill, IconCloseCircleFill } from '@arco-design/web-vue/es/icon';

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

const currentQuestion = computed(() => questions.value[currentIndex.value]);

onMounted(async () => {
  await fetchMistakes();
});

const fetchMistakes = async () => {
  try {
    const res = await axios.get(`/api/mistakes?user_id=${userId}`);
    questions.value = res.data.data || [];
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
  let correctArr = q.answer || [];
  let userArr = currentSelected.value;

  // 判断题特殊处理：如果后端存的是 ["对"]，前端选的是 ["A"]，需要映射
  if (q.question_type === 'judgment') {
    const map: Record<string, string> = { 'A': '对', 'B': '错' };
    userArr = userArr.map(k => map[k] || k);
  }

  const isRight = JSON.stringify(correctArr.sort()) === JSON.stringify(userArr.sort());
  
  showResult.value = true;
  isCorrect.value = isRight;

  if (isRight) {
    try {
      await axios.post('/api/mistakes/remove', {
        user_id: userId,
        question_id: q.id
      });
    } catch(e) {}
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
.mistake-container { max-width: 800px; margin: 0 auto; padding: 20px; }
.progress-bar { display: flex; margin-bottom: 20px; align-items: center; }
.q-card { min-height: 400px; border-radius: 8px; }
.q-header { margin-bottom: 20px; display: flex; align-items: flex-start; }
.type-tag { margin-right: 10px; flex-shrink: 0; }
.q-title { margin: 0; font-size: 18px; line-height: 1.5; }

.options-list { display: flex; flex-direction: column; gap: 12px; margin-bottom: 20px; }
.option-item {
  padding: 14px 20px; border: 2px solid #f2f3f5; border-radius: 8px; cursor: pointer; display: flex; align-items: center;
}
.option-item:hover { background: #f7f8fa; }
.opt-key { width: 30px; height: 30px; background: #f2f3f5; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-weight: bold; margin-right: 15px; }

.opt-active { border-color: #165dff; background: #e8f3ff; }
.opt-active .opt-key { background: #165dff; color: #fff; }
.opt-correct { border-color: #00b42a; background: #e8ffea; }
.opt-correct .opt-key { background: #00b42a; color: #fff; }
.opt-wrong { border-color: #f53f3f; background: #ffece8; }
.opt-wrong .opt-key { background: #f53f3f; color: #fff; }
.opt-disabled { opacity: 0.6; cursor: default; }

.analysis-box { margin-bottom: 20px; }
.nav-buttons { display: flex; justify-content: space-between; margin-top: auto; }

.empty-state { text-align: center; padding: 60px 0; }
.empty-icon { font-size: 60px; margin-bottom: 20px; }
</style>