package keyword_info

import (
	"fmt"

	"github.com/AnnonaOrg/annona_core/model"
)

func (u *KeyworldInfo) GetList() ([]*KeyworldInfo, int64, error) {
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

	list := make([]*KeyworldInfo, 0)

	switch {
	case len(u.ById) > 0:

		err := model.DB.Self.Model(&KeyworldInfo{}).
			Select("info_hash", "owner_info_hash", "owner_chat_id", "search_chat_id", "key_world").
			Limit(size).Offset(offset).
			Find(&list).
			Error
		model.DB.Self.Model(&KeyworldInfo{}).
			Count(&count)
		return list, count, err

	case len(u.OwnerInfoHash) > 0:
		err := model.DB.Self.Model(&KeyworldInfo{}).
			Select("id", "owner_info_hash", "owner_chat_id", "search_chat_id", "key_world").
			Where("owner_info_hash = ?", u.OwnerInfoHash).
			Limit(size).Offset(offset).
			Find(&list).
			Error
		model.DB.Self.Model(&KeyworldInfo{}).
			Where("owner_info_hash = ?", u.OwnerInfoHash).
			Count(&count)
		return list, count, err

	case u.OwnerChatId != 0:
		err := model.DB.Self.Model(&KeyworldInfo{}).
			Select("id", "owner_info_hash", "owner_chat_id", "search_chat_id", "key_world").
			Where("owner_chat_id = ?", u.OwnerChatId).
			Limit(size).Offset(offset).
			Find(&list).
			Error
		model.DB.Self.Model(&KeyworldInfo{}).
			Where("owner_chat_id = ?", u.OwnerChatId).
			Count(&count)
		return list, count, err

	default:
		return nil, 0, fmt.Errorf("未找到符合条件的列表")
	}
}
