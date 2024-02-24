package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Uid          uuid.UUID `json:"uid"`
	UserName     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	ProvinceId   string    `json:"province_id"`
	OccupationId string    `json:"occupation_id"`
	Photo        string    `json:"photo"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	TypeUser     string    `json:"type_user"`
	Balance      string    `json:"balance"`
	Savings      string    `json:"savings"`
	Cash         string    `json:"cash"`
	Debts        string    `json:"Debts"`
	Currency     string    `json:"currency"`
}

type CreateUserRequest struct {
	UserName string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	TypeUser string `json:"type_user" binding:"required"`
	Balance  string `json:"balance" inding:"required"`
	Savings  string `json:"savings" inding:"required"`
	Cash     string `json:"cash" inding:"required"`
	Debts    string `json:"Debts" inding:"required"`
	Currency string `json:"currency" inding:"required"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Uid      uuid.UUID `json:"uid"`
	UserName string    `json:"username"`
}

type TokenUserResponse struct {
	Token   string       `json:"token"`
	UserRes UserResponse `json:"user"`
}

type UserProfileResponse struct {
	Uid        uuid.UUID   `json:"uid"`
	UserName   string      `json:"username"`
	Email      string      `json:"email"`
	Regency    Regencies   `json:"regency"`
	Occupation Occupations `json:"occupation"`
	Photo      string      `json:"photo"`
	TypeUser   string      `json:"type_user"`
	Balance    string      `json:"balance"`
	Savings    string      `json:"savings"`
	Cash       string      `json:"cash"`
	Debts      string      `json:"Debts"`
	Currency   string      `json:"currency"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
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
