package account

import (
	"accounts.sidooh/db"
	"accounts.sidooh/util"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	viper.Set("APP_ENV", "TEST")

	conn := db.NewConnection()

	err := conn.AutoMigrate(&Model{})
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func createRandomAccount(t *testing.T, phone string) Model {
	arg := Model{
		Phone: phone,
	}

	account, err := Create(nil, arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Phone, account.Phone)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestAll(t *testing.T) {
	account1 := createRandomAccount(t, util.RandomPhone())
	account2 := createRandomAccount(t, util.RandomPhone())

	accounts, err := All()
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.GreaterOrEqual(t, len(accounts), 2)

	require.Equal(t, accounts[len(accounts)-2], account1)
	require.Equal(t, accounts[len(accounts)-1], account2)
}

func TestCreate(t *testing.T) {
	createRandomAccount(t, util.RandomPhone())
}

func TestById(t *testing.T) {
	account1 := createRandomAccount(t, util.RandomPhone())
	account2, err := ById(account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Phone, account2.Phone)
	require.Equal(t, account1.Pin, account2.Pin)
	require.Equal(t, account1.Active, account2.Active)
	require.Equal(t, account1.TelcoID, account2.TelcoID)

	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}

func TestByPhone(t *testing.T) {
	account1 := createRandomAccount(t, util.RandomPhone())
	account2, err := ByPhone(account1.Phone)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Phone, account2.Phone)
	require.Equal(t, account1.Pin, account2.Pin)
	require.Equal(t, account1.Active, account2.Active)
	require.Equal(t, account1.TelcoID, account2.TelcoID)

	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}
