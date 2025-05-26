package db

import "LiveLive/model"

func AddChoiceQuestion(choiceQuestion *model.ChoiceQuestion) error {
	err := Mysql.Create(choiceQuestion).Error
	return err
}
