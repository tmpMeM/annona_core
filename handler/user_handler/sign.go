package user_handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/AnnonaOrg/annona_core/handler"
	"github.com/AnnonaOrg/annona_core/internal/load_db2redis"
	"github.com/AnnonaOrg/annona_core/model/user_info"
	"github.com/AnnonaOrg/osenv"
	"github.com/AnnonaOrg/pkg/errno"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// 签到 并更新 用户信息
// update update exist user account info.
func Sign(c *gin.Context) {
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

	//解析用户信息
	var u user_info.UserInfo
	if err := c.ShouldBind(&u); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		log.Errorf("ShouldBind: %v", err)
		return
	}
	log.Debugf("ShouldBind: %+v", u)

	//通过平台id查询用户信息
	if userInfoTmp, err := u.GetInfo(); err != nil {
		handler.SendResponse(c, errno.ErrNotFound, nil)
		log.Errorf("GetInfo( %+v ): %v", u, err)
		return
	} else {
		// 计算有效期 时间戳 userInfoTmp.Exp 时间戳
		timeNow := time.Now()
		timeDate := timeNow.Format(time.DateOnly)
		if strings.EqualFold(timeDate, userInfoTmp.LastSignDate) {
			handler.SendResponse(c,
				fmt.Errorf("今天已经签过了: %s", userInfoTmp.LastSignDate),
				nil,
			)
			return
		}
		timeNowUnix := timeNow.Unix()
		if userInfoTmp.Exp < timeNowUnix {
			userInfoTmp.Exp =
				time.Now().Add(
					// time.Duration(cardInfo.Exp) * 24 * time.Hour,
					3 * time.Hour,
				).Unix()
		} else {
			userInfoTmp.Exp =
				userInfoTmp.Exp + 3*60*60
		}

		u.ID = userInfoTmp.ID
		u.Exp = userInfoTmp.Exp
		u.InfoHash = userInfoTmp.InfoHash
		u.TelegramStartBotId = userInfoTmp.TelegramStartBotId
		u.TelegramUsername = userInfoTmp.TelegramUsername
		u.LastSignDate = timeDate
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
	handler.SendResponse(c, nil, nil)
	return
}
