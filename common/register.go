package common

// IRegister common module register interface
type IRegister interface {
	Name() string
	CheckIn() error
}
