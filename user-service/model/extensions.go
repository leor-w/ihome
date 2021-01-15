package model

import (
	"time"

	"gorm.io/gorm"
)

// BeforeSave 用户保存前的钩子
func (model *User) BeforeSave(db *gorm.DB) error {
	model.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate 用户创建前的钩子
func (model *User) BeforeCreate(db *gorm.DB) error {
	timestamp := time.Now()
	model.CreatedAt = timestamp
	model.UpdatedAt = timestamp
	return nil
}
