package load_db2redis

import (
	"fmt"
	"time"

	"github.com/AnnonaOrg/annona_core/model/blockformsenderid_info"

	"github.com/AnnonaOrg/annona_core/internal/constvar"
	"github.com/AnnonaOrg/annona_core/internal/dbredis"
	"github.com/AnnonaOrg/annona_core/model/blockformchatid_info"
	"github.com/AnnonaOrg/annona_core/model/blockword_info"
	"github.com/AnnonaOrg/annona_core/model/keyword_info"
	"github.com/AnnonaOrg/annona_core/model/telebot_info"
	"github.com/AnnonaOrg/annona_core/model/user_info"
	log "github.com/sirupsen/logrus"
)

// 加载机器人信息
func LoadBot(botInfo *telebot_info.TeleBotInfo) {
	botToken := botInfo.TelegramBotToken
	botId := fmt.Sprintf("%d", botInfo.TelegramId)
	if err := dbredis.AddKeyValue(
		constvar.REDIS_KEY_PREFIX_BOT_TOKEN+botId,
		botToken,
	); err != nil {
		log.Errorf("AddKeyValue(%s,%s): %v",
			constvar.REDIS_KEY_PREFIX_BOT_TOKEN+botId,
			botToken,
			err,
		)
	}
}

// 加载用户信息
func LoadUser(userInfoA *user_info.UserInfo, isClearAll bool) {
	timeDuration := time.Unix(userInfoA.Exp, 0).Sub(time.Now())
	userInfoHash := userInfoA.InfoHash
	var userInfo user_info.UserInfo
	userInfo.AccoundPlatform = userInfoA.AccoundPlatform
	userInfo.AccoundPlatformId = userInfoA.AccoundPlatformId
	userInfo.Exp = userInfoA.Exp
	userInfo.InfoHash = userInfoA.InfoHash
	userInfo.LastCardUUID = userInfoA.LastCardUUID
	userInfo.TelegramChatId = userInfoA.TelegramChatId
	userInfo.TelegramStartBotId = userInfoA.TelegramStartBotId
	userInfo.TelegramUsername = userInfoA.TelegramUsername
	userInfo.TelegramNoticeChatId = userInfoA.TelegramNoticeChatId

	if timeDuration < time.Hour {
		log.Debugf("Ignore the user %s (%d) Exp: %d timeDuration: %d", userInfo.TelegramUsername, userInfo.TelegramChatId, userInfo.Exp, timeDuration)
		return
	}
	// 添加到所有用户集合
	if err := dbredis.AddToSetWithExpiration(
		constvar.REDIS_SET_KEY_ALL_USER,
		userInfoHash,
		timeDuration,
	); err != nil {
		log.Errorf("AddToSetWithExpiration(%s,%s): %v", constvar.REDIS_SET_KEY_ALL_USER, userInfoHash, err)
	} else {
		log.Debugf("AddToSetWithExpiration(%s,%s):ok", constvar.REDIS_SET_KEY_ALL_USER, userInfoHash)
	}
	// 添加到单用户键值对
	if err := dbredis.AddKeyValueWithExpiration(
		constvar.REDIS_KEY_PREFIX_USERINFO+userInfoHash,
		userInfo,
		timeDuration,
	); err != nil {
		log.Errorf("AddKeyValueWithExpiration(%s,%+v): %v",
			constvar.REDIS_KEY_PREFIX_USERINFO+userInfoHash, userInfo, err)
	} else {
		log.Debugf("AddKeyValueWithExpiration(%s,%+v)", constvar.REDIS_KEY_PREFIX_USERINFO+userInfoHash, userInfo)
		// log.Debugf("AddKeyValueWithExpiration*(%s,%+v)", constvar.REDIS_KEY_PREFIX_USERINFO+userInfoHash, *userInfo)
	}

	if list, count, err := keyword_info.GetAllByOwnerInfoHashToString(userInfoHash); err != nil || count == 0 {
		log.Debugf("the user (%s) keyword_info list is NULL ", userInfoHash)
		return
	} else {
		// 清空原有数据
		if isClearAll {
			DelAllKeyword(userInfoHash)
		}

		// 添加到单用户关键词集合
		if err := dbredis.AddMultipleToSetWithExpiration(
			constvar.REDIS_SET_KEY_USER_PREFIX_KEYWORD+userInfoHash,
			timeDuration,
			list,
		); err != nil {
			log.Errorf("AddMultipleToSetWithExpiration(%s,%+v): %v",
				constvar.REDIS_SET_KEY_USER_PREFIX_KEYWORD+userInfoHash, list, err,
			)
		} else {
			log.Debugf("AddMultipleToSetWithExpiration(%s,%+v): ok",
				constvar.REDIS_SET_KEY_USER_PREFIX_KEYWORD+userInfoHash, list,
			)
		}
		// 添加到所有关键词集合
		if err := dbredis.AddMultipleToSet(
			constvar.REDIS_SET_KEY_ALL_KEYWORD,
			list,
		); err != nil {
			log.Errorf("AddMultipleToSet(%s,%+v): %v",
				constvar.REDIS_SET_KEY_ALL_KEYWORD, list, err,
			)
		} else {
			log.Debugf("AddMultipleToSet(%s,%+v): ok",
				constvar.REDIS_SET_KEY_ALL_KEYWORD, list,
			)
		}
	}

	if list, count, err := blockword_info.GetAllByOwnerInfoHashToString(userInfoHash); err != nil || count == 0 {
		log.Debugf("the user (%s) blockword_info list is NULL ", userInfoHash)
		// return
	} else {
		// 清空原有数据
		if isClearAll {
			DelAllBlockword(userInfoHash)
		}

		// 添加到单用户屏蔽关键词集合
		if err := dbredis.AddMultipleToSetWithExpiration(
			constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKKEYWORD+userInfoHash,
			timeDuration,
			list,
		); err != nil {
			log.Errorf("AddMultipleToSetWithExpiration(%s,%+v): %v",
				constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKKEYWORD+userInfoHash,
				list, err,
			)
		} else {
			log.Debugf("AddMultipleToSetWithExpiration(%s,%+v): ok",
				constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKKEYWORD+userInfoHash,
				list,
			)
		}
		// 添加到所有屏蔽关键词集合
		if err := dbredis.AddMultipleToSet(
			constvar.REDIS_SET_KEY_ALL_BLOCKKEYWORD,
			list,
		); err != nil {
			log.Errorf("AddMultipleToSet(%s,%+v): %v",
				constvar.REDIS_SET_KEY_ALL_BLOCKKEYWORD,
				list, err,
			)
		} else {
			log.Debugf("AddMultipleToSet(%s,%+v): ok",
				constvar.REDIS_SET_KEY_ALL_BLOCKKEYWORD,
				list,
			)
		}
	}

	if list, count, err := blockformchatid_info.GetAllByOwnerInfoHashToString(userInfoHash); err != nil || count == 0 {
		log.Debugf("the user (%s) blockformchatid_info list is NULL ", userInfoHash)
		// return
	} else {
		// 清空原有数据
		if isClearAll {
			DelAllBlockformchatid(userInfoHash)
		}

		// 添加到单用户屏蔽来源会话id集合
		if err := dbredis.AddMultipleToSetWithExpiration(
			constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKFORMCHATID+userInfoHash,
			timeDuration,
			list,
		); err != nil {
			log.Errorf("AddMultipleToSetWithExpiration(%s,%+v): %v",
				constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKFORMCHATID+userInfoHash,
				list, err,
			)
		} else {
			log.Debugf("AddMultipleToSetWithExpiration(%s,%+v): ok",
				constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKFORMCHATID+userInfoHash,
				list,
			)
		}
		// 添加到所有屏蔽来源会话id集合
		if err := dbredis.AddMultipleToSet(
			constvar.REDIS_SET_KEY_ALL_BLOCKFORMCHATID,
			list,
		); err != nil {
			log.Errorf("AddMultipleToSet(%s,%+v): %v",
				constvar.REDIS_SET_KEY_ALL_BLOCKFORMCHATID,
				list, err,
			)
		} else {
			log.Debugf("AddMultipleToSet(%s,%+v): ok",
				constvar.REDIS_SET_KEY_ALL_BLOCKFORMCHATID,
				list,
			)
		}
	}

	if list, count, err := blockformsenderid_info.GetAllByOwnerInfoHashToString(userInfoHash); err != nil || count == 0 {
		log.Debugf("the user (%s) blockformsenderid_info list is NULL ", userInfoHash)
		// return
	} else {
		// 清空原有数据
		if isClearAll {
			DelAllBlockformsenderid(userInfoHash)
		}

		// 添加到单用户屏蔽来源Senderid集合
		if err := dbredis.AddMultipleToSetWithExpiration(
			constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKFORMSENDERID+userInfoHash,
			timeDuration,
			list,
		); err != nil {
			log.Errorf("AddMultipleToSetWithExpiration(%s,%+v): %v",
				constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKFORMSENDERID+userInfoHash,
				list, err,
			)
		} else {
			log.Debugf("AddMultipleToSetWithExpiration(%s,%+v): ok",
				constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKFORMSENDERID+userInfoHash,
				list,
			)
		}
		// 添加到所有屏蔽来源Senderid集合
		if err := dbredis.AddMultipleToSet(
			constvar.REDIS_SET_KEY_ALL_BLOCKFORMSENDERID,
			list,
		); err != nil {
			log.Errorf("AddMultipleToSet(%s,%+v): %v",
				constvar.REDIS_SET_KEY_ALL_BLOCKFORMSENDERID,
				list, err,
			)
		} else {
			log.Debugf("AddMultipleToSet(%s,%+v): ok",
				constvar.REDIS_SET_KEY_ALL_BLOCKFORMSENDERID,
				list,
			)
		}
	}
}

func LoadKeyword(userInfoHash, itemStr string) {
	dbredis.AddToSet(
		constvar.REDIS_SET_KEY_USER_PREFIX_KEYWORD+userInfoHash,
		itemStr,
	)
	dbredis.AddToSet(
		constvar.REDIS_SET_KEY_ALL_KEYWORD,
		itemStr,
	)
}
func DelKeyword(userInfoHash, itemStr string) {
	dbredis.RemoveFromSet(
		constvar.REDIS_SET_KEY_USER_PREFIX_KEYWORD+userInfoHash,
		itemStr,
	)
}
func DelAllKeyword(userInfoHash string) {
	dbredis.RemoveAllFromSet(
		constvar.REDIS_SET_KEY_USER_PREFIX_KEYWORD + userInfoHash,
	)
}

func LoadBlockword(userInfoHash, itemStr string) {
	dbredis.AddToSet(
		constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKKEYWORD+userInfoHash,
		itemStr,
	)
	dbredis.AddToSet(
		constvar.REDIS_SET_KEY_ALL_BLOCKKEYWORD,
		itemStr,
	)
}
func DelBlockword(userInfoHash, itemStr string) {
	dbredis.RemoveFromSet(
		constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKKEYWORD+userInfoHash,
		itemStr,
	)
}
func DelAllBlockword(userInfoHash string) {
	dbredis.RemoveAllFromSet(
		constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKKEYWORD + userInfoHash,
	)
}

func LoadBlockformchatid(userInfoHash, itemStr string) {
	dbredis.AddToSet(
		constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKFORMCHATID+userInfoHash,
		itemStr,
	)
	dbredis.AddToSet(
		constvar.REDIS_SET_KEY_ALL_BLOCKFORMCHATID,
		itemStr,
	)
}
func DelBlockformchatid(userInfoHash, itemStr string) {
	dbredis.RemoveFromSet(
		constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKFORMCHATID+userInfoHash,
		itemStr,
	)
}
func DelAllBlockformchatid(userInfoHash string) {
	dbredis.RemoveAllFromSet(
		constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKFORMCHATID + userInfoHash,
	)
}

func LoadBlockformsenderid(userInfoHash, itemStr string) {
	dbredis.AddToSet(
		constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKFORMSENDERID+userInfoHash,
		itemStr,
	)
	dbredis.AddToSet(
		constvar.REDIS_SET_KEY_ALL_BLOCKFORMSENDERID,
		itemStr,
	)
}
func DelBlockformsenderid(userInfoHash, itemStr string) {
	dbredis.RemoveFromSet(
		constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKFORMSENDERID+userInfoHash,
		itemStr,
	)
}
func DelAllBlockformsenderid(userInfoHash string) {
	dbredis.RemoveAllFromSet(
		constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKFORMSENDERID + userInfoHash,
	)
}
