package models

import "time"

// MaintenanceLog представляет запись о выполненном техническом обслуживании.
type MaintenanceLog struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	EquipmentID   uint      `json:"equipment_id" gorm:"not null"`
	Equipment     Equipment `json:"equipment" gorm:"foreignKey:EquipmentID"`
	WorkOrderID   *uint     `json:"work_order_id"`
	Type          string    `json:"type" gorm:"size:20"`
	Description   string    `json:"description" gorm:"type:text"`
	PerformedBy   uint      `json:"performed_by" gorm:"not null"`
	Performer     User      `json:"performer" gorm:"foreignKey:PerformedBy"`
	PerformedAt   time.Time `json:"performed_at"`
	DurationHours float64   `json:"duration_hours" gorm:"type:decimal(5,2)"`
	CreatedAt     time.Time `json:"created_at"`
}
