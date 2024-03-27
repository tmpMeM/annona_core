package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AnnonaOrg/annona_core/internal/dbredis"

	_ "github.com/AnnonaOrg/annona_core/internal/dotenv"
	_ "github.com/AnnonaOrg/annona_core/internal/log"

	"github.com/AnnonaOrg/annona_core/internal/constvar"
	"github.com/AnnonaOrg/annona_core/internal/kvstore"
	"github.com/AnnonaOrg/annona_core/model"
	"github.com/AnnonaOrg/annona_core/router"
	"github.com/AnnonaOrg/annona_core/router/middleware"
	"github.com/AnnonaOrg/osenv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	versionArg = pflag.BoolP("version", "v", false, "show version info.")
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("run time panic:%v\n", err)
		}
	}()

	pflag.Parse()
	if *versionArg {
		fmt.Println(constvar.APPAbout())
		return
	}

	fmt.Printf("%s %s\n%s\n",
		constvar.APPName(), constvar.APPVersion(), constvar.APPDesc(),
	)
	runAt := "运行在"
	if osenv.IsInDocker() {
		runAt = runAt + "(Docker)"
	}
	runAt = runAt + ": " + osenv.Getwd()
	fmt.Println(runAt)
	time.Sleep(time.Second * 3)

	// 数据库连接初始化
	if err := model.Init(); err != nil {
		log.Fatalf(
			"数据库( %s )连接初始化出错: %v",
			osenv.GetServerDbType(), err,
		)
	}
	defer model.Close()

	if err := kvstore.Init(); err != nil {
		log.Fatalf(
			"数据库(REDIS: %s)连接初始化出错: %v",
			osenv.GetServerDbRedisAddress(), err,
		)
	}
	defer kvstore.Close()
	if err := dbredis.Init(); err != nil {
		log.Fatalf(
			"数据库(REDIS: %s)连接初始化出错: %v",
			osenv.GetServerDbRedisAddress(), err,
		)
	}
	defer dbredis.Close()

	// Set gin mode.
	ginMode := osenv.GetServerGinRunode()
	if ginMode == "" {
		ginMode = "release"
	}
	gin.SetMode(ginMode)
	//Create the Gin engine.
	g := gin.New()
	//Routes.
	router.Load(
		g,
		middleware.Logging(),
		middleware.RequestId(),
	)

	go func() {
		if err := pingServer(); err != nil {
			log.Fatalf(
				"(%s)没有响应，请检查配置及网络状态: %v",
				constvar.APPName(), err,
			)
		}
		log.Infof("(%s)成功部署，服务地址:%s", constvar.APPName(), osenv.GetServerUrl())
	}()

	go mainTask()

	addr := ":" + osenv.GetServerPort()
	if err := http.ListenAndServe(addr, g); err != nil {
		log.Errorf(
			"(%s)出错了，需要重启: %v",
			constvar.APPName(), err,
		)
	}
}
