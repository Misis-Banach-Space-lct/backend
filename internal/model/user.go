package model

import (
	"context"
	"time"
)

var UsersTableName = "users"

type User struct {
	Id        int       `json:"id"`       // serial
	Username  string    `json:"username"` // unique
	Password  string    `json:"-"`        // not null
	Role      string    `json:"role"`     // admin, viewer
	Email     string    `json:"email"`    // unique
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	CreatedAt time.Time `json:"createdAt"` // default = current timestamp
	UpdatedAt time.Time `json:"updatedAt"` // default = current timestamp
}

type UserLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// accessible only for admins
type UserCreate struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Role      string `json:"role" validate:"default=viewer,oneof=admin viewer"`
	Email     string `json:"email" validate:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserGroupUpdate struct {
	Action  string `json:"action" validate:"required,oneof=add remove"`
	UserId  int    `json:"userId" validate:"required,gt=1"`
	GroupId int    `json:"groupId" validate:"gte=0,default=0"`
}

type UserRepository interface {
	InsertOne(c context.Context, userData UserCreate) error
	// InsertMany(c context.Context, usersData []UserCreate) error
	FindOne(c context.Context, filter string, value any) (User, error)
	FindMany(c context.Context, filter string, value any) ([]User, error)
	UpdateOne(c context.Context, userData User) error
	AddToGroup(c context.Context, userId, groupId int) error
	RemoveFromGroup(c context.Context, userId, groupId int) error
}
