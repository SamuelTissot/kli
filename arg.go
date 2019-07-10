package kli

import (
	"fmt"
	"reflect"
	"time"
)

type Arg struct {
	Kind reflect.Kind
	v    interface{}
}

func (a *Arg) ofKind(t reflect.Kind) error {
	if t != a.Kind {
		return fmt.Errorf("arg trying to get value as type %T but type is %T", t, a.Kind)
	}
	return nil
}

func (a *Arg) elem() reflect.Value {
	return reflect.ValueOf(a.v).Elem()
}

// Bool returns the boolean value of the Arg
// and error is return if trying to get a value of the
// wrong Type
func (a *Arg) Bool() (bool, error) {
	if err := a.ofKind(reflect.Bool); err != nil {
		return false, err
	}

	return a.elem().Bool(), nil
}

// Duration returns the duration value of Arg
func (a *Arg) Duration(unit time.Duration) (time.Duration, error) {
	if err := a.ofKind(reflect.Int64); err != nil {
		return 0, fmt.Errorf("arg trying to get value as type time.Duration but type is %T", a.Kind)
	}

	return time.Duration(a.elem().Int()), nil
}

// Float64 returns the fload64 value of Arg
func (a *Arg) Float64() (float64, error) {
	if err := a.ofKind(reflect.Float64); err != nil {
		return 0, err
	}

	return a.elem().Float(), nil
}

// Int returns the int value of Arg
func (a *Arg) Int() (int, error) {
	if err := a.ofKind(reflect.Int); err != nil {
		return 0, err
	}

	return int(a.elem().Int()), nil
}

// Int64 returns the int64 value of Arg
func (a *Arg) Int64() (int64, error) {
	if err := a.ofKind(reflect.Int64); err != nil {
		return 0, err
	}

	return a.elem().Int(), nil
}

// String returns the string value of Arg
func (a *Arg) String() (string, error) {
	if err := a.ofKind(reflect.String); err != nil {
		return "", err
	}

	return a.elem().String(), nil
}

// Uint returns the uint value of Arg
func (a *Arg) Uint() (uint, error) {
	if err := a.ofKind(reflect.Uint); err != nil {
		return 0, err
	}

	return uint(a.elem().Uint()), nil
}

// Uint64 returns the uint64 value of Arg
func (a *Arg) Uint64() (uint64, error) {
	if err := a.ofKind(reflect.Uint64); err != nil {
		return 0, err
	}

	return a.elem().Uint(), nil
}
