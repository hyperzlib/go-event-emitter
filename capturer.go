package eventemitter

import (
	"fmt"
	"reflect"
)

// Capturer is a container struct used to remove the capturer
type Capturer struct {
	handler  reflect.Value
	argTypes []reflect.Type
}

func (c *Capturer) Call(event EventType, args []interface{}) (ret []reflect.Value, err error) {
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
	if argLen < len(c.argTypes) {
		err = fmt.Errorf("missing args, required: %d, got: %d", argLen, len(c.argTypes))
		return
	}

	argValues := make([]reflect.Value, 0)
	argValues = append(argValues, reflect.ValueOf(event))

	for index, argument := range args {
		if index >= argLen {
			break
		}
		if c.argTypes[index] != reflect.TypeOf(argument) {
			err = fmt.Errorf("argument type mismatch at: %d", index)
		}
		argValues = append(argValues, reflect.ValueOf(argument))
	}

	ret = c.handler.Call(argValues)

	return
}

// newCapturerFunc from go-socket.io, get callback info
func newCapturerFunc(f interface{}) *Capturer {
	fv := reflect.ValueOf(f)

	if fv.Kind() != reflect.Func {
		panic("event handler must be a func.")
	}
	ft := fv.Type()

	if ft.NumIn() < 1 || ft.In(0) != reflect.TypeOf("") {
		panic("handler function should be like func(event EventType, ...)")
	}

	argTypes := make([]reflect.Type, ft.NumIn()-1)
	for i := range argTypes {
		argTypes[i] = ft.In(i + 1)
	}

	if len(argTypes) == 0 {
		argTypes = nil
	}

	return &Capturer{
		handler:  fv,
		argTypes: argTypes,
	}
}
