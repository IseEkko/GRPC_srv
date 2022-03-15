package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int32     `gorm:"primarykey"` //为什么使用int32，因为这个可以减少外键的失败
	CreateeAt time.Time `gorm:"column:add_time"`
	UpdateAt  time.Time `gorm:"column：update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}
