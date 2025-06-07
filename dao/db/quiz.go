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

func AddTrueOrFalseQuestion(trueOrFalseQuestion *model.TrueOrFalseQuestion) (error, uint) {
	err := Mysql.Create(trueOrFalseQuestion).Error
	return err, trueOrFalseQuestion.ID
}

func AddAnsweredTrueOrFalseQuestion(answeredTrueOrFalseQuestion *model.AnsweredTrueOrFalseQuestion) error {
	err := Mysql.Create(answeredTrueOrFalseQuestion).Error
	return err
}
