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

	conn := db.NewConnection()

	err := conn.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func createRandomUser(t *testing.T, password string) User {
	arg := User{
		Email:    util.RandomEmail(),
		Password: password,
	}

	user, err := CreateUser(arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt)

	return user
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

	user2, err := AuthUser(User{Email: user1.Email, Password: password})

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Password, user2.Password)
	require.Equal(t, user1.Status, user2.Status)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.EmailVerifiedAt.Time, user2.EmailVerifiedAt.Time, time.Second)
	require.WithinDuration(t, user1.CreatedAt.Time, user2.CreatedAt.Time, time.Second)
}
