package card_handler

import (
	"fmt"
	"strings"

	"github.com/AnnonaOrg/annona_core/handler"
	"github.com/AnnonaOrg/annona_core/model/card_info"
	"github.com/AnnonaOrg/osenv"
	"github.com/AnnonaOrg/pkg/errno"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

	var u card_info.CardInfo
	if err := c.ShouldBind(&u); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		log.Errorf("ShouldBind: %v", err)
		return
	}
	log.Debugf("ShouldBind: %+v", u)

	log.Debugf(
		"card_info list u.ById(%s) u.NDay(%s)",
		u.ById, u.NDay,
	)
	// 核验请求id
	if !osenv.IsBotManagerIDStr(u.ById) {
		handler.SendResponse(c, errno.ErrBadRequest, nil)
		return
	}

	if list, count, err := u.GetList(); err != nil {
		handler.SendResponse(c,
			fmt.Errorf("获取列表(%v)失败: %v %v", u, err, errno.ErrDatabase),
			nil,
		)
		log.Errorf("GetList: %v", err)
		return
	} else {
		handler.SendResponse(c, nil, handler.ListResponse{Items: list, Total: count})
		return
	}
}
