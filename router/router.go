package router

import (
	"github.com/AnnonaOrg/annona_core/handler/api_handler"
	"github.com/AnnonaOrg/annona_core/handler/blockformchatid_handler"
	"github.com/AnnonaOrg/annona_core/handler/blockformsenderid_handler"
	"github.com/AnnonaOrg/annona_core/handler/blockword_handler"
	"github.com/AnnonaOrg/annona_core/handler/card_handler"
	"github.com/AnnonaOrg/annona_core/handler/keyword_handler"
	"github.com/AnnonaOrg/annona_core/handler/redis_handler"
	"github.com/AnnonaOrg/annona_core/handler/telebot_handler"
	"github.com/AnnonaOrg/annona_core/handler/user_handler"
	"github.com/AnnonaOrg/annona_core/router/middleware"

	"github.com/gin-gonic/gin"
)

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(mw...)
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	g.MaxMultipartMemory = 8 << 20 // 8 MiB

	g.NoRoute(api_handler.ApiNotFound)
	g.GET("/", api_handler.ApiHello)
	g.GET("/ping", api_handler.ApiPing)

	// 通用请求接口
	apis := g.Group("/apis")
	{
		// 用户
		// 创建
		apis.POST("/v1/user/item/add", user_handler.Create)
		// 查询
		apis.POST("/v1/user/item/get", user_handler.Get)
		// 核销卡
		apis.POST("/v1/user/item/renew/:uuid", user_handler.Renew)
		// 列表
		apis.POST("/v1/user/list", user_handler.List)
		// 设置notice chatid
		apis.POST("/v1/user/item/updatenoticechatid/:chatid", user_handler.UpdateNoticeChatId)
		// 签到 sign
		apis.POST("/v1/user/item/sign", user_handler.Sign)

		// 卡
		// 创建
		apis.POST("/v1/card/item/convert", card_handler.Convert)
		apis.POST("/v1/card/item/add", card_handler.Create)
		// 查询
		apis.POST("/v1/card/item/get", card_handler.Get)
		// 列表
		apis.POST("/v1/card/list", card_handler.List)

		// Bot
		// 创建
		apis.POST("/v1/telebot/item/add", telebot_handler.Create)
		// 查询
		apis.POST("/v1/telebot/item/get", telebot_handler.Get)
		// 列表
		apis.POST("/v1/telebot/list", telebot_handler.List)

		// Keyword
		// 创建
		apis.POST("/v1/keyword/item/add", keyword_handler.Create)
		//删除
		apis.POST("/v1/keyword/item/del", keyword_handler.Delete)
		//删除all
		apis.POST("/v1/keyword/item/delall", keyword_handler.DeleteAll)
		// 列表
		apis.POST("/v1/keyword/list", keyword_handler.List)

		// blockword
		// 创建
		apis.POST("/v1/blockword/item/add", blockword_handler.Create)
		//删除
		apis.POST("/v1/blockword/item/del", blockword_handler.Delete)
		//删除all
		apis.POST("/v1/blockword/item/delall", blockword_handler.DeleteAll)
		// 列表
		apis.POST("/v1/blockword/list", blockword_handler.List)

		// blockformchatid
		// 创建
		apis.POST("/v1/blockformchatid/item/add", blockformchatid_handler.Create)
		//删除
		apis.POST("/v1/blockformchatid/item/del", blockformchatid_handler.Delete)
		//删除all
		apis.POST("/v1/blockformchatid/item/delall", blockformchatid_handler.DeleteAll)
		// 列表
		apis.POST("/v1/blockformchatid/list", blockformchatid_handler.List)

		// blockformsenderid
		// 创建
		apis.POST("/v1/blockformsenderid/item/add", blockformsenderid_handler.Create)
		//删除
		apis.POST("/v1/blockformsenderid/item/del", blockformsenderid_handler.Delete)
		//删除all
		apis.POST("/v1/blockformsenderid/item/delall", blockformsenderid_handler.DeleteAll)
		// 列表
		apis.POST("/v1/blockformsenderid/list", blockformsenderid_handler.List)

		// redis db
		// 根据key获取集合信息
		// 用户关键词列表 constvar.REDIS_SET_KEY_USER_PREFIX_KEYWORD + userInfoHash
		// 用户屏蔽来源会话id列表  constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKFORMCHATID + userInfoHash
		// 用户屏蔽词列表 constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKKEYWORD + userInfoHash
		// 用户屏蔽senderid列表 constvar.REDIS_SET_KEY_USER_PREFIX_BLOCKFORMSENDERID + userInfoHash
		apis.POST("/v1/db/redis/allbykey/list", redis_handler.GetAllByKey)
		// 获取所有关键词列表
		apis.POST("/v1/db/redis/allkeyword/list", redis_handler.GetAllKeyword)
		// 获取所有屏蔽词列表
		apis.POST("/v1/db/redis/allblockword/list", redis_handler.GetAllBlockword)
		// 获取所有用户信息hash
		apis.POST("/v1/db/redis/alluserinfohash/list", redis_handler.GetAllUserInfoHash)
		// 获取机器人token
		apis.POST("/v1/db/redis/bottoken/item/get", redis_handler.GetBotTokenByBotId)
		// 获取用户信息
		apis.POST("/v1/db/redis/userinfo/item/get", redis_handler.GetUserInfoByUserInfoHash)
	}

	return g
}
