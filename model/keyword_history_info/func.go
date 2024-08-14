package keyword_history_info

import (
	"time"

	"github.com/AnnonaOrg/annona_core/model"
)

func (c *KeyworldHistoryInfo) TableName() string {
	return "keyworld_history_info"
}

func (r *KeyworldHistoryInfo) Create() error {
	return model.DB.Self.Create(&r).Error
}

func DeleteKeyworldHistoryInfoBeforeOneDay() error {
	beforeHour := time.Now().Add(-24 * time.Hour).Format("2006-01-02 15:04:05")
	return model.DB.Self.Unscoped().
		Where("updated_at < ?", beforeHour).
		Delete(&KeyworldHistoryInfo{}).
		Error
}
