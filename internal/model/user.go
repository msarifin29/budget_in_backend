package model

import (
	"time"
)

type User struct {
	Uid       string     `json:"uid"`
	UserName  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Photo     string     `json:"photo"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	TypeUser  string     `json:"type_user"`
	Balance   float64    `json:"balance"`  // Depreceted 'old entity' will be remove later
	Savings   float64    `json:"savings"`  // Depreceted 'old entity' will be remove later
	Cash      float64    `json:"cash"`     // Depreceted 'old entity' will be remove later
	Debts     float64    `json:"debts"`    // Depreceted 'old entity' will be remove later
	Currency  string     `json:"currency"` // Depreceted 'old entity' will be remove later
	Status    string     `json:"status"`
}

type CreateUserRequest struct {
	UserName    string  `json:"username" binding:"required"`
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required,min=6"`
	AccountName string  `json:"account_name"`
	TypeUser    string  `json:"type_user" binding:"required,type_user"`
	Balance     float64 `json:"balance"`
	Savings     float64 `json:"savings"`
	Cash        float64 `json:"cash"`
	Debts       float64 `json:"debts"`
	Currency    string  `json:"currency"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserResponse struct {
	Uid       string `json:"uid"`
	UserName  string `json:"username"`
	AccountId string `json:"account_id"`
}

type TokenUserResponse struct {
	Token   string       `json:"token"`
	UserRes UserResponse `json:"user"`
}
type ResetPasswordRequest struct {
	Uid         string `json:"uid"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type UserProfileResponse struct {
	Uid         string     `json:"uid"`
	UserName    string     `json:"username"`
	Email       string     `json:"email"`
	AccountName string     `json:"account_name"`
	Photo       string     `json:"photo"`
	TypeUser    string     `json:"type_user"`
	Balance     float64    `json:"balance"`
	Savings     float64    `json:"savings"`
	Cash        float64    `json:"cash"`
	Debts       float64    `json:"debts"`
	Currency    string     `json:"currency"`
	Status      string     `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
type AccountUser struct {
	Uid         string    `json:"uid"`
	AccountId   string    `json:"account_id"`
	UserName    string    `json:"username"`
	AccountName string    `json:"account_name"`
	Email       string    `json:"email"`
	Photo       string    `json:"photo"`
	TypeUser    string    `json:"type_user"`
	Balance     float64   `json:"balance"`
	Savings     float64   `json:"savings"`
	Cash        float64   `json:"cash"`
	Debts       float64   `json:"debts"`
	Currency    string    `json:"currency"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateUserRequest struct {
	Uid      string `json:"uid" binding:"required"`
	UserName string `json:"username" binding:"required"`
}

type UserRequest struct {
	Uid string `uri:"uid"`
}
type CheckEmail struct {
	Email string `uri:"email"`
}

type NonActiveUserParams struct {
	Uid   string `json:"uid" binding:"required"`
	Email string `json:"email" binding:"required"`
}
type NonActiveUserRequest struct {
	Uid string `json:"uid"`
}

// type Regencies struct {
// 	ProvinceId string `json:"province_id"`
// 	Name       string `json:"name"`
// 	AltName    string `json:"alt_name"`
// 	Latitude   string `json:"latitude"`
// 	Longitude  string `json:"longitude"`
// }

// type Occupations struct {
// 	OccupationId string `json:"occupation_id"`
// 	Title        string `json:"title"`
// }

type EmailUserRequest struct {
	Email string `json:"email" binding:"required,email"`
}
