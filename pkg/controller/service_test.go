package controller

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sl4mmer/crawlAndBench/pkg/common"
)

func TestQuery(t *testing.T) {
	s := NewService(&common.Opts{
		Retries:           3,
		MaxConcurrency: 100,
		WorkerDelayMillis: 10,
		HostsParallel:     5,
		CacheTtlMillis:    3600 * 1000,
	})
	ctx,_:= context.WithTimeout(context.Background(),time.Second*30)
	results, err := s.Query(ctx, "запрос")
	assert.NoError(t, err)
	assert.True(t, len(results) > 0)
}
