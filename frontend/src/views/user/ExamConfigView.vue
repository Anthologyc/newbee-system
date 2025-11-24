<template>
  <div class="config-container">
    <a-card title="全真模拟考试设置" class="config-card">
      <a-form :model="form" layout="vertical">
        
        <a-form-item label="考试总题量">
          <a-slider v-model="form.totalCount" :min="10" :max="100" :step="10" show-input />
        </a-form-item>

        <a-divider orientation="center">题型与分值配置</a-divider>

        <a-row :gutter="20">
          <!-- 单选 -->
          <a-col :span="8">
            <a-form-item label="单选数量">
              <a-input-number v-model="form.singleCount" :min="0" />
            </a-form-item>
            <a-form-item label="单选分值">
              <a-input-number v-model="form.singleScore" :min="1" />
            </a-form-item>
          </a-col>
          
          <!-- 多选 -->
          <a-col :span="8">
             <a-form-item label="多选数量">
              <a-input-number v-model="form.multiCount" :min="0" />
            </a-form-item>
            <a-form-item label="多选分值">
              <a-input-number v-model="form.multiScore" :min="1" />
            </a-form-item>
          </a-col>
          
          <!-- 判断 -->
          <a-col :span="8">
             <a-form-item label="判断数量">
              <a-input-number v-model="form.judgeCount" :min="0" />
            </a-form-item>
            <a-form-item label="判断分值">
              <a-input-number v-model="form.judgeScore" :min="1" />
            </a-form-item>
          </a-col>
        </a-row>

        <div class="score-info">
          <p>总分预估：<span class="score">{{ estimatedScore }}</span> 分</p>
          <p v-if="!isCountValid" class="error-text">⚠️ 题型数量之和 ({{ sumCount }}) 必须等于总题量 ({{ form.totalCount }})</p>
        </div>

        <a-button type="primary" long size="large" @click="startExam" :disabled="!isCountValid">
          开始考试
        </a-button>
      </a-form>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { reactive, computed } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();

const form = reactive({
  totalCount: 50,
  singleCount: 20, singleScore: 2,
  multiCount: 20, multiScore: 4,
  judgeCount: 10, judgeScore: 2,
});

const sumCount = computed(() => form.singleCount + form.multiCount + form.judgeCount);
const isCountValid = computed(() => sumCount.value === form.totalCount);

const estimatedScore = computed(() => {
  return form.singleCount * form.singleScore + 
         form.multiCount * form.multiScore + 
         form.judgeCount * form.judgeScore;
});

const startExam = () => {
  router.push({
    name: 'exam-paper',
    query: {
      single: form.singleCount, sScore: form.singleScore,
      multi: form.multiCount, mScore: form.multiScore,
      judge: form.judgeCount, jScore: form.judgeScore
    }
  });
};
</script>

<style scoped>
.config-container { max-width: 800px; margin: 40px auto; padding: 20px; }
.score-info { text-align: center; margin: 20px 0; background: #f7f8fa; padding: 15px; border-radius: 8px; }
.score { font-size: 24px; color: #165dff; font-weight: bold; }
.error-text { color: #f53f3f; margin-top: 5px; }
</style>