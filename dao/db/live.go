package db

import "LiveLive/model"

func AddLive(live *model.LiveSession) error {
	err := Mysql.Create(live).Error
	return err
}
