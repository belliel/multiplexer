package services

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	httpRequestTimeout = time.Second
)

const MaxWorkers = 4

type UrlResponseStruct struct {
	data interface{}
	url string
	err error
}

func ProcessUrls(ctx context.Context, urls []string) (map[string]interface{}, error) {
	var (
		/**
			Default map can be used
			Because no concurrency jobs, only sequential write
		*/
		syncMap = sync.Map{}
		cancelChan = make(chan struct{})
		results = make(chan UrlResponseStruct, len(urls))
		jobs = make(chan string, len(urls))
	)

	for w := 0; w < MaxWorkers; w++ {
		go urlWorker(ctx, jobs, results, cancelChan)
	}
	for j := 0; j <= len(urls); j++ {
		jobs <- urls[j]
	}

	for result := range results {
		if result.err != nil {
			cancelChan <- struct{}{}
			return nil, result.err
		}
		syncMap.Store(result.url, result.data)
	}


	var result = make(map[string]interface{}, len(urls))

	syncMap.Range(func(key, value interface{}) bool {
		result[key.(string)] = value
		return true
	})


	return result, nil
}

func urlProcess(ctx context.Context, url string) UrlResponseStruct {
	result := UrlResponseStruct{}
	
	client := &http.Client{}

	requestContext, cancel := context.WithTimeout(ctx, httpRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(requestContext, "GET", url, nil)

	if err != nil {
		result.err = err
		return result
	}

	response, err := client.Do(req)
	select {
	case <-requestContext.Done():
		result.err = ctx.Err()
		return result
	default:
		break
	}

	if err != nil {
		result.err = err
		return result
	}
	defer response.Body.Close()


	if response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusMultipleChoices {
		result.data, result.err = io.ReadAll(response.Body)
		result.data = string(result.data.([]byte))
	}

	client.CloseIdleConnections()

	return result
}

func urlWorker(ctx context.Context, jobs <- chan string, results chan <- UrlResponseStruct, cancelChan <- chan struct{})  {
	for {
		select {
		case <- cancelChan:
			return
		case j := <- jobs:
			results <- urlProcess(ctx, j)
		}
	}
}