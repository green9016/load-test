package main

import (
	"os"
	"os/signal"
	"sync"
)

type ResultMap map[int]int

type Workers struct {
	systemCh   chan os.Signal
	resultCh   chan int
	resultMu   sync.RWMutex
	requesters []*Worker
	writer     *Worker
	reporter   *Worker
	result     ResultMap
	count      int
}

func (workers *Workers) stop() {
	workers.systemCh <- os.Kill
}

func (workers *Workers) start() {
	var wg sync.WaitGroup
	wg.Add(workers.count + 2)

	// running a writer
	go func() {
		workers.writer.runWriter(workers.result, workers.resultCh, &workers.resultMu)
		wg.Done()
	}()

	// running N requesters
	for i := 0; i < workers.count; i++ {
		go func(requester *Worker) {
			requester.runRequester(workers.resultCh)
			wg.Done()
		}(workers.requesters[i])
	}

	// running a reporter
	go func() {
		workers.reporter.runReporter(workers.result, &workers.resultMu)
		wg.Done()
	}()

	wg.Wait()
}

func newWorkers(count int) *Workers {
	systemCh := make(chan os.Signal, 1)
	signal.Notify(systemCh, os.Interrupt)

	resultCh := make(chan int, count*10)

	workers := &Workers{}
	workers.systemCh = systemCh
	workers.resultCh = resultCh
	workers.count = count
	workers.result = make(ResultMap)

	for i := 0; i < count; i++ {
		workers.requesters = append(workers.requesters, newWorker())
	}

	workers.writer = newWorker()
	workers.reporter = newWorker()

	go func() {
		<-systemCh

		for i := 0; i < count; i++ {
			workers.requesters[i].stop()
		}

		workers.writer.stop()
		workers.reporter.stop()
	}()

	return workers
}
