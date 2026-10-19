package model

import "time"

type Wallet struct {
	ID        string     `db:"id"`
	UserID    string     `db:"user_id"`
	Balance   float64    `db:"balance"`
	Currency  string     `db:"currency"`
	Status    string     `db:"status"`
	Version   int        `db:"version"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type CreateWalletRequest struct {
	UserID   string  `json:"user_id"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
	Status   string  `json:"status"`
	Version  int     `json:"version"`
}

type UpdatedWalletRequest struct {
	UserID   string  `db:"user_id"`
	Balance  float64 `db:"balance"`
	Currency string  `db:"currency"`
	Status   string  `db:"status"`
	Version  int     `db:"version"`
}
