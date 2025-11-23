<template>
  <div>
    <a-card title="题目管理">
      <div style="margin-bottom:16px;display:flex;gap:16px;">
        <a-input v-model="search.category" placeholder="按类别筛选" allow-clear style="width:200px;" />
        <a-select v-model="search.type" placeholder="按题型筛选" allow-clear style="width:200px;">
          <a-option value="">全部题型</a-option>
          <a-option value="single_choice">单选题</a-option>
          <a-option value="multiple_choice">多选题</a-option>
          <a-option value="judgment">判断题</a-option>
        </a-select>
        <a-button type="primary" @click="fetchQuestions">查询</a-button>
        <a-button type="outline" @click="goCreate">录入题目</a-button>
      </div>
      
      <a-table :columns="columns" :data="questions" row-key="id" :pagination="false">
        <!-- ⚠️ 注意：这里必须用 record，不能用 row -->
        <template #options="{ record }">
          <div v-if="record.options">
            <a-tag v-for="(text, key) in record.options" :key="key" style="margin-right: 4px; margin-bottom: 4px;">
              {{ key }}: {{ text }}
            </a-tag>
          </div>
          <div v-else style="color: #999; font-size: 12px;">
            无选项
          </div>
        </template>

        <!-- ⚠️ 注意：这里也改成 record -->
        <template #answer="{ record }">
          <div v-if="record.answer && record.answer.length">
            <a-tag v-for="item in record.answer" color="green" :key="item" style="margin-right: 4px;">
              {{ item }}
            </a-tag>
          </div>
          <div v-else style="color: #999; font-size: 12px;">
            无答案
          </div>
        </template>

        <!-- ⚠️ 注意：这里也改成 record -->
        <template #actions="{ record }">
          <a-button size="small" @click="edit(record.id)">编辑</a-button>
          <a-button size="small" status="danger" @click="del(record.id)">删除</a-button>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import axios from 'axios';
import { useRouter } from 'vue-router';
import { Message } from '@arco-design/web-vue';

const router = useRouter();

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '类别', dataIndex: 'category', width: 120 },
  { title: '题型', dataIndex: 'question_type', width: 110 },
  { title: '题目', dataIndex: 'question_text', ellipsis: true },
  { title: '选项', slotName: 'options', width: 200 },
  { title: '答案', slotName: 'answer', width: 120 },
  { title: '操作', slotName: 'actions', align: 'center', width: 150 }
];

const questions = ref([]);
const search = reactive({ category: '', type: '' });

const fetchQuestions = async () => {
  try {
    const params = { category: search.category, type: search.type };
    // 确保这里是相对路径 /api/questions
    const res = await axios.get('/api/questions', { params });
    questions.value = res.data.data || [];
  } catch (error) {
    console.error('获取题目失败', error);
    Message.error('获取题目列表失败，请检查网络或后台服务');
  }
};

const goCreate = () => router.push({ name: 'question-create' });
const edit = (id: number) => { /* TODO */ };
const del = async (id: number) => {
  await axios.delete(`/api/questions/${id}`);
  fetchQuestions();
};

onMounted(fetchQuestions);
</script>