package blockformchatid_handler

import (
	"fmt"
	"strings"

	"github.com/AnnonaOrg/annona_core/internal/load_db2redis"

	"github.com/AnnonaOrg/annona_core/model/user_info"

	"github.com/AnnonaOrg/osenv"

	"github.com/AnnonaOrg/annona_core/handler"
	model "github.com/AnnonaOrg/annona_core/model/blockformchatid_info"
	"github.com/AnnonaOrg/pkg/errno"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

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

	var u model.BlockformchatidInfo
	if err := c.ShouldBind(&u); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		log.Errorf("ShouldBind: %v", err)
		return
	}
	log.Debugf("ShouldBind: %+v", u)

	var userInfo user_info.UserInfo
	userInfo.TelegramChatId = u.OwnerChatId
	userInfo.AccoundPlatform = u.OwnerPlatform
	userInfo.AccoundPlatformId = fmt.Sprintf("%d", u.OwnerChatId)
	if userInfoTmp, err := userInfo.GetInfo(); err != nil {
		handler.SendResponse(c,
			fmt.Errorf("用户信息获取(%v)失败: %v %v", userInfo.AccoundPlatformId, err, errno.ErrDatabase),
			nil,
		)
		log.Errorf("用户信息获取失败: %v", err)
		return
	} else {
		u.OwnerInfoHash = userInfoTmp.InfoHash
	}

	log.Debugf(
		"info Add u.OwnerInfoHash(%s) u.OwnerChatId(%d) u.SearchChatId(%d) u.KeyWorld(%s)",
		u.OwnerInfoHash, u.OwnerChatId, u.SearchChatId, u.KeyWorld,
	)

	if err := u.Create(); err != nil {
		handler.SendResponse(c,
			fmt.Errorf("信息添加(%v)失败: %v %v", u.KeyWorld, err, errno.ErrDatabase),
			nil,
		)
		log.Errorf("信息添加失败: %v", err)
		return
	}
	handler.SendResponse(c, nil, u)
	go load_db2redis.LoadBlockformchatid(u.OwnerInfoHash, u.KeyWorld)
	return
}
