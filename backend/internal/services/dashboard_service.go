package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"toir-app/internal/models"

	"gorm.io/gorm"
)

// DashboardDBQuerier абстракция для запросов дашборда к БД.
type DashboardDBQuerier interface {
	CountOpenRequests(ctx context.Context) (int64, error)
	CountOverdueSchedules(ctx context.Context) (int64, error)
	AvgRepairTimeThisMonth(ctx context.Context) (float64, error)
	CountCompletedThisMonth(ctx context.Context) (int64, error)
	CountTotalEquipment(ctx context.Context) (int64, error)
	CountActiveEquipment(ctx context.Context) (int64, error)
}

// DashboardServiceImpl реализует DashboardService через запросы к БД.
type DashboardServiceImpl struct {
	querier DashboardDBQuerier
}

// NewDashboardService создаёт сервис дашборда.
func NewDashboardService(querier DashboardDBQuerier) *DashboardServiceImpl {
	return &DashboardServiceImpl{querier: querier}
}

// GetMetrics собирает агрегированные метрики для дашборда.
func (s *DashboardServiceImpl) GetMetrics(ctx context.Context) (*models.DashboardMetrics, error) {
	openRequests, err := s.querier.CountOpenRequests(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count open requests: %w", err)
	}

	overdueSchedules, err := s.querier.CountOverdueSchedules(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count overdue schedules: %w", err)
	}

	avgRepairTime, err := s.querier.AvgRepairTimeThisMonth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get avg repair time: %w", err)
	}

	completedThisMonth, err := s.querier.CountCompletedThisMonth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count completed this month: %w", err)
	}

	totalEquipment, err := s.querier.CountTotalEquipment(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count total equipment: %w", err)
	}

	activeEquipment, err := s.querier.CountActiveEquipment(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count active equipment: %w", err)
	}

	return &models.DashboardMetrics{
		OpenRequests:       openRequests,
		OverdueSchedules:   overdueSchedules,
		AvgRepairTimeHrs:   math.Round(avgRepairTime*100) / 100,
		CompletedThisMonth: completedThisMonth,
		TotalEquipment:     totalEquipment,
		ActiveEquipment:    activeEquipment,
	}, nil
}

// GormDashboardQuerier реализует DashboardDBQuerier через GORM.
type GormDashboardQuerier struct {
	db *gorm.DB
}

// NewGormDashboardQuerier создаёт querier на основе GORM.
func NewGormDashboardQuerier(db *gorm.DB) *GormDashboardQuerier {
	return &GormDashboardQuerier{db: db}
}

// CountOpenRequests считает заявки с открытыми статусами.
func (q *GormDashboardQuerier) CountOpenRequests(ctx context.Context) (int64, error) {
	var count int64
	err := q.db.WithContext(ctx).
		Table("repair_requests").
		Where("status NOT IN ?", []string{"completed", "closed"}).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count open requests: %w", err)
	}
	return count, nil
}

// CountOverdueSchedules считает просроченные расписания ТО.
func (q *GormDashboardQuerier) CountOverdueSchedules(ctx context.Context) (int64, error) {
	var count int64
	err := q.db.WithContext(ctx).
		Table("maintenance_schedules").
		Where("next_date < ? AND is_active = ?", time.Now(), true).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count overdue schedules: %w", err)
	}
	return count, nil
}

// AvgRepairTimeThisMonth возвращает среднее время ремонта за текущий месяц.
func (q *GormDashboardQuerier) AvgRepairTimeThisMonth(ctx context.Context) (float64, error) {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var avg *float64
	err := q.db.WithContext(ctx).
		Table("maintenance_logs").
		Where("performed_at >= ?", monthStart).
		Select("AVG(duration_hours)").
		Scan(&avg).Error
	if err != nil {
		return 0, fmt.Errorf("failed to get avg repair time: %w", err)
	}

	if avg == nil {
		return 0, nil
	}
	return *avg, nil
}

// CountCompletedThisMonth считает завершённые заявки за текущий месяц.
func (q *GormDashboardQuerier) CountCompletedThisMonth(ctx context.Context) (int64, error) {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	var count int64
	err := q.db.WithContext(ctx).
		Table("repair_requests").
		Where("status = ? AND completed_at >= ?", "completed", monthStart).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count completed this month: %w", err)
	}
	return count, nil
}

// CountTotalEquipment считает общее количество оборудования.
func (q *GormDashboardQuerier) CountTotalEquipment(ctx context.Context) (int64, error) {
	var count int64
	err := q.db.WithContext(ctx).
		Table("equipment").
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count total equipment: %w", err)
	}
	return count, nil
}

// CountActiveEquipment считает активное оборудование.
func (q *GormDashboardQuerier) CountActiveEquipment(ctx context.Context) (int64, error) {
	var count int64
	err := q.db.WithContext(ctx).
		Table("equipment").
		Where("status = ?", "active").
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count active equipment: %w", err)
	}
	return count, nil
}
