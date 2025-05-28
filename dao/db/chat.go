package db

import "LiveLive/model"

func AddChatMessageRecord(chatMsgRecord *model.ChatMsgRecord) error {
	err := Mysql.Create(chatMsgRecord).Error
	return err
}
