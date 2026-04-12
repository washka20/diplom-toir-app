import { defineStore } from 'pinia'
import { ref } from 'vue'
import client from '../api/client'

interface User {
  id: number
  username: string
  email: string
  full_name: string
  role: string
  is_active: boolean
}

export const useUsersStore = defineStore('users', () => {
  const items = ref<User[]>([])
  const loading = ref(false)

  async function fetchList() {
    loading.value = true
    try {
      const res = (await client.get('/users')) as unknown as User[]
      items.value = res
    } finally {
      loading.value = false
    }
  }

  async function create(data: {
    username: string
    email: string
    password: string
    full_name: string
    role: string
  }): Promise<User> {
    return (await client.post('/users', data)) as unknown as User
  }

  async function update(id: number | string, data: Partial<User>): Promise<User> {
    return (await client.put(`/users/${id}`, data)) as unknown as User
  }

  return { items, loading, fetchList, create, update }
})
