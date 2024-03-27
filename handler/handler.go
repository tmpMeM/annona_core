package handler

import (
	"net/http"

	"github.com/AnnonaOrg/pkg/errno"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
type ListResponse struct {
	Total  int64       `json:"total" form:"total"`
	Items  interface{} `json:"items" form:"items"`
	AdList interface{} `json:"ad_list,omitempty" form:"ad_list"`
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	code, message := errno.DecodeErr(err)

	//always return http.StatusOK
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func SendRedirect(c *gin.Context, data string) {
	c.Redirect(http.StatusMovedPermanently, data)
}

func SendRedirect302(c *gin.Context, data string) {
	c.Redirect(http.StatusFound, data)
}

type ResponseEx struct {
	Code    int         `json:"status"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func SendResponseEx(c *gin.Context, err error, data interface{}) {
	code, message := errno.DecodeErr(err)

	//always return http.StatusOK
	c.JSON(http.StatusOK, ResponseEx{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
