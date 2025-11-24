<template>
  <a-card :bordered="false" class="answer-sheet-card">
    <div class="sheet-header">
      <span class="title">答题卡</span>
      <div class="legend">
        <span 
          v-for="(dotClass, idx) in legendDots" 
          :key="idx"
          class="dot" 
          :class="dotClass"
        ></span>
      </div>
    </div>
    
    <div class="grid-container">
      <div
        v-for="(status, index) in statusList"
        :key="index"
        class="circle"
        :class="{
          active: currentIndex === index,
          correct: status === 'correct',
          wrong: status === 'wrong',
          done: status === 'done',
          'exam-filled': variant === 'exam' && status === 'done'
        }"
        @click="$emit('jump', index)"
      >
        {{ index + 1 }}
      </div>
    </div>
  </a-card>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = withDefaults(defineProps<{
  statusList: string[]; 
  currentIndex: number;
  variant?: 'practice' | 'exam' | 'mistake';
}>(), {
  variant: 'practice'
});

defineEmits(['jump']);

// Legend dots based on variant
const legendDots = computed(() => {
  if (props.variant === 'exam') {
    return ['exam-filled', 'unanswered'];
  }
  // practice and mistake use the same legend
  return ['wrong', 'correct', 'unanswered'];
});
</script>

<style scoped>
.answer-sheet-card {
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
  overflow: visible;
  position: relative;
}

.sheet-header {
  display: flex; justify-content: space-between; align-items: center;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #f2f3f5;
}
.title { font-weight: bold; font-size: 15px; }

.legend { display: flex; gap: 8px; }
.dot { width: 8px; height: 8px; border-radius: 50%; }
.dot.correct { background: #00b42a; }
.dot.wrong { background: #f53f3f; }
.dot.unanswered { background: #e5e6eb; }
.dot.exam-filled { background: var(--exam-blue, #0b63ff); }

.grid-container {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 10px;
  max-height: 300px;
  overflow-y: auto;
}

/* 题号圆圈 */
.circle {
  aspect-ratio: 1;
  border-radius: 50%;
  background: #f7f8fa;
  color: #86909c;
  display: flex; align-items: center; justify-content: center;
  cursor: pointer; font-size: 12px;
  transition: all 0.2s;
}

.circle.active { box-shadow: 0 0 0 2px #165dff; color: #165dff; background: #fff; font-weight: bold; }
.circle.correct { background-color: #00b42a; color: white; }
.circle.wrong { background-color: #f53f3f; color: white; }
.circle.exam-filled { background-color: var(--exam-blue, #0b63ff); color: white; }

/* 📱 移动端适配 */
@media (max-width: 768px) {
  .grid-container {
    grid-template-columns: repeat(8, 1fr); /* 手机横排更多 */
    gap: 6px;
  }
  .circle { font-size: 10px; }
}
</style>