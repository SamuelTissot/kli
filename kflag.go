package kli

import (
	"reflect"
	"time"
)

type KFlag interface {
	// Store returns a list of set flags with their reflect.Kind
	Store() map[string]reflect.Kind

	// Set sets a new flag
	SetFlag(name string, ptr interface{})

	// BoolFlag return the value of the flag "name"
	// ok is false if the flag does not exist or of wrong type
	BoolFlag(name string) (value, ok bool)

	// DurationFlag return the value of the flag "name"
	// ok is false if the flag does not exist or of wrong type
	DurationFlag(name string) (value time.Duration, ok bool)

	// Float64Flag return the value of the flag "name"
	// ok is false if the flag does not exist or of wrong type
	Float64Flag(name string) (value float64, ok bool)

	// IntFlag return the value of the flag "name"
	// ok is false if the flag does not exist or of wrong type
	IntFlag(name string) (value int, ok bool)

	// Int64Flag return the value of the flag "name"
	// ok is false if the flag does not exist or of wrong type
	Int64Flag(name string) (value int64, ok bool)

	// StringFlag return the value of the flag "name"
	// ok is false if the flag does not exist or of wrong type
	StringFlag(name string) (value string, ok bool)

	// UintFlag return the value of the flag "name"
	// ok is false if the flag does not exist or of wrong type
	UintFlag(name string) (value uint, ok bool)

	// Uint64Flag return the value of the flag "name"
	// ok is false if the flag does not exist or of wrong type
	Uint64Flag(name string) (value uint64, ok bool)
}

type FlagStore struct {
	f map[string]interface{}
}

func NewKflag() *FlagStore {
	return &FlagStore{map[string]interface{}{}}
}

func (a *FlagStore) Store() map[string]reflect.Kind {
	result := make(map[string]reflect.Kind)
	for name, f := range a.f {
		result[name] = reflect.TypeOf(f).Elem().Kind()
	}
	return result
}

func (a *FlagStore) SetFlag(name string, ptr interface{}) {
	a.f[name] = ptr
}

// flagElem returns the FlagStore for the given name
func (a *FlagStore) flagElem(name string) reflect.Value {
	if f, ok := a.f[name]; ok {
		e := reflect.ValueOf(f).Elem()
		return e
	}
	return reflect.Value{}
}

func (a *FlagStore) BoolFlag(name string) (value, ok bool) {
	f := a.flagElem(name)
	if f.Kind() != reflect.Bool {
		return
	}
	return f.Bool(), true
}

func (a *FlagStore) DurationFlag(name string) (value time.Duration, ok bool) {
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

func (a *FlagStore) Float64Flag(name string) (value float64, ok bool) {
	f := a.flagElem(name)
	if f.Kind() != reflect.Float64 {
		return
	}
	return f.Float(), true
}

func (a *FlagStore) IntFlag(name string) (value int, ok bool) {
	f := a.flagElem(name)
	if f.Kind() != reflect.Int {
		return
	}
	return int(f.Int()), true
}

func (a *FlagStore) Int64Flag(name string) (value int64, ok bool) {
	f := a.flagElem(name)
	if f.Kind() != reflect.Int64 {
		return
	}
	return f.Int(), true
}

func (a *FlagStore) StringFlag(name string) (value string, ok bool) {
	f := a.flagElem(name)
	if f.Kind() != reflect.String {
		return
	}

	return f.String(), true
}

func (a *FlagStore) UintFlag(name string) (value uint, ok bool) {
	f := a.flagElem(name)
	if f.Kind() != reflect.Uint {
		return
	}
	return uint(f.Uint()), true
}

func (a *FlagStore) Uint64Flag(name string) (value uint64, ok bool) {
	f := a.flagElem(name)
	if f.Kind() != reflect.Uint64 {
		return
	}
	return f.Uint(), true
}
