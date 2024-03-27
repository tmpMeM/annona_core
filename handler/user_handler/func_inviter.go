package user_handler

import (
	"fmt"

	"github.com/AnnonaOrg/annona_core/model/telebot_info"

	"github.com/AnnonaOrg/annona_core/internal/notice_message"
	"github.com/AnnonaOrg/annona_core/internal/redis_user"
	"github.com/AnnonaOrg/annona_core/model/user_info"
	log "github.com/sirupsen/logrus"
)

// userInfoA 邀请者
// userInfoB 被邀请者
func SendInviterNoticeMessage(userInfoA *user_info.UserInfo, userInfoB *user_info.UserInfo) {
	startBotId := fmt.Sprintf("%d", userInfoA.TelegramStartBotId)
	botToken := ""
	if botTokenTmp, err := redis_user.GetBotTokenByBotId(startBotId); err != nil {
		log.Errorf("GetBotTokenByBotId(%s)Fail: %v", startBotId, err)

	} else {
		botToken = botTokenTmp
		log.Debugf("GetBotTokenByBotId(%s)Success: %v", startBotId, botToken)
	}
	if len(botToken) == 0 {
		if telebotInfo, err := telebot_info.GetById(userInfoA.TelegramStartBotId); err != nil {
			log.Errorf("telebot_info.GetById(%d)Fail: %v", userInfoA.TelegramStartBotId, err)
		} else {
			botToken = telebotInfo.TelegramBotToken
		}
	}
	if len(botToken) == 0 {
		log.Errorf("Bot(%s) Token Is NULL", startBotId)
		return
	}

	messageText := fmt.Sprintf(
		"🎉 %s %s ID%s 已经接受你的邀请，你获得了 3小时 使用时间奖励！",
		userInfoB.TelegramFirstname, userInfoB.TelegramLasttname, userInfoB.AccoundPlatformId,
	)
	chatID := userInfoA.AccoundPlatformId
	topicID := new(int64)

	notice_message.SendNoticeMessage(
		messageText,
		botToken, chatID, topicID,
		true, true, true,
	)
}
