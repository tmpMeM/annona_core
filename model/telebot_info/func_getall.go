package telebot_info

import (
	"github.com/AnnonaOrg/annona_core/model"
)

func GetAll() ([]TeleBotInfo, int64, error) {
	list := make([]TeleBotInfo, 0)

	if err := model.DB.Self.Model(&TeleBotInfo{}).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, int64(len(list)), nil
}
