package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

const URI = "https://jsonplaceholder.typicode.com/todos/1"

type Worker struct {
	stopCh chan struct{}
}

func (worker *Worker) runRequester(resultCh chan int) {
	for {
		select {
		case <-worker.stopCh:
			return
		default:
			statusCode := checkUri()
			resultCh <- statusCode
		}

		time.Sleep(1)
	}
}

// write the result to the map
func (worker *Worker) runWriter(result ResultMap, resultCh chan int, mu *sync.RWMutex) {
	for {
		select {
		case statusCode := <-resultCh:
			mu.Lock()
			updateResult(result, statusCode)
			mu.Unlock()
		case <-worker.stopCh:
			return
		}

		time.Sleep(1)
	}
}

func (worker *Worker) runReporter(result ResultMap, mu *sync.RWMutex) {
	for {
		select {
		case <-worker.stopCh:
			return
		default:
			if mu.TryRLock() {
				reportResult(result)
				mu.RUnlock()

				time.Sleep(1 * time.Second)
			} else {
				time.Sleep(1)
			}
		}
	}
}

func (worker *Worker) stop() {
	worker.stopCh <- struct{}{}
}

func checkUri() int {
	res, err := http.Get(URI)

	if err != nil {
		log.Fatal(err.Error())
	}

	return res.StatusCode
}

func updateResult(result ResultMap, statusCode int) {
	if count, ok := result[statusCode]; ok {
		result[statusCode] = count + 1
	} else {
		result[statusCode] = 1
	}
}

func reportResult(result ResultMap) {
	s := "Load test:"
	for key := range result {
		s = s + fmt.Sprintf(" %d: %d", key, result[key])
	}

	fmt.Println(s)
}

func newWorker() *Worker {
	return &Worker{
		stopCh: make(chan struct{}),
	}
}
