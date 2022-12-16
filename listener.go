package eventemitter

import (
	"fmt"
	"reflect"
)

// Listener is a container struct used to remove the listener
type Listener struct {
	handler  reflect.Value
	argTypes []reflect.Type
}

func (l *Listener) Call(args []interface{}) (ret []reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("event call error: %s", r)
			}
		}
	}()

	argLen := len(args)
	if argLen < len(l.argTypes) {
		err = fmt.Errorf("missing args, required: %d, got: %d", argLen, len(l.argTypes))
		return
	}

	argValues := make([]reflect.Value, 0)

	for index, argument := range args {
		if index >= argLen {
			break
		}
		if l.argTypes[index] != reflect.TypeOf(argument) {
			err = fmt.Errorf("argument type mismatch at: %d", index)
		}
		argValues = append(argValues, reflect.ValueOf(argument))
	}

	ret = l.handler.Call(argValues)

	return
}

// newEventFunc from go-socket.io, get callback info
func newEventFunc(f interface{}) *Listener {
	fv := reflect.ValueOf(f)

	if fv.Kind() != reflect.Func {
		panic("event handler must be a func.")
	}
	ft := fv.Type()

	argTypes := make([]reflect.Type, ft.NumIn())
	for i := range argTypes {
		argTypes[i] = ft.In(i)
	}

	if len(argTypes) == 0 {
		argTypes = nil
	}

	return &Listener{
		handler:  fv,
		argTypes: argTypes,
	}
}
