package user

import (
	"accounts.sidooh/db"
	"accounts.sidooh/models"
	"accounts.sidooh/util"
	"errors"
	"fmt"
	"github.com/SamuelTissot/sqltime"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	models.ModelID

	Name            string        `json:"name" gorm:"size:64"`
	Username        string        `json:"username,omitempty" gorm:"uniqueIndex; size:16"`
	IdNumber        string        `json:"id_number,omitempty" gorm:"size:16"`
	Status          string        `json:"status,omitempty" gorm:"size:16"`
	Email           string        `json:"email" gorm:"uniqueIndex; size:256; not null"`
	EmailVerifiedAt *sqltime.Time `gorm:"type:timestamp" json:"-"`
	Password        string        `json:"-"`

	models.ModelTimeStamps
}

func (*Model) TableName() string {
	return "users"
}

func All() ([]Model, error) {
	conn := db.Connection()

	var users []Model
	result := conn.Find(&users)
	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}

func CreateUser(u Model) (Model, error) {
	conn := db.Connection()
	_, err := FindUserByEmail(u.Email)

	if err == nil {
		return Model{}, errors.New("email is already taken")
	}

	u.Password, _ = util.ToHash(u.Password)

	result := conn.Create(&u)
	if result.Error != nil {
		return Model{}, errors.New("error creating user")
	}

	return u, nil
}

func AuthUser(u Model) (Model, error) {
	user, err := FindUserByEmail(u.Email)

	if err != nil {
		return Model{}, errors.New("invalid credentials")
	}

	res := util.Compare(user.Password, u.Password)

	if !res {
		return Model{}, errors.New("invalid credentials")
	}

	return user, nil
}

func FindUserById(id uint) (Model, error) {
	return find("id = ?", id)
}

func FindUserByEmail(email string) (Model, error) {
	return find("email = ?", email)
}

func SearchByEmail(email string) ([]Model, error) {
	//%%  a literal percent sign; consumes no value
	return findAll("email LIKE ?", fmt.Sprintf("%%%s%%", email))
}

func findAll(query interface{}, args interface{}) ([]Model, error) {
	conn := db.Connection()

	var users []Model

	result := conn.Where(query, args).Find(&users)
	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}

func find(query interface{}, args interface{}) (Model, error) {
	conn := db.Connection()

	user := Model{}

	result := conn.Where(query, args).First(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func (u *Model) Save() *gorm.DB {
	return db.Connection().Save(&u)
}

func (u *Model) Update(column string, value string) *gorm.DB {
	return db.Connection().Model(&u).Update(column, value)
}

func TimeSeriesCount(limit int) (interface{}, error) {
	var users []struct {
		Date  int `json:"date"`
		Count int `json:"count"`
	}
	result := db.Connection().Raw(`
SELECT EXTRACT(YEAR_MONTH FROM created_at) as date, COUNT(id) as count
	FROM users
	GROUP BY date
	ORDER BY date DESC
	LIMIT ?`, limit).Scan(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func Summaries() (interface{}, error) {
	var users struct {
		Today int `json:"today"`
		Month int `json:"month"`
		Year  int `json:"year"`
		Total int `json:"total"`
	}
	now := time.Now().UTC()
	today := fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day())
	month := fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), 1)
	year := fmt.Sprintf("%d-%d-%d", now.Year(), 1, 1)

	result := db.Connection().Raw(`SELECT 
    	SUM(created_at > ?) as today,
    	SUM(created_at > ?) as month,
    	SUM(created_at > ?) as year,
       COUNT(created_at) as total
FROM users`, today, month, year).Scan(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}
