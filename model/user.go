package model

import (
	"time"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	ID           uint64 `gorm:"primary_key,AUTO_INCREMENT"`
	Email        string `gorm:"size:50;not null;unique_index" json:"email"`
	PasswordHash string `gorm:"size:100;not null"`
	Permission   string `gorm:"size:50;not null;default:'normal'" json:"permission"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	// DeletedAt for soft delete supported by gorm
	DeletedAt *time.Time
}

func newUserDao(db *gorm.DB) *gormUserDao {
	return &gormUserDao{
		db: db,
	}
}

type gormUserDao struct {
	db *gorm.DB
}

func (g *gormUserDao) schema() error {
	return g.db.AutoMigrate(user{}).Error
}

func (g *gormUserDao) index() error {
	return nil
}

func (g *gormUserDao) Create(email, password, permission string) (*User, error) {
	var output User

	ph, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	r := user{
		Email:        email,
		PasswordHash: string(ph),
		Permission:   permission,
	}
	if err = g.db.Model(&user{}).
		Where("email = ?", email).
		Assign(r).
		Create(&r).
		Error; err != nil {
		return nil, err
	}
	if err := copier.Copy(&output, &r); err != nil {
		return nil, err
	}
	return &output, nil
}

func (g *gormUserDao) GetByEmail(email string) (*User, error) {
	var rec user
	var output User
	query := g.db.
		Model(&user{}).
		Where("email = ?", email).
		First(&rec)
	if query.RecordNotFound() {
		// return nil if record is not found
		return nil, nil
	}
	if err := copier.Copy(&output, &rec); err != nil {
		return nil, err
	}
	return &output, query.Error
}

func (g *gormUserDao) GetByID(id string) (*User, error) {
	var rec user
	var output User
	query := g.db.
		Model(&user{}).
		Where("id = ?", id).
		First(&rec)
	if query.RecordNotFound() {
		// return nil if record is not found
		return nil, nil
	}
	if err := copier.Copy(&output, &rec); err != nil {
		return nil, err
	}
	return &output, query.Error
}
