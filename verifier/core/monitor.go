package core

type Monitor interface {
	Setup(*Verifier) error
	PreCheck(string, *Verifier) error
	Register(string, *Verifier) error
	Run(*Verifier)
}
