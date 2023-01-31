package entities

import (
	"accounts.sidooh/pkg/db"
)

type Invite struct {
	ModelID

	Phone     string `json:"phone" gorm:"uniqueIndex; size:16"`
	Status    string `json:"status" gorm:"size:16"`
	AccountID uint   `json:"account_id,omitempty"`
	InviterID uint   `json:"inviter_id"`

	ModelTimeStamps
}

type InviteWithAccount struct {
	Invite

	Account *Account `json:"account"`
}

type InviteWithInviter struct {
	Invite

	Inviter *AccountWithUser `json:"inviter"`
}

type InviteWithAccountAndInviter struct {
	Invite

	//TODO: Add a constraint to ensure these 2 can't have same values
	// 	i.e. a user can't invite themselves, obviously
	Account *Account `json:"account"`
	Inviter *Account `json:"inviter"`
}

func (*Invite) TableName() string {
	return "invites"
}
func (InviteWithAccount) TableName() string {
	return "invites"
}
func (InviteWithAccountAndInviter) TableName() string {
	return "invites"
}

func (r *Invite) Save() interface{} {
	return db.Connection().Save(&r)
}
