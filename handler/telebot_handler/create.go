package telebot_handler

import (
	"fmt"
	"strings"

	"github.com/AnnonaOrg/annona_core/internal/load_db2redis"

	"github.com/AnnonaOrg/osenv"

	"github.com/AnnonaOrg/annona_core/handler"
	model "github.com/AnnonaOrg/annona_core/model/telebot_info"
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

	var u model.TeleBotInfo
	if err := c.ShouldBind(&u); err != nil {
		handler.SendResponse(c, errno.ErrBind, nil)
		log.Errorf("ShouldBind: %v", err)
		return
	}
	log.Debugf("ShouldBind: %+v", u)

	log.Debugf(
		"info Add u.FirstName(%s) u.LastName(%s) u.TelegramUsername(%s) u.TelegramId(%d) u.TelegramBotToken(%s)",
		u.FirstName, u.LastName, u.TelegramUsername, u.TelegramId, u.TelegramBotToken,
	)

	if u.Count() > 0 {
		if err := u.Update(); err != nil {
			// handler.SendResponse(c,
			// 	fmt.Errorf("信息更新(%v)失败: %v %v", u, err, errno.ErrDatabase),
			// 	nil,
			// )
			log.Errorf("Update: %v", err)
		} else {
			go load_db2redis.LoadBot(&u)
			handler.SendResponse(c, nil, nil)
			return
		}
	}
	if err := u.Create(); err != nil {
		handler.SendResponse(c,
			fmt.Errorf("信息添加(%v)失败: %v %v", u, err, errno.ErrDatabase),
			nil,
		)
		log.Errorf("Create: %v", err)
		return
	} else {
		go load_db2redis.LoadBot(&u)
		handler.SendResponse(c, nil, nil)
		return
	}

}

// func load_db2redis
