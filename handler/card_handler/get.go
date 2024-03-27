package card_handler

import (
	"fmt"
	"strings"

	"github.com/AnnonaOrg/osenv"

	"github.com/AnnonaOrg/annona_core/handler"
	"github.com/AnnonaOrg/annona_core/model/card_info"
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

	var u card_info.CardInfo
	if err := c.ShouldBind(&u); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		return
	}
	log.Debugf("ShouldBind: %+v", u)

	if len(u.CardUUID) > 0 {
		cardInfo, err := card_info.GetCardInfoByUUID(u.CardUUID)
		if err != nil {
			handler.SendResponse(c,
				fmt.Errorf("查找卡信息(%s)失败: %v", u.CardUUID, errno.ErrDatabase),
				nil,
			)
			log.Errorf("查找卡信息: %+v,失败: %v", u, err)
			return
		}
		handler.SendResponse(c, nil, cardInfo)
		return
	} else {
		handler.SendResponse(c,
			fmt.Errorf("未知卡信息(%s): %v", u.CardUUID, errno.ErrValidation),
			nil,
		)
		return
	}
}
