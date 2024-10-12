package core

type CtStream interface {
	Init() error
	Start(Callback)
	Next(callback Callback)
	Run(Callback)
	Stop()
	Await()
}

type CTsStream[T CtStream] PararellCTsStream[T]
