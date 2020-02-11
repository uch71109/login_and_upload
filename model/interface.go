package model

import (
	"time"

	"github.com/gobuffalo/buffalo/binding"
)

const (
	// RoleAdmin defines admin user
	RoleAdmin = "admin"
	// RoleNormal defines normal user
	RoleNormal = "normal"
)

// DataOp is an set for all dao modules
type DataOp struct {
	User UserDao
}

type dao interface {
	schema() error
	index() error
}

// UserDao is an interface for User table operation
type UserDao interface {
	Create(email, password, permission string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByID(id string) (*User, error)
}

// User defines User model
type User struct {
	Email        string
	PasswordHash string
	Password     string
	Permission   string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// File defines File model
type File struct {
	UploadFile binding.File
}
