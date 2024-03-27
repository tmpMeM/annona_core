package user_handler

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/AnnonaOrg/annona_core/handler"
	model "github.com/AnnonaOrg/annona_core/model/user_info"
	"github.com/AnnonaOrg/osenv"
	"github.com/AnnonaOrg/pkg/errno"
	"github.com/gin-gonic/gin"
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

	// 解析用户信息
	var u model.UserInfo
	if err := c.ShouldBind(&u); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		log.Errorf("ShouldBind: %v", err)
		return
	}
	log.Debugf("ShouldBind: %+v", u)

	// 查询用户信息
	if userInfo, err := u.GetInfo(); err != nil {
		handler.SendResponse(c,
			fmt.Errorf("用户信息查询(%v)失败: %v %v", u.AccoundPlatformId, err, errno.ErrDatabase),
			nil,
		)
		log.Errorf("GetInfo( %+v ): %v", u, err)
		return
	} else {
		handler.SendResponse(c, nil, userInfo)
		log.Debugf("GetByAccoundPlatformId(%s): %+v", u.AccoundPlatformId, userInfo)
		return
	}
}
