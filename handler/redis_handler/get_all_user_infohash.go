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

// 获取所有用户列表
func GetAllUserInfoHash(c *gin.Context) {
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

	list, err := redis_user.GetAllUserInfoHashList()
	if err != nil {
		handler.SendResponse(c, errno.ErrDatabase, nil)
		log.Errorf("err: %v", err)
		return
	}
	handler.SendResponse(c,
		nil,
		handler.ListResponse{Items: list, Total: int64(len(list))},
	)
	return
}
