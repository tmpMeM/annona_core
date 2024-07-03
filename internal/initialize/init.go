package initialize

import (
	"github.com/AnnonaOrg/annona_core/internal/kvstore"
	"github.com/AnnonaOrg/annona_core/model/blockformchatid_info"
	"github.com/AnnonaOrg/annona_core/model/blockformsenderid_info"
	"github.com/AnnonaOrg/annona_core/model/blockword_info"
	"github.com/AnnonaOrg/annona_core/model/card_info"
	"github.com/AnnonaOrg/annona_core/model/keyword_history_info"
	"github.com/AnnonaOrg/annona_core/model/keyword_info"
	"github.com/AnnonaOrg/annona_core/model/telebot_info"
	"github.com/AnnonaOrg/annona_core/model/user_info"
)

func Init() {
	card_info.Init()
	keyword_info.Init()
	user_info.Init()
	telebot_info.Init()

	blockword_info.Init()
	blockformchatid_info.Init()
	blockformsenderid_info.Init()

	keyword_history_info.Init()

	kvstore.Init()
}
