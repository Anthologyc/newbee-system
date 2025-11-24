<template>
  <div class="filter-container">
    <a-card class="filter-card" title="定制你的练习" :bordered="false">
      <a-form :model="form" layout="vertical">
        
        <!-- 分类选择 -->
        <a-form-item label="题目分类">
          <a-select v-model="form.category" placeholder="不限 (所有分类)" allow-clear>
            <a-option v-for="cat in categories" :key="cat" :value="cat">{{ cat || '未分类' }}</a-option>
          </a-select>
        </a-form-item>

        <!-- 题型选择 -->
        <a-form-item label="题目类型">
          <a-radio-group v-model="form.type" type="button" size="large">
            <a-radio value="">全部</a-radio>
            <a-radio value="single_choice">单选</a-radio>
            <a-radio value="multiple_choice">多选</a-radio>
            <a-radio value="judgment">判断</a-radio>
          </a-radio-group>
        </a-form-item>

        <!-- 题目数量 -->
        <a-form-item label="练习题量">
          <a-slider v-model="form.count" :min="10" :max="100" :step="10" show-input />
        </a-form-item>

        <div class="actions">
          <a-button type="primary" size="large" long @click="startPractice">开始挑战</a-button>
        </div>
      </a-form>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { Message } from '@arco-design/web-vue';

const router = useRouter();
const categories = ref<string[]>([]);

const form = reactive({
  category: '',
  type: '',
  count: 50
});

onMounted(async () => {
  try {
    const res = await axios.get('/api/categories');
    categories.value = res.data.data || [];
  } catch (e) {
    console.error('Failed to fetch categories');
  }
});

const startPractice = () => {
  // 跳转到刷题页，带上参数
  router.push({
    name: 'practice-exam',
    params: { mode: 'random' }, // 只有随机模式才支持这种高级筛选
    query: {
      category: form.category,
      type: form.type,
      count: form.count
    }
  });
};
</script>

<style scoped>
.filter-container {
  max-width: 600px;
  margin: 40px auto;
  padding: 20px;
}
.filter-card {
  border-radius: 12px;
  box-shadow: 0 4px 10px rgba(0,0,0,0.05);
}
.actions { margin-top: 30px; }

@media (max-width: 768px) {
  .filter-container { margin: 10px auto; padding: 10px; }
}
</style>