<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import StatusBadge from '../components/StatusBadge.vue'
import { useRequestsStore } from '../stores/requests'
import { useAuthStore } from '../stores/auth'
import client from '../api/client'

interface RepairRequest {
  id: number
  equipment_id: number
  title: string
  description: string
  priority: string
  status: string
  created_by: number
  assigned_to: number | null
  created_at: string
  updated_at: string
  equipment?: { id: number; name: string; inventory_number: string }
  creator?: { id: number; full_name: string }
  assignee?: { id: number; full_name: string } | null
}

interface UserOption {
  id: number
  full_name: string
  role: string
}

const route = useRoute()
const router = useRouter()
const store = useRequestsStore()
const auth = useAuthStore()

const request = ref<RepairRequest | null>(null)
const loading = ref(true)
const error = ref('')
const actionLoading = ref(false)
const actionError = ref('')

const users = ref<UserOption[]>([])
const showAssignModal = ref(false)
const selectedUserId = ref<number>(0)

const isEngineerOrAdmin = computed(() =>
  auth.userRole === 'engineer' || auth.userRole === 'admin',
)

const statusTransitions: Record<string, { label: string; target: string }[]> = {
  new: [
    { label: 'Взять в работу', target: 'in_progress' },
    { label: 'Отменить', target: 'cancelled' },
  ],
  in_progress: [
    { label: 'Завершить', target: 'completed' },
    { label: 'Отменить', target: 'cancelled' },
  ],
  completed: [],
  cancelled: [],
}

const availableTransitions = computed(() => {
  if (!request.value || !isEngineerOrAdmin.value) return []
  return statusTransitions[request.value.status] ?? []
})

onMounted(async () => {
  try {
    request.value = (await store.fetchById(route.params.id as string)) as RepairRequest
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Ошибка загрузки'
  } finally {
    loading.value = false
  }
})

async function changeStatus(target: string) {
  if (!request.value) return
  actionLoading.value = true
  actionError.value = ''
  try {
    request.value = (await store.update(request.value.id, { status: target })) as RepairRequest
  } catch (e) {
    actionError.value = e instanceof Error ? e.message : 'Ошибка обновления статуса'
  } finally {
    actionLoading.value = false
  }
}

async function openAssignModal() {
  showAssignModal.value = true
  if (users.value.length === 0) {
    try {
      const res = (await client.get('/users')) as unknown as { data: UserOption[] }
      users.value = res.data
    } catch {
      actionError.value = 'Не удалось загрузить пользователей'
    }
  }
}

async function assignUser() {
  if (!request.value || !selectedUserId.value) return
  actionLoading.value = true
  actionError.value = ''
  try {
    request.value = (await store.update(request.value.id, {
      assigned_to: selectedUserId.value,
    })) as RepairRequest
    showAssignModal.value = false
    selectedUserId.value = 0
  } catch (e) {
    actionError.value = e instanceof Error ? e.message : 'Ошибка назначения'
  } finally {
    actionLoading.value = false
  }
}

function formatDate(dateStr: string | undefined): string {
  if (!dateStr) return '---'
  return new Date(dateStr).toLocaleString('ru-RU')
}

const statusSteps = ['new', 'in_progress', 'completed']
const statusLabels: Record<string, string> = {
  new: 'Новая',
  in_progress: 'В работе',
  completed: 'Завершена',
  cancelled: 'Отменена',
}

const priorityLabels: Record<string, string> = {
  low: 'Низкий',
  medium: 'Средний',
  high: 'Высокий',
  critical: 'Критический',
}
</script>

<template>
  <div class="detail-page">
    <button class="back-btn" @click="router.push('/requests')">
      &larr; К списку заявок
    </button>

    <div v-if="loading" class="detail-loading">Загрузка...</div>

    <div v-else-if="error" class="detail-error" role="alert">{{ error }}</div>

    <template v-else-if="request">
      <div class="detail-header">
        <h1>Заявка #{{ request.id }}</h1>
        <div class="detail-badges">
          <StatusBadge :status="request.status" />
          <StatusBadge :status="request.priority" />
        </div>
      </div>

      <div v-if="actionError" class="detail-error" role="alert">{{ actionError }}</div>

      <div class="status-timeline">
        <div
          v-for="(step, idx) in statusSteps"
          :key="step"
          class="timeline-step"
          :class="{
            active: request.status === step,
            done: statusSteps.indexOf(request.status) > idx,
            cancelled: request.status === 'cancelled',
          }"
        >
          <div class="timeline-dot" />
          <span class="timeline-label">{{ statusLabels[step] }}</span>
        </div>
        <div
          v-if="request.status === 'cancelled'"
          class="timeline-step cancelled active"
        >
          <div class="timeline-dot" />
          <span class="timeline-label">Отменена</span>
        </div>
      </div>

      <div class="detail-card">
        <h2 class="section-title">{{ request.title }}</h2>

        <div class="detail-grid">
          <div class="detail-field">
            <span class="detail-label">Оборудование</span>
            <span class="detail-value">
              {{ request.equipment?.name ?? '---' }}
              <span v-if="request.equipment" class="text-muted">
                ({{ request.equipment.inventory_number }})
              </span>
            </span>
          </div>
          <div class="detail-field">
            <span class="detail-label">Приоритет</span>
            <span class="detail-value">{{ priorityLabels[request.priority] ?? request.priority }}</span>
          </div>
          <div class="detail-field">
            <span class="detail-label">Создал</span>
            <span class="detail-value">{{ request.creator?.full_name ?? '---' }}</span>
          </div>
          <div class="detail-field">
            <span class="detail-label">Исполнитель</span>
            <span class="detail-value">{{ request.assignee?.full_name ?? 'Не назначен' }}</span>
          </div>
          <div class="detail-field">
            <span class="detail-label">Создана</span>
            <span class="detail-value">{{ formatDate(request.created_at) }}</span>
          </div>
          <div class="detail-field">
            <span class="detail-label">Обновлена</span>
            <span class="detail-value">{{ formatDate(request.updated_at) }}</span>
          </div>
        </div>

        <div v-if="request.description" class="description-block">
          <span class="detail-label">Описание</span>
          <p class="description-text">{{ request.description }}</p>
        </div>
      </div>

      <div v-if="isEngineerOrAdmin" class="actions-card">
        <h2 class="section-title">Действия</h2>
        <div class="actions-row">
          <button
            v-for="transition in availableTransitions"
            :key="transition.target"
            class="btn"
            :class="transition.target === 'cancelled' ? 'btn-danger' : 'btn-primary'"
            :disabled="actionLoading"
            @click="changeStatus(transition.target)"
          >
            {{ transition.label }}
          </button>
          <button
            class="btn btn-secondary"
            :disabled="actionLoading"
            @click="openAssignModal"
          >
            Назначить исполнителя
          </button>
        </div>
      </div>

      <div v-if="showAssignModal" class="modal-overlay" @click.self="showAssignModal = false">
        <div class="modal" role="dialog" aria-label="Назначить исполнителя">
          <h3 class="modal-title">Назначить исполнителя</h3>
          <div class="form-field">
            <label for="assign-user">Пользователь</label>
            <select id="assign-user" v-model.number="selectedUserId">
              <option :value="0" disabled>Выберите пользователя</option>
              <option
                v-for="u in users"
                :key="u.id"
                :value="u.id"
              >
                {{ u.full_name }} ({{ u.role }})
              </option>
            </select>
          </div>
          <div class="modal-actions">
            <button class="btn btn-secondary" @click="showAssignModal = false">
              Отмена
            </button>
            <button
              class="btn btn-primary"
              :disabled="!selectedUserId || actionLoading"
              @click="assignUser"
            >
              Назначить
            </button>
          </div>
        </div>
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
  margin-bottom: 16px;
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

.detail-badges {
  display: flex;
  gap: 8px;
}

.status-timeline {
  display: flex;
  align-items: center;
  gap: 0;
  margin-bottom: 20px;
  background: #fff;
  border-radius: 12px;
  padding: 18px 24px;
  box-shadow: 0 1px 6px rgba(0, 0, 0, 0.07);
}

.timeline-step {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  position: relative;
}

.timeline-step:not(:last-child)::after {
  content: '';
  flex: 1;
  height: 2px;
  background: #e2e8f0;
  margin: 0 8px;
}

.timeline-step.done:not(:last-child)::after {
  background: #16a34a;
}

.timeline-dot {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: #e2e8f0;
  flex-shrink: 0;
}

.timeline-step.done .timeline-dot {
  background: #16a34a;
}

.timeline-step.active .timeline-dot {
  background: #2563eb;
  box-shadow: 0 0 0 4px rgba(37, 99, 235, 0.2);
}

.timeline-step.cancelled .timeline-dot {
  background: #dc2626;
}

.timeline-label {
  font-size: 14px;
  color: #94a3b8;
  white-space: nowrap;
}

.timeline-step.active .timeline-label {
  color: #1e293b;
  font-weight: 600;
}

.timeline-step.done .timeline-label {
  color: #16a34a;
}

.detail-card,
.actions-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 20px;
  box-shadow: 0 1px 6px rgba(0, 0, 0, 0.07);
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 16px;
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

.text-muted {
  color: #94a3b8;
  font-size: 13px;
}

.description-block {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #f1f5f9;
}

.description-text {
  font-size: 14px;
  color: #334155;
  line-height: 1.6;
  margin-top: 4px;
  white-space: pre-wrap;
}

.actions-row {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
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

.btn-danger {
  background: #dc2626;
  color: #fff;
}

.btn-danger:hover:not(:disabled) {
  background: #b91c1c;
}

.btn-danger:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}

.modal {
  background: #fff;
  border-radius: 12px;
  padding: 28px;
  width: 100%;
  max-width: 440px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.18);
}

.modal-title {
  font-size: 18px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 16px;
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

.form-field select {
  width: 100%;
  height: 42px;
  padding: 10px 14px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 14px;
  color: #1e293b;
  background: #fff;
  box-sizing: border-box;
}

.form-field select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
