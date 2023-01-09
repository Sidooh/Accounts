package invites

import (
	"accounts.sidooh/pkg/db"
	"accounts.sidooh/pkg/entities"
	"accounts.sidooh/utils"
	"accounts.sidooh/utils/constants"
	"github.com/SamuelTissot/sqltime"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	utils.SetupConfig("../../")

	viper.Set("APP_ENV", "TEST")

	db.Init()
	conn := db.Connection()

	err := conn.AutoMigrate(&entities.Invite{}, &entities.Account{})
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func createRandomInvite(t *testing.T, phone string) entities.Invite {
	account := entities.Account{
		Phone: utils.RandomPhone(),
	}
	res := db.Connection().Omit("UserID").Create(&account)

	require.NoError(t, res.Error)

	arg := entities.Invite{
		InviterID: account.ID,
		Phone:     phone,
	}

	invite, err := Create(arg)
	require.NoError(t, err)
	require.NotEmpty(t, invite)

	require.Equal(t, arg.InviterID, account.ID)
	require.Equal(t, arg.Phone, invite.Phone)
	require.Equal(t, constants.PENDING, invite.Status)

	require.NotZero(t, invite.ID)
	require.NotZero(t, invite.CreatedAt)

	return invite
}

func refreshDatabase() {
	//Start clean slate
	conn := db.Connection()
	conn.Where("1 = 1").Delete(&entities.Invite{})
}

func TestAll(t *testing.T) {
	refreshDatabase()

	invite1 := createRandomInvite(t, utils.RandomPhone())
	invite2 := createRandomInvite(t, utils.RandomPhone())

	invites, err := ReadAll(constants.DEFAULT_QUERY_LIMIT)
	require.NoError(t, err)
	require.NotEmpty(t, invites)
	require.Equal(t, len(invites), 2)

	require.Equal(t, invites[1], invite1)
	require.Equal(t, invites[0], invite2)
}

func TestCreate(t *testing.T) {
	createRandomInvite(t, utils.RandomPhone())
}

func TestCreateWithInviteCode(t *testing.T) {
	phone := utils.RandomPhone()

	account := entities.Account{
		Phone:      phone,
		InviteCode: "TEST",
	}
	res := db.Connection().Omit("UserID").Create(&account)

	require.NoError(t, res.Error)
	require.NotEmpty(t, account)
}

func TestById(t *testing.T) {
	invite1 := createRandomInvite(t, utils.RandomPhone())
	invite2, err := ReadById(invite1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, invite2)

	require.Equal(t, invite1.AccountID, invite2.AccountID)
	require.Equal(t, invite1.InviterID, invite2.InviterID)
	require.Equal(t, invite1.Phone, invite2.Phone)
	require.Equal(t, invite1.Status, invite2.Status)

	require.WithinDuration(t, invite1.CreatedAt.Time, invite2.CreatedAt.Time, time.Second)
}

func TestByPhone(t *testing.T) {
	invite1 := createRandomInvite(t, utils.RandomPhone())
	invite2, err := ReadByPhone(invite1.Phone)
	require.NoError(t, err)
	require.NotEmpty(t, invite2)

	require.Equal(t, invite1.AccountID, invite2.AccountID)
	require.Equal(t, invite1.InviterID, invite2.InviterID)
	require.Equal(t, invite1.Phone, invite2.Phone)
	require.Equal(t, invite1.Status, invite2.Status)

	require.WithinDuration(t, invite1.CreatedAt.Time, invite2.CreatedAt.Time, time.Second)
}

func TestUnexpiredByPhone(t *testing.T) {
	//Start clean slate
	refreshDatabase()

	// Scenario 1 - active status
	phone := utils.RandomPhone()

	activeInvite := createRandomInvite(t, phone)
	activeInvite.Status = constants.ACTIVE
	activeInvite.Save()
	require.NotEmpty(t, activeInvite)

	invite, err := ReadUnexpiredByPhone(phone)
	require.Error(t, err)
	require.Empty(t, invite)

	// Scenario 2 - pending status
	phone = utils.RandomPhone()

	pendingInvite := createRandomInvite(t, phone)
	require.NotEmpty(t, pendingInvite)

	invite, err = ReadUnexpiredByPhone(phone)
	require.NoError(t, err)
	require.NotEmpty(t, invite)

	require.Equal(t, invite, pendingInvite)

	// Scenario 3 - time has passed
	phone = utils.RandomPhone()

	timeExpiredInvite := createRandomInvite(t, phone)
	timeExpiredInvite.CreatedAt = sqltime.Time{
		Time: time.Now().Add(-48 * time.Hour),
	}
	timeExpiredInvite.Save()
	require.NotEmpty(t, timeExpiredInvite)

	invite, err = ReadUnexpiredByPhone(phone)
	require.Error(t, err)
	require.Empty(t, invite)
}

func TestMarkExpired(t *testing.T) {
	//Start clean slate
	refreshDatabase()

	// We can have 3 states of invites
	// Active invite
	// pending but non time-expired invite
	// pending and time-expired <- to be removed
	activeInvite := createRandomInvite(t, utils.RandomPhone())
	activeInvite.Status = constants.ACTIVE
	activeInvite.Save()
	require.NotEmpty(t, activeInvite)

	pendingInvite := createRandomInvite(t, utils.RandomPhone())
	require.NotEmpty(t, pendingInvite)

	timeExpiredInvite := createRandomInvite(t, utils.RandomPhone())
	timeExpiredInvite.CreatedAt = sqltime.Time{
		Time: time.Now().Add(-48 * time.Hour),
	}
	timeExpiredInvite.Save()
	require.NotEmpty(t, timeExpiredInvite)

	err := MarkExpired()
	require.NoError(t, err)

	invites, err := ReadUnexpired()
	require.NoError(t, err)
	require.NotEmpty(t, invites)
	require.Equal(t, len(invites), 2)

	require.Equal(t, invites[0], activeInvite)
	require.Equal(t, invites[1], pendingInvite)
}
