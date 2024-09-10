package core

type Monitor interface {
	Setup(*Verifier) error
	Register(*Verifier) error
	Run(*Verifier)
}
