import { defineStore } from 'pinia'
import { ref } from 'vue'
import client from '../api/client'

interface MaintenanceSchedule {
  id: number
  equipment_id: number
  schedule_type: string
  interval_days: number
  description: string
  next_date: string
  status: string
  equipment?: { id: number; name: string; inventory_number: string }
}

interface PaginatedResponse {
  data: MaintenanceSchedule[]
  meta: { page: number; per_page: number; total: number }
}

export const useSchedulesStore = defineStore('schedules', () => {
  const items = ref<MaintenanceSchedule[]>([])
  const total = ref(0)
  const loading = ref(false)

  async function fetchList(page = 1) {
    loading.value = true
    try {
      const res = (await client.get('/maintenance-schedules', {
        params: { page, per_page: 20 },
      })) as unknown as PaginatedResponse
      items.value = res.data
      total.value = res.meta.total
    } finally {
      loading.value = false
    }
  }

  async function create(data: {
    equipment_id: number
    schedule_type: string
    interval_days: number
    description: string
  }): Promise<MaintenanceSchedule> {
    const res = (await client.post('/maintenance-schedules', data)) as unknown as { data: MaintenanceSchedule }
    return res.data
  }

  return { items, total, loading, fetchList, create }
})
