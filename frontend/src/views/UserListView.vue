<template>
  <a-card title="用户管理">
    <a-table :columns="columns" :data="users" row-key="ID" :pagination="{ pageSize: 10 }">
      <!-- 格式化时间 -->
      <template #createdAt="{ record }">
        {{ new Date(record.CreatedAt).toLocaleString() }}
      </template>
      
      <!-- 操作列 -->
      <template #actions="{ record }">
        <a-popconfirm content="确定要删除该用户吗？此操作不可恢复。" @ok="del(record.ID)">
          <a-button type="text" status="danger" size="small">
            <template #icon><icon-delete /></template>
            删除
          </a-button>
        </a-popconfirm>
      </template>
    </a-table>
  </a-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { Message } from '@arco-design/web-vue';
import { IconDelete } from '@arco-design/web-vue/es/icon';
import axios from 'axios';

const columns = [
  { title: 'ID', dataIndex: 'ID', width: 80 },
  { title: '用户名', dataIndex: 'username' },
  { title: '注册时间', slotName: 'createdAt' },
  { title: '操作', slotName: 'actions', width: 100, align: 'center' }
];

const users = ref([]);

const fetchUsers = async () => {
  try {
    const res = await axios.get('/api/users');
    users.value = res.data.data || [];
  } catch (e) {
    Message.error('获取用户列表失败');
  }
};

const del = async (id: number) => {
  try {
    await axios.delete(`/api/users/${id}`);
    Message.success('用户已删除');
    fetchUsers(); // 刷新列表
  } catch (e) {
    Message.error('删除失败');
  }
};

onMounted(fetchUsers);
</script>