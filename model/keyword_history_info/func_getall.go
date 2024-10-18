package keyword_history_info

import (
	"github.com/AnnonaOrg/annona_core/model"
)

func GetAll() ([]KeyworldHistoryInfo, int64, error) {

	list := make([]KeyworldHistoryInfo, 0)

	if err := model.DB.Self.Model(&KeyworldHistoryInfo{}).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, int64(len(list)), nil
}
