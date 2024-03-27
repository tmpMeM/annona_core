package blockword_info

import (
	"fmt"

	"github.com/AnnonaOrg/annona_core/model"
	"github.com/clin003/util"
)

func (c *BlockworldInfo) TableName() string {
	return "blockword_info"
}

func (r *BlockworldInfo) Create() error {
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
func (r *BlockworldInfo) Delete() error {

	switch {

	case len(r.KeyWorld) > 0 && r.OwnerChatId != 0:
		return model.DB.Self.
			Where("key_world = ? AND owner_chat_id = ?", r.KeyWorld, r.OwnerChatId).
			Delete(&BlockworldInfo{}).
			Error

	case len(r.OwnerInfoHash) > 0:
		return model.DB.Self.
			Where("owner_info_hash = ?", r.OwnerInfoHash).
			Delete(&BlockworldInfo{}).
			Error

	default:
		return fmt.Errorf("未知删除条件")
	}
}

func GetById(id uint64) (*BlockworldInfo, error) {
	uu := &BlockworldInfo{}
	d := model.DB.Self.Where("id = ?", id).First(&uu)
	return uu, d.Error
}

func (r *BlockworldInfo) Get() (*BlockworldInfo, error) {
	uu := &BlockworldInfo{}

	switch {
	case len(r.InfoHash) > 0:
		err := model.DB.Self.Model(&BlockworldInfo{}).
			Where("info_hash = ?", r.InfoHash).
			First(&uu).
			Error
		return uu, err

	case len(r.OwnerInfoHash) > 0 && len(r.KeyWorld) > 0:
		err := model.DB.Self.Model(&BlockworldInfo{}).
			Where("owner_info_hash = ? AND key_world = ?", r.OwnerInfoHash, r.KeyWorld).
			First(&uu).
			Error
		return uu, err

	case len(r.OwnerPlatform) > 0 && r.OwnerChatId != 0 && len(r.KeyWorld) > 0:
		err := model.DB.Self.Model(&BlockworldInfo{}).
			Where("owner_chat_id = ? AND owner_platform = ? AND key_world = ?", r.OwnerChatId, r.OwnerPlatform, r.KeyWorld).
			First(&uu).
			Error
		return uu, err
	case r.OwnerChatId != 0 && len(r.KeyWorld) > 0:
		err := model.DB.Self.Model(&BlockworldInfo{}).
			Where("owner_chat_id = ? AND key_world = ?", r.OwnerChatId, r.KeyWorld).
			First(&uu).
			Error
		return uu, err

	default:
		err := model.DB.Self.Model(&BlockworldInfo{}).
			Where("id = ?", r.ID).
			First(&uu).
			Error
		return uu, err
	}
}
