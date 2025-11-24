<template>
  <div class="answer-card-wrapper">
    <!-- Answer Sheet Grid -->
    <AnswerSheet 
      :statusList="statusList" 
      :currentIndex="currentIndex" 
      :variant="variant"
      @jump="handleJump" 
    />
    
    <!-- Statistics Card -->
    <a-card class="stats-card" :bordered="false">
      <div class="stats-content">
        <!-- Mode Badge -->
        <div class="stat-item mode-item">
          <span class="stat-label">模式：</span>
          <span class="stat-value mode-badge">{{ modeIcon }} {{ modeText }}</span>
        </div>
        
        <!-- Answered Count -->
        <div class="stat-item">
          <span class="stat-label">已答：</span>
          <span class="stat-value">{{ answeredCount }} 题</span>
        </div>
        
        <!-- Accuracy (optional - hidden for exam mode) -->
        <div v-if="showAccuracy" class="stat-item">
          <span class="stat-label">正确率：</span>
          <span class="stat-value">{{ accuracy }}%</span>
        </div>
      </div>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import AnswerSheet from './AnswerSheet.vue';

interface Props {
  mode: 'random' | 'sequential' | 'exam' | 'mistake';
  statusList: string[];
  currentIndex: number;
  showAccuracy?: boolean;
  variant?: 'practice' | 'exam' | 'mistake';
}

const props = withDefaults(defineProps<Props>(), {
  showAccuracy: true,
  variant: 'practice'
});

const emit = defineEmits<{
  jump: [index: number];
}>();

const modeIcon = computed(() => {
  const icons = {
    random: '🎲',
    sequential: '📄',
    exam: '📝',
    mistake: '❗'
  };
  return icons[props.mode] || '📚';
});

const modeText = computed(() => {
  const texts = {
    random: '随机练习',
    sequential: '顺序练习',
    exam: '全真模拟',
    mistake: '错题本'
  };
  return texts[props.mode] || '练习';
});

const answeredCount = computed(() => {
  return props.statusList.filter(s => s !== 'unanswered').length;
});

const accuracy = computed(() => {
  if (answeredCount.value === 0) return 0;
  const correctCount = props.statusList.filter(s => s === 'correct').length;
  return Math.round((correctCount / answeredCount.value) * 100);
});

const handleJump = (index: number) => {
  emit('jump', index);
};
</script>

<style scoped>
.answer-card-wrapper {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.stats-card {
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  overflow: visible;
}

.stats-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  font-size: 14px;
  line-height: 1.5;
}

.stat-label {
  color: #86909c;
  margin-right: 4px;
  flex-shrink: 0;
}

.stat-value {
  color: #1d2129;
  font-weight: 500;
}

.mode-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

/* Mobile responsive */
@media (max-width: 768px) {
  .stats-card {
    border-radius: 12px;
  }
  
  .stats-content {
    flex-direction: row;
    flex-wrap: wrap;
    gap: 8px 16px;
  }
  
  .stat-item {
    font-size: 13px;
  }
  
  .mode-item {
    width: 100%;
  }
}
</style>
