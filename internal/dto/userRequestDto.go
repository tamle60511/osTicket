package dto

import "ecommerce/internal/domain"

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignup struct {
	UserLogin
	Phone        string      `json:"phone"`
	UserName     string      `json:"user_name"`
	DepartmentID string      `json:"department_id"`
	Role         domain.Role `json:"role"`
}
