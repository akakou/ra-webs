package core

type Monitor interface {
	Setup(*Verifier) error
	Register(string, *Verifier) error
	Run(*Verifier)
}
