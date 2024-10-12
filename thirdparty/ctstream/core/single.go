package core

import (
	"context"
	"sync"
	"time"
)

type CTStream[T CtClient] struct {
	Client T
	Sleep  time.Duration
	Ctx    context.Context
	Wg     sync.WaitGroup
	stop   context.CancelFunc
}

func NewCTStream[T CtClient](
	client T,
	sleep time.Duration,
	Ctx context.Context,
) (*CTStream[T], error) {
	ctx, cancel := context.WithCancel(Ctx)
	return &CTStream[T]{
		Client: client,
		Sleep:  sleep,
		Ctx:    ctx,
		stop:   cancel,
		Wg:     sync.WaitGroup{},
	}, nil
}

func (stream *CTStream[T]) Init() error {
	return stream.Client.Init()
}

func (stream *CTStream[T]) Next(callback Callback) {
	stream.Client.Next(callback)
}

func (stream *CTStream[T]) start(callback Callback) {
	for {
		select {
		case <-stream.Ctx.Done():
			stream.Wg.Done()
			return
		default:
			stream.Next(callback)
			time.Sleep(stream.Sleep)
		}
	}
}

func (stream *CTStream[T]) Start(callback Callback) {
	stream.Wg.Add(1)
	go stream.start(callback)
}

func (stream *CTStream[T]) Run(callback Callback) {
	stream.Start(callback)
	stream.Await()
}

func (stream *CTStream[T]) Await() {
	stream.Wg.Wait()
}

func (stream *CTStream[T]) Stop() {
	stream.stop()
}
