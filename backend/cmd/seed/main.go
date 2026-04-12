package main

import (
	"fmt"
	"log"
	"time"

	"toir-app/internal/config"
	"toir-app/internal/database"
	"toir-app/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := autoMigrate(db); err != nil {
		log.Fatalf("failed to auto-migrate: %v", err)
	}

	if alreadySeeded(db) {
		log.Println("Database already seeded, skipping")
		return
	}

	if err := seed(db); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	log.Println("Seed completed successfully")
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Equipment{},
		&models.RepairRequest{},
		&models.MaintenanceSchedule{},
		&models.MaintenanceLog{},
		&models.Part{},
		&models.WorkOrder{},
		&models.WorkOrderPart{},
	)
}

func alreadySeeded(db *gorm.DB) bool {
	var count int64
	db.Model(&models.User{}).Count(&count)
	return count > 0
}

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}
	return string(hash)
}

func seed(db *gorm.DB) error {
	passwordHash := hashPassword("password123")

	// --- Users ---
	users := []models.User{
		{Username: "admin", Email: "admin@toir.local", PasswordHash: passwordHash, FullName: "Администратор Системы", Role: "admin"},
		{Username: "engineer", Email: "engineer@toir.local", PasswordHash: passwordHash, FullName: "Иванов Пётр Сергеевич", Role: "engineer"},
		{Username: "technician", Email: "technician@toir.local", PasswordHash: passwordHash, FullName: "Сидоров Алексей Николаевич", Role: "technician"},
		{Username: "operator", Email: "operator@toir.local", PasswordHash: passwordHash, FullName: "Козлова Мария Андреевна", Role: "operator"},
	}

	if err := db.Create(&users).Error; err != nil {
		return fmt.Errorf("failed to create users: %w", err)
	}
	log.Printf("Created %d users", len(users))

	// --- Equipment ---
	now := time.Now()
	installDate := now.AddDate(-2, 0, 0)
	lastMaintDate := now.AddDate(0, -1, 0)

	equipment := []models.Equipment{
		{Name: "Токарный станок 16К20", InventoryNumber: "ТС-001", Type: "Станок", Manufacturer: "СВЗС", Model: "16К20", SerialNumber: "SN-TC-001", Location: "Цех 1", Status: "active", InstallationDate: &installDate, LastMaintenanceDate: &lastMaintDate},
		{Name: "Фрезерный станок 6Р82", InventoryNumber: "ФС-001", Type: "Станок", Manufacturer: "ГФЗ", Model: "6Р82", SerialNumber: "SN-FS-001", Location: "Цех 1", Status: "active", InstallationDate: &installDate, LastMaintenanceDate: &lastMaintDate},
		{Name: "Компрессор Atlas Copco GA30", InventoryNumber: "КМ-001", Type: "Компрессор", Manufacturer: "Atlas Copco", Model: "GA30", SerialNumber: "SN-KM-001", Location: "Компрессорная", Status: "active", InstallationDate: &installDate, LastMaintenanceDate: &lastMaintDate},
		{Name: "Пресс гидравлический П6334", InventoryNumber: "ПГ-001", Type: "Пресс", Manufacturer: "ОЗПО", Model: "П6334", SerialNumber: "SN-PG-001", Location: "Цех 2", Status: "maintenance", InstallationDate: &installDate},
		{Name: "Конвейер ленточный КЛ-500", InventoryNumber: "КЛ-001", Type: "Конвейер", Manufacturer: "КонвейерМаш", Model: "КЛ-500", SerialNumber: "SN-KL-001", Location: "Склад", Status: "active", InstallationDate: &installDate, LastMaintenanceDate: &lastMaintDate},
		{Name: "Сварочный аппарат Lincoln Electric", InventoryNumber: "СА-001", Type: "Сварочное оборудование", Manufacturer: "Lincoln Electric", Model: "Invertec V350", SerialNumber: "SN-SA-001", Location: "Цех 3", Status: "active", InstallationDate: &installDate},
		{Name: "Насос центробежный К200", InventoryNumber: "НЦ-001", Type: "Насос", Manufacturer: "ГМС Насосы", Model: "К200", SerialNumber: "SN-NC-001", Location: "Насосная", Status: "active", InstallationDate: &installDate, LastMaintenanceDate: &lastMaintDate},
		{Name: "Электродвигатель АИР160", InventoryNumber: "ЭД-001", Type: "Электродвигатель", Manufacturer: "ВЭМЗ", Model: "АИР160", SerialNumber: "SN-ED-001", Location: "Цех 1", Status: "active", InstallationDate: &installDate},
		{Name: "Вентилятор промышленный ВЦ14", InventoryNumber: "ВП-001", Type: "Вентилятор", Manufacturer: "ВентПром", Model: "ВЦ14-46", SerialNumber: "SN-VP-001", Location: "Вентиляционная", Status: "maintenance", InstallationDate: &installDate},
		{Name: "Трансформатор силовой ТМ-630", InventoryNumber: "ТС-002", Type: "Трансформатор", Manufacturer: "ТрансЭнерго", Model: "ТМ-630", SerialNumber: "SN-TS-002", Location: "Подстанция", Status: "active", InstallationDate: &installDate, LastMaintenanceDate: &lastMaintDate},
		{Name: "Станок шлифовальный 3М132", InventoryNumber: "ШС-001", Type: "Станок", Manufacturer: "ЛШЗ", Model: "3М132", SerialNumber: "SN-SS-001", Location: "Цех 1", Status: "active", InstallationDate: &installDate},
		{Name: "Кран-балка подвесная КБП-5", InventoryNumber: "КБ-001", Type: "Грузоподъёмное", Manufacturer: "КранМаш", Model: "КБП-5", SerialNumber: "SN-KB-001", Location: "Цех 2", Status: "active", InstallationDate: &installDate, LastMaintenanceDate: &lastMaintDate},
	}

	if err := db.Create(&equipment).Error; err != nil {
		return fmt.Errorf("failed to create equipment: %w", err)
	}
	log.Printf("Created %d equipment items", len(equipment))

	// --- Repair Requests ---
	completedAt := now.AddDate(0, 0, -3)
	engineerID := users[1].ID
	technicianID := users[2].ID
	operatorID := users[3].ID

	repairRequests := []models.RepairRequest{
		{EquipmentID: equipment[0].ID, Title: "Вибрация шпинделя", Description: "Повышенная вибрация при работе на высоких оборотах", Priority: "high", Status: "new", CreatedBy: operatorID},
		{EquipmentID: equipment[3].ID, Title: "Утечка гидравлической жидкости", Description: "Обнаружена утечка в районе гидроцилиндра", Priority: "critical", Status: "new", CreatedBy: operatorID},
		{EquipmentID: equipment[4].ID, Title: "Проскальзывание ленты", Description: "Лента конвейера проскальзывает при нагрузке", Priority: "medium", Status: "assigned", CreatedBy: operatorID, AssignedTo: &technicianID},
		{EquipmentID: equipment[8].ID, Title: "Посторонний шум вентилятора", Description: "Появился гул при работе на полной мощности", Priority: "low", Status: "in_progress", CreatedBy: engineerID, AssignedTo: &technicianID},
		{EquipmentID: equipment[6].ID, Title: "Снижение производительности насоса", Description: "Расход снизился на 20% от нормы", Priority: "medium", Status: "completed", CreatedBy: operatorID, AssignedTo: &technicianID, CompletedAt: &completedAt},
		{EquipmentID: equipment[1].ID, Title: "Замена масла фрезерного станка", Description: "Плановая замена масла по регламенту", Priority: "low", Status: "closed", CreatedBy: engineerID, AssignedTo: &technicianID, CompletedAt: &completedAt},
	}

	if err := db.Create(&repairRequests).Error; err != nil {
		return fmt.Errorf("failed to create repair requests: %w", err)
	}
	log.Printf("Created %d repair requests", len(repairRequests))

	// --- Maintenance Schedules ---
	adminID := users[0].ID
	lastPerformed := now.AddDate(0, -1, 0)
	overdueDate := now.AddDate(0, 0, -5)

	schedules := []models.MaintenanceSchedule{
		{EquipmentID: equipment[0].ID, Type: "Плановое ТО", IntervalDays: 30, LastPerformed: &lastPerformed, NextDate: now.AddDate(0, 0, 25), Description: "Ежемесячное ТО токарного станка: смазка, проверка люфтов, замена фильтров", IsActive: true, CreatedBy: adminID},
		{EquipmentID: equipment[2].ID, Type: "Плановое ТО", IntervalDays: 90, LastPerformed: &lastPerformed, NextDate: now.AddDate(0, 2, 0), Description: "Квартальное ТО компрессора: замена масла, проверка клапанов, чистка радиатора", IsActive: true, CreatedBy: adminID},
		{EquipmentID: equipment[4].ID, Type: "Осмотр", IntervalDays: 7, LastPerformed: &lastPerformed, NextDate: overdueDate, Description: "Еженедельный осмотр конвейера: проверка натяжения ленты, состояние роликов", IsActive: true, CreatedBy: engineerID},
		{EquipmentID: equipment[9].ID, Type: "Плановое ТО", IntervalDays: 180, LastPerformed: &lastPerformed, NextDate: now.AddDate(0, 5, 0), Description: "Полугодовое ТО трансформатора: проверка масла, замер сопротивления изоляции", IsActive: true, CreatedBy: adminID},
	}

	if err := db.Create(&schedules).Error; err != nil {
		return fmt.Errorf("failed to create maintenance schedules: %w", err)
	}
	log.Printf("Created %d maintenance schedules", len(schedules))

	// --- Maintenance Logs ---
	logs := []models.MaintenanceLog{
		{EquipmentID: equipment[0].ID, Type: "repair", Description: "Замена подшипника шпинделя", PerformedBy: technicianID, PerformedAt: now.AddDate(0, 0, -10), DurationHours: 4.5},
		{EquipmentID: equipment[6].ID, Type: "repair", Description: "Замена уплотнительных колец насоса", PerformedBy: technicianID, PerformedAt: now.AddDate(0, 0, -3), DurationHours: 2.0},
		{EquipmentID: equipment[1].ID, Type: "maintenance", Description: "Плановая замена масла", PerformedBy: technicianID, PerformedAt: now.AddDate(0, 0, -5), DurationHours: 1.5},
	}

	if err := db.Create(&logs).Error; err != nil {
		return fmt.Errorf("failed to create maintenance logs: %w", err)
	}
	log.Printf("Created %d maintenance logs", len(logs))

	// --- Parts ---
	parts := []models.Part{
		{Name: "Подшипник 6208", ArticleNumber: "ПШ-6208", Quantity: 15, Unit: "шт", MinQuantity: 5, Location: "Стеллаж А-1"},
		{Name: "Ремень клиновой", ArticleNumber: "РК-А1500", Quantity: 8, Unit: "шт", MinQuantity: 3, Location: "Стеллаж А-2"},
		{Name: "Масло индустриальное И-40А", ArticleNumber: "МИ-40А-20", Quantity: 120, Unit: "л", MinQuantity: 40, Location: "Склад ГСМ"},
		{Name: "Фильтр масляный", ArticleNumber: "ФМ-001", Quantity: 10, Unit: "шт", MinQuantity: 4, Location: "Стеллаж Б-1"},
		{Name: "Прокладка уплотнительная", ArticleNumber: "ПУ-50", Quantity: 25, Unit: "шт", MinQuantity: 10, Location: "Стеллаж Б-2"},
		{Name: "Болт М12x40", ArticleNumber: "БМ-12-40", Quantity: 200, Unit: "шт", MinQuantity: 50, Location: "Стеллаж В-1"},
		{Name: "Электрод сварочный 3мм", ArticleNumber: "ЭС-3-5", Quantity: 50, Unit: "кг", MinQuantity: 10, Location: "Стеллаж В-2"},
		{Name: "Смазка ЛИТОЛ-24", ArticleNumber: "СМ-Л24-1", Quantity: 30, Unit: "кг", MinQuantity: 10, Location: "Склад ГСМ"},
	}

	if err := db.Create(&parts).Error; err != nil {
		return fmt.Errorf("failed to create parts: %w", err)
	}
	log.Printf("Created %d parts", len(parts))

	// --- Summary ---
	fmt.Println("=== Seed Summary ===")
	fmt.Printf("Users:                %d\n", len(users))
	fmt.Printf("Equipment:            %d\n", len(equipment))
	fmt.Printf("Repair Requests:      %d\n", len(repairRequests))
	fmt.Printf("Maintenance Schedules:%d\n", len(schedules))
	fmt.Printf("Maintenance Logs:     %d\n", len(logs))
	fmt.Printf("Parts:                %d\n", len(parts))

	return nil
}
