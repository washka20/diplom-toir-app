import axios from 'axios'
import type { AxiosResponse } from 'axios'

interface ApiEnvelope<T = unknown> {
  success: boolean
  data: T
  error: string | null
  meta: Record<string, unknown> | null
}

const client = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

client.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

client.interceptors.response.use(
  (response: AxiosResponse<ApiEnvelope>) => {
    return response.data.data as AxiosResponse
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      window.location.href = '/login'
    }
    const message = error.response?.data?.error || error.message || 'Network error'
    return Promise.reject(new Error(message))
  },
)

export default client
