package user_handler

import (
	"strings"

	"github.com/AnnonaOrg/annona_core/internal/load_db2redis"

	log "github.com/sirupsen/logrus"

	"github.com/AnnonaOrg/annona_core/handler"
	"github.com/AnnonaOrg/annona_core/model/user_info"
	"github.com/AnnonaOrg/osenv"
	"github.com/AnnonaOrg/pkg/errno"

	"github.com/gin-gonic/gin"
)

// 创建用户
func Create(c *gin.Context) {
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

	var u user_info.UserInfo
	if err := c.ShouldBind(&u); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		log.Errorf("ShouldBind: %+v", err)
		return
	}
	log.Debugf("ShouldBind: %+v", u)

	var uc user_info.UserInfo
	if len(u.AccoundPlatformId) > 0 && len(u.AccoundPlatform) > 0 {
		uc.AccoundPlatform = u.AccoundPlatform
		uc.AccoundPlatformId = u.AccoundPlatformId
		uc.TelegramChatId = u.TelegramChatId
	} else {
		handler.SendResponse(c, errno.ErrValidation, nil)
		return
	}

	// 检查是否已经注册
	if uc.Count() > 0 {
		if err := u.Update(); err != nil {
			handler.SendResponse(c, errno.ErrDatabase, nil)
			log.Errorf("Update: %v", err)
			return
		} else {
			handler.SendResponse(c, nil, nil)
		}
		return
	}

	inviterUserHash := ""
	inviterCode := u.Inviter // strings.TrimPrefix(u.Inviter, "vo_")
	inviterUser := &user_info.UserInfo{}
	var err error
	if len(inviterCode) > 0 {
		inviterUser, err = user_info.GetByInviterCode(inviterCode, u.AccoundPlatform)
		if err != nil {
		} else {
			inviterUserHash = inviterUser.InfoHash
		}
	}

	u.Inviter = inviterUserHash
	u.InviterCode = u.AccoundPlatformId
	if err := u.Create(); err != nil {
		handler.SendResponse(c, errno.ErrDatabase, nil)
		log.Errorf("Create: %v", err)
		return
	} else {
		handler.SendResponse(c, nil, u)
		go func() {
			// 更新用户信息
			load_db2redis.LoadUser(&u, true)
			// 更新邀请者信息
			if inviterUser != nil && len(inviterUserHash) > 0 {
				if len(inviterUser.InviterCode) == 0 {
					inviterUser.InviterCode = inviterUser.AccoundPlatformId
				}
				inviterUser.Exp = addUserExp(inviterUser.Exp, 3)
				if err := inviterUser.Update(); err != nil {
					log.Errorf("Update inviterUser(%+v): %v", inviterUser, err)
				} else {
					load_db2redis.LoadUser(inviterUser, true)
					// 发用邀请成功通知
					SendInviterNoticeMessage(inviterUser, &u)
				}
			}
		}()
		// // 更新用户信息
		// go load_db2redis.LoadUser(&u, true)
		// // 更新邀请者信息
		// if inviterUser != nil && len(inviterUserHash) > 0 {
		// 	if len(inviterUser.InviterCode) == 0 {
		// 		inviterUser.InviterCode = inviterUser.AccoundPlatformId
		// 	}
		// 	inviterUser.Exp = addUserExp(inviterUser.Exp, 3)
		// 	if err := inviterUser.Update(); err != nil {
		// 		log.Errorf("Update inviterUser(%+v): %v", inviterUser, err)
		// 	} else {
		// 		go load_db2redis.LoadUser(inviterUser, true)
		// 	}
		// }
		return
	}
}
