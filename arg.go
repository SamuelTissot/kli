package kli

import (
	"reflect"
	"time"
)

type KFlag struct {
	f map[string]interface{}
}

func NewArg() *KFlag {
	return &KFlag{map[string]interface{}{}}
}

// Flags returns a list of flags
func (a *KFlag) Flags() map[string]reflect.Kind {
	result := make(map[string]reflect.Kind)
	for name, f := range a.f {
		result[name] = reflect.TypeOf(f).Elem().Kind()
	}
	return result
}

// flagElem returns the KFlag for the given name
func (a *KFlag) flagElem(name string) reflect.Value {
	if f, ok := a.f[name]; ok {
		e := reflect.ValueOf(f).Elem()
		return e
	}
	return reflect.Value{}
}

// BoolFlag return the value of the flag "name"
// ok is false if the flag does not exist or of wrong type
func (a *KFlag) BoolFlag(name string) (value, ok bool) {
	f := a.flagElem(name)
	if f.Kind() != reflect.Bool {
		return
	}
	return f.Bool(), true
}

// DurationFlag return the value of the flag "name"
// ok is false if the flag does not exist or of wrong type
func (a *KFlag) DurationFlag(name string) (value time.Duration, ok bool) {
	f := a.flagElem(name)
	if f.Kind() != reflect.String {
		return
	}

	d, e := time.ParseDuration(f.String())
	if e != nil {
		return
	}

	return d, true
}

// Float64Flag return the value of the flag "name"
// ok is false if the flag does not exist or of wrong type
func (a *KFlag) Float64Flag(name string) (value float64, ok bool) {
	f := a.flagElem(name)
	if f.Kind() != reflect.Float64 {
		return
	}
	return f.Float(), true
}

// IntFlag return the value of the flag "name"
// ok is false if the flag does not exist or of wrong type
func (a *KFlag) IntFlag(name string) (value int, ok bool) {
	f := a.flagElem(name)
	if f.Kind() != reflect.Int {
		return
	}
	return int(f.Int()), true
}

// Int64Flag return the value of the flag "name"
// ok is false if the flag does not exist or of wrong type
func (a *KFlag) Int64Flag(name string) (value int64, ok bool) {
	f := a.flagElem(name)
	if f.Kind() != reflect.Int64 {
		return
	}
	return f.Int(), true
}

// StringFlag return the value of the flag "name"
// ok is false if the flag does not exist or of wrong type
func (a *KFlag) StringFlag(name string) (value string, ok bool) {
	f := a.flagElem(name)
	if f.Kind() != reflect.String {
		return
	}

	return f.String(), true
}

// UintFlag return the value of the flag "name"
// ok is false if the flag does not exist or of wrong type
func (a *KFlag) UintFlag(name string) (value uint, ok bool) {
	f := a.flagElem(name)
	if f.Kind() != reflect.Uint {
		return
	}
	return uint(f.Uint()), true
}

// Uint64Flag return the value of the flag "name"
// ok is false if the flag does not exist or of wrong type
func (a *KFlag) Uint64Flag(name string) (value uint64, ok bool) {
	f := a.flagElem(name)
	if f.Kind() != reflect.Uint64 {
		return
	}
	return f.Uint(), true
}
