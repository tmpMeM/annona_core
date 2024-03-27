package api_handler

import (
	"github.com/AnnonaOrg/annona_core/handler"
	"github.com/AnnonaOrg/annona_core/internal/constvar"
	"github.com/AnnonaOrg/pkg/errno"
	"github.com/gin-gonic/gin"
)

// 404 Not found
func ApiNotFound(c *gin.Context) {
	handler.SendResponse(c, errno.Err404, constvar.APPDesc404())
}

// API Hello
func ApiHello(c *gin.Context) {
	handler.SendResponse(c, errno.SayHello, constvar.APPDesc())
}

// ping
func ApiPing(c *gin.Context) {
	handler.SendResponse(c, errno.PONG, constvar.APPVersion())
}
