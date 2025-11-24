<template>
  <div class="exam-paper-container">
    <!-- 顶部状态栏 -->
    <div class="status-bar">
      <div class="left-info">
        <span class="mode-badge">全真模拟</span>
        <div class="timer" :class="{ 'danger': timeLeft < 300 }">
          <icon-clock-circle /> {{ formatTime(timeLeft) }}
        </div>
      </div>
      <div class="right-info">
        已答 {{ answeredCount }} / {{ questions.length }}
      </div>
    </div>

    <div class="main-layout" v-if="questions.length > 0">
      <!-- 答题区 -->
      <div class="question-area">
        <a-card class="q-card" :bordered="false">
          <div class="q-header">
            <a-tag color="arcoblue" size="large" class="type-tag">
              {{ getTypeName(currentQuestion?.question_type) }}
            </a-tag>
            <span class="score-tag">({{ getScore(currentQuestion?.question_type) }}分)</span>
            <h3 class="q-title">{{ currentIndex + 1 }}. {{ currentQuestion?.question_text }}</h3>
          </div>

          <!-- 🚀 关键修复：这里必须用 getDisplayOptions(currentQuestion) -->
          <div class="options-list">
            <div v-for="(optContent, optKey) in getDisplayOptions(currentQuestion)" :key="optKey" class="option-item"
              :class="{ 'selected': isSelected(currentQuestion?.id, String(optKey)) }"
              @click="selectOption(currentQuestion?.id, String(optKey), currentQuestion?.question_type)">
              <div class="opt-key">{{ optKey }}</div>
              <div class="opt-val">{{ optContent }}</div>
            </div>
          </div>

          <!-- 底部导航 -->
          <div class="nav-buttons">
            <a-button @click="prev" :disabled="currentIndex === 0">
              <template #icon><icon-left /></template> 上一题
            </a-button>

            <a-button v-if="currentIndex < questions.length - 1" type="primary" @click="next">
              下一题 <template #icon><icon-right /></template>
            </a-button>

            <a-button v-else type="primary" status="danger" @click="() => handleSubmitConfirm(false)">
              <template #icon><icon-send /></template> 交卷
            </a-button>
          </div>
        </a-card>
      </div>

      <!-- 侧边答题卡 -->
      <div class="sidebar">
        <AnswerCard 
          mode="exam" 
          variant="exam"
          :statusList="statusListForCard" 
          :currentIndex="currentIndex" 
          :showAccuracy="false"
          @jump="handleJump" 
        />
        <a-button type="primary" status="danger" long class="submit-btn-bottom" @click="() => handleSubmitConfirm(false)">
          <template #icon><icon-send /></template> 交卷
        </a-button>
      </div>
    </div>

    <div v-else class="loading-box">
      <a-spin dot /> 正在分发试卷...
    </div>

    <!-- 结果弹窗 -->
    <a-modal v-model:visible="showResultModal" title="考试成绩单" fullscreen :footer="false">
      <div class="result-page">
        <div class="score-header">
          <div class="total-score">{{ examResult.score }} <span class="unit">分</span></div>
          <div class="sub-info">
            答对: <span class="green">{{ examResult.correct }}</span>
            答错: <span class="red">{{ examResult.wrong }}</span>
          </div>
        </div>

        <!-- 只显示错题 -->
        <div class="detail-list">
          <h3>错题解析</h3>
          <div v-if="wrongList.length === 0" class="no-wrong">🎉 太棒了！全对！</div>

          <div v-for="(res, idx) in wrongList" :key="idx" class="detail-item wrong">
            <div class="d-title">
              第 {{ getQuestionIndex(res.question_id) }} 题
              <span style="margin-left:10px; font-weight:normal">{{ getQuestionText(res.question_id) }}</span>
            </div>

            <div class="d-ans">
              <div class="ans-row user-err">
                <span class="label">你的选择:</span>
                {{ formatAnswer(res.user_answer) }}
              </div>
              <div class="ans-row corr">
                <span class="label">正确答案:</span>
                {{ formatAnswer(res.correct_answer) }}
              </div>
            </div>

            <!-- 🚀 修复：补上选项展示 -->
            <div class="d-options-review">
              <div v-for="(val, key) in getQuestionOptions(res.question_id)" :key="key" class="review-opt" :class="{
                'is-user': res.user_answer.includes(String(key)),
                'is-corr': isResultOptionCorrect(res.correct_answer, String(key), getQuestionType(res.question_id))
              }">
                <span class="opt-letter">{{ key }}.</span> {{ val }}
              </div>
            </div>
          </div>
        </div>

        <div class="result-actions">
          <a-button type="primary" @click="$router.push({ name: 'practice-home' })">退出考场</a-button>
          <a-button @click="$router.push({ name: 'mistake-book' })">查看错题本</a-button>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRoute } from 'vue-router';
import axios from 'axios';
import { IconClockCircle, IconLeft, IconRight, IconSend } from '@arco-design/web-vue/es/icon';
import { Message, Modal } from '@arco-design/web-vue';
import AnswerCard from '../../components/AnswerCard.vue';

const route = useRoute();

const questions = ref<any[]>([]);
const userAnswers = ref<Record<number, string[]>>({});
const currentIndex = ref(0);
const timeLeft = ref(3600);
let timer: any = null;
const showResultModal = ref(false);
const examResult = ref<any>({ score: 0, correct: 0, wrong: 0, details: [] });

const sScore = Number(route.query.sScore) || 2;
const mScore = Number(route.query.mScore) || 4;
const jScore = Number(route.query.jScore) || 2;

const currentQuestion = computed(() => questions.value[currentIndex.value]);
const answeredCount = computed(() => Object.keys(userAnswers.value).length);

const wrongList = computed(() => {
  return (examResult.value.details || []).filter((item: any) => !item.is_correct);
});

// Status list for AnswerCard component
const statusListForCard = computed(() => {
  return questions.value.map(q => {
    return isAnswered(q.id) ? 'done' : 'unanswered';
  });
});

const handleJump = (index: number) => {
  currentIndex.value = index;
};
// 🚀 辅助函数：判断结果页的选项是否正确（用于高亮）
const isResultOptionCorrect = (correctArr: string[], key: string, type: string) => {
  if (type === 'judgment') {
    // 如果是判断题，correctArr 是 ["对"]，而 key 是 "A"，需要映射
    const val = key === 'A' ? '对' : '错';
    return correctArr.includes(val);
  }
  // 其他题型直接比对
  return correctArr.includes(key);
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
  const { single, multi, judge } = route.query;
  // 🚀 关键修复：调用新的 generate 接口
  try {
    const res = await axios.get(`/api/exam/generate?single=${single || 0}&multi=${multi || 0}&judge=${judge || 0}`);
    questions.value = res.data.data || [];
    startTimer();
  } catch (error) {
    console.error('Failed to generate exam paper:', error);
    Message.error('试卷生成失败');
  }

  // Register keyboard shortcuts
  window.addEventListener('keydown', onGlobalKeydown);
});

onUnmounted(() => { 
  if (timer) clearInterval(timer); 
  // Unregister keyboard shortcuts
  window.removeEventListener('keydown', onGlobalKeydown);
});

const startTimer = () => {
  timer = setInterval(() => {
    if (timeLeft.value > 0) timeLeft.value--;
    else { clearInterval(timer); handleSubmitConfirm(true); }
  }, 1000);
};

const formatTime = (seconds: number) => {
  const m = Math.floor(seconds / 60);
  const s = seconds % 60;
  return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
};

// 🚀 关键修复：统一处理选项显示 (判断题兜底)
const getDisplayOptions = (q: any) => {
  if (!q) return {};
  if (q.question_type === 'judgment') return { 'A': '对', 'B': '错' };
  return q.options || {}; // 加上 || {} 防止报错
};

const selectOption = (qID: number, key: string, type: string) => {
  let current = userAnswers.value[qID] || [];
  if (type === 'multiple_choice') {
    if (current.includes(key)) current = current.filter(k => k !== key);
    else { current.push(key); current.sort(); }
  } else {
    current = [key];
  }
  userAnswers.value = { ...userAnswers.value, [qID]: current };
};

const isSelected = (qID: number, key: string) => {
  return userAnswers.value[qID]?.includes(key);
};

const isAnswered = (qID: number) => {
  const answer = userAnswers.value?.[qID];
  return Array.isArray(answer) && answer.length > 0;
};

const handleSubmitConfirm = (force: boolean = false) => {
  if (typeof force !== 'boolean') force = false;
  if (!force) {
    Modal.confirm({
      title: '确认交卷',
      content: `还有 ${questions.value.length - answeredCount.value} 道题未做，确定要交卷吗？`,
      onOk: submitPaper
    });
  } else {
    submitPaper();
  }
};

const submitPaper = async () => {
  if (timer) clearInterval(timer);
  const userId = Number(localStorage.getItem('user_id') || 0);

  // 🚀 处理判断题 A/B -> 对/错 映射
  const finalAnswers: Record<number, string[]> = {};
  for (const [qidStr, ans] of Object.entries(userAnswers.value)) {
    const qid = Number(qidStr);
    const q = questions.value.find(item => item.id === qid);
    if (q && q.question_type === 'judgment') {
      const map: Record<string, string> = { 'A': '对', 'B': '错' };
      finalAnswers[qid] = ans.map(k => map[k] || k);
    } else {
      finalAnswers[qid] = ans;
    }
  }

  const payload = {
    user_id: userId,
    answers: finalAnswers,
    single_score: sScore,
    multi_score: mScore,
    judge_score: jScore
  };

  try {
    const res = await axios.post('/api/exam/submit', payload);
    examResult.value = res.data;
    showResultModal.value = true;
  } catch (error) {
    console.error('Failed to submit exam:', error);
    Message.error('交卷失败');
  }
};

const getTypeName = (type: string) => {
  const map: Record<string, string> = { 'single_choice': '单选', 'multiple_choice': '多选', 'judgment': '判断' };
  return map[type] || '题';
};

const getScore = (type: string) => {
  if (type === 'single_choice') return sScore;
  if (type === 'multiple_choice') return mScore;
  if (type === 'judgment') return jScore;
  return 2;
};

// 结果页辅助函数
const getQuestionIndex = (qid: number) => questions.value.findIndex(q => q.id === qid) + 1;
const getQuestionText = (qid: number) => questions.value.find(q => q.id === qid)?.question_text || '';
const getQuestionType = (qid: number) => questions.value.find(q => q.id === qid)?.question_type || '';


const getQuestionOptions = (qid: number) => {
  const q = questions.value.find(q => q.id === qid);
  if (!q) return {};
  if (q.question_type === 'judgment') return { 'A': '对', 'B': '错' };
  return q.options;
};

const formatAnswer = (ans: string[]) => {
  if (!ans || ans.length === 0) return '未答';
  return ans.join('、');
};

const next = () => { if (currentIndex.value < questions.value.length - 1) currentIndex.value++; };
const prev = () => { if (currentIndex.value > 0) currentIndex.value--; };
</script>

<style scoped>
.exam-paper-container {
  max-width: 1200px;
  margin: 0 auto;
  padding-top: 20px;
  padding-bottom: 40px;
}

/* ... 样式和之前一致，省略部分 ... */
.status-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #fff;
  padding: 15px 25px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  margin-bottom: 20px;
}

.left-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.mode-badge {
  background: #165dff;
  color: #fff;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.timer {
  display: flex;
  align-items: center;
  gap: 5px;
  font-weight: bold;
  font-size: 16px;
}

.timer.danger {
  color: #f53f3f;
}

.main-layout {
  display: flex;
  gap: 24px;
}

.question-area {
  flex: 1;
}

.sidebar {
  width: 280px;
  flex-shrink: 0;
  position: sticky;
  top: 20px;
  align-self: flex-start;
}

.q-card {
  min-height: 500px;
  border-radius: 8px;
}

.q-header {
  margin-bottom: 20px;
  display: flex;
  align-items: flex-start;
}

.type-tag {
  margin-right: 8px;
  flex-shrink: 0;
}

.score-tag {
  color: #86909c;
  margin-right: 8px;
  flex-shrink: 0;
  font-size: 12px;
  margin-top: 4px;
}

.q-title {
  margin: 0;
  font-size: 18px;
  line-height: 1.5;
}

.options-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 30px;
}

.option-item {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  border: 2px solid #f2f3f5;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.option-item:hover {
  background: #f7f8fa;
}

.option-item.selected {
  border-color: #165dff;
  background: #e8f3ff;
}

.option-item.selected .opt-key {
  background: #165dff;
  color: #fff;
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
  margin-right: 15px;
}

.nav-buttons {
  display: flex;
  justify-content: space-between;
  margin-top: auto;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.submit-btn-bottom {
  margin-top: 15px;
  border-radius: 8px;
}

.result-page {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.score-header {
  text-align: center;
  margin-bottom: 40px;
}

.total-score {
  font-size: 48px;
  font-weight: bold;
  color: #165dff;
}

.unit {
  font-size: 16px;
  color: #86909c;
}

.sub-info {
  margin-top: 10px;
  font-size: 16px;
}

.green {
  color: #00b42a;
  font-weight: bold;
  margin-right: 15px;
}

.red {
  color: #f53f3f;
  font-weight: bold;
}

.detail-list {
  background: #f7f8fa;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 30px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  padding: 12px;
  border-bottom: 1px solid #eee;
  background: #fff;
  margin-bottom: 10px;
  border-radius: 4px;
}

.detail-item.wrong {
  border-left: 4px solid #f53f3f;
}

.badge {
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
  margin-right: 10px;
  color: #fff;
}

.result-actions {
  display: flex;
  justify-content: center;
  gap: 20px;
}

@media (max-width: 768px) {
  .exam-paper-container {
    padding-bottom: 20px;
  }

  .main-layout {
    flex-direction: column;
    gap: 15px;
  }

  .question-area {
    width: 100%;
  }

  .sidebar {
    width: 100%;
    position: static;
    order: 2;
  }

  .status-bar {
    padding: 10px;
    font-size: 12px;
  }
}

/* ... (其他样式) */

.d-options-review {
  margin-top: 15px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.review-opt {
  padding: 8px 12px;
  background: #fff;
  border: 1px solid #eee;
  border-radius: 4px;
  font-size: 14px;
  color: #4e5969;
}

.review-opt.is-user {
  background: #ffece8;
  border-color: #f53f3f;
  color: #f53f3f;
  text-decoration: line-through; /* 选错的划掉 */
}

.review-opt.is-corr {
  background: #e8ffea;
  border-color: #00b42a;
  color: #00b42a;
  font-weight: bold;
  text-decoration: none;
}

/* 如果既选了又是对的（在错题列表里通常不会出现全对的情况，除非部分对） */
/* 这里简化处理：正确答案优先级更高 */
</style>