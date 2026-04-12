<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DataTable from '../components/DataTable.vue'
import { useUsersStore } from '../stores/users'

interface UserForm {
  username: string
  email: string
  password: string
  full_name: string
  role: string
}

const store = useUsersStore()

const showCreateForm = ref(false)
const editingId = ref<number | null>(null)

const createForm = ref<UserForm>({
  username: '',
  email: '',
  password: '',
  full_name: '',
  role: 'operator',
})

const editForm = ref({
  email: '',
  full_name: '',
  role: '',
  is_active: true,
})

const formError = ref('')
const formLoading = ref(false)

const columns = [
  { key: 'username', label: 'Логин', width: '140px' },
  { key: 'full_name', label: 'ФИО' },
  { key: 'email', label: 'Email' },
  { key: 'role', label: 'Роль', width: '120px' },
  { key: 'is_active', label: 'Активен', width: '100px' },
  { key: 'actions', label: '', width: '140px' },
]

onMounted(() => {
  store.fetchList()
})

async function handleCreate() {
  formError.value = ''
  if (!createForm.value.username || !createForm.value.password || !createForm.value.full_name) {
    formError.value = 'Заполните обязательные поля'
    return
  }
  formLoading.value = true
  try {
    await store.create(createForm.value)
    showCreateForm.value = false
    createForm.value = { username: '', email: '', password: '', full_name: '', role: 'operator' }
    store.fetchList()
  } catch (e) {
    formError.value = e instanceof Error ? e.message : 'Ошибка создания'
  } finally {
    formLoading.value = false
  }
}

function startEdit(row: Record<string, unknown>) {
  editingId.value = row.id as number
  editForm.value = {
    email: String(row.email ?? ''),
    full_name: String(row.full_name ?? ''),
    role: String(row.role ?? 'operator'),
    is_active: row.is_active as boolean,
  }
  formError.value = ''
}

function cancelEdit() {
  editingId.value = null
  formError.value = ''
}

async function saveEdit() {
  if (editingId.value === null) return
  formError.value = ''
  formLoading.value = true
  try {
    await store.update(editingId.value, editForm.value)
    editingId.value = null
    store.fetchList()
  } catch (e) {
    formError.value = e instanceof Error ? e.message : 'Ошибка сохранения'
  } finally {
    formLoading.value = false
  }
}

async function toggleActive(row: Record<string, unknown>) {
  try {
    await store.update(row.id as number, { is_active: !(row.is_active as boolean) })
    store.fetchList()
  } catch {
    // silently fail, user sees no change
  }
}
</script>

<template>
  <div class="users-page">
    <div class="page-header">
      <h1>Пользователи</h1>
      <button class="btn btn-primary" @click="showCreateForm = !showCreateForm">
        {{ showCreateForm ? 'Отмена' : 'Добавить пользователя' }}
      </button>
    </div>

    <div v-if="showCreateForm" class="form-card">
      <h2 class="form-title">Новый пользователь</h2>
      <div v-if="formError && editingId === null" class="form-error" role="alert">{{ formError }}</div>
      <form class="form-grid" @submit.prevent="handleCreate">
        <div class="form-field">
          <label for="user-login">Логин *</label>
          <input id="user-login" v-model="createForm.username" type="text" required />
        </div>
        <div class="form-field">
          <label for="user-name">ФИО *</label>
          <input id="user-name" v-model="createForm.full_name" type="text" required />
        </div>
        <div class="form-field">
          <label for="user-email">Email</label>
          <input id="user-email" v-model="createForm.email" type="email" />
        </div>
        <div class="form-field">
          <label for="user-pass">Пароль *</label>
          <input id="user-pass" v-model="createForm.password" type="password" required />
        </div>
        <div class="form-field">
          <label for="user-role">Роль</label>
          <select id="user-role" v-model="createForm.role">
            <option value="admin">admin</option>
            <option value="engineer">engineer</option>
            <option value="technician">technician</option>
            <option value="operator">operator</option>
          </select>
        </div>
        <div class="form-actions">
          <button type="submit" class="btn btn-primary" :disabled="formLoading">
            {{ formLoading ? 'Сохранение...' : 'Создать' }}
          </button>
        </div>
      </form>
    </div>

    <div v-if="editingId !== null" class="form-card">
      <h2 class="form-title">Редактирование пользователя</h2>
      <div v-if="formError" class="form-error" role="alert">{{ formError }}</div>
      <form class="form-grid" @submit.prevent="saveEdit">
        <div class="form-field">
          <label for="edit-name">ФИО</label>
          <input id="edit-name" v-model="editForm.full_name" type="text" />
        </div>
        <div class="form-field">
          <label for="edit-email">Email</label>
          <input id="edit-email" v-model="editForm.email" type="email" />
        </div>
        <div class="form-field">
          <label for="edit-role">Роль</label>
          <select id="edit-role" v-model="editForm.role">
            <option value="admin">admin</option>
            <option value="engineer">engineer</option>
            <option value="technician">technician</option>
            <option value="operator">operator</option>
          </select>
        </div>
        <div class="form-field">
          <label class="checkbox-label">
            <input type="checkbox" v-model="editForm.is_active" />
            Активен
          </label>
        </div>
        <div class="form-actions">
          <button type="button" class="btn btn-secondary" @click="cancelEdit">
            Отмена
          </button>
          <button type="submit" class="btn btn-primary" :disabled="formLoading">
            {{ formLoading ? 'Сохранение...' : 'Сохранить' }}
          </button>
        </div>
      </form>
    </div>

    <div v-if="store.loading" class="table-loading">Загрузка...</div>
    <DataTable
      v-else
      :columns="columns"
      :rows="(store.items as unknown as Record<string, unknown>[])"
    >
      <template #is_active="{ value }">
        <span :class="value ? 'active-badge' : 'inactive-badge'">
          {{ value ? 'Да' : 'Нет' }}
        </span>
      </template>
      <template #actions="{ row }">
        <div class="row-actions">
          <button class="btn-small" @click.stop="startEdit(row)">Изменить</button>
          <button
            class="btn-small"
            :class="(row as Record<string, unknown>).is_active ? 'btn-small-danger' : 'btn-small-success'"
            @click.stop="toggleActive(row)"
          >
            {{ (row as Record<string, unknown>).is_active ? 'Откл.' : 'Вкл.' }}
          </button>
        </div>
      </template>
    </DataTable>
  </div>
</template>

<style scoped>
.users-page {
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

.active-badge {
  color: #166534;
  font-weight: 500;
  font-size: 13px;
  background: #dcfce7;
  padding: 3px 10px;
  border-radius: 10px;
}

.inactive-badge {
  color: #991b1b;
  font-weight: 500;
  font-size: 13px;
  background: #fee2e2;
  padding: 3px 10px;
  border-radius: 10px;
}

.row-actions {
  display: flex;
  gap: 6px;
}

.btn-small {
  padding: 6px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  background: #fff;
  font-size: 13px;
  color: #374151;
  cursor: pointer;
  transition: background 0.15s;
}

.btn-small:hover {
  background: #f3f4f6;
}

.btn-small-danger {
  color: #dc2626;
  border-color: #fecaca;
}

.btn-small-danger:hover {
  background: #fef2f2;
}

.btn-small-success {
  color: #16a34a;
  border-color: #bbf7d0;
}

.btn-small-success:hover {
  background: #f0fdf4;
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

.form-field input[type="text"],
.form-field input[type="email"],
.form-field input[type="password"],
.form-field select {
  height: 42px;
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 14px;
  color: #1e293b;
  background: #fff;
}

.form-field input:focus,
.form-field select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-top: 20px;
  cursor: pointer;
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
