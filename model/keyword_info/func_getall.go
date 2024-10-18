package keyword_info

import (
	"fmt"

	"github.com/AnnonaOrg/annona_core/model"
)

func GetAll() ([]KeyworldInfo, int64, error) {

	list := make([]KeyworldInfo, 0)

	if err := model.DB.Self.Model(&KeyworldInfo{}).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, int64(len(list)), nil
}

func GetAllByOwnerInfoHash(ownerInfoHash string) ([]KeyworldInfo, int64, error) {
	var err error
	list := make([]KeyworldInfo, 0)
	rows, err := model.DB.Self.Model(&KeyworldInfo{}).
		Where("owner_info_hash = ?", ownerInfoHash).
		Rows()
	defer rows.Close()
	for rows.Next() {
		var item KeyworldInfo
		if err1 := model.DB.Self.ScanRows(rows, &item); err1 != nil {
			err = err1
			continue
		} else {
			list = append(list, item)
		}
	}
	if len(list) > 0 {
		return list, int64(len(list)), nil
	}
	return list, 0, err
}

func GetAllByOwnerInfoHashToString(ownerInfoHash string) ([]string, int64, error) {
	var err error
	list := make([]string, 0)

	listTmp, count, err := GetAllByOwnerInfoHash(ownerInfoHash)
	if err != nil {
		return nil, count, err
	} else if count == 0 {
		return nil, 0, fmt.Errorf("List is NULL")
	}

	for _, v := range listTmp {
		keywordStr := v.KeyWorld
		list = append(list, keywordStr)
	}
	return list, int64(len(list)), nil
}
