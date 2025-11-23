<template>
  <div>
    <a-card title="公告管理">
      <template #extra>
        <a-button type="primary" @click="visible = true">
          <template #icon><icon-plus /></template>
          发布新公告
        </a-button>
      </template>

      <a-table :columns="columns" :data="data" :pagination="false">
        <template #createdAt="{ record }">
          {{ new Date(record.created_at).toLocaleString() }}
        </template>
        <template #actions="{ record }">
          <a-popconfirm content="确定删除这条公告吗？" @ok="del(record.id)">
            <a-button type="text" status="danger" size="small">删除</a-button>
          </a-popconfirm>
        </template>
      </a-table>
    </a-card>

    <!-- 发布公告弹窗 -->
    <a-modal v-model:visible="visible" title="发布新公告" @ok="submit" :ok-loading="loading">
      <a-form :model="form">
        <a-form-item label="标题">
          <a-input v-model="form.title" placeholder="请输入公告标题" />
        </a-form-item>
        <a-form-item label="内容">
          <a-textarea v-model="form.content" placeholder="请输入公告内容" :auto-size="{minRows:4,maxRows:6}" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { Message } from '@arco-design/web-vue';
import { IconPlus } from '@arco-design/web-vue/es/icon';
import axios from 'axios';

const columns = [
  { title: 'ID', dataIndex: 'id', width: 80 },
  { title: '标题', dataIndex: 'title', width: 200 },
  { title: '内容', dataIndex: 'content', ellipsis: true }, // 内容过长自动省略
  { title: '发布时间', slotName: 'createdAt', width: 180 },
  { title: '操作', slotName: 'actions', width: 100, align: 'center' }
];

const data = ref([]);
const visible = ref(false);
const loading = ref(false);
const form = reactive({ title: '', content: '' });

const fetchList = async () => {
  const res = await axios.get('/api/announcements');
  data.value = res.data.data || [];
};

const submit = async () => {
  if (!form.title || !form.content) {
    Message.warning('标题和内容不能为空');
    return false; // 阻止弹窗关闭
  }
  loading.value = true;
  try {
    await axios.post('/api/announcements', form);
    Message.success('发布成功');
    visible.value = false;
    form.title = '';
    form.content = '';
    fetchList(); // 刷新列表
  } catch (e) {
    Message.error('发布失败');
  } finally {
    loading.value = false;
  }
};

const del = async (id: number) => {
  try {
    await axios.delete(`/api/announcements/${id}`);
    Message.success('删除成功');
    fetchList();
  } catch (e) {
    Message.error('删除失败');
  }
};

onMounted(fetchList);
</script>