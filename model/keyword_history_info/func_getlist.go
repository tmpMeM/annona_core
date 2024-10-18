package keyword_history_info

import (
	"fmt"
	"strings"
	"time"

	"github.com/AnnonaOrg/annona_core/internal/utils"
	"github.com/AnnonaOrg/annona_core/model"
)

func (u *KeyworldHistoryInfo) GetList() ([]KeyworldHistoryInfo, int64, error) {
	list := make([]KeyworldHistoryInfo, 0)
	// var count int64
	var err error
	keyworld := strings.TrimSpace(u.KeyWorld)
	switch {
	case len(keyworld) > 0:
		list, _, err = u.GetListByKeyworld()
		if err != nil {
			return nil, 0, err
		}

	case u.SenderId != 0:
		list, _, err = u.GetListBySenderID()
		if err != nil {
			return nil, 0, err
		}

	default:
		return nil, 0, fmt.Errorf("未找到符合条件的列表")
	}

	return GetKeyworldHistoryInfoWithContentJoinToNote(list)
}

func (u *KeyworldHistoryInfo) GetListByKeyworld() ([]KeyworldHistoryInfo, int64, error) {
	var count int64

	list := make([]KeyworldHistoryInfo, 0)
	keyworld := strings.TrimSpace(u.KeyWorld)
	beforeHour := time.Now().Add(-24 * time.Hour).Format("2006-01-02 15:04:05")
	if len(keyworld) == 0 {
		return nil, 0, fmt.Errorf("keyworld is NULL")
	}

	keyworldLike := "%" + keyworld + "%"

	err := model.DB.Self.Model(&KeyworldHistoryInfo{}).
		Where("key_world LIKE ? AND updated_at > ?", keyworldLike, beforeHour).
		Order("id DESC").
		Find(&list).
		Error
	model.DB.Self.Model(&KeyworldHistoryInfo{}).
		Where("key_world LIKE ? AND updated_at > ?", keyworldLike, beforeHour).
		Count(&count)
	return list, count, err
}

// xxx
func (u *KeyworldHistoryInfo) GetListByKeyworldEx() ([]KeyworldHistoryInfo, int64, error) {
	list, _, err := u.GetListByKeyworld()
	if err != nil {
		return nil, 0, err
	}

	var listSender []int64
	listMap := make(map[int64]KeyworldHistoryInfo, 0)
	for _, v := range list {
		vc := v

		vc.Note = utils.GetStringRuneN(vc.MessageContentText, 30) + " " + vc.MessageLink + "\n"
		if len(vc.MessageLink) == 0 {
			vc.NoteHtml = utils.GetStringRuneN(vc.MessageContentText, 18)
		} else {
			vc.NoteHtml = utils.GetStringRuneN(vc.MessageContentText, 18) +
				" <a href=\"" + vc.MessageLink + "\">来源</a>" + "\n"
		}
		if _, isAdd := listMap[vc.SenderId]; isAdd {
			vcM := listMap[vc.SenderId]

			vcM.Note = vcM.Note + vc.Note
			vcM.NoteHtml = vcM.NoteHtml + vc.NoteHtml
			listMap[vc.SenderId] = vcM
		} else {
			listMap[vc.SenderId] = vc
			listSender = append(listSender, vc.SenderId)
		}
	}
	newList := make([]KeyworldHistoryInfo, 0)
	for _, v := range listSender {
		newList = append(newList, listMap[v])
	}

	return newList, int64(len(newList)), nil
}

func (u *KeyworldHistoryInfo) GetListBySenderID() ([]KeyworldHistoryInfo, int64, error) {
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

	list := make([]KeyworldHistoryInfo, 0)
	beforeHour := time.Now().Add(-24 * time.Hour).Format("2006-01-02 15:04:05")

	err := model.DB.Self.Model(&KeyworldHistoryInfo{}).
		Where("sender_id = ? AND updated_at > ?", u.SenderId, beforeHour).
		Limit(size).Offset(offset).
		Order("id DESC").
		Find(&list).
		Error
	model.DB.Self.Model(&KeyworldHistoryInfo{}).
		Where("sender_id = ? AND updated_at > ?", u.SenderId, beforeHour).
		Count(&count)
	return list, count, err
}

// xxx
func (u *KeyworldHistoryInfo) GetListBySenderIDEx() ([]KeyworldHistoryInfo, int64, error) {
	list, _, err := u.GetListBySenderID()
	if err != nil {
		return nil, 0, err
	}

	// retList := make([]*KeyworldHistoryInfo, 0)
	var listSender []int64
	listMap := make(map[int64]KeyworldHistoryInfo, 0)
	for _, v := range list {
		vc := v

		vc.Note = utils.GetStringRuneN(vc.MessageContentText, 30) + " " + vc.MessageLink + "\n"
		if len(vc.MessageLink) == 0 {
			vc.NoteHtml = utils.GetStringRuneN(vc.MessageContentText, 18)
		} else {
			vc.NoteHtml = utils.GetStringRuneN(vc.MessageContentText, 18) +
				" <a href=\"" + vc.MessageLink + "\">来源</a>" + "\n"
		}

		if _, isAdd := listMap[vc.SenderId]; isAdd {
			vcM := listMap[vc.SenderId]

			vcM.Note = vcM.Note + vc.Note
			vcM.NoteHtml = vcM.NoteHtml + vc.NoteHtml
			listMap[vc.SenderId] = vcM
		} else {
			listMap[vc.SenderId] = vc
			listSender = append(listSender, vc.SenderId)
		}
	}
	newList := make([]KeyworldHistoryInfo, 0)
	for _, v := range listSender {
		newList = append(newList, listMap[v])
	}

	return newList, int64(len(newList)), nil
}

// 将 KeyworldHistoryInfo 中的 ContentText 拼接到 Note中
func GetKeyworldHistoryInfoWithContentJoinToNote(list []KeyworldHistoryInfo) ([]KeyworldHistoryInfo, int64, error) {
	var listSender []int64
	listMap := make(map[int64]KeyworldHistoryInfo, 0)
	for _, v := range list {
		vc := v

		vc.Note = utils.GetStringRuneN(vc.MessageContentText, 30) + " " + vc.MessageLink + "\n"
		if len(vc.MessageLink) == 0 {
			vc.NoteHtml = utils.GetStringRuneN(vc.MessageContentText, 18)
		} else {
			vc.NoteHtml = utils.GetStringRuneN(vc.MessageContentText, 18) +
				" <a href=\"" + vc.MessageLink + "\">来源</a>" + "\n"
		}

		if _, isAdd := listMap[vc.SenderId]; isAdd {
			vcM := listMap[vc.SenderId]

			vcM.Note = vcM.Note + vc.Note
			vcM.NoteHtml = vcM.NoteHtml + vc.NoteHtml
			listMap[vc.SenderId] = vcM
		} else {
			listMap[vc.SenderId] = vc
			listSender = append(listSender, vc.SenderId)
		}
	}
	newList := make([]KeyworldHistoryInfo, 0)
	for _, v := range listSender {
		newList = append(newList, listMap[v])
	}

	return newList, int64(len(newList)), nil
}
