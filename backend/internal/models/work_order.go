package models

import "time"

// WorkOrder представляет наряд-заказ на выполнение работ.
type WorkOrder struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	RepairRequestID *uint          `json:"repair_request_id"`
	RepairRequest   *RepairRequest `json:"repair_request" gorm:"foreignKey:RepairRequestID"`
	ScheduleID      *uint          `json:"schedule_id"`
	Description     string         `json:"description" gorm:"type:text"`
	PlannedStart    *time.Time     `json:"planned_start"`
	PlannedEnd      *time.Time     `json:"planned_end"`
	ActualStart     *time.Time     `json:"actual_start"`
	ActualEnd       *time.Time     `json:"actual_end"`
	Status          string         `json:"status" gorm:"size:20;default:pending"`
	AssignedTo      uint           `json:"assigned_to" gorm:"not null"`
	Assignee        User           `json:"assignee" gorm:"foreignKey:AssignedTo"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}
