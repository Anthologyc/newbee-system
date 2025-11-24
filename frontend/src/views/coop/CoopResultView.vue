<template>
  <div class="result-container">
    <div class="header">
      <h2>🏆 合作考试结算</h2>
      <p>协同作战，其利断金！</p>
    </div>

    <a-card class="score-card">
      <div class="total-score-box">
        <div class="label">合作总分</div>
        <div class="score">{{ averageScore }}</div>
        <div class="desc">两人的平均分</div>
      </div>
      
      <a-divider />

      <div class="pk-area">
        <!-- 我 -->
        <div class="player-res mine">
          <a-avatar :size="64" :style="{ backgroundColor: '#165dff' }">我</a-avatar>
          <div class="p-name">我</div>
          <div class="p-score">{{ myData?.Score || 0 }} 分</div>
          <div class="p-detail">对 {{ myData?.Correct }} / 错 {{ myData?.Wrong }}</div>
        </div>

        <div class="vs">VS</div>

        <!-- 队友 -->
        <div class="player-res opponent">
          <a-avatar :size="64" :style="{ backgroundColor: '#ff7d00' }">友</a-avatar>
          <div class="p-name">队友</div>
          <div class="p-score">{{ opponentData?.Score || 0 }} 分</div>
          <div class="p-detail">对 {{ opponentData?.Correct }} / 错 {{ opponentData?.Wrong }}</div>
        </div>
      </div>
    </a-card>

    <div class="actions">
      <a-button type="primary" size="large" @click="backToLobby">返回大厅</a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const myId = Number(localStorage.getItem('user_id') || 0);

const report = ref<any>({});

onMounted(() => {
  const data = localStorage.getItem('coop_report');
  if (data) {
    report.value = JSON.parse(data);
  }
});

const myData = computed(() => report.value[myId]);
const opponentData = computed(() => {
  const uid = Object.keys(report.value).find(id => Number(id) !== myId);
  return uid ? report.value[uid] : null;
});

const averageScore = computed(() => {
  const s1 = myData.value?.Score || 0;
  const s2 = opponentData.value?.Score || 0;
  return Math.round((s1 + s2) / 2);
});

const backToLobby = () => {
  router.push({ name: 'coop-lobby' });
};
</script>

<style scoped>
.result-container { max-width: 600px; margin: 40px auto; text-align: center; padding: 20px; }
.header { margin-bottom: 30px; }
.score-card { padding: 20px; border-radius: 12px; }
.total-score-box .score { font-size: 48px; font-weight: bold; color: #f53f3f; line-height: 1.2; }
.total-score-box .label { color: #86909c; }

.pk-area { display: flex; justify-content: space-around; align-items: center; margin-top: 30px; }
.player-res { display: flex; flex-direction: column; align-items: center; gap: 8px; }
.p-score { font-size: 24px; font-weight: bold; color: #165dff; }
.p-detail { font-size: 12px; color: #86909c; }
.vs { font-size: 24px; color: #ccc; font-weight: bold; }

.actions { margin-top: 40px; }
</style>