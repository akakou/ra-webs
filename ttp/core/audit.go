package core

type Audit interface {
	Setup(*TTP) error
	Run(*TTP)
}
