package core

type Monitor interface {
	Setup(*Verifier) error
	Run(*Verifier)
}
