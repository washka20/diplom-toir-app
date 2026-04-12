<script setup lang="ts">
interface Props {
  title: string
  value: string | number
  color?: string
}

withDefaults(defineProps<Props>(), {
  color: '#3b82f6',
})
</script>

<template>
  <div class="dashboard-card">
    <div class="dashboard-card-indicator" :style="{ background: color }" />
    <div class="dashboard-card-body">
      <div class="dashboard-card-title">{{ title }}</div>
      <div class="dashboard-card-value">{{ value }}</div>
    </div>
  </div>
</template>

<style scoped>
.dashboard-card {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
  display: flex;
  overflow: hidden;
}

.dashboard-card-indicator {
  width: 4px;
  flex-shrink: 0;
}

.dashboard-card-body {
  padding: 16px 20px;
}

.dashboard-card-title {
  font-size: 13px;
  color: #64748b;
  margin-bottom: 6px;
}

.dashboard-card-value {
  font-size: 28px;
  font-weight: 700;
  color: #1e293b;
}
</style>
