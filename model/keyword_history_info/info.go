package keyword_history_info

import (
	"github.com/AnnonaOrg/annona_core/model"
)

type KeyworldHistoryInfo struct {
	model.Model

	ChatId             int64  `json:"chat_id" form:"chat_id" gorm:"column:chat_id;"`
	SenderId           int64  `json:"sender_id" form:"sender_id" gorm:"column:sender_id;"`
	SenderUsername     string `json:"sender_username" form:"sender_username" gorm:"column:sender_username;"`
	MessageId          int64  `json:"message_id" form:"message_id" gorm:"column:message_id;"`
	MessageLink        string `json:"message_link" form:"message_link" gorm:"column:message_link;"`
	MessageContentText string `json:"message_content_text" form:"message_content_text" gorm:"column:message_content_text;"`

	KeyWorld string `json:"key_world" form:"key_world" gorm:"column:key_world;"`
	Total    int64  `json:"total" form:"total" gorm:"column:total;"`

	// 核验请求id
	ById string `json:"by_id" form:"by_id" gorm:"-"`

	Page   int    `json:"page" form:"page" gorm:"-"`
	Size   int    `json:"size" form:"size" gorm:"-"`
	Filter string `json:"filter" form:"filter" gorm:"-"`

	Note     string `json:"note" form:"-" gorm:"-"`
	NoteHtml string `json:"note_html" form:"-" gorm:"-"`
}
