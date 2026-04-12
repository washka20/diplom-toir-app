<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  status: string
}

const props = defineProps<Props>()

const colorMap: Record<string, { bg: string; text: string }> = {
  new: { bg: '#dbeafe', text: '#1d4ed8' },
  assigned: { bg: '#e0e7ff', text: '#4338ca' },
  in_progress: { bg: '#fef3c7', text: '#92400e' },
  waiting_parts: { bg: '#fef9c3', text: '#854d0e' },
  completed: { bg: '#dcfce7', text: '#166534' },
  closed: { bg: '#f1f5f9', text: '#475569' },
  cancelled: { bg: '#f3f4f6', text: '#6b7280' },
  active: { bg: '#dcfce7', text: '#166534' },
  inactive: { bg: '#f3f4f6', text: '#6b7280' },
  maintenance: { bg: '#fef3c7', text: '#92400e' },
  decommissioned: { bg: '#fee2e2', text: '#991b1b' },
  pending: { bg: '#e0e7ff', text: '#3730a3' },
  overdue: { bg: '#fee2e2', text: '#991b1b' },
  scheduled: { bg: '#dbeafe', text: '#1d4ed8' },
  low: { bg: '#f0fdf4', text: '#166534' },
  medium: { bg: '#fefce8', text: '#854d0e' },
  high: { bg: '#fff7ed', text: '#c2410c' },
  critical: { bg: '#fef2f2', text: '#991b1b' },
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
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 13px;
  font-weight: 500;
  white-space: nowrap;
}
</style>
