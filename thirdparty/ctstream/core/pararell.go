package core

import (
	"time"
)

type PararellCTsStream[T CtStream] struct {
	Streams []T
	Sleep   time.Duration
}

func NewPararellCTsStream[T CtStream](streams []T, sleep time.Duration) (*PararellCTsStream[T], error) {
	return &PararellCTsStream[T]{
		Streams: streams,
		Sleep:   sleep,
	}, nil
}

func (stream *PararellCTsStream[T]) Init() error {
	for _, s := range stream.Streams {
		err := s.Init()
		if err != nil {
			return err
		}
	}

	return nil
}

func (stream *PararellCTsStream[T]) Start(callback Callback) {
	for _, s := range stream.Streams {
		go s.Run(callback)
	}
}

func (stream *PararellCTsStream[T]) Next(callback Callback) {
}

func (stream PararellCTsStream[T]) Await() {
	for _, s := range stream.Streams {
		time.Sleep(stream.Sleep)
		s.Await()
	}
}

func (stream *PararellCTsStream[T]) Run(callback Callback) {
	stream.Start(callback)
	stream.Await()
}

func (stream *PararellCTsStream[T]) Stop() {
	for _, s := range stream.Streams {
		s.Stop()
	}
}
