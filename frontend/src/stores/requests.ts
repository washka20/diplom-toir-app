import { defineStore } from 'pinia'
import { ref } from 'vue'
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

interface PaginatedResponse {
  data: RepairRequest[]
  meta: { page: number; per_page: number; total: number }
}

export const useRequestsStore = defineStore('requests', () => {
  const items = ref<RepairRequest[]>([])
  const total = ref(0)
  const loading = ref(false)

  async function fetchList(page = 1, filters: { status?: string; priority?: string } = {}) {
    loading.value = true
    try {
      const params: Record<string, string | number> = { page, per_page: 20 }
      if (filters.status) params.status = filters.status
      if (filters.priority) params.priority = filters.priority

      const res = (await client.get('/repair-requests', { params })) as unknown as PaginatedResponse
      items.value = res.data
      total.value = res.meta.total
    } finally {
      loading.value = false
    }
  }

  async function fetchById(id: number | string): Promise<RepairRequest> {
    return (await client.get(`/repair-requests/${id}`)) as unknown as RepairRequest
  }

  async function create(data: {
    equipment_id: number
    title: string
    description: string
    priority: string
  }): Promise<RepairRequest> {
    return (await client.post('/repair-requests', data)) as unknown as RepairRequest
  }

  async function update(
    id: number | string,
    data: { status?: string; assigned_to?: number },
  ): Promise<RepairRequest> {
    return (await client.put(`/repair-requests/${id}`, data)) as unknown as RepairRequest
  }

  return { items, total, loading, fetchList, fetchById, create, update }
})
