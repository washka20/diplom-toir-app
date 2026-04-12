<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  status: string
}

const props = defineProps<Props>()

const colorMap: Record<string, { bg: string; text: string }> = {
  new: { bg: '#dbeafe', text: '#1e40af' },
  in_progress: { bg: '#fef3c7', text: '#92400e' },
  completed: { bg: '#d1fae5', text: '#065f46' },
  cancelled: { bg: '#f3f4f6', text: '#6b7280' },
  active: { bg: '#d1fae5', text: '#065f46' },
  inactive: { bg: '#f3f4f6', text: '#6b7280' },
  maintenance: { bg: '#fef3c7', text: '#92400e' },
  decommissioned: { bg: '#fee2e2', text: '#991b1b' },
  pending: { bg: '#e0e7ff', text: '#3730a3' },
  overdue: { bg: '#fee2e2', text: '#991b1b' },
  scheduled: { bg: '#dbeafe', text: '#1e40af' },
  low: { bg: '#f0fdf4', text: '#166534' },
  medium: { bg: '#fefce8', text: '#854d0e' },
  high: { bg: '#fff7ed', text: '#c2410c' },
  critical: { bg: '#fef2f2', text: '#dc2626' },
}

const defaultColor = { bg: '#f3f4f6', text: '#6b7280' }

const colors = computed(() => colorMap[props.status] ?? defaultColor)
</script>

<template>
  <span
    class="status-badge"
    :style="{ background: colors.bg, color: colors.text }"
  >
    {{ status }}
  </span>
</template>

<style scoped>
.status-badge {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
  white-space: nowrap;
}
</style>
