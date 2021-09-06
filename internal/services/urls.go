package services

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	httpRequestTimeout = 1 * time.Second
)

const MaxWorkers = 4

type UrlResponseStruct struct {
	data interface{}
	url  string
	err  error
}

func ProcessUrls(ctx context.Context, urls []string) ([]map[string]interface{}, error) {
	var (
		cancelChan = make(chan struct{})
		processed  = make([]map[string]interface{}, 0, len(urls))
		results    = make(chan UrlResponseStruct)
		jobs       = make(chan string, len(urls))
		counter    = 0
	)

	for w := 0; w < MaxWorkers; w++ {
		go urlWorker(ctx, jobs, results, cancelChan)
	}
	for j := 0; j < len(urls); j++ {
		jobs <- urls[j]
	}

	for counter < len(urls) {
		select {
		case result := <-results:
			if result.err != nil {
				go func() { cancelChan <- struct{}{} }()
				return nil, result.err
			}
			processed = append(processed, map[string]interface{}{
				"counter": counter,
				"url":     result.url,
				"data":    result.data,
			})
			counter++
		}
	}
	cancelChan <- struct{}{}
	return processed, nil
}

func urlProcess(ctx context.Context, url string) UrlResponseStruct {
	result := UrlResponseStruct{
		url: url,
	}

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
		result.err = requestContext.Err()
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
	} else {
		result.err = errors.New(response.Status)
	}

	client.CloseIdleConnections()

	return result
}

func urlWorker(ctx context.Context, jobs <-chan string, results chan<- UrlResponseStruct, cancelChan <-chan struct{}) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-cancelChan:
			return
		case j := <-jobs:
			results <- urlProcess(ctx, j+"?q="+strconv.Itoa(int(time.Now().Unix())))
		}
	}
}
