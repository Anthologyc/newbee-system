<template>
  <a-card title="录入题目">
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
            <a-input v-model="form.question_text" />
          </a-form-item>
          <a-form-item label="选项" v-if="form.question_type !== 'judgment'">
            <div v-for="k in ['A','B','C','D']" :key="k" style="display:flex;align-items:center;gap:8px;">
              <span>{{ k }}：</span>
              <a-input v-model="form.options[k]" placeholder="选项内容"/>
            </div>
          </a-form-item>
          <a-form-item label="答案">
            <a-input v-model="answerStr" placeholder="支持多选，逗号分隔（如 A,B）"/>
          </a-form-item>
          <a-button type="primary" long @click="submit">提交</a-button>
        </a-form>
      </a-tab-pane>
      <a-tab-pane key="json" title="批量 JSON 导入">
        <a-textarea v-model="jsonStr" rows="8" placeholder="贴上题目 JSON数组，如：[{}]" />
        <a-button type="primary" long @click="importJson" style="margin-top:8px;">导入</a-button>
      </a-tab-pane>
    </a-tabs>
  </a-card>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { Message } from '@arco-design/web-vue';
import axios from 'axios';
import { useRouter } from 'vue-router';
const router = useRouter();

const tab = ref('form');

const form = reactive({
  category: '',
  question_type: 'single_choice',
  question_text: '',
  options: { A: '', B: '', C: '', D: '' } as Record<string, string>,
  answer: [] as string[]
});
const answerStr = ref('');

const submit = async () => {
  form.answer = answerStr.value.split(',').map(i => i.trim());
  try {
    // await axios.post('http://localhost:8080/api/questions', form);
    await axios.post('/api/questions', form);
    Message.success('提交成功');
    router.push({ name: 'question-list' });
  } catch (e) {
    Message.error('提交失败');
  }
};

// JSON 批量导入
const jsonStr = ref('');
const importJson = async () => {
  try {
    const arr = JSON.parse(jsonStr.value);
    if (!Array.isArray(arr)) throw new Error('不是数组');
    // 逐条入库，实际可优化为后端批量接口
    for (const q of arr) {
      await axios.post('http://localhost:8080/api/questions', q);
    }
    Message.success('批量导入成功');
    router.push({ name: 'question-list' });
  } catch (e: any) {
    Message.error(e.message || '格式错误');
  }
};
</script>