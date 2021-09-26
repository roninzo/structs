package utils

import (
	"reflect"
	"strings"
	"time"
)

// CanString returns true if reflect value is of type string.
func CanString(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return true
	}
	return false
}

// CanBool returns true if reflect value is of type bool.
func CanBool(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return true
	}
	return false
}

// CanInt returns true if reflect value is of type int.
func CanInt(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	}
	return false
}

// CanUint returns true if reflect value is of type uint.
func CanUint(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	}
	return false
}

// CanFloat returns true if reflect value is of type float.
func CanFloat(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

// CanComplex returns true if reflect value is of type float.
func CanComplex(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Complex64, reflect.Complex128:
		return true
	}
	return false
}

// CanBytes returns true if reflect value is of type []byte.
//
// NOTE: []byte == []uint8
func CanBytes(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return true
		}
	}
	return false
}

// CanTime returns true if reflect value is of type time.Time, else returns false.
func CanTime(v reflect.Value) bool {
	if v.Type().String() == "time.Time" {
		return true
	}
	_, ok := v.Interface().(time.Time)
	return ok
}

// CanDuration returns true if reflect value is of type time.Duration, else
// returns false.
//
// NOTE: time.Duration <=> int64
func CanDuration(v reflect.Value) bool {
	if v.Type().String() == "time.Duration" {
		return true
	}
	_, ok := v.Interface().(time.Duration)
	return ok
}

// CanError returns true if reflect value implements the error interface.
func CanError(v reflect.Value) bool {
	if v.Type().String() == "error" {
		return true
	}
	_, ok := v.Interface().(interface{ Error() string })
	return ok
}

// CanSlice reports whether reflect Slice can be used on reflect value
// without panicking.
func CanSlice(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return true
	}
	return false
}

// CanMap reports whether reflect Map can be used on reflect value
// without panicking.
func CanMap(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Map:
		return true
	}
	return false
}

// CanInterface reports whether reflect Interface can be used on reflect value
// without panicking.
func CanInterface(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Interface:
		return true
	}
	return false
}

// CanStruct returns true if reflect value represents a nested struct,
// else returns false.
func CanStruct(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Struct:
		if !CanTime(v) {
			return true
		}
	case reflect.Ptr:
		v, _, err := StructValueElem(v, v.Type())
		if err != nil {
			return false
		}
		if v.Kind() == reflect.Struct {
			if !CanTime(v) {
				return true
			}
		}
	}
	return false
}

// CanNil returns true if reflect value represents a type that can be set
// to nil, else returns false.
func CanNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return true
	}
	return false
}

// CanPtr returns true if reflect value represents a pointer,
// else returns false.
func CanPtr(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr: // , reflect.UnsafePointer:
		return true
	}
	return false
}

// Zero returns the zero-value of reflect value v.
func Zero(v reflect.Value) reflect.Value {
	return reflect.Zero(v.Type())
}

// Preset returns concrete value of v in a way that it is ready to be set.
// If v is a pointer to another type, then it returns the concrete value that the pointer
// points to. In the case that the pointer is nil (not pointing to anything yet), Preset
// first has to create zero-value of the type and change the value of the pointer from nil
// to the address of that zero-value. Then it can returns the concrete value that the pointer
// points to as well.
//
// Makes sure no pointers points to something concrete, instead of nil.
func Preset(v reflect.Value) reflect.Value {
	t := v.Type()
	if v.Kind() == reflect.Ptr {
		e := v.Elem()
		t = t.Elem()
		if !e.IsValid() { // Check if the pointer is nil
			v.Set(reflect.New(t))
		}
		return v.Elem()
	}
	return v
}

// Time returns time.Time reflect value, else returns zero value.
func Time(v reflect.Value) (t time.Time) {
	i := v.Interface()
	t, _ = i.(time.Time)
	return t
}

// Duration returns time.Duration reflect value, else returns zero value.
func Duration(v reflect.Value) (d time.Duration) {
	i := v.Interface()
	d, _ = i.(time.Duration)
	return d
}

// Error returns the error reflect value, else returns nil.
func Error(v reflect.Value) error {
	i := v.Interface()
	if err, ok := i.(interface{ Error() string }); ok {
		return err
	}
	return nil
}

// Kinds returns a comma separated string from list of reflect.Kind's.
func Kinds(kinds ...reflect.Kind) string {
	l := make([]string, 0)
	for _, k := range kinds {
		l = append(l, k.String())
	}
	return strings.Join(l, ",")
}
