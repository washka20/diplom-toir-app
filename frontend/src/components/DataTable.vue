<script setup lang="ts">
interface Column {
  key: string
  label: string
  width?: string
}

interface Props {
  columns: Column[]
  rows: Record<string, unknown>[]
  currentPage?: number
  totalPages?: number
}

const props = withDefaults(defineProps<Props>(), {
  currentPage: 1,
  totalPages: 1,
})

const emit = defineEmits<{
  'page-change': [page: number]
  'row-click': [row: Record<string, unknown>]
}>()

function prevPage() {
  if (props.currentPage > 1) {
    emit('page-change', props.currentPage - 1)
  }
}

function nextPage() {
  if (props.currentPage < props.totalPages) {
    emit('page-change', props.currentPage + 1)
  }
}
</script>

<template>
  <div class="data-table-wrapper">
    <table class="data-table" role="table">
      <thead>
        <tr>
          <th
            v-for="col in columns"
            :key="col.key"
            :style="col.width ? { width: col.width } : {}"
          >
            {{ col.label }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(row, idx) in rows"
          :key="idx"
          class="data-table-row"
          @click="emit('row-click', row)"
        >
          <td v-for="col in columns" :key="col.key">
            <slot :name="col.key" :row="row" :value="row[col.key]">
              {{ row[col.key] }}
            </slot>
          </td>
        </tr>
        <tr v-if="rows.length === 0">
          <td :colspan="columns.length" class="data-table-empty">
            Нет данных
          </td>
        </tr>
      </tbody>
    </table>

    <div v-if="totalPages > 1" class="data-table-pagination">
      <button
        class="pagination-btn"
        :disabled="currentPage <= 1"
        @click="prevPage"
      >
        &larr; Назад
      </button>
      <span class="pagination-info">{{ currentPage }} / {{ totalPages }}</span>
      <button
        class="pagination-btn"
        :disabled="currentPage >= totalPages"
        @click="nextPage"
      >
        Вперёд &rarr;
      </button>
    </div>
  </div>
</template>

<style scoped>
.data-table-wrapper {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
  overflow: hidden;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th {
  text-align: left;
  padding: 12px 16px;
  font-size: 12px;
  font-weight: 600;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
}

.data-table td {
  padding: 12px 16px;
  font-size: 14px;
  color: #334155;
  border-bottom: 1px solid #f1f5f9;
}

.data-table-row {
  cursor: pointer;
  transition: background 0.1s;
}

.data-table-row:hover {
  background: #f8fafc;
}

.data-table-empty {
  text-align: center;
  color: #94a3b8;
  padding: 32px 16px;
}

.data-table-pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 12px 16px;
  border-top: 1px solid #e2e8f0;
}

.pagination-btn {
  padding: 6px 14px;
  background: #fff;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 13px;
  color: #374151;
  cursor: pointer;
  transition: background 0.15s;
}

.pagination-btn:hover:not(:disabled) {
  background: #f3f4f6;
}

.pagination-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.pagination-info {
  font-size: 13px;
  color: #64748b;
}
</style>
