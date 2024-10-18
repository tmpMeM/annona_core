package load_db2redis

import (
	"fmt"
	"sync"

	"github.com/AnnonaOrg/annona_core/model/telebot_info"

	"github.com/AnnonaOrg/annona_core/model/user_info"
)

var isReady bool
var lockReady sync.Mutex

func IsReady() bool {
	lockReady.Lock()
	defer lockReady.Unlock()

	return isReady
}
func setReady() {
	lockReady.Lock()
	defer lockReady.Unlock()
	isReady = true
	fmt.Println("LoadAll is ready")
}

func LoadAll(isClearAll bool) {
	LoadAllUser(isClearAll)
	LoadAllBot()
	setReady()
}

func LoadAllUser(isClearAll bool) {
	allUser, count, err := user_info.GetAll()
	if err != nil {

		return
	} else if count == 0 {

		return
	}

	for _, v := range allUser {
		user := v
		LoadUser(&user, isClearAll)
	}
}

func LoadAllBot() {
	allBot, count, err := telebot_info.GetAll()
	if err != nil {
		return
	} else if count == 0 {
		return
	}
	for _, v := range allBot {
		bot := v
		LoadBot(&bot)
	}

}
