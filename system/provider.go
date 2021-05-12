package system

import (
	"log"

	"github.com/powerpuffpenguin/sessionid"
	"github.com/powerpuffpenguin/sessionid/provider/bolt"
	"github.com/powerpuffpenguin/sessionid/provider/redis"
	"github.com/powerpuffpenguin/sessionid_server/configure"
)

var provider sessionid.Provider

func DefaultProvider() sessionid.Provider {
	return provider
}
func initProvider() {
	cnf := configure.DefaultConfigure().Provider
	switch cnf.Backend {
	case `Memory`:
		initProviderMemory(&cnf.Memory)
	case `Redis`:
		initProviderRedis(&cnf.Redis)
	case `Bolt`:
		initProviderBolt(&cnf.Bolt)
	default:
		panic(`unknow Provider.Backend : ` + cnf.Backend)
	}
}
func initProviderMemory(cnf *configure.ProviderMemory) {
	provider = sessionid.NewProvider(
		sessionid.WithProviderAccess(cnf.Access),
		sessionid.WithProviderRefresh(cnf.Refresh),
		sessionid.WithProviderMaxSize(cnf.MaxSize),
		sessionid.WithProviderCheckBatch(cnf.Batch),
		sessionid.WithProviderClear(cnf.Clear),
	)
}
func initProviderRedis(cnf *configure.ProviderRedis) {
	var e error
	provider, e = redis.New(
		redis.WithURL(cnf.URL),
		redis.WithAccess(cnf.Access),
		redis.WithRefresh(cnf.Refresh),
		redis.WithCheckBatch(cnf.Batch),
		redis.WithKeyPrefix(cnf.KeyPrefix),
		redis.WithMetadataKey(cnf.MetadataKey),
	)
	if e != nil {
		log.Fatalln(e)
	}
}
func initProviderBolt(cnf *configure.ProviderBolt) {
	var e error
	provider, e = bolt.New(
		bolt.WithFilename(cnf.Filename),
		bolt.WithAccess(cnf.Access),
		bolt.WithRefresh(cnf.Refresh),
		bolt.WithMaxSize(cnf.MaxSize),
		bolt.WithCheckBatch(cnf.Batch),
		bolt.WithClear(cnf.Clear),
	)
	if e != nil {
		log.Fatalln(e)
	}
}
