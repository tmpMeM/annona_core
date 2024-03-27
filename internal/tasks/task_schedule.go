package tasks

import (
	"sync"
	"time"

	"github.com/AnnonaOrg/annona_core/internal/load_db2redis"
)

func taskSchedule() {
	taskSchedule_1()
}

func taskSchedule_1() {
	isFirstRun := true
	var wgTask sync.WaitGroup
	for {
		wgTask.Add(1) //为计数器设置值

		go func() {

			if timeNow := time.Now().Hour(); timeNow < 6 && timeNow > 0 {
				load_db2redis.LoadAll(true)
				time.Sleep(3 * time.Hour)
			} else {
				load_db2redis.LoadAll(isFirstRun)
				time.Sleep(1 * time.Hour)
			}

			isFirstRun = false
			wgTask.Done()
		}()

		wgTask.Wait() //阻塞到计数器的值为0
		time.Sleep(time.Minute)
	}
}
