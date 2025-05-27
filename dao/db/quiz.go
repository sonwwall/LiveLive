package db

import "LiveLive/model"

func AddChoiceQuestion(choiceQuestion *model.ChoiceQuestion) (error, uint) {
	err := Mysql.Create(choiceQuestion).Error
	return err, choiceQuestion.ID
}

func AddAnsweredChoiceQuestion(answeredChoiceQuestion *model.AnsweredChoiceQuestion) error {
	err := Mysql.Create(answeredChoiceQuestion).Error
	return err
}
