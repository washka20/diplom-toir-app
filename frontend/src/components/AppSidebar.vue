<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()

const isEngineerOrAdmin = computed(() =>
  auth.userRole === 'engineer' || auth.userRole === 'admin',
)
const isAdmin = computed(() => auth.userRole === 'admin')

function handleLogout() {
  auth.logout()
}
</script>

<template>
  <aside class="sidebar">
    <div class="sidebar-header">
      <div class="sidebar-logo">ТОиР</div>
      <div v-if="auth.user" class="sidebar-user">
        <div class="sidebar-user-name">{{ auth.user.full_name }}</div>
        <div class="sidebar-user-role">{{ auth.user.role }}</div>
      </div>
    </div>

    <nav class="sidebar-nav">
      <router-link to="/" class="sidebar-link" exact-active-class="active">
        <span class="sidebar-link-icon">&#9632;</span>
        Дашборд
      </router-link>
      <router-link to="/equipment" class="sidebar-link" active-class="active">
        <span class="sidebar-link-icon">&#9881;</span>
        Оборудование
      </router-link>
      <router-link to="/requests" class="sidebar-link" active-class="active">
        <span class="sidebar-link-icon">&#9998;</span>
        Заявки на ремонт
      </router-link>
      <router-link
        v-if="isEngineerOrAdmin"
        to="/schedules"
        class="sidebar-link"
        active-class="active"
      >
        <span class="sidebar-link-icon">&#128197;</span>
        Графики ТО
      </router-link>
      <router-link
        v-if="isAdmin"
        to="/users"
        class="sidebar-link"
        active-class="active"
      >
        <span class="sidebar-link-icon">&#128101;</span>
        Пользователи
      </router-link>
    </nav>

    <button class="sidebar-logout" @click="handleLogout">
      Выход
    </button>
  </aside>
</template>

<style scoped>
.sidebar {
  width: 240px;
  min-height: 100vh;
  background: #1e293b;
  color: #e2e8f0;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.sidebar-header {
  padding: 20px 16px;
  border-bottom: 1px solid #334155;
}

.sidebar-logo {
  font-size: 20px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 12px;
}

.sidebar-user-name {
  font-size: 14px;
  font-weight: 500;
  color: #f1f5f9;
}

.sidebar-user-role {
  font-size: 12px;
  color: #94a3b8;
  margin-top: 2px;
}

.sidebar-nav {
  flex: 1;
  padding: 12px 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.sidebar-link {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  color: #cbd5e1;
  text-decoration: none;
  font-size: 14px;
  transition: background 0.15s, color 0.15s;
}

.sidebar-link:hover {
  background: #334155;
  color: #fff;
}

.sidebar-link.active {
  background: #3b82f6;
  color: #fff;
}

.sidebar-link-icon {
  width: 20px;
  text-align: center;
  font-size: 16px;
}

.sidebar-logout {
  margin: 16px;
  padding: 10px;
  background: transparent;
  border: 1px solid #475569;
  color: #94a3b8;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: background 0.15s, color 0.15s;
}

.sidebar-logout:hover {
  background: #dc2626;
  color: #fff;
  border-color: #dc2626;
}
</style>
