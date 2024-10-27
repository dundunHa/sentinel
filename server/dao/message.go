package dao

import (
	"sentinel/server/model"
	"sentinel/server/storage"
)

func GetLastProcessedID(appID int) (int, error) {
	var gotifyMessage model.GotifyMessage
	if err := storage.DB.Where("app_id = ?", appID).First(&gotifyMessage).Error; err != nil {
		return 0, err
	}
	return gotifyMessage.LastProcessedID, nil
}

func UpdateLastProcessedID(appID, id int) error {
	gotifyMessage := model.GotifyMessage{
		AppID:           appID,
		LastProcessedID: id,
	}

	return storage.DB.
		Where("app_id = ?", appID).
		Assign(gotifyMessage).
		FirstOrCreate(&gotifyMessage).Error
}
