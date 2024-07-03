package keyword_history_info

import (
	"fmt"

	"github.com/AnnonaOrg/annona_core/model"
)

func (u *KeyworldHistoryInfo) GetList() ([]*KeyworldHistoryInfo, int64, error) {
	var count int64
	if u.Size <= 0 || u.Size > 100 {
		u.Size = 100
	}
	size := u.Size
	if u.Page <= 0 {
		u.Page = -1
	}
	offset := u.Page - 1
	if offset > 0 {
		offset = offset * u.Size
	}

	list := make([]*KeyworldHistoryInfo, 0)

	switch {
	case len(u.KeyWorld) > 0:
		err := model.DB.Self.Model(&KeyworldHistoryInfo{}).
			Where("key_world = ?", u.KeyWorld).
			Limit(size).Offset(offset).
			Find(&list).
			Error
		model.DB.Self.Model(&KeyworldHistoryInfo{}).
			Where("key_world = ?", u.KeyWorld).
			Count(&count)
		return list, count, err

	case u.SenderId != 0:
		err := model.DB.Self.Model(&KeyworldHistoryInfo{}).
			Where("sender_id = ?", u.SenderId).
			Limit(size).Offset(offset).
			Find(&list).
			Error
		model.DB.Self.Model(&KeyworldHistoryInfo{}).
			Where("sender_id = ?", u.SenderId).
			Count(&count)
		return list, count, err

	default:
		return nil, 0, fmt.Errorf("未找到符合条件的列表")
	}
}
