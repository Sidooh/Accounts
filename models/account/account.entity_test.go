package account

import (
	"accounts.sidooh/pkg/db"
	"accounts.sidooh/utils"
	"database/sql"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	viper.Set("APP_ENV", "TEST")

	db.Init()
	conn := db.Connection()

	err := conn.AutoMigrate(&Model{})
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func CreateRandomAccount(t *testing.T, phone string) Model {
	arg := Model{
		Phone: phone,
	}

	account, err := Create(arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Phone, account.Phone)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func refreshDatabase() {
	//Start clean slate
	conn := db.Connection()
	conn.Where("1 = 1").Delete(&Model{})
}

func TestAll(t *testing.T) {
	account1 := CreateRandomAccount(t, utils.RandomPhone())
	account2 := CreateRandomAccount(t, utils.RandomPhone())

	accounts, err := All()
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.GreaterOrEqual(t, len(accounts), 2)

	require.Equal(t, accounts[1], account1)
	require.Equal(t, accounts[0], account2)
}

func TestCreate(t *testing.T) {
	CreateRandomAccount(t, utils.RandomPhone())
}

func TestById(t *testing.T) {
	account1 := CreateRandomAccount(t, utils.RandomPhone())
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
	account1 := CreateRandomAccount(t, utils.RandomPhone())
	account2, err := ByPhone(account1.Phone)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Phone, account2.Phone)
	require.Equal(t, account1.Pin, account2.Pin)
	require.Equal(t, account1.Active, account2.Active)
	require.Equal(t, account1.TelcoID, account2.TelcoID)

	require.WithinDuration(t, account1.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}

func TestUpdate(t *testing.T) {
	account := CreateRandomAccount(t, utils.RandomPhone())
	require.NotEmpty(t, account)

	result := account.Update("pin", "new_pin")
	require.NoError(t, result.Error)

	require.Equal(t, account.Pin, sql.NullString{String: "new_pin", Valid: true})
	require.WithinDuration(t, account.UpdatedAt.Time, account.UpdatedAt.Time, time.Second)

	result = account.Update("pins", "new_pin")
	require.Error(t, result.Error)
}

func TestSearchByPhone(t *testing.T) {
	refreshDatabase()

	account1 := CreateRandomAccount(t, "714611696")
	account2 := CreateRandomAccount(t, "780611696")

	accounts, err := SearchByPhone("714")
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Equal(t, 1, len(accounts))

	require.Equal(t, accounts[0], account1)

	accounts, err = SearchByPhone("6116")
	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Equal(t, 2, len(accounts))

	require.Equal(t, accounts[0], account2)

	accounts, err = SearchByPhone("3")
	require.NoError(t, err)
	require.Empty(t, accounts)

}
