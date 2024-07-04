package keyword_history_handler

import (
	"strings"

	"github.com/AnnonaOrg/annona_core/handler"
	model "github.com/AnnonaOrg/annona_core/model/keyword_history_info"
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

	var u model.KeyworldHistoryInfo
	if err := c.ShouldBind(&u); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		log.Errorf("ShouldBind: %v", err)
		return
	}
	log.Debugf("ShouldBind: %+v", u)

	if list, count, err := u.GetList(); err != nil {
		handler.SendResponse(c, errno.ErrDatabase, nil)
		log.Errorf("GetList: %v", err)
		return
	} else {
		handler.SendResponse(c,
			nil,
			handler.ListResponse{Items: list, Total: count},
		)
		// log.Debugf("GetList: %+v", list)
		return
	}
}
