package core

type Monitor interface {
	Setup(*Verifier) error
	PreCheck(string, bool, *Verifier) error
	Register(string, bool, *Verifier) error
	Run(*Verifier)
}
