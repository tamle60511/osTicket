package domain

import "time"

type Role string

const (
	Admin   Role = "admin"
	Staff   Role = "staff"
	Manager Role = "manager"
	IT      Role = "it"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserName     string    `json:"user_name"`
	Email        string    `json:"email" gorm:"index;unique; not null"`
	Password     string    `json:"password"`
	Phone        string    `json:"phone"`
	DepartmentID string    `json:"department_id"`
	Expiry       time.Time `json:"expiry"`
	Verified     bool      `json:"verified" gorm:"default:false"`
	Role         Role      `json:"role" gorm:"default:staff"`
	CreatedAt    time.Time `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"default:current_timestamp"`
}
