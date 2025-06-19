package dto

import "time"

type SignUpRequest struct {
	Email        string   `json:"email" validate:"required,email"`
	Username     string   `json:"username" validate:"required"`
	LocationId   uint64   `json:"locationId" validate:"required"`
	RoleId       uint64   `json:"roleId" validate:"required"`
	PlacementIds []uint64 `json:"placementIds" validate:"required"`
	PhotoProfile string   `json:"photoProfile"`
	PhoneNumber  string   `json:"phoneNumber" validate:"required"`
	Address      string   `json:"address" validate:"required"`
	Salary       string   `json:"salary" validate:"required"`
	Name         string   `json:"name" validate:"required"`
	Password     string   `json:"password" validate:"required,min=4"`
}

type SignUpResponse struct {
	Id           string           `json:"id"`
	PhotoProfile string           `json:"photoProfile"`
	Email        string           `json:"email"`
	Role         RoleResponse     `json:"role"`
	Location     LocationResponse `json:"location"`
}

type SignInRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=4"`
}

type SignInResponse struct {
	Id           string           `json:"id"`
	TokenType    string           `json:"tokenType"`
	Role         RoleResponse     `json:"role"`
	Location     LocationResponse `json:"location"`
	Name         string           `json:"name"`
	PhotoProfile string           `json:"photoProfile"`
	AccessToken  string           `json:"accessToken"`
	ExpiredAt    time.Duration    `json:"expiredAt"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"oldPassword" validate:"required,min=4"`
	NewPassword     string `json:"newPassword" validate:"required,min=4"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=4"`
}

type ChangePasswordResponse struct {
	Id           string           `json:"id"`
	PhotoProfile string           `json:"photoProfile"`
	Email        string           `json:"email"`
	Location     LocationResponse `json:"location"`
	Role         RoleResponse     `json:"role"`
}
