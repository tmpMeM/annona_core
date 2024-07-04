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
