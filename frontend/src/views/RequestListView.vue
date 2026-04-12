<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import DataTable from '../components/DataTable.vue'
import StatusBadge from '../components/StatusBadge.vue'
import { useRequestsStore } from '../stores/requests'

const router = useRouter()
const store = useRequestsStore()

const page = ref(1)
const statusFilter = ref('')
const priorityFilter = ref('')

const columns = [
  { key: 'id', label: 'ID', width: '70px' },
  { key: 'equipment_name', label: 'Оборудование' },
  { key: 'title', label: 'Заголовок' },
  { key: 'priority', label: 'Приоритет', width: '120px' },
  { key: 'status', label: 'Статус', width: '130px' },
  { key: 'created_at', label: 'Дата создания', width: '140px' },
]

const totalPages = computed(() => Math.ceil(store.total / 20) || 1)

const tableRows = computed(() =>
  store.items.map((item) => ({
    ...item,
    equipment_name: item.equipment?.name ?? '---',
    created_at: item.created_at ? new Date(item.created_at).toLocaleDateString('ru-RU') : '---',
  })),
)

function load() {
  store.fetchList(page.value, {
    status: statusFilter.value || undefined,
    priority: priorityFilter.value || undefined,
  })
}

onMounted(load)
watch([page, statusFilter, priorityFilter], load)

function handleRowClick(row: Record<string, unknown>) {
  router.push(`/requests/${row.id}`)
}

function resetFilters() {
  statusFilter.value = ''
  priorityFilter.value = ''
  page.value = 1
}
</script>

<template>
  <div class="requests-page">
    <div class="page-header">
      <h1>Заявки на ремонт</h1>
      <button class="btn btn-primary" @click="router.push('/requests/create')">
        Создать заявку
      </button>
    </div>

    <div class="filter-bar">
      <select v-model="statusFilter" aria-label="Фильтр по статусу">
        <option value="">Все статусы</option>
        <option value="new">new</option>
        <option value="in_progress">in_progress</option>
        <option value="completed">completed</option>
        <option value="cancelled">cancelled</option>
      </select>
      <select v-model="priorityFilter" aria-label="Фильтр по приоритету">
        <option value="">Все приоритеты</option>
        <option value="low">low</option>
        <option value="medium">medium</option>
        <option value="high">high</option>
        <option value="critical">critical</option>
      </select>
      <button class="btn btn-secondary" @click="resetFilters">Сбросить</button>
    </div>

    <div v-if="store.loading" class="table-loading">Загрузка...</div>
    <DataTable
      v-else
      :columns="columns"
      :rows="(tableRows as Record<string, unknown>[])"
      :current-page="page"
      :total-pages="totalPages"
      @row-click="handleRowClick"
      @page-change="(p: number) => (page = p)"
    >
      <template #status="{ value }">
        <StatusBadge :status="String(value)" />
      </template>
      <template #priority="{ value }">
        <StatusBadge :status="String(value)" />
      </template>
    </DataTable>
  </div>
</template>

<style scoped>
.requests-page {
  width: 100%;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.page-header h1 {
  margin-bottom: 0;
}

.filter-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.filter-bar select {
  height: 40px;
  padding: 8px 14px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 14px;
  color: #1e293b;
  background: #fff;
}

.filter-bar select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
}

.table-loading {
  color: #64748b;
  padding: 40px 0;
  text-align: center;
}

.btn {
  height: 40px;
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.15s;
}

.btn-primary {
  background: #2563eb;
  color: #fff;
}

.btn-primary:hover {
  background: #1d4ed8;
}

.btn-secondary {
  background: #fff;
  color: #374151;
  border: 1px solid #d1d5db;
}

.btn-secondary:hover {
  background: #f3f4f6;
}
</style>
