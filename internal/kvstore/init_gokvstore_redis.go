package kvstore

import (
	"github.com/AnnonaOrg/annona_core/internal/repository"
	"github.com/AnnonaOrg/gokv"
	"github.com/AnnonaOrg/gokv/redis"
	"github.com/AnnonaOrg/osenv"
)

var (
	kvStore gokv.Store
)

func GetRedisOptions() redis.Options {
	options := redis.DefaultOptions
	options.Address = osenv.GetServerDbRedisAddress()
	if pw := osenv.GetServerDbRedisPassword(); len(pw) > 0 {
		options.Password = pw
	}
	return options
}

func Init() error {
	var err error
	options := GetRedisOptions()
	kvStore, err = redis.NewClient(options)
	if err != nil {
		return err
	}
	// defer client.Close()
	// KVStore = client
	repository.KVStore = kvStore
	return nil
}

func NewClient() (gokv.Store, error) {
	options := GetRedisOptions()
	return redis.NewClient(options)
}

func Client() gokv.Store {
	return kvStore
}

func Close() error {
	return kvStore.Close()
}
