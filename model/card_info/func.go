package card_info

import (
	"fmt"

	"github.com/AnnonaOrg/annona_core/model"
	"github.com/clin003/util"
	"github.com/google/uuid"
)

func (c *CardInfo) TableName() string {
	return "card_info"
}

// 生成单张卡
func (r *CardInfo) Create() error {
	if r.CardUUID == "" || r.CardHash == "" {
		r.CardUUID = uuid.New().String()
		r.CardHash = util.EncryptMd5(
			r.CardUUID +
				fmt.Sprintf("%d", r.Exp) + r.Note,
		)
	}
	r.Stat = 1
	return model.DB.Self.Create(&r).Error
}

// 批量生卡
func (r *CardInfo) CreateN(nNum int64) ([]*CardInfo, int64, error) {
	var n int
	if nNum > 50 {
		n = 50
	} else if nNum > 0 {
		n = int(nNum)
	} else {
		n = 1
	}

	list := make([]*CardInfo, 0)

	for i := 0; i < n; i++ {
		var cardInfo CardInfo
		cardInfo.Exp = r.Exp
		cardInfo.Note = r.Note
		cardInfo.Stat = 1
		cardInfo.CardUUID = uuid.New().String()
		cardInfo.CardHash = util.EncryptMd5(
			cardInfo.CardUUID +
				fmt.Sprintf("%d", cardInfo.Exp) + cardInfo.Note,
		)
		list = append(list, &cardInfo)
	}
	err := model.DB.Self.Create(&list).Error

	return list, int64(len(list)), err
}

// 删除
func (r *CardInfo) Delete() error {
	switch {
	case len(r.CardUUID) > 0:
		return model.DB.Self.
			Where("card_uuid = ?", r.CardUUID).
			Delete(&CardInfo{}).
			Error
	case r.ID > 0:
		return model.DB.Self.
			Where("id = ?", r.ID).
			Delete(&CardInfo{}).
			Error
	default:
		return fmt.Errorf("未找到符合条件的待删除对象: %v", r)
	}
}

func (r *CardInfo) Update() error {
	switch {
	case r.CardUUID != "":
		return model.DB.Self.Model(&CardInfo{}).
			Omit("card_uuid", "info_hash", "exp").
			Select("stat").
			Where("card_uuid = ?", r.CardUUID).
			Updates(&r).
			Error
	default:
		return fmt.Errorf("未找到符合条件的更新对象: %v", r)
	}
}

func GetCardInfoByUUID(cardUUID string) (*CardInfo, error) {
	uu := &CardInfo{}
	d := model.DB.Self.Model(&CardInfo{}).
		Where("card_uuid = ?", cardUUID).
		First(&uu)
	return uu, d.Error
}

// 通过 gid goods_id  获取相关信息
func (u *CardInfo) Get() (*CardInfo, error) {
	uu := &CardInfo{}
	switch {
	case len(u.CardUUID) > 0:
		d := model.DB.Self.Model(&CardInfo{}).
			Where("card_uuid = ?", u.CardUUID).
			First(&uu)
		return uu, d.Error
	default:
		return nil, fmt.Errorf("未知卡信息: %v", u)
	}
}

func (u *CardInfo) GetList() ([]*CardInfo, int64, error) {
	var count int64
	if u.Size <= 0 || u.Size > 100 {
		u.Size = 10
	}
	size := u.Size
	if u.Page <= 0 {
		u.Page = -1
	}
	offset := u.Page - 1
	if offset > 0 {
		offset = offset * u.Size
	}

	list := make([]*CardInfo, 0)

	switch {
	case u.Stat > 0:
		err := model.DB.Self.Model(&CardInfo{}).
			Where("stat = ?", u.Stat).
			Order("updated_at desc").
			Limit(size).Offset(offset).Find(&list).Error
		model.DB.Self.Model(&CardInfo{}).
			Where("stat = ?", u.Stat).
			Count(&count)
		return list, count, err

	case len(u.Filter) > 0:
		// filter := u.Filter
		likeFilter := fmt.Sprintf("%s%s%s", "%", u.Filter, "%")
		err := model.DB.Self.Model(&CardInfo{}).
			Where("note LIKE ?", likeFilter).
			Order("stat desc").Order("updated_at desc").
			Limit(size).Offset(offset).Find(&list).Error
		model.DB.Self.Model(&CardInfo{}).
			Where("note LIKE ?", likeFilter).
			Count(&count)
		return list, count, err

	default:
		err := model.DB.Self.Model(&CardInfo{}).
			Order("stat").
			Limit(size).Offset(offset).Find(&list).Error
		model.DB.Self.Model(&CardInfo{}).Count(&count)
		return list, count, err
	}
}
