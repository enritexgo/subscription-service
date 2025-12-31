package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"subscription-service/internal/model"
	"subscription-service/internal/repository"
)

type SubscriptionService struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func parseMonthYear(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, errors.New("пустая дата")
	}
	return time.Parse("01-2006", s)
}

func (s *SubscriptionService) Create(req *model.CreateSubscriptionRequest) (*model.Subscription, error) {
	startDate, err := parseMonthYear(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("неверный формат start_date: %v", err)
	}

	var endDate sql.NullTime
	if req.EndDate != "" {
		endDate.Time, err = parseMonthYear(req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("неверный формат end_date: %v", err)
		}
		endDate.Valid = true
	}

	sub := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	err = s.repo.Create(sub)
	if err != nil {
		return nil, fmt.Errorf("ошибка сохранения в БД: %v", err)
	}

	return sub, nil
}

func (s *SubscriptionService) GetByID(id int) (*model.Subscription, error) {
	return s.repo.GetByID(id)
}

func (s *SubscriptionService) Update(id int, req *model.CreateSubscriptionRequest) (*model.Subscription, error) {
	startDate, err := parseMonthYear(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("неверный формат start_date: %v", err)
	}

	var endDate sql.NullTime
	if req.EndDate != "" {
		endDate.Time, err = parseMonthYear(req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("неверный формат end_date: %v", err)
		}
		endDate.Valid = true
	}

	sub := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	err = s.repo.Update(id, sub)
	if err != nil {
		return nil, fmt.Errorf("ошибка обновления: %v", err)
	}

	updated, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения обновленной записи: %v", err)
	}

	return updated, nil
}

func (s *SubscriptionService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *SubscriptionService) List(filters map[string]interface{}) ([]*model.Subscription, error) {
	return s.repo.List(filters)
}

func (s *SubscriptionService) CalculateTotalCost(
	userID, serviceName, periodStart, periodEnd string,
) (int, error) {
	startDate, err := parseMonthYear(periodStart)
	if err != nil {
		return 0, fmt.Errorf("неверный формат period_start: %v", err)
	}

	endDate, err := parseMonthYear(periodEnd)
	if err != nil {
		return 0, fmt.Errorf("неверный формат period_end: %v", err)
	}

	return s.repo.CalculateTotalCost(userID, serviceName, startDate, endDate)
}
