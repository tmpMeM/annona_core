package telebot_info

import (
	"fmt"

	"github.com/AnnonaOrg/annona_core/model"
)

func (c *TeleBotInfo) TableName() string {
	return "telebot_info"
}

func (r *TeleBotInfo) Create() error {
	return model.DB.Self.Create(&r).Error
}

// 删除 无条件，默认删除30天过期用户
func (r *TeleBotInfo) Delete() error {

	switch {

	case r.TelegramId != 0:
		return model.DB.Self.
			Unscoped().
			Where("telegram_id = ?", r.TelegramId).
			Delete(&TeleBotInfo{}).
			Error
	case len(r.TelegramUsername) > 0:
		return model.DB.Self.
			Unscoped().
			Where("telegram_username = ?", r.TelegramUsername).
			Delete(&TeleBotInfo{}).
			Error
	default:
		return fmt.Errorf("未知信息，删除出错: %v", r)
	}
}

func GetById(teleId int64) (*TeleBotInfo, error) {
	uu := &TeleBotInfo{}
	d := model.DB.Self.
		Select("telegram_bot_token", "telegram_id", "telegram_username").
		Where("telegram_id = ?", teleId).
		Last(&uu)
	return uu, d.Error
}
func GetByAccoundPlatformId(teleUsername string) (*TeleBotInfo, error) {
	uu := &TeleBotInfo{}
	d := model.DB.Self.
		Select("telegram_bot_token", "telegram_id", "telegram_username").
		Where("telegram_username = ?", teleUsername).
		Last(&uu)
	return uu, d.Error
}

func (r *TeleBotInfo) GetInfo() (*TeleBotInfo, error) {

	uu := &TeleBotInfo{}
	switch {
	case r.TelegramId != 0:
		err := model.DB.Self.
			Select("telegram_bot_token", "telegram_id", "telegram_username").
			Where("telegram_id = ?", r.TelegramId).
			First(&uu).
			Error
		return uu, err

	case len(r.TelegramUsername) > 0:
		err := model.DB.Self.
			Select("telegram_bot_token", "telegram_id", "telegram_username").
			Where("telegram_username = ?", r.TelegramUsername).
			First(&uu).
			Error
		return uu, err

	default:
		return nil, fmt.Errorf("未知信息，查询出错: %v", r)
	}
}
func (r *TeleBotInfo) Count() int64 {

	var count int64

	switch {
	case r.TelegramId != 0:

		model.DB.Self.Model(&TeleBotInfo{}).
			Where("telegram_id = ?", r.TelegramId).
			Count(&count)
		return count

	case len(r.TelegramUsername) > 0:

		model.DB.Self.Model(&TeleBotInfo{}).
			Where("telegram_username = ?", r.TelegramUsername).
			Count(&count)
		return count

	default:
		return count
	}
}

func (r *TeleBotInfo) Update() error {
	switch {

	case r.ID > 0:
		return model.DB.Self.Model(&r).
			Omit("telegram_id", "telegram_bot_token").
			Updates(&r).
			Error

	case r.TelegramId != 0 && len(r.TelegramBotToken) > 0 && len(r.TelegramUsername) > 0:
		err := model.DB.Self.Model(&r).
			Select("telegram_bot_token", "telegram_username").
			Omit("telegram_id").
			Where("telegram_id = ?", r.TelegramId).
			Updates(&r).
			Error
		return err
	case r.TelegramId != 0 && len(r.TelegramBotToken) > 0:
		err := model.DB.Self.Model(&r).
			Select("telegram_bot_token").
			Omit("telegram_id").
			Where("telegram_id = ?", r.TelegramId).
			Updates(&r).
			Error
		return err

	default:
		return fmt.Errorf("更新信息失败: %v", r)
	}
}

// 需鉴权使用
func (u *TeleBotInfo) GetList() ([]*TeleBotInfo, int64, error) {
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

	list := make([]*TeleBotInfo, 0)

	switch {
	case len(u.Filter) > 0:
		err := model.DB.Self.Model(&TeleBotInfo{}).
			Select("first_name", "last_name", "is_bot", "is_premium", "is_forum").
			Limit(size).Offset(offset).Find(&list).Error
		model.DB.Self.Model(&TeleBotInfo{}).
			Select("first_name", "last_name", "is_bot", "is_premium", "is_forum").
			Count(&count)
		return list, count, err

	default:
		return nil, count, fmt.Errorf("filter is NULL")
	}
}
