// Copyright 2021 Roninzo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package structs

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/pkg/errors"
)

/*   S t r u c t   d e f i n i t i o n   */

// StructField represents a single struct field that encapsulates high level
// functions around the field.
type StructField struct {
	index   int                 // sequence index of field.
	indexes []int               // absolute indexes of field inside struct.
	value   reflect.Value       // struct field value.
	field   reflect.StructField // struct field definition.
	Parent  *StructValue        // field's own struct reference.
}

/*   C o n s t r u c t o r   */

// Field returns nil or one of the fields of the struct that matches argument dest.
// Its argument dest can be the name or the index of the field.
// Field(nil) returns nil and adds an error to StructValue.
//
// NOTE: Field is an alias to either the getFieldByName or the getFieldByIndex method.
//
// TODO: getFieldByIndex ability to access nested field names using []int as index.
//   e.g.: [1, 3, 1] <=> v.Field(1).Field(3).Field(1)
//
// TODO: getFieldByName ability to parse nested field names inside n.
//   e.g.: "Struct.Nested.String" <=> v.FieldByName("Struct").FieldByName("Nested").FieldByName("String")
func (s *StructValue) Field(dest interface{}) *StructField {
	if dest == nil {
		s.setError("invalid nil argument")
		return nil
	}
	switch arg := dest.(type) {
	// case []int:
	// 	return s.getFieldByIndexes(arg)
	// case []string:
	// 	return s.getFieldByNames(arg)
	case int:
		return s.getFieldByIndex(arg)
	case string:
		return s.getFieldByName(arg)
	}
	t := reflect.TypeOf(dest)
	s.setErrorf("invalid argument type; want: 'string' or 'int', got: '%s'", t.Kind())
	return nil
}

/*   I m p l e m e n t a t i o n   */

// IsValid returns true if StructField has been loaded successfully.
// Useful for checking StructField is valid before use.
func (f *StructField) IsValid() bool {
	if f.Parent != nil && f.index != OutOfRange {
		return f.value.IsValid()
	}
	return false
}

// Index returns the struct index of the given field.
// If field is not valid, Index returns the OutOfRange constant, i.e. -1.
func (f *StructField) Index() int {
	return f.index
}

// Name returns returns the name of StructField, unless it was invalid.
// In which case, Name returns zero-value string.
func (f *StructField) Name() string {
	return f.field.Name
}

// Namespace is similar to the Name method, except that it includes its related struct names
// all the way to the top level struct (in a dot separated string).
func (f *StructField) Namespace() (n string) {
	s := f.Parent
	n = s.Name()
	for {
		if s.IsNested() {
			s = s.Parent
			n = fmt.Sprintf("%s.%s", s.Name(), n)
		} else {
			break
		}
	}
	if n == "" {
		return f.Name()
	}
	n = fmt.Sprintf("%s.%s", n, f.Name())
	return n
}

// Type returns the underlying type of the field.
func (f *StructField) Type() reflect.Type {
	return f.value.Type()
}

// Kind returns the fields kind, such as "string", "int", "bool", etc ..
func (f *StructField) Kind() reflect.Kind {
	return f.value.Type().Kind()
}

// Tag returns the value associated with key in the tag string.
// If the key is present in the tag the value (which may be empty) is returned.
// Otherwise the returned value will be the empty string. The ok return value
// reports whether the value was explicitly set in the tag string.
func (f *StructField) Tag(key string) (string, bool) {
	return f.field.Tag.Lookup(key)
}

// IsAnonymous returns true if the given field is an anonymous field, meaning a field
// having no name. This obviously related to the use of the Name method.
func (f *StructField) IsAnonymous() bool {
	return f.field.Anonymous
}

// IsEmbedded is a alias to the IsAnonymous method.
// An embedded field can be an anonymous nested struct field.
func (f *StructField) IsEmbedded() bool {
	return f.IsAnonymous()
}

// Interface returns true if underlying value of the field is modifiable.
func (f *StructField) CanSet() bool {
	return f.value.CanSet() // Unexported struct fields will be neglected.
}

// ... NOTE: Review/improve doc.
// IsExported returns true if the given field is exported and its json tag is
// not equal to "-". Those fields are neglected for getter and setter methods.
func (f *StructField) IsExported() bool {
	if ok := f.field.PkgPath == ""; ok {
		return true // if f.IsHidden() { return false }
	}
	return false
}

// ... NOTE: Review/improve doc.
// IsHidden returns true if the given field is exported and its json tag is
// not equal to "-". Those fields are neglected for getter and setter methods.
func (f *StructField) IsHidden() bool {
	if !f.IsExported() {
		return true
	}
	if val, ok := f.Tag("json"); ok {
		if val == "-" {
			return true
		}
		if strings.Contains(val, "omitempty") {
			if f.IsZero() {
				return true
			}
		}
	}
	return false
}

// Zero returns field's type specific zero value. For instance, the zero-value
// of a string field is "", of an int is 0, and so on.
func (f *StructField) Zero() reflect.Value {
	v := f.value
	return Zero(v)
}

// IsZero returns true if the given field is a zero-value, i.e. not initialized.
// Unexported struct fields will be neglected.
func (f *StructField) IsZero() bool {
	if f.IsExported() {
		return reflect.DeepEqual(f.Interface(), f.Zero().Interface()) // v := f.value; z := Zero(v); return v == z
	}
	return false
}

// IsNil reports whether its argument f is nil. The argument must be a chan, func,
// interface, map, pointer, or slice value; if it is not, IsNil returns nil.
// Unexported struct fields will be neglected.
func (f *StructField) IsNil() bool {
	v := f.value
	if canNil(v) {
		return v.IsNil()
	}
	return false
}

// ... NOTE: Review/improve doc.
// Equal compares field value with reflect value argument and returns true
func (f *StructField) Equal(x *StructField) bool {
	if !f.IsExported() {
		return false
	}
	if x == nil {
		return false
	}
	return f.equal(x.value) != OutOfRange
}

// G e t t e r s
// Getter methods return f's underlying value for a specifc type T. It panics if f's kind is not T.

// Get returns the value of the field as interface.
// Unexported struct fields will be neglected.
func (f *StructField) Get() (interface{}, error) {
	v := f.Value()
	err := f.Parent.Err()
	if err != nil {
		return nil, err
	}
	return v.Interface(), nil
}

func (f *StructField) Time() time.Time         { v := f.value; return Time(v) }
func (f *StructField) Duration() time.Duration { v := f.value; return Duration(v) }
func (f *StructField) Error() error            { v := f.value; return Error(v) }
func (f *StructField) String() string          { v := f.value; return v.String() }
func (f *StructField) Bool() bool              { v := f.value; return v.Bool() }
func (f *StructField) Int() int64              { v := f.value; return v.Int() }
func (f *StructField) Uint() uint64            { v := f.value; return v.Uint() }
func (f *StructField) Float() float64          { v := f.value; return v.Float() }
func (f *StructField) Complex() complex128     { v := f.value; return v.Complex() }
func (f *StructField) Bytes() []byte           { v := f.value; return v.Bytes() }
func (f *StructField) Interface() interface{}  { v := f.value; return v.Interface() }

// Struct returns nested struct from field or nil if f is not a nested struct.
func (f *StructField) Struct() *StructValue {
	v := f.value
	s, err := New(v.Interface(), f.Parent)
	if err != nil {
		f.Parent.setErr(err)
		return nil
	}
	f.Parent.Error = nil
	return s
}

// C h e c k e r s
// Checker methods report whether type requested can be used without panicking.

func (f *StructField) CanNil() bool       { v := f.value; return canNil(v) }
func (f *StructField) CanPtr() bool       { v := f.value; return canPtr(v) }
func (f *StructField) CanTime() bool      { v := f.value; return canTime(v) }
func (f *StructField) CanDuration() bool  { v := f.value; return canDuration(v) }
func (f *StructField) CanError() bool     { v := f.value; return canError(v) }
func (f *StructField) CanString() bool    { v := f.value; return canString(v) }
func (f *StructField) CanBool() bool      { v := f.value; return canBool(v) }
func (f *StructField) CanInt() bool       { v := f.value; return canInt(v) }
func (f *StructField) CanUint() bool      { v := f.value; return canUint(v) }
func (f *StructField) CanFloat() bool     { v := f.value; return canFloat(v) }
func (f *StructField) CanComplex() bool   { v := f.value; return canComplex(v) }
func (f *StructField) CanBytes() bool     { v := f.value; return canBytes(v) }
func (f *StructField) CanInterface() bool { v := f.value; return canInterface(v) }
func (f *StructField) CanStruct() bool    { v := f.value; return canStruct(v) }

// S e t t e r s
// Setter methods assigns x to the field f. no assignment is carried out if CanSet
// returns false. As in Go, x's value must be assignable to f's type.

// SetZero sets the field to its zero value.
// Unsettable struct fields will return an error.
func (f *StructField) SetZero() error {
	v, ctx := f.value, fmt.Sprintf("could not set field %s to zero-value", f.Namespace())
	if !v.CanSet() {
		return errors.Wrap(ErrNotSettable, ctx)
	}
	v.Set(Zero(v))
	return nil
}

// SetNil sets the field to its zero value.
// Unsettable/Un-nillable struct fields will return an error.
func (f *StructField) SetNil() error {
	v, ctx := f.value, fmt.Sprintf("could not set field %s to nil", f.Namespace())
	if !v.CanSet() {
		return errors.Wrap(ErrNotSettable, ctx)
	}
	if !canNil(v) {
		return errors.Wrap(ErrNotNillable, ctx)
	}
	v.Set(Zero(v))
	return nil
}

// SetTime sets the field to the time.Time value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetTime(x time.Time) {
	v := f.value
	if v.CanSet() {
		v := Preset(v)
		v.Set(reflect.ValueOf(x))
	}
}

// SetDuration sets the field to the time.Duration value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetDuration(x time.Duration) {
	v := f.value
	if v.CanSet() {
		v := Preset(v)
		v.Set(reflect.ValueOf(x))
	}
}

// SetError sets the field to the error value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetError(x error) {
	v := f.value
	if v.CanSet() {
		v := Preset(v)
		v.Set(reflect.ValueOf(x))
	}
}

// SetString sets the field to the string value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetString(x string) {
	v := f.value
	if v.CanSet() {
		v := Preset(v)
		v.SetString(x)
	}
}

// SetBool sets the field to the bool value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetBool(x bool) {
	v := f.value
	if v.CanSet() {
		v := Preset(v)
		v.SetBool(x)
	}
}

// SetInt sets the field to the int64 value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetInt(x int64) {
	v := f.value
	if v.CanSet() {
		v := Preset(v)
		v.SetInt(x)
	}
}

// SetUint sets the field to the uint64 value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetUint(x uint64) {
	v := f.value
	if v.CanSet() {
		v := Preset(v)
		v.SetUint(x)
	}
}

// SetFloat sets the field to the float64 value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetFloat(x float64) {
	v := f.value
	if v.CanSet() {
		v := Preset(v)
		v.SetFloat(x)
	}
}

// SetComplex sets the field to the complex128 value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetComplex(x complex128) {
	v := f.value
	if v.CanSet() {
		v := Preset(v)
		v.SetComplex(x)
	}
}

// SetBytes sets the field to the slice of bytes value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetBytes(x []byte) {
	v := f.value
	if v.CanSet() {
		v := Preset(v)
		v.SetBytes(x)
	}
}

// SetInterface sets the field to the interface value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetInterface(x interface{}) { // NOTE: Not used!
	v := f.value
	if v.CanSet() {
		v := Preset(v)
		v.Set(reflect.ValueOf(x))
	}
}

// SetStruct sets the field to the StructValue value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetStruct(x *StructValue) {
	v := f.value
	if v.CanSet() {
		v := Preset(v)
		v.Set(x.value)
	}
}

// Value returns the underlying value of the field.
// Unexported struct fields will be neglected.
func (f *StructField) Value() reflect.Value {
	v := f.value
	if f.IsExported() {
		return v
	}
	f.Parent.setErrorsf(ErrNotExported, "could not get value of field %s", f.Namespace())
	t := v.Type()
	return reflect.New(t).Elem()
}

// PtrValue returns the underlying value of the field and handles struct fields that
// are pointer to data types by returning the element of the pointer instead.
// Unexported struct fields will be neglected.
func (f *StructField) PtrValue() reflect.Value {
	v := f.value
	t := v.Type()
	if f.IsExported() {
		if v.Kind() == reflect.Ptr {
			e := v.Elem()
			t := t.Elem()
			if !e.IsValid() { // Check if the pointer is nil
				return reflect.New(t).Elem()
			}
			return e
		}
		return v
	}
	f.Parent.setErrorsf(ErrNotExported, "cannot get field value of %s", f.Namespace())
	return reflect.New(t).Elem()
}

// AssignableTo reports whether a field value is assignable to interface dest.
func (f *StructField) AssignableTo(dest interface{}) bool {
	x := reflect.ValueOf(dest)
	return f.assignableTo(x)
}

// ... NOTE: Review/improve doc.
// Set sets the field to a given value dest. It returns an error if the field is not
// settable (not addressable or not exported) or if the given value's type doesn't
// match the fields type.
//
// The pointers are not expected for field ...
//
// NOTE: Set might benefit from using reflect.Type.AssignableTo() or ConvertibleTo().
func (f *StructField) Set(dest interface{}) error {
	fieldName := f.Namespace()
	if !f.CanSet() {
		return errors.Wrapf(ErrNotSettable, "could not set field %s", fieldName)
	}
	if dest == nil {
		return f.SetNil()
	}
	v := f.value
	x := reflect.ValueOf(dest)
	if canPtr(v) && !canPtr(x) {
		v = Preset(v)
	}
	//
	// Assignables
	switch {
	//
	// case canPtr(v) && canPtr(x):
	//     v.SetPointer(x.Pointer()); return nil
	case canTime(v) && canTime(x), canDuration(v) && canDuration(x), canError(v) && canError(x):
		v.Set(x)
		return nil
	case canString(v) && canString(x):
		v.SetString(x.String())
		return nil
	case canBool(v) && canBool(x):
		v.SetBool(x.Bool())
		return nil
	case canInt(v) && canInt(x):
		int64X := x.Int()
		if v.OverflowInt(int64X) {
			return errors.Errorf("field %s(%s) could not represent int64", fieldName, x.Type())
		}
		v.SetInt(int64X)
		return nil
	case canUint(v) && canUint(x):
		uint64X := x.Uint()
		if v.OverflowUint(uint64X) {
			return errors.Errorf("field %s(%s) could not represent uint64", fieldName, x.Type())
		}
		v.SetUint(uint64X)
		return nil
	case canFloat(v) && canFloat(x):
		float64X := x.Float()
		if v.OverflowFloat(float64X) {
			return errors.Errorf("field %s(%s) could not represent float64", fieldName, x.Type())
		}
		v.SetFloat(float64X)
		return nil
	case canComplex(v) && canComplex(x):
		complex128X := x.Complex()
		if v.OverflowComplex(complex128X) {
			return errors.Errorf("field %s(%s) could not represent complex128", fieldName, x.Type())
		}
		v.SetComplex(complex128X)
		return nil
	case canBytes(v) && canBytes(x):
		v.SetBytes(x.Bytes())
		return nil
	case f.assignableTo(x):
		v.Set(x)
		return nil
	}
	//
	// Semi-Assignables
	switch {
	//
	// - text     <- bool
	//   text     <- number
	//   text     <- []byte
	//   text     <- date
	// - bool     <- text
	// - number   <- bool
	//   number   <- float (losing decimal point value)
	// - float    <- number
	// - []byte   <- text
	// - date     <- text
	// - duration <- text
	//   duration <- number
	// - error    <- text
	case canString(v):
		switch {
		case canBool(x):
			v.SetString(fmt.Sprintf("%t", x.Bool()))
			return nil
		case canInt(x):
			v.SetString(fmt.Sprintf("%d", x.Int()))
			return nil
		case canUint(x):
			v.SetString(fmt.Sprintf("%d", x.Uint()))
			return nil
		case canFloat(x):
			v.SetString(fmt.Sprintf("%f", x.Float()))
			return nil
		case canBytes(x):
			v.SetString(string(x.Bytes()))
			return nil
		case canTime(x):
			v.SetString(Time(x).Format(time.RFC3339))
			return nil
		}
	case canBool(v):
		switch {
		case canString(x):
			switch strings.ToLower(x.String()) {
			case "true", "yes", "y", "ok", "1":
				v.SetBool(true)
			default:
				v.SetBool(false)
			}
			return nil
		}
	case canInt(v):
		switch {
		case canBool(x):
			b := x.Bool()
			if b {
				v.SetInt(1) // var i int64 = 1; v.SetInt(i)
			} else {
				v.SetInt(0) // var i int64 = 0; v.SetInt(i)
			}
			return nil
		case canFloat(x):
			v.SetInt(int64(x.Float()))
			return nil
		}
	case canUint(v):
		switch {
		case canBool(x):
			b := x.Bool()
			if b {
				v.SetUint(1) // var i uint64 = 1; v.SetUint(i)
			} else {
				v.SetUint(0) // var i uint64 = 0; v.SetUint(i)
			}
			return nil
		case canFloat(x):
			v.SetUint(uint64(x.Float()))
			return nil
		}
	case canFloat(v):
		switch {
		case canInt(x):
			v.SetFloat(float64(x.Int()))
		case canUint(x):
			v.SetFloat(float64(x.Uint()))
		}
	case canBytes(v):
		switch {
		case canString(x):
			v.SetBytes([]byte(x.String()))
			return nil
		}
	case canTime(v):
		switch {
		case canString(x):
			txt := x.String()
			t, found := time.Now(), false
			layouts := []string{
				"2006-01-02 15:04:05",    // MySQL DATETIME
				"2006-01-02",             //       DATE
				"2006/01/02",             // EXCEL DATE
				"02-Jan-2006",            //       DATE
				"01-02-2006 03:04:05 PM", // CSV   DATETIME1
				"02/01/2006  15:04:05",   //       DATETIME2
				"02/01/2006 15:04",       //       DATETIME3
			}
			var errs error
			for _, layout := range layouts {
				date, err := time.Parse(layout, txt)
				if err != nil {
					if errs == nil {
						errs = err
					} else {
						errs = errors.Wrap(errs, err.Error())
					}
				}
				t = date
				found = true
				break
			}
			if found {
				x := reflect.ValueOf(t)
				v.Set(x)
			} else {
				return errors.Wrapf(errs, "Invalid date value. found: %s; formats expected: %v", txt, layouts)
			}
			return nil
		}
	case canDuration(v):
		switch {
		case canString(x):
			t := x.String()
			if !strings.ContainsAny(t, "nsuµmh") {
				t = t + "ns"
			}
			d, err := time.ParseDuration(t)
			if err != nil {
				return errors.Wrapf(err, "Invalid duration value. found: %s; want: [1s, 3h, ... ]", t)
			}
			v.Set(reflect.ValueOf(d))
			return nil
		case canInt(x):
			v.Set(reflect.ValueOf(time.Duration(x.Int())))
			return nil
		case canUint(x):
			v.Set(reflect.ValueOf(time.Duration(int64(x.Uint()))))
			return nil
		case canFloat(x):
			v.Set(reflect.ValueOf(time.Duration(int64(x.Float()))))
			return nil
		}
	case canError(v):
		switch {
		case canString(x):
			v.Set(reflect.ValueOf(errors.New(x.String())))
			return nil
		}
	}
	//
	// Non-Assignables
	return errors.Errorf("wrong kind of value for field %s. got: '%s' want: '%s'", fieldName, x.Type(), v.Type())
}

// func (f *StructField) SetElem(i int, dest interface{}) error {
// 	fieldName := f.Namespace()
// 	if !f.CanSet() {
// 		return errors.Wrapf(ErrNotSettable, "could not set field %s", fieldName)
// 	}
// 	v := f.value
// 	if dest == nil {
// 		return f.SetNil()
// 	}
// 	// v = Preset(v) // make sure no pointers points to something concrete, instead of nil
// 	x := reflect.ValueOf(dest)
// 	// vi := v.Interface()
// 	// xi := x.Interface()
// 	// fmt.Printf("%s: vi: %v; xi: %v.\n", fieldName, vi, xi)
// 	// switch {
// 	// case canSlice(v):
// 	// 	if canSlice(x) {
// 	// 		v.Set(x)
// 	// 		return nil
// 	// 	}
// 	// case canMap(v):
// 	// 	if canMap(x) {
// 	// 		v.Set(x)
// 	// 		return nil
// 	// 	}
// 	// case v.Kind() == x.Kind():
// 	// 	v.Set(x)
// 	// 	return nil
// 	// }
// 	return errors.Errorf("wrong kind of value for field %s. got: '%s' want: '%s'", fieldName, x.Kind(), v.Kind())
// }

// func GetUnexportedField(field reflect.Value) interface{} {
//     return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
// }

// func SetUnexportedField(field reflect.Value, value interface{}) {
//     reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
//         Elem().
//         Set(reflect.ValueOf(value))
// }

/*   U n e x p o r t e d   */

// // value is the reflection interface to a Go value.
// func (f *StructField) value() reflect.Value {
// 	return f.Parent.value.Field(f.index)
// }

// // field describes a single field in a struct.
// func (f *StructField) field() reflect.StructField {
// 	return f.Parent.value.Type().Field(f.index)
// }

// equal compares field value with reflect value argument and returns field index
// if they are equal, else returns OutOfRange, i.e. -1.
//
// NOTE: Equal might benefit from using reflect.Type.Comparable().
// ... NOTE: Review/improve doc.
func (f *StructField) equal(x reflect.Value) int {
	if !f.IsExported() {
		return OutOfRange
	}
	v := f.value
	// if !v.Type().Comparable() {
	// 	return OutOfRange
	// }
	switch {
	case canDuration(v) && canDuration(x):
		if Duration(v) == Duration(x) {
			return f.Index()
		}
	case canTime(v) && canTime(x):
		if Time(v) == Time(x) {
			return f.Index()
		}
	case canError(v) && canError(x):
		errV := Error(v)
		errX := Error(x)
		if errV.Error() == errX.Error() { // if errors.Is(errV, errX) {
			return f.Index()
		}
	case canString(v) && canString(x):
		if v.String() == x.String() {
			return f.Index()
		}
	case canBool(v) && canBool(x):
		if v.Bool() == x.Bool() {
			return f.Index()
		}
	case canInt(v) && canInt(x):
		int64V := v.Int()
		int64X := x.Int()
		if v.OverflowInt(int64X) {
			return OutOfRange
		}
		if int64V == int64X {
			return f.Index()
		}
	case canUint(v) && canUint(x):
		uint64V := v.Uint()
		uint64X := x.Uint()
		if v.OverflowUint(uint64X) {
			return OutOfRange
		}
		if uint64V == uint64X {
			return f.Index()
		}
	case canFloat(v) && canFloat(x):
		float64V := v.Float()
		float64X := x.Float()
		if v.OverflowFloat(float64X) {
			return OutOfRange
		}
		if float64V == float64X {
			return f.Index()
		}
	case canComplex(v) && canComplex(x):
		complex128V := v.Complex()
		complex128X := x.Complex()
		if v.OverflowComplex(complex128X) {
			return OutOfRange
		}
		if complex128V == complex128X {
			return f.Index()
		}
	case canBytes(v) && canBytes(x):
		if bytes.Equal(v.Bytes(), x.Bytes()) {
			return f.Index()
		}
	case f.assignableTo(x), v.CanInterface(), f.CanStruct():
		if reflect.DeepEqual(v.Interface(), x.Interface()) {
			return f.Index()
		}
	default:
		// case canPtr(v) && canPtr(x): v.SetPointer(x.Pointer()); return nil
		// case reflect.Invalid:
		// case reflect.Slice:
		// case reflect.Array:
		// case reflect.Map:
		// case reflect.Func:
		// case reflect.Chan:
		// case reflect.Ptr:
		// case reflect.Uintptr:
		// case reflect.UnsafePointer:
	}
	return OutOfRange
}

// assignableTo reports whether a field value is assignable to reflect value x.
func (f *StructField) assignableTo(x reflect.Value) bool {
	v := f.value
	vt := v.Type()
	xt := x.Type()
	// fmt.Printf("%s(%s) with dest(%s) == %v\n", f.Namespace(), vt, xt, x.Interface())
	switch {
	case vt.AssignableTo(xt):
		// fmt.Println("vt.AssignableTo(xt)")
		return true
	case vt.ConvertibleTo(xt): // Not Used?
		// fmt.Println("vt.ConvertibleTo(xt)")
		return true
	case vt.String() == xt.String():
		// fmt.Println("vt.String() == xt.String()")
		return true
	// case canTime(v):
	// 	fmt.Println("canTime(v)")
	// 	if _, ok := x.Interface().(time.Time); ok {
	// 		return true
	// 	}
	// case canDuration(v):
	// 	fmt.Println("canDuration(v)")
	// 	if _, ok := x.Interface().(time.Duration); ok {
	// 		return true
	// 	}
	// case canStruct(v):
	// 	fmt.Println("canStruct(v)")
	// 	s, err := New(x.Interface(), f.Parent)
	// 	if err != nil {
	// 		return errors.Wrapf(err, "could not load field %s's as nested struct", f.Namespace())
	// 	}
	// 	if Type(v) != s.Name() {
	// 		return errors.Errorf("wrong struct name for field %s. got: '%s' want: '%s'", f.Namespace(), s.Name(), Type(v))
	// 	}
	// 	return true
	// case canSlice(v):
	// 	fmt.Println("canSlice(v)")
	// 	if canSlice(x) {
	// 		return true
	// 	}
	// 	// Value struct
	// 	//
	// 	// // Copy copies the contents of src into dst until either
	// 	// // dst has been filled or src has been exhausted.
	// 	// // It returns the number of elements copied.
	// 	// // Dst and src each must have kind Slice or Array, and
	// 	// // dst and src must have the same element type.
	// 	// //
	// 	// // As a special case, src can have kind String if the element type of dst is kind Uint8.
	// 	// func Copy(dst, src Value) int
	// case canMap(v):
	// 	// fmt.Println("canMap(v)")
	// 	if canMap(x) {
	// 		return true
	// 	}
	case canError(v) && canError(v): // vt.String() == "error":
		// fmt.Println("canError(v)")
		// if vt.Implements(xt) {
		// 	fmt.Println("Type.Implements.Type(error)")
		// if _, ok := x.Interface().(interface{ Error() string }); ok {
		return true
	case canInterface(v):
		// fmt.Println("canInterface(v)")
		// canInterface(x): no need
		// x came from interface{} dest
		return true
	}
	return false
}
