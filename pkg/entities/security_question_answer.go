package entities

type QuestionAnswer struct {
	ModelID

	Answer string `json:"-"`

	QuestionID uint `json:"-"`
	AccountID  uint `json:"-"`

	ModelTimeStamps
}

type QuestionAnswerWithQuestion struct {
	// TODO: Can we flatten the results from here? (like array flatten),
	// 	related: Since this only brings back id, is it necessary?
	QuestionAnswer /*`json:"-"`*/

	Question Question `json:"question"`
}

type QuestionAnswerWithAccountAndQuestion struct {
	QuestionAnswerWithQuestion

	Account Account `json:"account"`
}

func (QuestionAnswer) TableName() string {
	return "security_question_answers"
}

func (QuestionAnswerWithQuestion) TableName() string {
	return "security_question_answers"
}

func (QuestionAnswerWithAccountAndQuestion) TableName() string {
	return "security_question_answers"
}
