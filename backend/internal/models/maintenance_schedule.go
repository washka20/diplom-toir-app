package models

import "time"

// MaintenanceSchedule представляет расписание планового ТО для оборудования.
type MaintenanceSchedule struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	EquipmentID   uint       `json:"equipment_id" gorm:"not null"`
	Equipment     Equipment  `json:"equipment" gorm:"foreignKey:EquipmentID"`
	Type          string     `json:"type" gorm:"size:100"`
	IntervalDays  int        `json:"interval_days" gorm:"not null"`
	LastPerformed *time.Time `json:"last_performed"`
	NextDate      time.Time  `json:"next_date" gorm:"not null"`
	Description   string     `json:"description" gorm:"type:text"`
	IsActive      bool       `json:"is_active" gorm:"default:true"`
	CreatedBy     uint       `json:"created_by" gorm:"not null"`
	Creator       User       `json:"creator" gorm:"foreignKey:CreatedBy"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
