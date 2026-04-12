import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import { useAuthStore } from '../stores/auth'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'login',
    component: LoginView,
    meta: { public: true },
  },
  {
    path: '/',
    name: 'dashboard',
    component: () => import('../views/DashboardView.vue'),
  },
  {
    path: '/equipment',
    name: 'equipment-list',
    component: () => import('../views/EquipmentListView.vue'),
  },
  {
    path: '/equipment/:id',
    name: 'equipment-detail',
    component: () => import('../views/EquipmentDetailView.vue'),
  },
  {
    path: '/requests',
    name: 'request-list',
    component: () => import('../views/RequestListView.vue'),
  },
  {
    path: '/requests/create',
    name: 'request-create',
    component: () => import('../views/RequestCreateView.vue'),
  },
  {
    path: '/requests/:id',
    name: 'request-detail',
    component: () => import('../views/RequestDetailView.vue'),
  },
  {
    path: '/schedules',
    name: 'schedule-list',
    component: () => import('../views/ScheduleListView.vue'),
  },
  {
    path: '/users',
    name: 'users',
    component: () => import('../views/UsersView.vue'),
    meta: { roles: ['admin'] },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to) => {
  const auth = useAuthStore()

  if (to.meta.public) {
    return true
  }

  if (!auth.isAuthenticated) {
    return { name: 'login' }
  }

  if (to.meta.roles && Array.isArray(to.meta.roles)) {
    if (!auth.userRole || !to.meta.roles.includes(auth.userRole)) {
      return { name: 'dashboard' }
    }
  }

  return true
})

export default router
