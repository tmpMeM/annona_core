package main

import (
	"fmt"
	"time"

	"github.com/AnnonaOrg/annona_core/internal/constvar"
	"github.com/AnnonaOrg/annona_core/internal/initialize"
	"github.com/AnnonaOrg/annona_core/internal/tasks"
	"github.com/AnnonaOrg/osenv"
	"github.com/clin003/util"
	log "github.com/sirupsen/logrus"
)

func mainTask() {
	initialize.Init()

	go doTask()
}

// 自检openAPI服务是否正常运行
func pingServer() error {
	apiURL := osenv.GetServerUrl()
	for i := 0; i < 10; i++ {

		if util.CheckPingBaseURL(apiURL) {
			return nil
		}

		log.Debugf(
			"(%s)等待自检, 1秒后重试(%d) %s.",
			constvar.APPName(), i, apiURL,
		)
		time.Sleep(time.Second * 2)
	}
	return fmt.Errorf(
		"(%s)自检失败 %s.",
		constvar.APPName(), apiURL,
	)
}

func doTask() {
	tasks.Init()
}
