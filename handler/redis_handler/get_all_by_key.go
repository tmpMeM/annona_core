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

// 获取信息列表
// 用户关键词列表 constvar.REDIS_SET_KEY_USER_PREFIX_KEYWORD + userInfoHash
// 用户屏蔽来源会话id列表  constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKFORMCHATID + userInfoHash
// 用户屏蔽词列表 constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKKEYWORD + userInfoHash
// 用户屏蔽senderid列表 constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKFORMSENDERID + userInfoHash
func GetAllByKey(c *gin.Context) {
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

	list, err := redis_user.GetAllByKey(u.Redis_Key)
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
