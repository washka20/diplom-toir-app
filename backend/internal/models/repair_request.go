package models

import "time"

// RepairRequest представляет заявку на ремонт оборудования.
type RepairRequest struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	EquipmentID uint       `json:"equipment_id" gorm:"not null"`
	Equipment   Equipment  `json:"equipment" gorm:"foreignKey:EquipmentID"`
	Title       string     `json:"title" gorm:"size:255;not null"`
	Description string     `json:"description" gorm:"type:text"`
	Priority    string     `json:"priority" gorm:"size:20;not null"`
	Status      string     `json:"status" gorm:"size:20;default:new"`
	CreatedBy   uint       `json:"created_by" gorm:"not null"`
	Creator     User       `json:"creator" gorm:"foreignKey:CreatedBy"`
	AssignedTo  *uint      `json:"assigned_to"`
	Assignee    *User      `json:"assignee" gorm:"foreignKey:AssignedTo"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at"`
}
