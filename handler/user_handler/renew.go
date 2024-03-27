package user_handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/AnnonaOrg/annona_core/internal/load_db2redis"

	log "github.com/sirupsen/logrus"

	"github.com/AnnonaOrg/annona_core/model/card_info"

	"github.com/AnnonaOrg/annona_core/handler"
	"github.com/AnnonaOrg/annona_core/model/user_info"
	"github.com/AnnonaOrg/osenv"
	"github.com/AnnonaOrg/pkg/errno"
	"github.com/gin-gonic/gin"
)

// 核销卡 并更新 用户信息
// update update exist user account info.
func Renew(c *gin.Context) {
	coreApiToken := osenv.GetCoreApiToken()
	if len(coreApiToken) > 0 {
		if !strings.EqualFold(coreApiToken, c.Request.Header.Get("Apiclient")) {
			handler.SendResponse(c, errno.ErrBadRequest, nil)
			return
		}
	} else {
		handler.SendResponse(c, errno.ErrBadRequest, nil)
		return
	}

	cardUUID := c.Param("uuid")
	if cardUUID == "" {
		handler.SendResponse(c,
			fmt.Errorf("卡信息获取(%s)失败: %v", cardUUID, errno.ErrValidation),
			nil,
		)
	}
	// 判断卡状态
	if checkCard(cardUUID) {
		handler.SendResponse(c,
			fmt.Errorf("当前卡(%s)状态异常，请检查是否有正在进行的其他操作，或联系管理员操作: %v", cardUUID, errno.ErrValidation),
			nil,
		)
		log.Debugf("checkCard(%s): 当前卡状态异常，请检查是否有正在进行的其他操作，或联系管理员操作", cardUUID)
		return
	}
	// 锁卡
	if err := lockCard(cardUUID); err != nil {
		handler.SendResponse(c,
			fmt.Errorf("锁卡(%s)失败，请检查是否有正在进行的其他操作: %v", cardUUID, errno.ErrValidation),
			nil,
		)
		log.Errorf("lockCard: %v", err)
		return
	}
	// 解卡
	defer unlockCard(cardUUID)

	//解析用户信息
	var u user_info.UserInfo
	if err := c.ShouldBind(&u); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		log.Errorf("ShouldBind: %v", err)
		return
	}
	log.Debugf("ShouldBind: %+v", u)

	// 查询卡信息
	cardInfo, err := card_info.GetCardInfoByUUID(cardUUID)
	if err != nil {
		handler.SendResponse(c,
			fmt.Errorf("卡信息获取(%s)失败%v %v", cardUUID, err, errno.ErrDatabase),
			nil,
		)
		log.Errorf("GetCardInfoByUUID: %v", err)
		return
	}
	if cardInfo.Stat > 2 {
		handler.SendResponse(c,
			fmt.Errorf("检查卡(%s)已核销,不可重复使用: %v", cardUUID, errno.ErrDatabase),
			nil,
		)
		log.Debugf("检查卡(%s)已核销,不可重复使用", cardUUID)
		return
	}

	//通过平台id查询用户信息
	if userInfoTmp, err := u.GetInfo(); err != nil {
		handler.SendResponse(c, errno.ErrNotFound, nil)
		log.Errorf("GetInfo( %+v ): %v", u, err)
		return
	} else {
		// 计算有效期 时间戳 userInfoTmp.Exp 时间戳
		timeNowUnix := time.Now().Unix()
		if userInfoTmp.Exp < timeNowUnix {
			userInfoTmp.Exp =
				time.Now().Add(
					time.Duration(cardInfo.Exp) * 24 * time.Hour,
				).Unix()
		} else {
			userInfoTmp.Exp =
				userInfoTmp.Exp + cardInfo.Exp*24*60*60 + 24*60*60
		}

		u.ID = userInfoTmp.ID
		u.Exp = userInfoTmp.Exp
		u.InfoHash = userInfoTmp.InfoHash
		u.TelegramStartBotId = userInfoTmp.TelegramStartBotId
		u.TelegramUsername = userInfoTmp.TelegramUsername
		u.LastCardUUID = cardUUID
	}

	//更新用户信息
	if err := u.Update(); err != nil {
		handler.SendResponse(c, errno.ErrDatabase, nil)
		log.Errorf("Update: %v", err)
		return
	}
	go load_db2redis.LoadUser(&u, true)
	log.Debugf("Update: %+v", u)
	// 用户信息更新完成
	// 标记卡状态
	cardInfo.Stat = 3
	if err := cardInfo.Update(); err != nil {
		handler.SendResponse(c,
			fmt.Errorf("卡状态更新(%v)失败: %v", cardInfo, err, errno.ErrDatabase),
			nil,
		)
		log.Errorf("Update: %v", err)
		return
	}
	handler.SendResponse(c, nil, nil)
	return
}
