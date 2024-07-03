package keyword_history_info

import (
	"github.com/AnnonaOrg/annona_core/model"
)

func (c *KeyworldHistoryInfo) TableName() string {
	return "keyworld_history_info"
}

func (r *KeyworldHistoryInfo) Create() error {
	return model.DB.Self.Create(&r).Error
}
