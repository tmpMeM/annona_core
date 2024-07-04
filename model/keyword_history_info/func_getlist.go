package keyword_history_info

import (
	"fmt"
	"strings"

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
	keyworld := strings.TrimSpace(u.KeyWorld)
	switch {
	case len(keyworld) > 0:
		keyworldLike := "%" + keyworld + "%"
		// beforeHour := time.Now().Add(-24 * 7 * time.Hour).Format("2006-01-02 15:04:05")
		// Where("updated_at > ?", beforeHour).
		err := model.DB.Self.Model(&KeyworldHistoryInfo{}).
			Select("sender_username, count(*) as total").
			Where("key_world LIKE ?", keyworldLike).
			Or("message_content_text LIKE ?", keyworldLike).
			Group("sender_username").
			Limit(size).Offset(offset).
			// Order("id DESC").
			Find(&list).
			Error
		model.DB.Self.Model(&KeyworldHistoryInfo{}).
			Select("sender_username").
			Where("key_world LIKE ?", keyworldLike).
			Or("message_content_text LIKE ?", keyworldLike).
			Group("sender_username").
			Count(&count)
		return list, count, err

	case u.SenderId != 0:
		err := model.DB.Self.Model(&KeyworldHistoryInfo{}).
			Where("sender_id = ?", u.SenderId).
			Limit(size).Offset(offset).
			Order("id DESC").
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
