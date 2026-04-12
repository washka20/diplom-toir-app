package models

// DashboardMetrics содержит агрегированные метрики для дашборда.
type DashboardMetrics struct {
	OpenRequests       int64   `json:"open_requests"`
	OverdueSchedules   int64   `json:"overdue_schedules"`
	AvgRepairTimeHrs   float64 `json:"avg_repair_time_hrs"`
	CompletedThisMonth int64   `json:"completed_this_month"`
	TotalEquipment     int64   `json:"total_equipment"`
	ActiveEquipment    int64   `json:"active_equipment"`
}
