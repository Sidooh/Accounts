package invites

import (
	"accounts.sidooh/utils"
	"strings"
)

func GetInviteById(id uint, with string) (interface{}, error) {
	relations := strings.Split(with, ",")

	if utils.InArray("account", relations) && utils.InArray("inviter", relations) {
		return ReadWithAccountAndInviter(id)
	} else if utils.InArray("account", relations) {
		return ReadWithAccount(id)
	} else {
		return ReadById(id)
	}
}
