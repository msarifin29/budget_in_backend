package model

import (
	"time"
)

type User struct {
	Uid       string    `json:"uid"`
	UserName  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Photo     string    `json:"photo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TypeUser  string    `json:"type_user"`
	Balance   float64   `json:"balance"`
	Savings   float64   `json:"savings"`
	Cash      float64   `json:"cash"`
	Debts     float64   `json:"debts"`
	Currency  string    `json:"currency"`
}

type CreateUserRequest struct {
	UserName string  `json:"username" binding:"required"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	TypeUser string  `json:"type_user" binding:"required,type_user"`
	Balance  float64 `json:"balance" binding:"required"`
	Savings  float64 `json:"savings"`
	Cash     float64 `json:"cash" binding:"required"`
	Debts    float64 `json:"debts"`
	Currency string  `json:"currency"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
	Uid      string `json:"uid"`
	UserName string `json:"username"`
}

type TokenUserResponse struct {
	Token   string       `json:"token"`
	UserRes UserResponse `json:"user"`
}

type UserProfileResponse struct {
	Uid      string `json:"uid"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	// Regency    Regencies   `json:"regency"`
	// Occupation Occupations `json:"occupation"`
	Photo     string    `json:"photo"`
	TypeUser  string    `json:"type_user"`
	Balance   float64   `json:"balance"`
	Savings   float64   `json:"savings"`
	Cash      float64   `json:"cash"`
	Debts     float64   `json:"debts"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateUserRequest struct {
	Uid      string `json:"uid" binding:"required"`
	UserName string `json:"username"`
}

type UserRequest struct {
	Uid string `uri:"uid"`
}

type Regencies struct {
	ProvinceId string `json:"province_id"`
	Name       string `json:"name"`
	AltName    string `json:"alt_name"`
	Latitude   string `json:"latitude"`
	Longitude  string `json:"longitude"`
}

type Occupations struct {
	OccupationId string `json:"occupation_id"`
	Title        string `json:"title"`
}

type EmailUserRequest struct {
	Email string `json:"email" binding:"required,email"`
}
