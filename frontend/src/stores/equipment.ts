import { defineStore } from 'pinia'
import { ref } from 'vue'
import client from '../api/client'

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

interface PaginatedResponse {
  data: Equipment[]
  meta: { page: number; per_page: number; total: number }
}

export const useEquipmentStore = defineStore('equipment', () => {
  const items = ref<Equipment[]>([])
  const total = ref(0)
  const loading = ref(false)

  async function fetchList(page = 1, filters: { status?: string; location?: string } = {}) {
    loading.value = true
    try {
      const params: Record<string, string | number> = { page, per_page: 20 }
      if (filters.status) params.status = filters.status
      if (filters.location) params.location = filters.location

      const res = (await client.get('/equipment', { params })) as unknown as PaginatedResponse
      items.value = res.data
      total.value = res.meta.total
    } finally {
      loading.value = false
    }
  }

  async function fetchById(id: number | string): Promise<Equipment> {
    return (await client.get(`/equipment/${id}`)) as unknown as Equipment
  }

  async function create(data: Partial<Equipment>): Promise<Equipment> {
    return (await client.post('/equipment', data)) as unknown as Equipment
  }

  async function update(id: number | string, data: Partial<Equipment>): Promise<Equipment> {
    return (await client.put(`/equipment/${id}`, data)) as unknown as Equipment
  }

  return { items, total, loading, fetchList, fetchById, create, update }
})
