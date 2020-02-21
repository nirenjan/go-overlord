package module

import (
	"errors"
	"fmt"
)

type Callback int

const (
	BuildCommandTree Callback = iota
	ModuleInit
	Backup
	Restore
	maxCallbacks
)

type Module struct {
	// Module name, this is used in error messages to narrow down issues
	Name string

	// A list of callbacks that are supported by the module
	Callbacks [maxCallbacks]func() error

	// A list of callbacks that accept byte arrays as input and return
	// byte arrays as output
	DataCallbacks [maxCallbacks]func([]byte) ([]byte, error)
}

type CallbackIterator struct {
	// Module name
	Name string

	// Callback
	Callback func([]byte) ([]byte, error)
}

var modules []Module

func RegisterModule(mod Module) {
	modules = append(modules, mod)
}

func RunCallback(cb Callback) error {
	for _, mod := range modules {
		f := mod.Callbacks[cb]
		if f != nil {
			err := f()
			if err != nil {
				errmsg := fmt.Sprintf("Callback %v:%v failed with error %v",
					mod.Name, cb, err)
				return errors.New(errmsg)
			}
		}
	}

	return nil
}

// IterateCallback will return a list of modules that have the callback
// function set
func IterateCallback(cb Callback) []CallbackIterator {
	list := make([]CallbackIterator, len(modules))
	i := 0
	for _, mod := range modules {
		f := mod.DataCallbacks[cb]
		if f != nil {
			list[i].Name = mod.Name
			list[i].Callback = mod.DataCallbacks[cb]
			i++
		}
	}

	list = list[:i]
	return list
}
