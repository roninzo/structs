// Copyright 2021 Roninzo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package structs

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
)

/*   F u n c t i o n s   */

// Elem follows the pointer v and returns the element it points to,
// for borth reflect value and reflect type, as well as nil error.
// If v is not a pointer or not a supported pointer, it returns
// zero-value reflect value and reflect type as well as an error.
func Elem(v reflect.Value, t reflect.Type) (rv reflect.Value, rt reflect.Type, err error) {
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

// Dump returns a MarshalIndent string.
//
// BUG(roninzo): Dump uses json marshaling which  does not support complex
// types (complex64 and complex128).
func Dump(dest interface{}) string {
	// s := fmt.Sprintf("%#v", t)
	// m := make(map[string]interface{}) // convert dest to m first?
	// http://choly.ca/post/go-json-marshalling/
	// https://www.py4u.net/discuss/1206302
	// https://play.golang.org/p/MuW6gwSAKi
	// https://attilaolah.eu/2013/11/29/json-decoding-in-go/
	// https://mariadesouza.com/2017/09/07/custom-unmarshal-json-in-golang/
	j, err := json.MarshalIndent(dest, " ", "   ")
	if err != nil {
		return err.Error()
	}
	return string(j)
}

// Kinds returns a comma separated string from list of reflect.Kind's.
func Kinds(kinds ...reflect.Kind) string {
	l := make([]string, 0)
	for _, k := range kinds {
		l = append(l, k.String())
	}
	return strings.Join(l, ",")
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

// Interface returns interface reflect value, else returns zero value.
func Interface(v reflect.Value) interface{} {
	return v.Interface()
}

// Error returns the error reflect value, else returns nil.
func Error(v reflect.Value) error {
	i := v.Interface()
	if err, ok := i.(interface{ Error() string }); ok {
		return err
	}
	return nil
}

// Struct returns the StructValue object or panics (returns nil).
func Struct(v reflect.Value) *StructValue {
	i := v.Interface()
	// Is dest already a *StructValue? ...
	if s, ok := i.(*StructValue); ok {
		return s
	}
	//... else, create it!
	s, err := New(i)
	if err != nil {
		panic(err) // return nil
	}
	return s
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

// Zero returns the zero-value of reflect value v.
func Zero(v reflect.Value) reflect.Value {
	return reflect.Zero(v.Type())
}

// PtrString, PtrBool, etc. returns a pointer to type value x.
func PtrString(x string) *string                 { return &x }
func PtrBool(x bool) *bool                       { return &x }
func PtrInt(x int) *int                          { return &x }
func PtrInt8(x int8) *int8                       { return &x }
func PtrInt16(x int16) *int16                    { return &x }
func PtrInt32(x int32) *int32                    { return &x }
func PtrInt64(x int64) *int64                    { return &x }
func PtrUint(x uint) *uint                       { return &x }
func PtrUint8(x uint8) *uint8                    { return &x }
func PtrUint16(x uint16) *uint16                 { return &x }
func PtrUint32(x uint32) *uint32                 { return &x }
func PtrUint63(x uint64) *uint64                 { return &x }
func PtrFloat32(x float32) *float32              { return &x }
func PtrFloat64(x float64) *float64              { return &x }
func PtrComplex64(x complex64) *complex64        { return &x }
func PtrComplex128(x complex128) *complex128     { return &x }
func PtrBytes(x []byte) *[]byte                  { return &x }
func PtrTime(x time.Time) *time.Time             { return &x }
func PtrDuration(x time.Duration) *time.Duration { return &x }
func PtrError(x error) *error                    { return &x }

/*   U n e x p o r t e d   */

// canString returns true if reflect value is of type string.
func canString(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return true
	}
	return false
}

// canBool returns true if reflect value is of type bool.
func canBool(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return true
	}
	return false
}

// canInt returns true if reflect value is of type int.
func canInt(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	}
	return false
}

// canUint returns true if reflect value is of type uint.
func canUint(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	}
	return false
}

// canFloat returns true if reflect value is of type float.
func canFloat(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

// canComplex returns true if reflect value is of type float.
func canComplex(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Complex64, reflect.Complex128:
		return true
	}
	return false
}

// canBytes returns true if reflect value is of type []byte.
//
// NOTE: []byte == []uint8
func canBytes(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return true
		}
	}
	return false
}

// canTime returns true if reflect value is of type time.Time, else returns false.
func canTime(v reflect.Value) bool {
	if v.Type().String() == "time.Time" {
		return true
	}
	_, ok := v.Interface().(time.Time)
	return ok
}

// canDuration returns true if reflect value is of type time.Duration, else
// returns false.
// NOTE: time.Duration <=> int64
func canDuration(v reflect.Value) bool {
	if v.Type().String() == "time.Duration" {
		return true
	}
	_, ok := v.Interface().(time.Duration)
	return ok
}

// canError returns true if reflect value implements the error interface.
func canError(v reflect.Value) bool {
	if v.Type().String() == "error" {
		return true
	}
	_, ok := v.Interface().(interface{ Error() string })
	return ok
}

// canSlice reports whether reflect Slice can be used on reflect value
// without panicking.
func canSlice(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return true
	}
	return false
}

// canMap reports whether reflect Map can be used on reflect value
// without panicking.
func canMap(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Map:
		return true
	}
	return false
}

// canInterface reports whether reflect Interface can be used on reflect value
// without panicking.
func canInterface(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Interface:
		return true
	}
	return false
}

// canStruct returns true if reflect value represents a nested struct,
// else returns false.
func canStruct(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Struct:
		if !canTime(v) {
			return true
		}
	case reflect.Ptr:
		v, _, err := Elem(v, v.Type())
		if err != nil {
			return false
		}
		if v.Kind() == reflect.Struct {
			if !canTime(v) {
				return true
			}
		}
	}
	return false
}

// canNil returns true if reflect value represents a type that can be set
// to nil, else returns false.
func canNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return true
	}
	return false
}

// canPtr returns true if reflect value represents a pointer,
// else returns false.
func canPtr(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.UnsafePointer:
		return true
	}
	return false
}
