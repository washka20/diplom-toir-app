<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DashboardCard from '../components/DashboardCard.vue'
import client from '../api/client'

interface DashboardData {
  open_requests: number
  overdue_schedules: number
  avg_repair_time_hrs: number
  completed_this_month: number
  total_equipment: number
  active_equipment: number
}

const data = ref<DashboardData | null>(null)
const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    const res = (await client.get('/dashboard')) as unknown as { data: DashboardData }
    data.value = res.data
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Ошибка загрузки данных'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="dashboard-page">
    <h1>Дашборд</h1>

    <div v-if="loading" class="dashboard-loading">Загрузка...</div>

    <div v-else-if="error" class="dashboard-error" role="alert">{{ error }}</div>

    <div v-else-if="data" class="dashboard-grid">
      <DashboardCard
        title="Открытые заявки"
        :value="data.open_requests"
        color="#f59e0b"
      />
      <DashboardCard
        title="Просроченные ТО"
        :value="data.overdue_schedules"
        color="#dc2626"
      />
      <DashboardCard
        title="Среднее время ремонта"
        :value="`${data.avg_repair_time_hrs.toFixed(1)} ч`"
        color="#2563eb"
      />
      <DashboardCard
        title="Выполнено за месяц"
        :value="data.completed_this_month"
        color="#16a34a"
      />
      <DashboardCard
        title="Всего оборудования"
        :value="data.total_equipment"
        color="#6b7280"
      />
      <DashboardCard
        title="Активное оборудование"
        :value="data.active_equipment"
        color="#16a34a"
      />
    </div>
  </div>
</template>

<style scoped>
.dashboard-page {
  width: 100%;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
}

@media (max-width: 768px) {
  .dashboard-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 480px) {
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
}

.dashboard-loading {
  color: #64748b;
  padding: 40px 0;
  text-align: center;
}

.dashboard-error {
  background: #fef2f2;
  color: #dc2626;
  padding: 12px 16px;
  border-radius: 8px;
  border: 1px solid #fecaca;
}
</style>
