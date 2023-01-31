package entities

import (
	"accounts.sidooh/pkg/db"
	"github.com/SamuelTissot/sqltime"
	"gorm.io/gorm"
)

type User struct {
	ModelID

	Name            string        `json:"name" gorm:"size:64"`
	Username        string        `json:"username,omitempty" gorm:"uniqueIndex; size:16"`
	IdNumber        string        `json:"id_number,omitempty" gorm:"size:16"`
	Status          string        `json:"status,omitempty" gorm:"size:16"`
	Email           string        `json:"email" gorm:"uniqueIndex; size:256; not null"`
	EmailVerifiedAt *sqltime.Time `gorm:"type:timestamp" json:"-"`
	Password        string        `json:"-"`

	ModelTimeStamps
}

func (*User) TableName() string {
	return "users"
}

func (u *User) Save() *gorm.DB {
	return db.Connection().Save(&u)
}

func (u *User) Update(column string, value string) *gorm.DB {
	return db.Connection().Model(&u).Update(column, value)
}
