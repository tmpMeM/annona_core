package blockformsenderid_info

import (
	"fmt"

	"github.com/AnnonaOrg/annona_core/model"
	"github.com/clin003/util"
)

func (c *BlockformsenderidInfo) TableName() string {
	return "blockformsenderid_info"
}

func (r *BlockformsenderidInfo) Create() error {
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
func (r *BlockformsenderidInfo) Delete() error {

	switch {

	case len(r.KeyWorld) > 0 && r.OwnerChatId != 0:
		return model.DB.Self.
			Where("key_world = ? AND owner_chat_id = ?", r.KeyWorld, r.OwnerChatId).
			Delete(&BlockformsenderidInfo{}).
			Error

	case len(r.OwnerInfoHash) > 0:
		return model.DB.Self.
			Where("owner_info_hash = ?", r.OwnerInfoHash).
			Delete(&BlockformsenderidInfo{}).
			Error

	default:
		return fmt.Errorf("未知删除条件")
	}
}

func GetById(id uint64) (*BlockformsenderidInfo, error) {
	uu := &BlockformsenderidInfo{}
	d := model.DB.Self.Where("id = ?", id).First(&uu)
	return uu, d.Error
}

func (r *BlockformsenderidInfo) Get() (*BlockformsenderidInfo, error) {
	uu := &BlockformsenderidInfo{}

	switch {
	case len(r.InfoHash) > 0:
		err := model.DB.Self.Model(&BlockformsenderidInfo{}).
			Where("info_hash = ?", r.InfoHash).
			First(&uu).
			Error
		return uu, err

	case len(r.OwnerInfoHash) > 0 && len(r.KeyWorld) > 0:
		err := model.DB.Self.Model(&BlockformsenderidInfo{}).
			Where("owner_info_hash = ? AND key_world = ?", r.OwnerInfoHash, r.KeyWorld).
			First(&uu).
			Error
		return uu, err

	case len(r.OwnerPlatform) > 0 && r.OwnerChatId != 0 && len(r.KeyWorld) > 0:
		err := model.DB.Self.Model(&BlockformsenderidInfo{}).
			Where("owner_chat_id = ? AND owner_platform = ? AND key_world = ?", r.OwnerChatId, r.OwnerPlatform, r.KeyWorld).
			First(&uu).
			Error
		return uu, err
	case r.OwnerChatId != 0 && len(r.KeyWorld) > 0:
		err := model.DB.Self.Model(&BlockformsenderidInfo{}).
			Where("owner_chat_id = ? AND key_world = ?", r.OwnerChatId, r.KeyWorld).
			First(&uu).
			Error
		return uu, err

	default:
		err := model.DB.Self.Model(&BlockformsenderidInfo{}).
			Where("id = ?", r.ID).
			First(&uu).
			Error
		return uu, err
	}
}
