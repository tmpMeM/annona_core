package repository

import (
	// "github.com/shitpostingio/autopostingbot/config/structs"
	// "github.com/zelenin/go-tdlib/client"
	"github.com/AnnonaOrg/gokv"

	"github.com/redis/go-redis/v9"
)

var (
	// // // Config contains all the configuration structures.
	// // Config *structs.Config

	// // Tdlib is the Telegram client instance.
	// Tdlib *client.Client

	// // Me represents the current bot as a Telegram client.User.
	// Me *client.User

	KVStore gokv.Store

	DBRedis *redis.Client
)
