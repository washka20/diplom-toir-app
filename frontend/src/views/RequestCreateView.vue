<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useRequestsStore } from '../stores/requests'
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

const router = useRouter()
const requestsStore = useRequestsStore()

const equipmentList = ref<EquipmentOption[]>([])
const form = ref({
  equipment_id: 0,
  title: '',
  description: '',
  priority: 'medium',
})
const loading = ref(false)
const error = ref('')

onMounted(async () => {
  try {
    const res = (await client.get('/equipment', {
      params: { per_page: 200 },
    })) as unknown as EquipmentResponse
    equipmentList.value = res.data
  } catch {
    error.value = 'Не удалось загрузить список оборудования'
  }
})

async function handleSubmit() {
  error.value = ''
  if (!form.value.equipment_id) {
    error.value = 'Выберите оборудование'
    return
  }
  if (!form.value.title.trim()) {
    error.value = 'Введите заголовок'
    return
  }
  loading.value = true
  try {
    const result = await requestsStore.create(form.value)
    router.push(`/requests/${result.id}`)
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Ошибка создания заявки'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="create-page">
    <button class="back-btn" @click="router.push('/requests')">
      &larr; К списку заявок
    </button>

    <h1>Новая заявка на ремонт</h1>

    <div class="form-card">
      <div v-if="error" class="form-error" role="alert">{{ error }}</div>

      <form @submit.prevent="handleSubmit">
        <div class="form-field">
          <label for="req-equipment">Оборудование *</label>
          <select id="req-equipment" v-model.number="form.equipment_id" required>
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
          <label for="req-title">Заголовок *</label>
          <input
            id="req-title"
            v-model="form.title"
            type="text"
            placeholder="Краткое описание проблемы"
            required
          />
        </div>

        <div class="form-field">
          <label for="req-description">Описание</label>
          <textarea
            id="req-description"
            v-model="form.description"
            rows="4"
            placeholder="Подробное описание проблемы"
          />
        </div>

        <div class="form-field">
          <label for="req-priority">Приоритет</label>
          <select id="req-priority" v-model="form.priority">
            <option value="low">Низкий (low)</option>
            <option value="medium">Средний (medium)</option>
            <option value="high">Высокий (high)</option>
            <option value="critical">Критический (critical)</option>
          </select>
        </div>

        <div class="form-actions">
          <button
            type="button"
            class="btn btn-secondary"
            @click="router.push('/requests')"
          >
            Отмена
          </button>
          <button type="submit" class="btn btn-primary" :disabled="loading">
            {{ loading ? 'Создание...' : 'Создать заявку' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.create-page {
  max-width: 640px;
}

.back-btn {
  background: none;
  border: none;
  color: #2563eb;
  cursor: pointer;
  font-size: 14px;
  padding: 0;
  margin-bottom: 16px;
}

.back-btn:hover {
  text-decoration: underline;
}

.form-card {
  background: #fff;
  border-radius: 8px;
  padding: 24px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
}

.form-error {
  background: #fef2f2;
  color: #dc2626;
  padding: 10px 14px;
  border-radius: 6px;
  font-size: 14px;
  margin-bottom: 16px;
  border: 1px solid #fecaca;
}

.form-field {
  margin-bottom: 16px;
}

.form-field label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 6px;
}

.form-field input,
.form-field select,
.form-field textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  color: #1e293b;
  background: #fff;
  box-sizing: border-box;
  font-family: inherit;
}

.form-field input:focus,
.form-field select:focus,
.form-field textarea:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 20px;
}

.btn {
  padding: 10px 20px;
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
</style>
