package user_info

import (
	"github.com/AnnonaOrg/annona_core/model"
)

func GetAll() ([]UserInfo, int64, error) {
	// var err error
	// list := make([]*UserInfo, 0)
	// rows, err := model.DB.Self.Model(&UserInfo{}).Rows()
	// defer rows.Close()
	// for rows.Next() {
	// 	var item UserInfo
	// 	if err1 := model.DB.Self.ScanRows(rows, &item); err1 != nil {
	// 		err = err1
	// 		continue
	// 	} else {
	// 		list = append(list, &item)
	// 	}
	// }
	var list []UserInfo
	err := model.DB.Self.Model(&UserInfo{}).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	if len(list) > 0 {
		return list, int64(len(list)), nil
	}
	return list, 0, err
}
