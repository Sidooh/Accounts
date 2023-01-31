package invites

import (
	"golang.org/x/exp/slices"
	"strings"
)

func GetInvites(with string, limit int) (interface{}, error) {
	relations := strings.Split(with, ",")

	if slices.Contains(relations, "inviter") {
		return ReadAllWithInviter(limit)
	} else {
		return ReadAll(limit)
	}
}

func GetInviteById(id uint, with string) (interface{}, error) {
	relations := strings.Split(with, ",")

	if slices.Contains(relations, "account") && slices.Contains(relations, "inviter") {
		return ReadWithAccountAndInviter(id)
	} else if slices.Contains(relations, "account") {
		return ReadWithAccount(id)
	} else {
		return ReadById(id)
	}
}
