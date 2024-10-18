package user_info

import (
	"github.com/AnnonaOrg/annona_core/model"
)

func GetAll() ([]UserInfo, int64, error) {

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
