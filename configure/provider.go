package configure

import "time"

type Provider struct {
	Backend string
	Memory  ProviderMemory
	Redis   ProviderRedis
	Bolt    ProviderBolt
}
type ProviderMemory struct {
	Access  time.Duration
	Refresh time.Duration
	MaxSize int
	Batch   int
	Clear   time.Duration
}
type ProviderRedis struct {
	URL     string
	Access  time.Duration
	Refresh time.Duration

	Batch       int
	KeyPrefix   string
	MetadataKey string
}
type ProviderBolt struct {
	Filename string

	Access  time.Duration
	Refresh time.Duration
	MaxSize int
	Batch   int
	Clear   time.Duration
}
