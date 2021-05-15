local def = import "def.libsonnet";
local access = def.Duration.Hour*1;
local refresh = def.Duration.Hour*12*3;
local clear = def.Duration.Minute*30;
local provider = def.Provider;
{
    Backend: provider.Memory,
    Memory: {
            Access: access,
            Refresh: refresh,
            MaxSize: 1000,
            Batch: 128,
            Clear: clear,
    },
    Redis: {
        	URL: 'redis://redis:6379/0',
            Access: access,
            Refresh: refresh,
            Batch: 128,
            KeyPrefix: 'sessionid.provider.redis.',
            MetadataKey: '__private_provider_redis',
    },
    Bolt: {
        Filename: '/opt/server/data/db/sessionid.db',

        Access: access,
        Refresh: refresh,
        MaxSize: 10000,
        Batch: 128,
        Clear: clear,
    },
}