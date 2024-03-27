package user_handler

import (
	"github.com/AnnonaOrg/annona_core/internal/kvstore"

	log "github.com/sirupsen/logrus"
)

type cardNull struct {
	IsLock bool
}

// 如果kvStore 服务故障 返回 true ，避免重复兑换。 可锁定状态(没找到存储对象)返回 false ， 锁定状态(找到存储对象) 返回 true
func checkCard(cardUUID string) bool {
	keyCard := "check_" + cardUUID
	var n cardNull
	if isFound, err := kvstore.Client().Get(keyCard, n); err != nil {
		log.Errorf("卡状态检测(%s)异常 : %v", cardUUID, err)
		return true
	} else {
		return isFound
	}
}

func lockCard(cardUUID string) error {
	keyCard := "check_" + cardUUID
	return kvstore.Client().Set(keyCard, cardNull{IsLock: true})
}

func unlockCard(cardUUID string) error {
	keyCard := "check_" + cardUUID
	return kvstore.Client().Delete(keyCard)
}
