package common

// IRegister common module regist interface
type IRegister interface {
	Name() string
	Regist() error
}
