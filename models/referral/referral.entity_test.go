package referral

import (
	"accounts.sidooh/db"
	Account "accounts.sidooh/models/account"
	"accounts.sidooh/util"
	"accounts.sidooh/util/constants"
	"github.com/SamuelTissot/sqltime"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	util.SetupConfig("../../")

	viper.Set("APP_ENV", "TEST")

	db.Init()
	conn := db.Connection()

	err := conn.AutoMigrate(&Model{}, &Account.Model{})
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func createRandomReferral(t *testing.T, phone string) Model {
	account, err := Account.Create(Account.Model{
		Phone: util.RandomPhone(),
	})
	require.NoError(t, err)

	arg := Model{
		AccountID:    account.ID,
		RefereePhone: phone,
	}

	referral, err := Create(arg)
	require.NoError(t, err)
	require.NotEmpty(t, referral)

	require.Equal(t, arg.AccountID, account.ID)
	require.Equal(t, arg.RefereePhone, referral.RefereePhone)
	require.Equal(t, constants.PENDING, referral.Status)

	require.NotZero(t, referral.ID)
	require.NotZero(t, referral.CreatedAt)

	return referral
}

func refreshDatabase() {
	//Start clean slate
	conn := db.Connection()
	conn.Where("1 = 1").Delete(&Model{})
}

func TestAll(t *testing.T) {
	refreshDatabase()

	referral1 := createRandomReferral(t, util.RandomPhone())
	referral2 := createRandomReferral(t, util.RandomPhone())

	referrals, err := All()
	require.NoError(t, err)
	require.NotEmpty(t, referrals)
	require.Equal(t, len(referrals), 2)

	require.Equal(t, referrals[0], referral1)
	require.Equal(t, referrals[1], referral2)
}

func TestCreate(t *testing.T) {
	createRandomReferral(t, util.RandomPhone())
}

func TestById(t *testing.T) {
	referral1 := createRandomReferral(t, util.RandomPhone())
	referral2, err := ById(referral1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, referral2)

	require.Equal(t, referral1.AccountID, referral2.AccountID)
	require.Equal(t, referral1.RefereeID, referral2.RefereeID)
	require.Equal(t, referral1.RefereePhone, referral2.RefereePhone)
	require.Equal(t, referral1.Status, referral2.Status)

	require.WithinDuration(t, referral1.CreatedAt.Time, referral2.CreatedAt.Time, time.Second)
}

func TestByPhone(t *testing.T) {
	referral1 := createRandomReferral(t, util.RandomPhone())
	referral2, err := ByPhone(referral1.RefereePhone)
	require.NoError(t, err)
	require.NotEmpty(t, referral2)

	require.Equal(t, referral1.AccountID, referral2.AccountID)
	require.Equal(t, referral1.RefereeID, referral2.RefereeID)
	require.Equal(t, referral1.RefereePhone, referral2.RefereePhone)
	require.Equal(t, referral1.Status, referral2.Status)

	require.WithinDuration(t, referral1.CreatedAt.Time, referral2.CreatedAt.Time, time.Second)
}

func TestUnexpiredByPhone(t *testing.T) {
	//Start clean slate
	refreshDatabase()

	// Scenario 1 - active status
	phone := util.RandomPhone()

	activeReferral := createRandomReferral(t, phone)
	activeReferral.Status = constants.ACTIVE
	activeReferral.Save()
	require.NotEmpty(t, activeReferral)

	referral, err := UnexpiredByPhone(phone)
	require.Error(t, err)
	require.Empty(t, referral)

	// Scenario 2 - pending status
	phone = util.RandomPhone()

	pendingReferral := createRandomReferral(t, phone)
	require.NotEmpty(t, pendingReferral)

	referral, err = UnexpiredByPhone(phone)
	require.NoError(t, err)
	require.NotEmpty(t, referral)

	require.Equal(t, referral, pendingReferral)

	// Scenario 3 - time has passed
	phone = util.RandomPhone()

	timeExpiredReferral := createRandomReferral(t, phone)
	timeExpiredReferral.CreatedAt = sqltime.Time{
		Time: time.Now().Add(-48 * time.Hour),
	}
	timeExpiredReferral.Save()
	require.NotEmpty(t, timeExpiredReferral)

	referral, err = UnexpiredByPhone(phone)
	require.Error(t, err)
	require.Empty(t, referral)
}

func TestMarkExpired(t *testing.T) {
	//Start clean slate
	refreshDatabase()

	// We can have 3 states of referrals
	// Active referral
	// pending but non time-expired referral
	// pending and time-expired <- to be removed
	activeReferral := createRandomReferral(t, util.RandomPhone())
	activeReferral.Status = constants.ACTIVE
	activeReferral.Save()
	require.NotEmpty(t, activeReferral)

	pendingReferral := createRandomReferral(t, util.RandomPhone())
	require.NotEmpty(t, pendingReferral)

	timeExpiredReferral := createRandomReferral(t, util.RandomPhone())
	timeExpiredReferral.CreatedAt = sqltime.Time{
		Time: time.Now().Add(-48 * time.Hour),
	}
	timeExpiredReferral.Save()
	require.NotEmpty(t, timeExpiredReferral)

	err := MarkExpired()
	require.NoError(t, err)

	referrals, err := Unexpired()
	require.NoError(t, err)
	require.NotEmpty(t, referrals)
	require.Equal(t, len(referrals), 2)

	require.Equal(t, referrals[0], activeReferral)
	require.Equal(t, referrals[1], pendingReferral)
}
