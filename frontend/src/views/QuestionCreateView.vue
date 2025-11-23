<template>
  <a-card :title="isEditMode ? '编辑题目' : '录入题目'">
    <a-tabs v-model="tab">
      <a-tab-pane key="form" title="表单录入">
        <a-form :model="form" style="max-width:500px;">
          <a-form-item label="类别">
            <a-input v-model="form.category" />
          </a-form-item>
          <a-form-item label="题型">
            <a-select v-model="form.question_type">
              <a-option value="single_choice">单选题</a-option>
              <a-option value="multiple_choice">多选题</a-option>
              <a-option value="judgment">判断题</a-option>
            </a-select>
          </a-form-item>
          <a-form-item label="题目内容">
            <a-textarea v-model="form.question_text" :auto-size="{minRows:3,maxRows:5}"/>
          </a-form-item>
          
          <a-form-item label="选项" v-if="form.question_type !== 'judgment'">
            <div v-for="k in ['A','B','C','D']" :key="k" style="display:flex;align-items:center;gap:8px;margin-bottom:8px;width:100%">
              <span style="font-weight:bold">{{ k }}.</span>
              <a-input v-model="form.options[k]" placeholder="选项内容" />
            </div>
          </a-form-item>
          
          <a-form-item label="答案">
            <a-input v-model="answerStr" placeholder="多选用逗号分隔，如 A,B" />
          </a-form-item>
          
          <div style="display:flex; justify-content:flex-end; gap: 10px;">
            <a-button @click="router.back()">取消</a-button>
            <a-button type="primary" :loading="loading" @click="submit">
              {{ isEditMode ? '保存修改' : '立即创建' }}
            </a-button>
          </div>
        </a-form>
      </a-tab-pane>
      
      <!-- 编辑模式下隐藏导入功能 -->
      <a-tab-pane key="json" title="批量 JSON 导入" v-if="!isEditMode">
        <a-textarea v-model="jsonStr" rows="8" placeholder="贴上题目 JSON数组，如：[{}]" />
        <a-button type="primary" long @click="importJson" style="margin-top:8px;">导入</a-button>
      </a-tab-pane>
    </a-tabs>
  </a-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue';
import { Message } from '@arco-design/web-vue';
import axios from 'axios';
import { useRouter, useRoute } from 'vue-router';

const router = useRouter();
const route = useRoute();

const tab = ref('form');
const loading = ref(false);
const jsonStr = ref('');

// 判断是否为编辑模式
const isEditMode = computed(() => !!route.params.id);
const questionId = route.params.id;

const form = reactive({
  category: '',
  question_type: 'single_choice',
  question_text: '',
  options: { A: '', B: '', C: '', D: '' } as Record<string, string>,
  answer: [] as string[]
});
const answerStr = ref('');

// 初始化：如果是编辑模式，拉取数据
onMounted(async () => {
  if (isEditMode.value) {
    try {
      const res = await axios.get(`/api/questions/${questionId}`);
      const data = res.data;
      
      // 回显数据
      form.category = data.category;
      form.question_type = data.question_type;
      form.question_text = data.question_text;
      // 保证 options 存在
      form.options = data.options || { A: '', B: '', C: '', D: '' };
      form.answer = data.answer || [];
      
      // 处理答案显示字符串
      answerStr.value = form.answer.join(',');
    } catch (e) {
      Message.error('获取题目详情失败');
    }
  }
});

const submit = async () => {
  loading.value = true;
  // 处理答案数组
  form.answer = answerStr.value.split(/[,，]/).map(i => i.trim().toUpperCase()).filter(i => i);
  
  try {
    if (isEditMode.value) {
      // 编辑 -> PUT
      await axios.put(`/api/questions/${questionId}`, form);
      Message.success('修改成功');
    } else {
      // 新增 -> POST
      await axios.post('/api/questions', form);
      Message.success('创建成功');
    }
    router.push({ name: 'question-list' });
  } catch (e) {
    Message.error('操作失败');
  } finally {
    loading.value = false;
  }
};

const importJson = async () => {
  try {
    const arr = JSON.parse(jsonStr.value);
    if (!Array.isArray(arr)) throw new Error('不是数组');
    for (const q of arr) {
      await axios.post('/api/questions', q);
    }
    Message.success('批量导入成功');
    router.push({ name: 'question-list' });
  } catch (e: any) {
    Message.error(e.message || '格式错误');
  }
};
</script>