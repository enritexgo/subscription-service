package model

import (
	"database/sql"
	"time"
)

type Subscription struct {
	ID          int          `json:"id"`
	ServiceName string       `json:"service_name"`
	Price       int          `json:"price"`
	UserID      string       `json:"user_id"`
	StartDate   time.Time    `json:"start_date"`
	EndDate     sql.NullTime `json:"end_date,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name" binding:"required,min=1,max=255"`
	Price       int    `json:"price" binding:"required,min=0"`
	UserID      string `json:"user_id" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"` // формат "MM-YYYY"
	EndDate     string `json:"end_date"`                      // опционально, формат "MM-YYYY"
}
