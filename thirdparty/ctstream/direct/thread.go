package direct

import (
	"fmt"
	"sync"
	"time"
)

var FetchingThread = 0
var MaxThread = 50

type threadManager struct {
	FetchingThread int
	MaxThread      int
	Mutex          sync.Mutex
}

var ThreadManager = threadManager{
	FetchingThread: FetchingThread,
	MaxThread:      MaxThread,
	Mutex:          sync.Mutex{},
}

type Args []any

func (tm *threadManager) Run(f func(Args), param Args) {
	tm.SleepToAdjust()
	go tm.run(f, param)
}

func (tm *threadManager) run(f func(Args), param Args) {
	tm.Mutex.Lock()
	tm.FetchingThread++
	tm.Mutex.Unlock()

	f(param)

	tm.Mutex.Lock()
	tm.FetchingThread--
	tm.Mutex.Unlock()
}

func (tm *threadManager) SleepToAdjust() {
	for FetchingThread > MaxThread {
		fmt.Printf("Sleeping to adjust THREAD(%v) > MAX(%v)\n", FetchingThread, MaxThread)
		time.Sleep(DefaultSleep)
	}
}
