<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()

const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

async function handleSubmit() {
  error.value = ''
  loading.value = true
  try {
    await auth.login(username.value, password.value)
    router.push('/')
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Ошибка авторизации'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <form class="login-card" @submit.prevent="handleSubmit">
      <div class="login-icon">&#9881;</div>
      <h1 class="login-title">ТОиР</h1>
      <p class="login-subtitle">Система управления техническим обслуживанием и ремонтом</p>

      <div v-if="error" class="login-error" role="alert">
        {{ error }}
      </div>

      <div class="login-field">
        <label for="username" class="login-label">Имя пользователя</label>
        <input
          id="username"
          v-model="username"
          type="text"
          class="login-input"
          autocomplete="username"
          required
        />
      </div>

      <div class="login-field">
        <label for="password" class="login-label">Пароль</label>
        <input
          id="password"
          v-model="password"
          type="password"
          class="login-input"
          autocomplete="current-password"
          required
        />
      </div>

      <button type="submit" class="login-button" :disabled="loading">
        {{ loading ? 'Вход...' : 'Войти' }}
      </button>
    </form>
  </div>
</template>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f1f5f9;
}

.login-card {
  background: #fff;
  padding: 40px;
  border-radius: 12px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 420px;
}

.login-icon {
  text-align: center;
  font-size: 48px;
  margin-bottom: 8px;
}

.login-title {
  font-size: 28px;
  font-weight: 700;
  color: #1e293b;
  margin: 0 0 8px;
  text-align: center;
}

.login-subtitle {
  font-size: 14px;
  color: #64748b;
  margin: 0 0 28px;
  text-align: center;
}

.login-error {
  background: #fef2f2;
  color: #dc2626;
  padding: 10px 14px;
  border-radius: 6px;
  font-size: 14px;
  margin-bottom: 16px;
  border: 1px solid #fecaca;
}

.login-field {
  margin-bottom: 18px;
}

.login-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 6px;
}

.login-input {
  width: 100%;
  height: 44px;
  padding: 10px 14px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 15px;
  color: #1e293b;
  background: #fff;
  transition: border-color 0.15s;
  box-sizing: border-box;
}

.login-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
}

.login-button {
  width: 100%;
  height: 46px;
  padding: 12px;
  background: #2563eb;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  margin-top: 8px;
  transition: background 0.15s;
}

.login-button:hover:not(:disabled) {
  background: #1d4ed8;
}

.login-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
