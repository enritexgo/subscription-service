package repository

import (
	"database/sql"
	"fmt"
	"time"

	"subscription-service/internal/model"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(sub *model.Subscription) error {
	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	var endDate *time.Time
	if sub.EndDate.Valid {
		endDate = &sub.EndDate.Time
	}

	err := r.db.QueryRow(query,
		sub.ServiceName, sub.Price, sub.UserID,
		sub.StartDate, endDate,
	).Scan(&sub.ID, &sub.CreatedAt, &sub.UpdatedAt)

	return err
}

func (r *SubscriptionRepository) GetByID(id int) (*model.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions WHERE id = $1`

	row := r.db.QueryRow(query, id)

	var sub model.Subscription
	var endDate sql.NullTime

	err := row.Scan(
		&sub.ID, &sub.ServiceName, &sub.Price,
		&sub.UserID, &sub.StartDate, &endDate,
		&sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	sub.EndDate = endDate

	return &sub, nil
}

func (r *SubscriptionRepository) Update(id int, sub *model.Subscription) error {
	query := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, user_id = $3,
		    start_date = $4, end_date = $5, updated_at = NOW()
		WHERE id = $6`

	var endDate *time.Time
	if sub.EndDate.Valid {
		endDate = &sub.EndDate.Time
	}

	_, err := r.db.Exec(query,
		sub.ServiceName, sub.Price, sub.UserID,
		sub.StartDate, endDate, id)

	return err
}

func (r *SubscriptionRepository) Delete(id int) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *SubscriptionRepository) List(filters map[string]interface{}) ([]*model.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at FROM subscriptions WHERE true`
	var args []interface{}
	argIndex := 1

	if userID, ok := filters["user_id"].(string); ok {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, userID)
		argIndex++
	}

	if serviceName, ok := filters["service_name"].(string); ok {
		query += fmt.Sprintf(" AND service_name = $%d", argIndex)
		args = append(args, serviceName)
		argIndex++
	}

	if startAfter, ok := filters["start_after"]; ok {
		query += fmt.Sprintf(" AND start_date >= $%d", argIndex)
		args = append(args, startAfter)
		argIndex++
	}

	if endBefore, ok := filters["end_before"]; ok {
		query += fmt.Sprintf(" AND (end_date IS NULL OR end_date <= $%d)", argIndex)
		args = append(args, endBefore)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*model.Subscription
	for rows.Next() {
		var sub model.Subscription
		var endDate sql.NullTime
		err := rows.Scan(
			&sub.ID, &sub.ServiceName, &sub.Price,
			&sub.UserID, &sub.StartDate, &endDate,
			&sub.CreatedAt, &sub.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sub.EndDate = endDate
		subscriptions = append(subscriptions, &sub)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) CalculateTotalCost(
	userID string,
	serviceName string,
	startDate time.Time,
	endDate time.Time,
) (int, error) {
	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE user_id = $1
		  AND service_name = $2
		  AND start_date <= $4
		  AND (end_date IS NULL OR end_date >= $3)`

	var total int
	err := r.db.QueryRow(query, userID, serviceName, startDate, endDate).Scan(&total)
	return total, err
}
