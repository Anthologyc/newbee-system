<template>
  <div class="dashboard-container">
    <a-card title="工作台数据">
      <a-grid :cols="3" :colGap="20">
        <a-grid-item>
          <a-statistic title="总题目数" :value="120" show-group-separator />
        </a-grid-item>
        <a-grid-item>
          <a-statistic title="注册用户" :value="56" show-group-separator />
        </a-grid-item>
        <a-grid-item>
          <a-statistic title="今日活跃" :value="12" show-group-separator />
        </a-grid-item>
      </a-grid>
    </a-card>
    
    <a-card title="系统公告" style="margin-top: 20px;">
      <p>欢迎使用 NewBee 后台管理系统！</p>
      <p>当前版本：v0.0.1 - 开发中</p>
    </a-card>
  </div>
</template>

<script setup lang="ts">
// 这里暂时不需要逻辑
</script>

<style scoped>
</style>