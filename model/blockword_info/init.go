package blockword_info

import (
	"fmt"
	"sync"

	"github.com/AnnonaOrg/annona_core/model"
)

var once sync.Once

func Init() {
	if !model.DBIsReady {
		return
	}
	once.Do(func() {

		if err := model.DB.Self.AutoMigrate(&BlockworldInfo{}); err != nil {
			fmt.Println(err)
			return
		}

	})
}
