<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import DataTable from '../components/DataTable.vue'
import StatusBadge from '../components/StatusBadge.vue'
import { useEquipmentStore } from '../stores/equipment'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const store = useEquipmentStore()
const auth = useAuthStore()

const page = ref(1)
const statusFilter = ref('')
const locationFilter = ref('')
const showForm = ref(false)

const form = ref({
  inventory_number: '',
  name: '',
  equipment_type: '',
  model: '',
  manufacturer: '',
  serial_number: '',
  location: '',
  status: 'active',
  notes: '',
})
const formError = ref('')
const formLoading = ref(false)

const isEngineerOrAdmin = computed(() =>
  auth.userRole === 'engineer' || auth.userRole === 'admin',
)

const columns = [
  { key: 'inventory_number', label: 'Инв. номер', width: '130px' },
  { key: 'name', label: 'Наименование' },
  { key: 'equipment_type', label: 'Тип', width: '140px' },
  { key: 'location', label: 'Местоположение', width: '160px' },
  { key: 'status', label: 'Статус', width: '130px' },
]

const totalPages = computed(() => Math.ceil(store.total / 20) || 1)

function load() {
  store.fetchList(page.value, {
    status: statusFilter.value || undefined,
    location: locationFilter.value || undefined,
  })
}

onMounted(load)
watch([page, statusFilter, locationFilter], load)

function handleRowClick(row: Record<string, unknown>) {
  router.push(`/equipment/${row.id}`)
}

function resetFilters() {
  statusFilter.value = ''
  locationFilter.value = ''
  page.value = 1
}

async function handleCreate() {
  formError.value = ''
  if (!form.value.name || !form.value.inventory_number) {
    formError.value = 'Заполните обязательные поля'
    return
  }
  formLoading.value = true
  try {
    await store.create(form.value)
    showForm.value = false
    form.value = {
      inventory_number: '',
      name: '',
      equipment_type: '',
      model: '',
      manufacturer: '',
      serial_number: '',
      location: '',
      status: 'active',
      notes: '',
    }
    load()
  } catch (e) {
    formError.value = e instanceof Error ? e.message : 'Ошибка создания'
  } finally {
    formLoading.value = false
  }
}
</script>

<template>
  <div class="equipment-page">
    <div class="page-header">
      <h1>Оборудование</h1>
      <button
        v-if="isEngineerOrAdmin"
        class="btn btn-primary"
        @click="showForm = !showForm"
      >
        {{ showForm ? 'Отмена' : 'Добавить' }}
      </button>
    </div>

    <div v-if="showForm" class="form-card">
      <h2 class="form-title">Новое оборудование</h2>
      <div v-if="formError" class="form-error" role="alert">{{ formError }}</div>
      <form class="form-grid" @submit.prevent="handleCreate">
        <div class="form-field">
          <label for="eq-inv">Инв. номер *</label>
          <input id="eq-inv" v-model="form.inventory_number" type="text" required />
        </div>
        <div class="form-field">
          <label for="eq-name">Наименование *</label>
          <input id="eq-name" v-model="form.name" type="text" required />
        </div>
        <div class="form-field">
          <label for="eq-type">Тип</label>
          <input id="eq-type" v-model="form.equipment_type" type="text" />
        </div>
        <div class="form-field">
          <label for="eq-model">Модель</label>
          <input id="eq-model" v-model="form.model" type="text" />
        </div>
        <div class="form-field">
          <label for="eq-manufacturer">Производитель</label>
          <input id="eq-manufacturer" v-model="form.manufacturer" type="text" />
        </div>
        <div class="form-field">
          <label for="eq-serial">Серийный номер</label>
          <input id="eq-serial" v-model="form.serial_number" type="text" />
        </div>
        <div class="form-field">
          <label for="eq-location">Местоположение</label>
          <input id="eq-location" v-model="form.location" type="text" />
        </div>
        <div class="form-field">
          <label for="eq-status">Статус</label>
          <select id="eq-status" v-model="form.status">
            <option value="active">active</option>
            <option value="maintenance">maintenance</option>
            <option value="decommissioned">decommissioned</option>
          </select>
        </div>
        <div class="form-field form-field-full">
          <label for="eq-notes">Примечания</label>
          <textarea id="eq-notes" v-model="form.notes" rows="2" />
        </div>
        <div class="form-actions">
          <button type="submit" class="btn btn-primary" :disabled="formLoading">
            {{ formLoading ? 'Сохранение...' : 'Создать' }}
          </button>
        </div>
      </form>
    </div>

    <div class="filter-bar">
      <select v-model="statusFilter" aria-label="Фильтр по статусу">
        <option value="">Все статусы</option>
        <option value="active">active</option>
        <option value="maintenance">maintenance</option>
        <option value="decommissioned">decommissioned</option>
      </select>
      <input
        v-model="locationFilter"
        type="text"
        placeholder="Местоположение"
        aria-label="Фильтр по местоположению"
      />
      <button class="btn btn-secondary" @click="resetFilters">Сбросить</button>
    </div>

    <div v-if="store.loading" class="table-loading">Загрузка...</div>
    <DataTable
      v-else
      :columns="columns"
      :rows="(store.items as Record<string, unknown>[])"
      :current-page="page"
      :total-pages="totalPages"
      @row-click="handleRowClick"
      @page-change="(p: number) => (page = p)"
    >
      <template #status="{ value }">
        <StatusBadge :status="String(value)" />
      </template>
    </DataTable>
  </div>
</template>

<style scoped>
.equipment-page {
  max-width: 1100px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.page-header h1 {
  margin-bottom: 0;
}

.filter-bar {
  display: flex;
  gap: 10px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.filter-bar select,
.filter-bar input {
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  color: #1e293b;
  background: #fff;
}

.filter-bar select:focus,
.filter-bar input:focus {
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
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
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

.btn-secondary {
  background: #f3f4f6;
  color: #374151;
  border: 1px solid #d1d5db;
}

.btn-secondary:hover {
  background: #e5e7eb;
}

.form-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px 24px;
  margin-bottom: 16px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
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
  gap: 12px;
}

@media (max-width: 600px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.form-field label {
  font-size: 13px;
  font-weight: 500;
  color: #374151;
}

.form-field input,
.form-field select,
.form-field textarea {
  padding: 8px 10px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  color: #1e293b;
  background: #fff;
}

.form-field input:focus,
.form-field select:focus,
.form-field textarea:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
}

.form-field-full {
  grid-column: 1 / -1;
}

.form-actions {
  grid-column: 1 / -1;
  display: flex;
  justify-content: flex-end;
  padding-top: 4px;
}
</style>
