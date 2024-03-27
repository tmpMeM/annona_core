package user_handler

import (
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/AnnonaOrg/annona_core/handler"
	"github.com/AnnonaOrg/annona_core/model/user_info"
	"github.com/AnnonaOrg/osenv"
	"github.com/AnnonaOrg/pkg/errno"

	"github.com/gin-gonic/gin"
)

// 获取信息列表
func List(c *gin.Context) {
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

	// log.Info("List function called.")
	var u user_info.UserInfo
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

	if list, count, err := u.GetList(); err != nil {
		handler.SendResponse(c, errno.ErrDatabase, nil)
		log.Errorf("GetList: %v", err)
		return
	} else {
		handler.SendResponse(c, nil, handler.ListResponse{Items: list, Total: count})
		return
	}
}
