package models

// WorkOrderPart представляет использованную запчасть в наряд-заказе.
type WorkOrderPart struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	WorkOrderID  uint      `json:"work_order_id" gorm:"not null"`
	WorkOrder    WorkOrder `json:"work_order" gorm:"foreignKey:WorkOrderID"`
	PartID       uint      `json:"part_id" gorm:"not null"`
	Part         Part      `json:"part" gorm:"foreignKey:PartID"`
	QuantityUsed int       `json:"quantity_used" gorm:"not null"`
}
