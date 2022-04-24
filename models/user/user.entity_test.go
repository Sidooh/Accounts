package user

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

	db.Init()
	conn := db.Connection()

	err := conn.AutoMigrate(&Model{})
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func createUser(arg Model) (Model, error) {
	user, err := CreateUser(arg)
	return user, err
}

func createRandomUser(t *testing.T, password string) Model {
	arg := Model{
		Email:    util.RandomEmail(),
		Password: password,
	}

	user, err := createUser(arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt)

	return user
}

func refreshDatabase() {
	//Start clean slate
	conn := db.Connection()
	conn.Where("1 = 1").Delete(&Model{})
}

func TestAll(t *testing.T) {
	user1 := createRandomUser(t, util.RandomString(6))
	user2 := createRandomUser(t, util.RandomString(6))

	users, err := All()
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.GreaterOrEqual(t, len(users), 2)

	require.Equal(t, users[len(users)-2], user1)
	require.Equal(t, users[len(users)-1], user2)
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t, util.RandomString(6))
}

func TestFindUserById(t *testing.T) {
	user1 := createRandomUser(t, util.RandomString(6))
	user2, err := FindUserById(user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.Status, user2.Status)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.EmailVerifiedAt.Time, user2.EmailVerifiedAt.Time, time.Second)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}

func TestFindUserByEmail(t *testing.T) {
	user1 := createRandomUser(t, util.RandomString(6))
	user2, err := FindUserByEmail(user1.Email)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.Status, user2.Status)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.EmailVerifiedAt.Time, user2.EmailVerifiedAt.Time, time.Second)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}

func TestAuthUser(t *testing.T) {
	password := util.RandomString(6)
	user1 := createRandomUser(t, password)

	user2, err := AuthUser(Model{Email: user1.Email, Password: password})

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.Status, user2.Status)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.EmailVerifiedAt.Time, user2.EmailVerifiedAt.Time, time.Second)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}

func TestSearchByEmail(t *testing.T) {
	refreshDatabase()

	password := util.RandomString(6)
	arg := Model{
		Email:    "ab@a.a",
		Password: password,
	}
	user1, err := createUser(arg)

	arg = Model{
		Email:    "a@a.a",
		Password: password,
	}
	user2, err := createUser(arg)

	users, err := SearchByEmail("ab")
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, 1, len(users))

	require.Equal(t, users[0], user1)

	users, err = SearchByEmail("a")
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, 2, len(users))

	require.Equal(t, users[1], user2)

	users, err = SearchByEmail("c")
	require.NoError(t, err)
	require.Empty(t, users)

}
