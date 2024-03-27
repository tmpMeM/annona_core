package redis_handler

import (
	"strings"

	"github.com/AnnonaOrg/annona_core/handler"
	"github.com/AnnonaOrg/annona_core/internal/redis_user"
	"github.com/AnnonaOrg/osenv"
	"github.com/AnnonaOrg/pkg/errno"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// 根据用户infohash 获取用户信息
func GetUserInfoByUserInfoHash(c *gin.Context) {
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
	var u RedisKey
	if err := c.ShouldBind(&u); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		log.Errorf("ShouldBind: %v", err)
		return
	}
	log.Debugf("ShouldBind: %+v", u)

	ret, err := redis_user.GetUserInfoByUserInfoHash(u.Redis_Key)
	if err != nil {
		handler.SendResponse(c, errno.ErrDatabase, nil)
		log.Errorf("err: %v", err)
		return
	}
	handler.SendResponse(c,
		nil,
		ret,
	)
	return
}
