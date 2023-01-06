package models

import "github.com/SamuelTissot/sqltime"

type ModelID struct {
	ID uint `gorm:"primaryKey" json:"id"`
}

type ModelTimeStamps struct {
	CreatedAt sqltime.Time `gorm:"type:timestamp" json:"created_at"`
	UpdatedAt sqltime.Time `gorm:"type:timestamp" json:"-"`
}
