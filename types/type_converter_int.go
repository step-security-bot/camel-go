package types

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

// ==========================
//
// Int converter
//
// ==========================

// ToInt --
type ToInt interface {
	ToInt() (int, error)
}

// ToInt8 --
type ToInt8 interface {
	ToInt8() (int8, error)
}

// ToInt16 --
type ToInt16 interface {
	ToInt16() (int16, error)
}

// ToInt32 --
type ToInt32 interface {
	ToInt32() (int32, error)
}

// ToInt64 --
type ToInt64 interface {
	ToInt64() (int64, error)
}

// ==========================
//
//
//
// ==========================

// ToInt ..
func (target Integer) ToInt() int {
	return int(target)
}

// ToInt8 ..
func (target Integer) ToInt8() int8 {
	return int8(target)
}

// ToInt16 ..
func (target Integer) ToInt16() int16 {
	return int16(target)
}

// ToInt32 ..
func (target Integer) ToInt32() int32 {
	return int32(target)
}

// ToInt64 ..
func (target Integer) ToInt64() int64 {
	return int64(target)
}

// ==========================
//
// ToIntConverter
//
// ==========================

// ToIntConverter --
func ToIntConverter(source interface{}, targetType reflect.Type) (interface{}, error) {
	if !IsInt(targetType) {
		return nil, errors.New("unsupported")
	}

	var answer interface{}
	var err error

	sourceType := reflect.TypeOf(source)
	sourceKind := sourceType.Kind()

	switch targetType.Kind() {
	case reflect.Int:
		if sourceKind == reflect.Struct {
			if v, ok := source.(ToInt); ok {
				answer, err = v.ToInt()
			} else {
				err = fmt.Errorf("unable to convert struct:%T to:%v", source, targetType)
			}
		} else {
			answer, err = cast.ToIntE(source)
		}
	case reflect.Int8:
		if sourceKind == reflect.Struct {
			if v, ok := source.(ToInt8); ok {
				answer, err = v.ToInt8()
			} else {
				err = fmt.Errorf("unable to convert struct:%T to:%v", source, targetType)
			}
		} else {
			answer, err = cast.ToInt8E(source)
		}
	case reflect.Int16:
		if sourceKind == reflect.Struct {
			if v, ok := source.(ToInt16); ok {
				answer, err = v.ToInt16()
			} else {
				err = fmt.Errorf("unable to convert struct:%T to:%v", source, targetType)
			}
		} else {
			answer, err = cast.ToInt16E(source)
		}
	case reflect.Int32:
		if sourceKind == reflect.Struct {
			if v, ok := source.(ToInt32); ok {
				answer, err = v.ToInt32()
			} else {
				err = fmt.Errorf("unable to convert struct:%T to:%v", source, targetType)
			}
		} else {
			answer, err = cast.ToInt32E(source)
		}
	case reflect.Int64:
		if sourceKind == reflect.Struct {
			if v, ok := source.(ToInt64); ok {
				answer, err = v.ToInt64()
			} else {
				err = fmt.Errorf("unable to convert struct:%T to:%v", source, targetType)
			}
		} else {
			answer, err = cast.ToInt64E(source)
		}
	}

	return answer, err
}