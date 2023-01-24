package special

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	Id        uint           `gorm:"column:id;primaryKey" json:"id" form:"id"`
	CreatedAt time.Time      `gorm:"column:created_at;" json:"created_at" form:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;" json:"updated_at" form:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;" json:"deleted_at" form:"deleted_at"`
}
