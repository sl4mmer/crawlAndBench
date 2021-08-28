package bench

import (
	"context"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/sl4mmer/crawlAndBench/pkg/common"
)

func Run(ctx context.Context, url string, concurrency int, delayMillis int) int {
	var (
		counter int
		wg      sync.WaitGroup
	)

	results := make(chan *Task,concurrency)
	go func() {
		for task := range results {
			if !task.isError {
				counter++
			}
		}
	}()
	// вообще надо бы логарифмически наращивать число потоков
	for i := 0; i < concurrency; i++ {
		task := Task{
			url:     url,
			num:     i,
			isError: false,
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- runTask(ctx, task)
		}()
		time.Sleep(time.Duration(delayMillis) * time.Millisecond)
	}
	wg.Wait()
	close(results)
	return counter
}

func runTask(ctx context.Context, task Task) *Task {
	client := http.Client{}
	requestCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(requestCtx, "GET", task.url, nil)
	req.Header.Set("User-Agent", common.DefaultUA)

	startTime := time.Now()
	resp, err := client.Do(req)
	task.duration = time.Now().Sub(startTime)
	// редиректы по идее тоже успешный запрос
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 400 {
		task.isError = true
		return &task
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil || len(b) == 0 {
		task.isError = true
	}
	return &task
}
