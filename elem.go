package utils

import (
	"reflect"

	"github.com/pkg/errors"
)

// StructValueElem follows the pointer v and returns the element it points to,
// for borth reflect value and reflect type, as well as nil error.
// If v is not a pointer or not a supported pointer, it returns
// zero-value reflect value and reflect type as well as an error.
func StructValueElem(v reflect.Value, t reflect.Type) (rv reflect.Value, rt reflect.Type, err error) {
	switch t.Kind() {
	case reflect.Ptr:
		rv, rt = v.Elem(), t.Elem()
	case reflect.Slice, reflect.Array:
		rt = t.Elem()
		if i := 0; i < v.Len() {
			rv = v.Index(i)
		} else {
			rv = reflect.New(rt).Elem()
		}
	case reflect.Map, reflect.Chan, reflect.Func:
		err = errors.Errorf("'%s' is an unsupported pointer to a struct", t.Kind())
	default:
		err = errors.Errorf("'%s' is not a pointer", t.Kind())
	}
	return rv, rt, err
}
