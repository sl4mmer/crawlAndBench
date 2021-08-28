package controller

import (
	"context"
	"errors"
	"sync"
	"time"

	lwg "github.com/IAD/go-limited-waitgroup"

	"github.com/sl4mmer/crawlAndBench/pkg/bench"
	"github.com/sl4mmer/crawlAndBench/pkg/common"
	"github.com/sl4mmer/crawlAndBench/pkg/yandex"
)



type item struct {
	value int
	born  time.Time
}

type Service struct {
	sync.RWMutex
	opts  *common.Opts
	cache map[string]*item
}

func NewService(o *common.Opts) *Service {
	return &Service{
		opts:  o,
		cache: make(map[string]*item),
	}
}

func (s *Service) CacheClearanceLoop() {
	for {
		time.Sleep(100 * time.Millisecond)
		s.killSomebody()
	}
}

func (s *Service) killSomebody() {
	now := time.Now()
	s.Lock()
	for key, val := range s.cache {
		if now.Sub(val.born) > time.Duration(s.opts.CacheTtlMillis)*time.Millisecond {
			delete(s.cache, key)
		}
	}
	s.Unlock()
}

func (s *Service) Query(ctx context.Context, q string) (map[string]int, error) {
	var (
		yandexRetries int
		links         []*yandex.Item
		err           error
	)
	for {
		links, err = yandex.Search(q)
		if errors.Is(err, yandex.ErrRequestFailed) {
			if yandexRetries >= s.opts.Retries {
				return nil, err
			}
			yandexRetries++
			continue
		}
		if err != nil {
			return nil, err
		}
		break
	}

	result := make(map[string]int, len(links))
	wg := lwg.NewLimitedWaitGroup(s.opts.HostsParallel)
	for _, link := range links {
		s.RLock()
		cached, ok := s.cache[link.Host]
		s.RUnlock()
		if ok {
			result[link.Host] = cached.value
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			max := bench.Run(ctx, link.Url, s.opts.MaxConcurrency, s.opts.WorkerDelayMillis)
			s.Lock()
			s.cache[link.Host] = &item{
				value: max,
				born:  time.Now(),
			}
			result[link.Host] = max
			s.Unlock()
		}()

	}
	wg.Wait()
	return result, nil
}
