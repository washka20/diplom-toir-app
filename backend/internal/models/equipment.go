package models

import "time"

// Equipment представляет единицу оборудования.
type Equipment struct {
	ID                  uint       `json:"id" gorm:"primaryKey"`
	Name                string     `json:"name" gorm:"size:255;not null"`
	InventoryNumber     string     `json:"inventory_number" gorm:"uniqueIndex;size:50;not null"`
	Type                string     `json:"type" gorm:"size:100"`
	Manufacturer        string     `json:"manufacturer" gorm:"size:255"`
	Model               string     `json:"model" gorm:"size:255"`
	SerialNumber        string     `json:"serial_number" gorm:"size:100"`
	Location            string     `json:"location" gorm:"size:255"`
	Status              string     `json:"status" gorm:"size:20;default:active"`
	InstallationDate    *time.Time `json:"installation_date"`
	LastMaintenanceDate *time.Time `json:"last_maintenance_date"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}
