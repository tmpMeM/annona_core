package keyword_history_info

import (
	"github.com/AnnonaOrg/annona_core/model"
)

func GetAll() ([]*KeyworldHistoryInfo, int64, error) {
	var err error
	list := make([]*KeyworldHistoryInfo, 0)
	rows, err := model.DB.Self.Model(&KeyworldHistoryInfo{}).Rows()
	defer rows.Close()
	for rows.Next() {
		var item KeyworldHistoryInfo
		if err1 := model.DB.Self.ScanRows(rows, &item); err1 != nil {
			err = err1
			continue
		} else {
			list = append(list, &item)
		}
	}
	if len(list) > 0 {
		return list, int64(len(list)), nil
	}
	return list, 0, err
}

// func (u *KeyworldHistoryInfo) GetListBySenderID() ([]*KeyworldHistoryInfo, int64, error) {

// 	list := make([]*KeyworldHistoryInfo, 0)

// 	switch {

// 	case u.SenderId != 0:
// 		err := model.DB.Self.Model(&KeyworldHistoryInfo{}).
// 			Select("sender_id", "count(tag) as total", "coalesce(max(stat), 0) as max_stat").
// 			Where("tag IN ?", filterList).
// 			//Or("tag LIKE ?", likeFilter).
// 			Group("chat_id").
// 			Where("sender_id = ?", u.SenderId).
// 			Limit(size).Offset(offset).
// 			Find(&list).
// 			Error
// 		model.DB.Self.Model(&KeyworldHistoryInfo{}).
// 			Where("sender_id = ?", u.SenderId).
// 			Count(&count)
// 		return list, count, err

// 	default:
// 		return nil, 0, fmt.Errorf("未找到符合条件的列表")
// 	}
// }
