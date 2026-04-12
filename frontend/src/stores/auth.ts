import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import client from '../api/client'

interface User {
  id: number
  username: string
  full_name: string
  role: string
}

interface LoginResponse {
  access_token: string
  refresh_token: string
}

function decodeJwtPayload(token: string): Record<string, unknown> {
  const base64Url = token.split('.')[1]
  const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
  const jsonPayload = decodeURIComponent(
    atob(base64)
      .split('')
      .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
      .join(''),
  )
  return JSON.parse(jsonPayload)
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(localStorage.getItem('access_token'))
  const refreshToken = ref<string | null>(localStorage.getItem('refresh_token'))

  const isAuthenticated = computed(() => !!accessToken.value)
  const userRole = computed(() => user.value?.role ?? null)

  function restoreUserFromToken() {
    if (accessToken.value && !user.value) {
      try {
        const payload = decodeJwtPayload(accessToken.value)
        user.value = {
          id: payload.user_id as number,
          username: payload.username as string,
          full_name: payload.full_name as string,
          role: payload.role as string,
        }
      } catch {
        clearAuth()
      }
    }
  }

  async function login(username: string, password: string): Promise<void> {
    const data = (await client.post('/auth/login', { username, password })) as unknown as LoginResponse
    accessToken.value = data.access_token
    refreshToken.value = data.refresh_token
    localStorage.setItem('access_token', data.access_token)
    localStorage.setItem('refresh_token', data.refresh_token)

    const payload = decodeJwtPayload(data.access_token)
    user.value = {
      id: payload.user_id as number,
      username: payload.username as string,
      full_name: payload.full_name as string,
      role: payload.role as string,
    }
  }

  function logout() {
    clearAuth()
    window.location.href = '/login'
  }

  function clearAuth() {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
  }

  restoreUserFromToken()

  return {
    user,
    accessToken,
    refreshToken,
    isAuthenticated,
    userRole,
    login,
    logout,
  }
})
