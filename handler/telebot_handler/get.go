package telebot_handler

import (
	"fmt"
	"strings"

	"github.com/AnnonaOrg/osenv"

	"github.com/AnnonaOrg/annona_core/handler"
	model "github.com/AnnonaOrg/annona_core/model/telebot_info"
	"github.com/AnnonaOrg/pkg/errno"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Get an user by the user identifier
func Get(c *gin.Context) {
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

	var u model.TeleBotInfo
	if err := c.ShouldBind(&u); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		log.Errorf("ShouldBind: %v", err)
		return
	}
	log.Debugf("ShouldBind: %+v", u)
	// 核验请求id
	if !osenv.IsBotManagerIDStr(u.ById) {
		handler.SendResponse(c, errno.ErrBadRequest, nil)
		return
	}

	if uu, err := u.GetInfo(); err != nil {
		handler.SendResponse(c,
			fmt.Errorf("查询(%v)失败: %v %v", u, err, errno.ErrDatabase),
			nil,
		)
		log.Errorf("GetInfo: %v", err)
	} else {
		handler.SendResponse(c, nil, uu)
	}
	return
}
