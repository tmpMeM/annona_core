package keyword_info

import (
	"fmt"

	"github.com/AnnonaOrg/annona_core/model"
	"github.com/clin003/util"
)

func (c *KeyworldInfo) TableName() string {
	return "keyworld_info"
}

func (r *KeyworldInfo) Create() error {
	if len(r.KeyWorld) > 0 && len(r.OwnerInfoHash) > 0 {
		r.InfoHash = util.EncryptMd5(
			fmt.Sprintf("%s_%s_%d",
				r.KeyWorld, r.OwnerInfoHash, r.SearchChatId,
			))
	} else {
		err := fmt.Errorf("InfoHash:%s or OwnerInfoHash:%s is Null", r.InfoHash, r.OwnerInfoHash)
		return err
	}
	return model.DB.Self.Create(&r).Error
}

// 删除
func (r *KeyworldInfo) Delete() error {

	switch {

	case len(r.KeyWorld) > 0 && r.OwnerChatId != 0:
		return model.DB.Self.
			Where("key_world = ? AND owner_chat_id = ?", r.KeyWorld, r.OwnerChatId).
			Delete(&KeyworldInfo{}).
			Error

	case len(r.OwnerInfoHash) > 0:
		return model.DB.Self.
			Where("owner_info_hash = ?", r.OwnerInfoHash).
			Delete(&KeyworldInfo{}).
			Error

	default:
		return fmt.Errorf("未知删除条件")
	}
}

func GetById(id uint64) (*KeyworldInfo, error) {
	uu := &KeyworldInfo{}
	d := model.DB.Self.Where("id = ?", id).First(&uu)
	return uu, d.Error
}

func (r *KeyworldInfo) Get() (*KeyworldInfo, error) {
	uu := &KeyworldInfo{}

	switch {
	case len(r.InfoHash) > 0:
		err := model.DB.Self.Model(&KeyworldInfo{}).
			Where("info_hash = ?", r.InfoHash).
			First(&uu).
			Error
		return uu, err

	case len(r.OwnerInfoHash) > 0 && len(r.KeyWorld) > 0:
		err := model.DB.Self.Model(&KeyworldInfo{}).
			Where("owner_info_hash = ? AND key_world = ?", r.OwnerInfoHash, r.KeyWorld).
			First(&uu).
			Error
		return uu, err

	case len(r.OwnerPlatform) > 0 && r.OwnerChatId != 0 && len(r.KeyWorld) > 0:
		err := model.DB.Self.Model(&KeyworldInfo{}).
			Where("owner_chat_id = ? AND owner_platform = ? AND key_world = ?", r.OwnerChatId, r.OwnerPlatform, r.KeyWorld).
			First(&uu).
			Error
		return uu, err
	case r.OwnerChatId != 0 && len(r.KeyWorld) > 0:
		err := model.DB.Self.Model(&KeyworldInfo{}).
			Where("owner_chat_id = ? AND key_world = ?", r.OwnerChatId, r.KeyWorld).
			First(&uu).
			Error
		return uu, err

	default:
		err := model.DB.Self.Model(&KeyworldInfo{}).
			Where("id = ?", r.ID).
			First(&uu).
			Error
		return uu, err
	}
}
