package blockformsenderid_handler

import (
	"fmt"
	"strings"

	"github.com/AnnonaOrg/annona_core/internal/load_db2redis"

	"github.com/AnnonaOrg/annona_core/handler"
	model "github.com/AnnonaOrg/annona_core/model/blockformsenderid_info"
	"github.com/AnnonaOrg/annona_core/model/user_info"
	"github.com/AnnonaOrg/osenv"
	"github.com/AnnonaOrg/pkg/errno"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// 删除
func Delete(c *gin.Context) {
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

	var u model.BlockformsenderidInfo
	if err := c.ShouldBind(&u); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		log.Errorf("ShouldBind: %v", err)
		return
	}
	log.Debugf("ShouldBind: %+v", u)

	var userInfoHash string
	var userInfo user_info.UserInfo
	userInfo.AccoundPlatform = u.OwnerPlatform
	userInfo.AccoundPlatformId = fmt.Sprintf("%d", u.OwnerChatId)
	if userInfoTmp, err := userInfo.GetInfo(); err != nil {
		handler.SendResponse(c,
			fmt.Errorf("用户信息获取(%v)失败: %v %v", userInfo.AccoundPlatformId, err, errno.ErrDatabase),
			nil,
		)
		log.Errorf("GetInfo: %v", err)
		return
	} else {
		// u.OwnerInfoHash = userInfoTmp.InfoHash
		userInfoHash = userInfoTmp.InfoHash
	}

	if err := u.Delete(); err != nil {
		handler.SendResponse(c,
			fmt.Errorf("删除(%v)失败: %v %v", u.KeyWorld, err, errno.ErrDatabase),
			nil,
		)
		log.Errorf("Delete: %v", err)
		return
	}
	handler.SendResponse(c, nil, nil)

	// 同步删除缓存
	if len(userInfoHash) > 0 {
		go load_db2redis.DelBlockformsenderid(userInfoHash, u.KeyWorld)
	}
}

// 删除all
func DeleteAll(c *gin.Context) {
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

	var u model.BlockformsenderidInfo
	if err := c.ShouldBind(&u); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		log.Errorf("ShouldBind: %v", err)
		return
	}

	var userInfoHash string
	var userInfo user_info.UserInfo
	userInfo.AccoundPlatform = u.OwnerPlatform
	userInfo.AccoundPlatformId = fmt.Sprintf("%d", u.OwnerChatId)
	if userInfoTmp, err := userInfo.GetInfo(); err != nil {
		handler.SendResponse(c,
			fmt.Errorf("用户信息获取(%v)失败: %v %v", userInfo.AccoundPlatformId, err, errno.ErrDatabase),
			nil,
		)
		log.Errorf("GetInfo: %v", err)
		return
	} else {
		u.OwnerInfoHash = userInfoTmp.InfoHash
		userInfoHash = userInfoTmp.InfoHash
	}

	if err := u.Delete(); err != nil {
		handler.SendResponse(c,
			fmt.Errorf("删除(%v)失败: %v %v", userInfo.AccoundPlatformId, err, errno.ErrDatabase),
			nil,
		)
		log.Errorf("Delete: %v", err)
		return
	}
	handler.SendResponse(c, nil, nil)

	// 同步删除缓存
	if len(userInfoHash) > 0 {
		go load_db2redis.DelAllBlockformsenderid(userInfoHash)
	}
}
