package tasks

import (
	"sync"
)

// var mutex sync.Mutex
var once sync.Once

func Init() {
	once.Do(func() {
		go taskSchedule()
	})
}
