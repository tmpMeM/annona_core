package telebot_info

import (
	"github.com/AnnonaOrg/annona_core/model"
)

func GetAll() ([]*TeleBotInfo, int64, error) {
	var err error
	list := make([]*TeleBotInfo, 0)
	rows, err := model.DB.Self.Model(&TeleBotInfo{}).Rows()
	defer rows.Close()
	for rows.Next() {
		var item TeleBotInfo
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
