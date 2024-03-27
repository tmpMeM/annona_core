package user_handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AnnonaOrg/annona_core/handler"
	"github.com/AnnonaOrg/annona_core/internal/load_db2redis"
	"github.com/AnnonaOrg/annona_core/model/user_info"
	"github.com/AnnonaOrg/osenv"
	"github.com/AnnonaOrg/pkg/errno"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// 设置通知会话id
// update update exist user account info.
func UpdateNoticeChatId(c *gin.Context) {
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

	noticeChatId, err := strconv.ParseInt(c.Param("chatid"), 10, 64)
	if err != nil || noticeChatId == 0 {
		handler.SendResponse(c,
			fmt.Errorf("通知会话ID获取(%s)失败: %v", c.Param("chatid"), errno.ErrValidation),
			nil,
		)
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

	u.TelegramNoticeChatId = noticeChatId

	//更新用户信息
	if err := u.Update(); err != nil {
		handler.SendResponse(c, errno.ErrDatabase, nil)
		log.Errorf("Update: %v", err)
		return
	}

	log.Debugf("Update: %+v", u)
	// 用户信息更新完成
	handler.SendResponse(c, nil, nil)
	// go load_db2redis.LoadUser(&u, false)
	go func() {
		if userTmp, err := u.GetInfo(); err != nil {
		} else {
			load_db2redis.LoadUser(userTmp, false)
		}
	}()
	return
}
