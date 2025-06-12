package models

import "gorm.io/gorm"

type SessionLogin struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	SessionID string `json:"session_id"`
}
