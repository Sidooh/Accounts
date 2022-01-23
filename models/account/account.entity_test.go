package account

import (
	"accounts.sidooh/db"
	"accounts.sidooh/util"
	"database/sql"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var datastore = new(db.DB)

func TestMain(m *testing.M) {
	viper.Set("APP_ENV", "TEST")

	datastore = db.NewConnection()

	err := datastore.Conn.AutoMigrate(&Model{})
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func CreateRandomAccount(t *testing.T, phone string) Model {
	arg := Model{
		Phone: phone,
	}

	account, err := Create(datastore, arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Phone, account.Phone)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestAll(t *testing.T) {
	account1 := CreateRandomAccount(t, util.RandomPhone())
	account2 := CreateRandomAccount(t, util.RandomPhone())

	accounts, err := All(datastore)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.GreaterOrEqual(t, len(accounts), 2)

	require.Equal(t, accounts[len(accounts)-2], account1)
	require.Equal(t, accounts[len(accounts)-1], account2)
}

func TestCreate(t *testing.T) {
	CreateRandomAccount(t, util.RandomPhone())
}

func TestById(t *testing.T) {
	account1 := CreateRandomAccount(t, util.RandomPhone())
	account2, err := ById(datastore, account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Phone, account2.Phone)
	require.Equal(t, account1.Pin, account2.Pin)
	require.Equal(t, account1.Active, account2.Active)
	require.Equal(t, account1.TelcoID, account2.TelcoID)

	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}

func TestByPhone(t *testing.T) {
	account1 := CreateRandomAccount(t, util.RandomPhone())
	account2, err := ByPhone(datastore, account1.Phone)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Phone, account2.Phone)
	require.Equal(t, account1.Pin, account2.Pin)
	require.Equal(t, account1.Active, account2.Active)
	require.Equal(t, account1.TelcoID, account2.TelcoID)

	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}

func TestUpdate(t *testing.T) {
	account := CreateRandomAccount(t, util.RandomPhone())
	require.NotEmpty(t, account)

	result := account.Update(datastore, "pin", "new_pin")
	require.NoError(t, result.Error)

	require.Equal(t, account.Pin, sql.NullString{String: "new_pin", Valid: true})
	require.WithinDuration(t, account.UpdatedAt.Time, account.UpdatedAt.Time, time.Second)

	result = account.Update(datastore, "pins", "new_pin")
	require.Error(t, result.Error)
}
