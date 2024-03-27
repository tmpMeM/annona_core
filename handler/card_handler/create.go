package card_handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AnnonaOrg/osenv"

	"github.com/AnnonaOrg/annona_core/handler"
	"github.com/AnnonaOrg/annona_core/model/card_info"
	"github.com/AnnonaOrg/pkg/errno"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Create(c *gin.Context) {
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
		"card_info Add u.ById(%s) u.NDay(%s) u.NNum(%d)",
		u.ById, u.NDay, u.NNum,
	)
	// 核验请求id
	if !osenv.IsBotManagerIDStr(u.ById) {
		handler.SendResponse(c, errno.ErrBadRequest, nil)
		return
	}

	var expTime int64
	if strings.HasSuffix(u.NDay, "d") {
		u.NDay = strings.ReplaceAll(u.NDay, "d", "h")
	}
	// if nDay, err := time.ParseDuration(u.NDay); err != nil {
	// 	log.Errorf("time.ParseDuration(%s) %v", u.NDay, err)
	// 	expTime = int64(24 * 30 * time.Hour) // time.Now().Add(24 * 30 * time.Hour).Unix()
	// } else {
	// 	expTime = int64(24 * nDay * time.Hour) // time.Now().Add(24 * nDay).Unix()
	// }
	// 有效期 单位: 天 = 24 * time.Hour
	if nDay, err := strconv.ParseInt(u.NDay, 10, 64); err != nil {
		expTime = 30
	} else if nDay > 0 {
		expTime = nDay
	} else {
		expTime = 1
	}

	u.Exp = expTime
	//新增记录
	if list, count, err := u.CreateN(u.NNum); err != nil {
		handler.SendResponse(c,
			fmt.Errorf("新卡登记失败: %v %v", err, errno.ErrDatabase), //
			nil,
		)
		log.Errorf("新卡登记失败: %v", err)
		return
	} else {
		handler.SendResponse(c, nil, handler.ListResponse{Items: list, Total: count})
		return
		// handler.SendResponse(c, nil, u)
	}
}
