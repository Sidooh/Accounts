package entities

import (
	"accounts.sidooh/pkg/db"
	"database/sql"
	"gorm.io/gorm"
)

type Account struct {
	ModelID

	Phone     string         `json:"phone" gorm:"uniqueIndex; size:16"`
	Active    bool           `json:"active"`
	Pin       sql.NullString `json:"-"`
	TelcoID   int            `json:"-"`
	InviterID uint           `json:"inviter_id,omitempty"`
	UserID    uint           `json:"user_id,omitempty"`

	InviteCode string `json:"invite_code,omitempty"`

	ModelTimeStamps
}

type AccountWithUser struct {
	Account

	User *User `json:"user"`
}

type AccountWithUserAndInviter struct {
	Account

	User    *User            `json:"user"`
	Inviter *AccountWithUser `json:"inviter"`
}

type InviteModel struct {
	Account

	Level int `json:"level"`
}

func (*Account) TableName() string {
	return "accounts"
}
func (AccountWithUser) TableName() string {
	return "accounts"
}
func (AccountWithUserAndInviter) TableName() string {
	return "accounts"
}

func (a *Account) Save() *gorm.DB {
	return db.Connection().Save(&a)
}

func (a *Account) Update(column string, value string) *gorm.DB {
	return db.Connection().Model(&a).Update(column, value)
}
