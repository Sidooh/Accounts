package entities

type Question struct {
	ModelID

	Question string `json:"question" gorm:"unique"`
	Status   string `json:"status"`

	ModelTimeStamps
}

func (Question) TableName() string {
	return "security_questions"
}
