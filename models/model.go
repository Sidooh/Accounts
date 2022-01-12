package models

import (
	"github.com/SamuelTissot/sqltime"
)

type Model struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	CreatedAt sqltime.Time `gorm:"type:timestamp" json:"-"`
	UpdatedAt sqltime.Time `gorm:"type:timestamp" json:"-"`
}
