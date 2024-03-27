package user_info

import (
	"fmt"
	"time"

	"github.com/AnnonaOrg/annona_core/model"
	"github.com/clin003/util"
)

func (c *UserInfo) TableName() string {
	return "user_info"
}

func (r *UserInfo) Create() error {
	if len(r.AccoundPlatform) > 0 && len(r.AccoundPlatformId) > 0 {
		r.InfoHash = util.EncryptMd5(
			fmt.Sprintf("%s_%s",
				r.AccoundPlatform, r.AccoundPlatformId,
			))
		if r.Exp < 10 {
			if len(r.Inviter) > 0 {
				r.Exp = time.Now().Add(8 * time.Hour).Unix()
			} else {
				r.Exp = time.Now().Add(10 * time.Hour).Unix()
			}
		}
	} else {
		err := fmt.Errorf(
			"AccoundPlatform:%s , AccoundPlatformId:%s 未知",
			r.AccoundPlatform, r.AccoundPlatformId,
		)
		return err
	}
	return model.DB.Self.Create(&r).Error
}

// 删除 无条件，默认删除30天过期用户
func (r *UserInfo) Delete() error {

	if r.InfoHash == "" {
		if len(r.AccoundPlatform) > 0 && len(r.AccoundPlatformId) > 0 {
			r.InfoHash = util.EncryptMd5(
				fmt.Sprintf("%s_%s",
					r.AccoundPlatform, r.AccoundPlatformId,
				))
		}
	}

	switch {
	case len(r.InfoHash) > 0:
		return model.DB.Self.
			Unscoped().
			Where("info_hash = ?", r.InfoHash).
			Delete(&UserInfo{}).
			Error

	case r.TelegramChatId > 0:
		return model.DB.Self.
			Unscoped().
			Where("telegram_chat_id = ?", r.TelegramChatId).
			Delete(&UserInfo{}).
			Error
	default:
		expTime := time.Now().Add(720 * time.Hour).Unix()
		return model.DB.Self.
			Unscoped().
			Where("exp < ?", expTime).
			Delete(&r).
			Error
		return nil
	}
}

func GetById(id uint64) (*UserInfo, error) {
	uu := &UserInfo{}
	d := model.DB.Self.Where("id = ?", id).First(&uu)
	return uu, d.Error
}
func GetByInviterCode(inviterCode string, inviterPlatform string) (*UserInfo, error) {
	uu := &UserInfo{}
	d := model.DB.Self.
		Where("inviter_code = ?", inviterCode).
		Or("accound_platform = ? AND accound_platform_id =?", inviterPlatform, inviterCode).
		First(&uu)
	return uu, d.Error
}

func (r *UserInfo) GetInfo() (*UserInfo, error) {
	if r.InfoHash == "" {
		if len(r.AccoundPlatform) > 0 && len(r.AccoundPlatformId) > 0 {
			r.InfoHash = util.EncryptMd5(
				fmt.Sprintf("%s_%s",
					r.AccoundPlatform, r.AccoundPlatformId,
				))
		}
	}

	uu := &UserInfo{}
	switch {
	case len(r.InfoHash) > 0:
		err := model.DB.Self.Model(&UserInfo{}).
			// Select("info_hash", "accound_platform", "accound_platform_id", "telegram_chat_id", "telegram_start_bot_id").
			Where("info_hash = ?", r.InfoHash).
			First(&uu).
			Error
		return uu, err

	case len(r.AccoundPlatformId) > 0:
		err := model.DB.Self.Model(&UserInfo{}).
			Select("accound_platform_id", "telegram_chat_id").
			Where("accound_platform_id = ?", r.AccoundPlatformId).
			First(&uu).
			Error
		return uu, err

	default:
		err := model.DB.Self.Model(&UserInfo{}).
			Select("accound_platform_id", "telegram_chat_id").
			Where("id = ?", r.ID).
			First(&uu).
			Error
		return uu, err
	}
}
func (r *UserInfo) Count() int64 {
	if r.InfoHash == "" {
		if len(r.AccoundPlatform) > 0 && len(r.AccoundPlatformId) > 0 {
			r.InfoHash = util.EncryptMd5(
				fmt.Sprintf("%s_%s",
					r.AccoundPlatform, r.AccoundPlatformId,
				))
		}
	}

	var count int64

	switch {
	case len(r.InfoHash) > 0:

		model.DB.Self.Model(&UserInfo{}).
			Where("info_hash = ?", r.InfoHash).
			Count(&count)
		return count

	case len(r.AccoundPlatformId) > 0:

		model.DB.Self.Model(&UserInfo{}).
			Where("accound_platform_id = ?", r.AccoundPlatformId).
			Count(&count)
		return count

	default:
		return count
	}
}

func (r *UserInfo) Update() error {
	if r.InfoHash == "" {
		if len(r.AccoundPlatform) > 0 && len(r.AccoundPlatformId) > 0 {
			r.InfoHash = util.EncryptMd5(
				fmt.Sprintf("%s_%s",
					r.AccoundPlatform, r.AccoundPlatformId,
				))
		}
	}

	switch {

	case r.ID > 0:
		return model.DB.Self.Model(&UserInfo{}).
			Omit("id", "info_hash", "accound_platform", "accound_platform_id", "telegram_chat_id", "inviter").
			Where("id = ?", r.ID).
			Updates(&r).
			Error

	case len(r.InfoHash) > 0 && r.TelegramNoticeChatId != 0:
		return model.DB.Self.Model(&UserInfo{}).
			Select("telegram_notice_chat_id").
			Omit("id", "info_hash", "accound_platform", "accound_platform_id", "telegram_chat_id", "inviter").
			Where("info_hash = ?", r.InfoHash).
			Updates(&r).
			Error

	case len(r.InfoHash) > 0:
		return model.DB.Self.Model(&UserInfo{}).
			// Select("telegram_start_bot_id").
			Omit("id", "info_hash", "accound_platform", "accound_platform_id", "telegram_chat_id", "inviter").
			Where("info_hash = ?", r.InfoHash).
			Updates(&r).
			Error

	case r.TelegramChatId != 0:
		return model.DB.Self.Model(&UserInfo{}).
			Select("telegram_start_bot_id").
			Omit("id", "info_hash", "accound_platform", "accound_platform_id", "telegram_chat_id", "inviter").
			Where("telegram_chat_id = ?", r.TelegramChatId).
			Updates(&r).
			Error

	default:
		return fmt.Errorf("更新用户信息失败: %v", r)
	}
}

// 需鉴权使用
func (u *UserInfo) GetList() ([]*UserInfo, int64, error) {
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

	list := make([]*UserInfo, 0)

	switch {
	case len(u.Filter) > 0:
		// filter := u.Filter
		// likeFilter := fmt.Sprintf("%s%s%s", "%", u.Filter, "%")
		expTime := time.Now().Unix()
		err := model.DB.Self.Model(&UserInfo{}).
			Select("info_hash", "accound_platform", "accound_platform_id", "telegram_chat_id").
			Where("exp < ?", expTime).
			Limit(size).Offset(offset).Find(&list).Error
		model.DB.Self.Model(&UserInfo{}).
			Select("info_hash", "accound_platform", "accound_platform_id", "telegram_chat_id").
			Where("exp < ?", expTime).
			Count(&count)
		return list, count, err

	default:
		return nil, count, fmt.Errorf("filter is NULL")
	}
}
