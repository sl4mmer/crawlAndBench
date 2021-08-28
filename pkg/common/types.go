package common

import "context"

type Opts struct {
	Retries           int `mapstructure:"retries"`
	MaxConcurrency    int `mapstructure:"max_concurrency"`
	WorkerDelayMillis int `mapstructure:"worker_delay_millis"`
	HostsParallel     int `mapstructure:"hosts_parallel"`
	CacheTtlMillis    int `mapstructure:"cache_ttl_millis"`
}

type Queerer interface {
	Query(ctx context.Context, q string) (map[string]int, error)
}
