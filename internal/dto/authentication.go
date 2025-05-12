package dto

import "time"

type SignUpRequest struct {
	Email        string `json:"email" validate:"required,email"`
	PhotoProfile string `json:"photoProfile"`
	PhoneNumber  string `json:"phoneNumber" validate:"required"`
	Address      string `json:"address" validate:"required"`
	Salary       string `json:"salary" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Password     string `json:"password" validate:"required,min=8"`
	RoleId       uint64 `json:"roleId" validate:"required"`
}

type SignUpResponse struct {
	Id           string       `json:"id"`
	PhotoProfile string       `json:"photoProfile"`
	Email        string       `json:"email"`
	Role         RoleResponse `json:"role"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type SignInResponse struct {
	TokenType    string        `json:"tokenType"`
	Role         RoleResponse  `json:"role"`
	Name         string        `json:"name"`
	PhotoProfile string        `json:"photoProfile"`
	AccessToken  string        `json:"accessToken"`
	ExpiredAt    time.Duration `json:"expiredAt"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"oldPassword" validate:"required,min=8"`
	NewPassword     string `json:"newPassword" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8"`
}

type UpdateAccountRequest struct {
	Email        string `json:"email" validate:"required,email"`
	RoleId       uint64 `json:"roleId" validate:"required"`
	PhotoProfile string `json:"photoProfile" validate:"required"`
}

type ChangePasswordResponse struct {
	Id           string       `json:"id"`
	PhotoProfile string       `json:"photoProfile"`
	Email        string       `json:"email"`
	Role         RoleResponse `json:"role"`
}

type AccountResponse struct {
	Id           string       `json:"id"`
	Email        string       `json:"email"`
	PhotoProfile string       `json:"photoProfile"`
	Role         RoleResponse `json:"role"`
}
