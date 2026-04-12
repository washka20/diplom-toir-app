<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import StatusBadge from '../components/StatusBadge.vue'
import { useEquipmentStore } from '../stores/equipment'
import { useAuthStore } from '../stores/auth'

interface Equipment {
  id: number
  inventory_number: string
  name: string
  equipment_type: string
  location: string
  status: string
  model: string
  manufacturer: string
  serial_number: string
  installation_date: string
  last_maintenance_date: string
  notes: string
}

const route = useRoute()
const router = useRouter()
const store = useEquipmentStore()
const auth = useAuthStore()

const equipment = ref<Equipment | null>(null)
const loading = ref(true)
const error = ref('')
const editing = ref(false)
const editForm = ref<Partial<Equipment>>({})
const saveLoading = ref(false)
const saveError = ref('')

const isEngineerOrAdmin = computed(() =>
  auth.userRole === 'engineer' || auth.userRole === 'admin',
)

onMounted(async () => {
  try {
    equipment.value = (await store.fetchById(route.params.id as string)) as Equipment
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Ошибка загрузки'
  } finally {
    loading.value = false
  }
})

function startEdit() {
  if (!equipment.value) return
  editForm.value = { ...equipment.value }
  editing.value = true
  saveError.value = ''
}

function cancelEdit() {
  editing.value = false
  saveError.value = ''
}

async function saveEdit() {
  if (!equipment.value) return
  saveLoading.value = true
  saveError.value = ''
  try {
    equipment.value = (await store.update(equipment.value.id, editForm.value)) as Equipment
    editing.value = false
  } catch (e) {
    saveError.value = e instanceof Error ? e.message : 'Ошибка сохранения'
  } finally {
    saveLoading.value = false
  }
}

const fields: { key: keyof Equipment; label: string }[] = [
  { key: 'inventory_number', label: 'Инвентарный номер' },
  { key: 'name', label: 'Наименование' },
  { key: 'equipment_type', label: 'Тип' },
  { key: 'model', label: 'Модель' },
  { key: 'manufacturer', label: 'Производитель' },
  { key: 'serial_number', label: 'Серийный номер' },
  { key: 'location', label: 'Местоположение' },
  { key: 'installation_date', label: 'Дата установки' },
  { key: 'last_maintenance_date', label: 'Последнее ТО' },
  { key: 'notes', label: 'Примечания' },
]
</script>

<template>
  <div class="detail-page">
    <button class="back-btn" @click="router.push('/equipment')">
      &larr; К списку оборудования
    </button>

    <div v-if="loading" class="detail-loading">Загрузка...</div>

    <div v-else-if="error" class="detail-error" role="alert">{{ error }}</div>

    <template v-else-if="equipment">
      <div class="detail-header">
        <h1>{{ equipment.name }}</h1>
        <div class="detail-header-actions">
          <StatusBadge :status="equipment.status" />
          <button
            v-if="isEngineerOrAdmin && !editing"
            class="btn btn-primary"
            @click="startEdit"
          >
            Редактировать
          </button>
        </div>
      </div>

      <div v-if="!editing" class="detail-card">
        <div class="detail-grid">
          <div v-for="field in fields" :key="field.key" class="detail-field">
            <span class="detail-label">{{ field.label }}</span>
            <span class="detail-value">{{ equipment[field.key] || '---' }}</span>
          </div>
        </div>
      </div>

      <div v-else class="detail-card">
        <h2 class="form-title">Редактирование</h2>
        <div v-if="saveError" class="form-error" role="alert">{{ saveError }}</div>
        <form class="form-grid" @submit.prevent="saveEdit">
          <div class="form-field">
            <label for="edit-name">Наименование</label>
            <input id="edit-name" v-model="editForm.name" type="text" />
          </div>
          <div class="form-field">
            <label for="edit-type">Тип</label>
            <input id="edit-type" v-model="editForm.equipment_type" type="text" />
          </div>
          <div class="form-field">
            <label for="edit-model">Модель</label>
            <input id="edit-model" v-model="editForm.model" type="text" />
          </div>
          <div class="form-field">
            <label for="edit-manufacturer">Производитель</label>
            <input id="edit-manufacturer" v-model="editForm.manufacturer" type="text" />
          </div>
          <div class="form-field">
            <label for="edit-location">Местоположение</label>
            <input id="edit-location" v-model="editForm.location" type="text" />
          </div>
          <div class="form-field">
            <label for="edit-status">Статус</label>
            <select id="edit-status" v-model="editForm.status">
              <option value="active">active</option>
              <option value="maintenance">maintenance</option>
              <option value="decommissioned">decommissioned</option>
            </select>
          </div>
          <div class="form-field form-field-full">
            <label for="edit-notes">Примечания</label>
            <textarea id="edit-notes" v-model="editForm.notes" rows="3" />
          </div>
          <div class="form-actions">
            <button type="button" class="btn btn-secondary" @click="cancelEdit">
              Отмена
            </button>
            <button type="submit" class="btn btn-primary" :disabled="saveLoading">
              {{ saveLoading ? 'Сохранение...' : 'Сохранить' }}
            </button>
          </div>
        </form>
      </div>

      <div class="detail-card">
        <h2 class="section-title">История обслуживания</h2>
        <p class="placeholder-text">История обслуживания будет отображена здесь</p>
      </div>
    </template>
  </div>
</template>

<style scoped>
.detail-page {
  width: 100%;
}

.back-btn {
  background: none;
  border: none;
  color: #2563eb;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  padding: 0;
  margin-bottom: 16px;
}

.back-btn:hover {
  text-decoration: underline;
}

.detail-loading {
  color: #64748b;
  padding: 40px 0;
  text-align: center;
}

.detail-error {
  background: #fef2f2;
  color: #dc2626;
  padding: 12px 16px;
  border-radius: 8px;
  border: 1px solid #fecaca;
}

.detail-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 20px;
}

.detail-header h1 {
  margin-bottom: 0;
}

.detail-header-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.detail-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 20px;
  box-shadow: 0 1px 6px rgba(0, 0, 0, 0.07);
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

@media (max-width: 600px) {
  .detail-grid {
    grid-template-columns: 1fr;
  }
}

.detail-field {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-label {
  font-size: 12px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.detail-value {
  font-size: 15px;
  color: #1e293b;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 12px;
}

.placeholder-text {
  color: #94a3b8;
  font-size: 14px;
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

.form-field-full {
  grid-column: 1 / -1;
}

.form-actions {
  grid-column: 1 / -1;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
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

.btn-secondary {
  background: #fff;
  color: #374151;
  border: 1px solid #d1d5db;
}

.btn-secondary:hover {
  background: #f3f4f6;
}
</style>
