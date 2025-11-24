<template>
  <div class="mode-container">
    <h2 class="page-title">选择练习模式</h2>
    
    <div class="mode-cards">
      <!-- 1. 随机练习 -->
      <div class="mode-card" @click="go('random')">
        <div class="icon-box random-icon">
          <icon-dice size="40" />
        </div>
        <h3>随机练习</h3>
        <p>从题库中随机抽取题目，支持筛选分类和题型，适合碎片时间刷题。</p>
        <a-button type="primary" shape="round">开始挑战</a-button>
      </div>

      <!-- 2. 顺序练习 -->
      <div class="mode-card" @click="go('sequential')">
        <div class="icon-box seq-icon">
          <icon-ordered-list size="40" />
        </div>
        <h3>顺序练习</h3>
        <p>按照题库顺序逐一练习，适合系统性地复习所有知识点。</p>
        <a-button type="outline" shape="round">开始学习</a-button>
      </div>

      <!-- 3. 全真模拟 -->
      <div class="mode-card" @click="go('exam')">
        <div class="icon-box exam-icon">
          <icon-schedule size="40" />
        </div>
        <h3>全真模拟</h3>
        <p>模拟真实考试环境，自定义题量与分值，交卷后统一判分。</p>
        <a-button type="primary" status="danger" shape="round">进入考场</a-button>
      </div>

      <!-- 4. 错题本 (带红点) -->
      <div class="mode-card mistake-card-wrapper" @click="go('mistake')">
        <a-badge :count="mistakeCount" :offset="[10, -10]">
          <div class="icon-box mistake-icon">
            <icon-book size="40" />
          </div>
        </a-badge>
        <h3>错题本</h3>
        <p>自动记录你的每一次错误，重点攻克薄弱环节，消灭错题。</p>
        <a-button type="primary" status="warning" shape="round">复习错题</a-button>
      </div>
      <!-- 5. 多人协作 -->
      <div class="mode-card" @click="go('coop')">
        <div class="icon-box coop-icon">
          <icon-user-group size="40" />
        </div>
        <h3>多人协作</h3>
        <p>邀请好友一起实时刷题，PK 进度与正确率，共同进步。</p>
        <a-button type="primary" status="success" shape="round">进入大厅</a-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { IconUserGroup } from '@arco-design/web-vue/es/icon';
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { 
  IconDice, 
  IconOrderedList, 
  IconSchedule, 
  IconBook 
} from '@arco-design/web-vue/es/icon';

const router = useRouter();
const mistakeCount = ref(0);
const userId = Number(localStorage.getItem('user_id') || 0);

onMounted(async () => {
  if (userId) {
    try {
      // 复用获取错题列表接口，后端最好有个 count 接口，这里先取列表长度
      const res = await axios.get(`/api/mistakes?user_id=${userId}`);
      const list = res.data.data || [];
      mistakeCount.value = list.length;
    } catch (e) {
      console.error(e);
    }
  }
});

const go = (mode: string) => {
  switch (mode) {
    case 'random': router.push({ name: 'practice-filter' }); break;
    case 'sequential': router.push({ name: 'practice-exam', params: { mode: 'sequential' } }); break;
    case 'exam': router.push({ name: 'exam-config' }); break;
    case 'mistake': router.push({ name: 'mistake-book' }); break;
    case 'coop':router.push({ name: 'coop-lobby' });break;
  }
};
</script>

<style scoped>
/* ... 样式保持不变，只增加一个 overflow visible 以免红点被切 */
.mistake-card-wrapper { overflow: visible !important; }

/* 保持原有的样式 */
.mode-container {
  max-width: 1000px; margin: 40px auto; text-align: center;
}
.page-title { margin-bottom: 40px; color: #1d2129; font-size: 28px; }
.mode-cards { display: grid; grid-template-columns: repeat(2, 1fr); gap: 30px; padding: 0 20px; }
.mode-card {
  background: #fff; border-radius: 12px; padding: 40px 20px;
  box-shadow: 0 4px 10px rgba(0,0,0,0.05); cursor: pointer;
  transition: all 0.3s; border: 1px solid transparent;
  display: flex; flex-direction: column; align-items: center;
}
.mode-card:hover { transform: translateY(-5px); box-shadow: 0 8px 20px rgba(0,0,0,0.1); border-color: #165dff; }
.icon-box {
  width: 70px; height: 70px; border-radius: 50%; margin-bottom: 20px;
  display: flex; align-items: center; justify-content: center; color: #fff;
}
.random-icon { background: linear-gradient(135deg, #165dff 0%, #722ed1 100%); }
.seq-icon { background: linear-gradient(135deg, #00b42a 0%, #86909c 100%); }
.exam-icon { background: linear-gradient(135deg, #f53f3f 0%, #ff7d00 100%); }
.mistake-icon { background: linear-gradient(135deg, #ff7d00 0%, #f7ba1e 100%); }
h3 { margin-bottom: 10px; font-size: 20px; color: #1d2129; }
p { color: #86909c; margin-bottom: 25px; height: 40px; font-size: 14px; max-width: 300px; }

@media (max-width: 768px) {
  .mode-container { margin: 20px auto; }
  .mode-cards { grid-template-columns: 1fr; gap: 20px; }
  .mode-card { padding: 30px 20px; flex-direction: row; text-align: left; align-items: center; }
  .icon-box { width: 50px; height: 50px; margin-bottom: 0; margin-right: 15px; }
  h3 { font-size: 18px; margin-bottom: 5px; }
  p { display: none; }
  .mode-card button { display: none; }
  .mode-card > div:not(.icon-box) { flex: 1; }
}
.coop-icon { background: linear-gradient(135deg, #00d0b6 0%, #00b42a 100%); }
</style>