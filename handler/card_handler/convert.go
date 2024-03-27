package card_handler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/AnnonaOrg/annona_core/internal/load_db2redis"

	"github.com/AnnonaOrg/annona_core/handler"
	"github.com/AnnonaOrg/annona_core/model/card_info"
	"github.com/AnnonaOrg/annona_core/model/user_info"
	"github.com/AnnonaOrg/osenv"
	"github.com/AnnonaOrg/pkg/errno"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Convert(c *gin.Context) {
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
	// // 核验请求id
	// if !osenv.IsBotManagerIDStr(u.ById) {
	// 	handler.SendResponse(c, errno.ErrBadRequest, nil)
	// 	return
	// }
	// 获取用户信息
	userInfoTmp := &user_info.UserInfo{}
	userInfoTmp.AccoundPlatform = u.AccoundPlatform
	userInfoTmp.AccoundPlatformId = u.AccoundPlatformId
	userInfo, err := userInfoTmp.GetInfo()
	if err != nil {
		handler.SendResponse(c, errno.ErrDatabase, nil)
		log.Errorf("user GetInfo(%+v): %v", userInfoTmp, err)
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
	nDay := expTime
	nNum := u.NNum
	u.Exp = expTime
	// 检查用户Exp是否足够转换
	needExpTime := u.NNum * u.Exp * 24 * 60 * 60
	needTime := time.Now().Unix() + needExpTime
	if needTime > userInfo.Exp {
		needTimeStr := time.Unix(needTime, 0).Format(time.DateTime)
		userTimeStr := time.Unix(userInfo.Exp, 0).Format(time.DateTime)
		handler.SendResponse(c,
			fmt.Errorf("时间兑换卡失败: 现有到期时间(%s),需要到期时间(%s)", userTimeStr, needTimeStr),
			nil,
		)
		log.Errorf("convert Time To Card err: have(%s) need(%s)", userTimeStr, needTimeStr)
		return
	}
	userExpOld := userInfo.Exp
	userInfo.Exp = userInfo.Exp - needExpTime
	if err := userInfo.Update(); err != nil {
		handler.SendResponse(c,
			fmt.Errorf("更新用户信息失败: %v %v", errno.ErrDatabase, err),
			nil,
		)
		log.Errorf("Update Userinfo err: %v", err)
		return
	}
	//新增记录
	if list, count, err := u.CreateN(u.NNum); err != nil {
		userInfo.Exp = userExpOld
		if err := userInfo.Update(); err != nil {
			handler.SendResponse(c,
				fmt.Errorf("新卡登记失败 & 用户信息回退失败: %v %v with: nDay(%d) nNum(%d)", err, errno.ErrDatabase, nDay, nNum), //
				nil,
			)
			log.Errorf("遇到个倒霉蛋,新卡登记失败 & 用户信息回退失败: %v %v with: nDay(%d) nNum(%d)", err, errno.ErrDatabase, nDay, nNum)
			return
		}
		handler.SendResponse(c,
			fmt.Errorf("新卡登记失败(请重试): %v %v", err, errno.ErrDatabase), //
			nil,
		)
		log.Errorf("新卡登记失败: %v", err)
		return
	} else {
		handler.SendResponse(c, nil, handler.ListResponse{Items: list, Total: count})
		go func() {
			load_db2redis.LoadUser(userInfo, true)
			log.Debugf("Update: %+v", userInfo)
		}()
		return
	}
}
