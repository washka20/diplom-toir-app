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
  border-radius: 12px;
  box-shadow: 0 1px 6px rgba(0, 0, 0, 0.07);
  display: flex;
  overflow: hidden;
  transition: box-shadow 0.15s;
}

.dashboard-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.dashboard-card-indicator {
  width: 4px;
  flex-shrink: 0;
}

.dashboard-card-body {
  padding: 24px;
}

.dashboard-card-title {
  font-size: 14px;
  color: #64748b;
  margin-bottom: 8px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.dashboard-card-value {
  font-size: 36px;
  font-weight: 700;
  color: #1e293b;
}
</style>
