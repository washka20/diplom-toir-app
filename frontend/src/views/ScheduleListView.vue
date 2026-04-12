<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import DataTable from '../components/DataTable.vue'
import StatusBadge from '../components/StatusBadge.vue'
import { useSchedulesStore } from '../stores/schedules'
import client from '../api/client'

interface EquipmentOption {
  id: number
  name: string
  inventory_number: string
}

interface EquipmentResponse {
  data: EquipmentOption[]
  meta: { page: number; per_page: number; total: number }
}

const store = useSchedulesStore()

const page = ref(1)
const showForm = ref(false)
const equipmentList = ref<EquipmentOption[]>([])

const form = ref({
  equipment_id: 0,
  schedule_type: '',
  interval_days: 30,
  description: '',
})
const formError = ref('')
const formLoading = ref(false)

const columns = [
  { key: 'equipment_name', label: 'Оборудование' },
  { key: 'type_name', label: 'Тип ТО', width: '150px' },
  { key: 'interval_days', label: 'Интервал (дни)', width: '130px' },
  { key: 'next_date', label: 'Следующая дата', width: '150px' },
  { key: 'status', label: 'Статус', width: '120px' },
]

const totalPages = computed(() => Math.ceil(store.total / 20) || 1)

const today = new Date().toISOString().split('T')[0]

const tableRows = computed(() =>
  store.items.map((item: Record<string, unknown>) => ({
    ...item,
    equipment_name: (item.equipment as Record<string, unknown>)?.name ?? '---',
    type_name: item.type ?? '---',
    next_date: item.next_date ? String(item.next_date).split('T')[0] : '---',
    status: item.is_active ? 'active' : 'inactive',
    _overdue: item.next_date ? String(item.next_date).split('T')[0] < today && item.is_active : false,
  })),
)

function load() {
  store.fetchList(page.value)
}

onMounted(load)
watch(page, load)

async function openForm() {
  showForm.value = true
  if (equipmentList.value.length === 0) {
    try {
      const res = (await client.get('/equipment', {
        params: { per_page: 200 },
      })) as unknown as EquipmentResponse
      equipmentList.value = res.data
    } catch {
      formError.value = 'Не удалось загрузить список оборудования'
    }
  }
}

async function handleCreate() {
  formError.value = ''
  if (!form.value.equipment_id) {
    formError.value = 'Выберите оборудование'
    return
  }
  if (!form.value.schedule_type.trim()) {
    formError.value = 'Укажите тип ТО'
    return
  }
  formLoading.value = true
  try {
    await store.create(form.value)
    showForm.value = false
    form.value = { equipment_id: 0, schedule_type: '', interval_days: 30, description: '' }
    load()
  } catch (e) {
    formError.value = e instanceof Error ? e.message : 'Ошибка создания'
  } finally {
    formLoading.value = false
  }
}
</script>

<template>
  <div class="schedules-page">
    <div class="page-header">
      <h1>Графики ТО</h1>
      <button class="btn btn-primary" @click="showForm ? (showForm = false) : openForm()">
        {{ showForm ? 'Отмена' : 'Создать график' }}
      </button>
    </div>

    <div v-if="showForm" class="form-card">
      <h2 class="form-title">Новый график обслуживания</h2>
      <div v-if="formError" class="form-error" role="alert">{{ formError }}</div>
      <form class="form-grid" @submit.prevent="handleCreate">
        <div class="form-field">
          <label for="sch-equipment">Оборудование *</label>
          <select id="sch-equipment" v-model.number="form.equipment_id" required>
            <option :value="0" disabled>Выберите оборудование</option>
            <option
              v-for="eq in equipmentList"
              :key="eq.id"
              :value="eq.id"
            >
              {{ eq.inventory_number }} - {{ eq.name }}
            </option>
          </select>
        </div>
        <div class="form-field">
          <label for="sch-type">Тип ТО *</label>
          <input id="sch-type" v-model="form.schedule_type" type="text" placeholder="Ежемесячный осмотр" required />
        </div>
        <div class="form-field">
          <label for="sch-interval">Интервал (дни)</label>
          <input id="sch-interval" v-model.number="form.interval_days" type="number" min="1" />
        </div>
        <div class="form-field">
          <label for="sch-desc">Описание</label>
          <textarea id="sch-desc" v-model="form.description" rows="2" />
        </div>
        <div class="form-actions">
          <button type="submit" class="btn btn-primary" :disabled="formLoading">
            {{ formLoading ? 'Сохранение...' : 'Создать' }}
          </button>
        </div>
      </form>
    </div>

    <div v-if="store.loading" class="table-loading">Загрузка...</div>
    <DataTable
      v-else
      :columns="columns"
      :rows="(tableRows as Record<string, unknown>[])"
      :current-page="page"
      :total-pages="totalPages"
      @page-change="(p: number) => (page = p)"
    >
      <template #status="{ value }">
        <StatusBadge :status="String(value)" />
      </template>
      <template #next_date="{ row }">
        <span :class="{ overdue: (row as Record<string, unknown>)._overdue }">
          {{ (row as Record<string, unknown>).next_date }}
        </span>
      </template>
    </DataTable>
  </div>
</template>

<style scoped>
.schedules-page {
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

.table-loading {
  color: #64748b;
  padding: 40px 0;
  text-align: center;
}

.overdue {
  color: #dc2626;
  font-weight: 600;
  background: #fef2f2;
  padding: 2px 8px;
  border-radius: 4px;
}

.form-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 20px;
  box-shadow: 0 1px 6px rgba(0, 0, 0, 0.07);
}

.form-title {
  font-size: 18px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 16px;
}

.form-error {
  background: #fef2f2;
  color: #dc2626;
  padding: 10px 14px;
  border-radius: 6px;
  font-size: 14px;
  margin-bottom: 12px;
  border: 1px solid #fecaca;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

@media (max-width: 600px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-field label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.form-field input,
.form-field select,
.form-field textarea {
  height: 42px;
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 14px;
  color: #1e293b;
  background: #fff;
  font-family: inherit;
}

.form-field textarea {
  height: auto;
  min-height: 80px;
}

.form-field input:focus,
.form-field select:focus,
.form-field textarea:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
}

.form-actions {
  grid-column: 1 / -1;
  display: flex;
  justify-content: flex-end;
  padding-top: 8px;
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

.btn-primary:hover:not(:disabled) {
  background: #1d4ed8;
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
