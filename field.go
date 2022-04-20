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
	"github.com/roninzo/structs/utils"
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
	s.setErrorf("invalid argument type; want: %q or %q, got: %q", reflect.String, reflect.Int, t.Kind())
	return nil
}

/*   I m p l e m e n t a t i o n   */

// IsValid returns true if StructField has been loaded successfully.
// Useful for checking StructField is valid before use.
func (f *StructField) IsValid() bool {
	if f.Parent != nil && f.index != OutOfRange {
		v := f.value
		return v.IsValid()
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

// NameJson returns returns the string name of StructField
// defined in its related json struct tag, else it generates it.
func (f *StructField) NameJson() string {
	tag, ok := f.Tag("json")
	if ok && tag != "-" {
		tag = strings.TrimSuffix(tag, ",omitempty")
		return tag
	}
	return utils.CamelCaseToUnderscore(f.field.Name)
}

// Default returns returns the string default value of StructField
// defined in its related default struct tag, else returns empty string.
func (f *StructField) Default() string {
	tag, ok := f.Tag("default")
	if ok {
		// string:   Roninzo
		// number:   12
		// decimal:  3.99
		// bool:     false
		// datetime: 2020-12-11T01:00:00+02:00
		skip := strings.Contains(tag, "(") &&
			strings.Contains(tag, ")") || strings.ToLower(tag) == "null" || tag == ""
		if !skip {
			return strings.TrimSpace(tag)
		}
	}
	return ""
}

// FullName is similar to the Name method, except that it includes its related struct names
// all the way to the top level struct (in a dot separated string).
func (f *StructField) FullName() (n string) {
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

// TO REVISIT

// IsExported returns true if the given field is exported and its json tag is
// not equal to "-". Those fields are neglected for getter and setter methods.
func (f *StructField) IsExported() bool {
	if ok := f.field.PkgPath == ""; ok {
		return true
	}
	return false
}

// TO REVISIT

// IsHidden returns true if the given field is exported and its json tag is
// not equal to "-". Those fields are neglected for getter and setter methods.
func (f *StructField) IsHidden() bool {
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

// Interface returns true if underlying value of the field is modifiable.
func (f *StructField) CanSet() bool {
	return f.value.CanSet()
}

// Zero returns field's type specific zero value. For instance, the zero-value
// of a string field is "", of an int is 0, and so on.
func (f *StructField) Zero() reflect.Value {
	v := f.value
	return utils.Zero(v)
}

// IsZero returns true if the given field is a zero-value, i.e. not initialized.
// Unexported struct fields will be neglected.
func (f *StructField) IsZero() bool {
	if f.IsExported() {
		return reflect.DeepEqual(f.Interface(), f.Zero().Interface()) // v := f.value; z := utils.Zero(v); return v == z
	}
	return false
}

// IsNil reports whether its argument f is nil. The argument must be a chan, func,
// interface, map, pointer, or slice value; if it is not, IsNil returns nil.
// Unexported struct fields will be neglected.
func (f *StructField) IsNil() bool {
	v := f.value
	if utils.CanNil(v) {
		return v.IsNil()
	}
	return false
}

// TO REVISIT

// Equal compares field value with reflect value argument and returns true
func (f *StructField) Equal(x *StructField) bool {
	if x == nil {
		return false
	}
	return f.equal(x.value) != OutOfRange
}

// TO REVISIT

// equal compares field value with reflect value argument and returns field index
// if they are equal, else returns OutOfRange, i.e. -1.
//
// NOTE: Equal might benefit from using reflect.Type.Comparable().
// if !v.Type().Comparable() {
// 	return OutOfRange
// }
//
// case utils.CanPtr(v) && utils.CanPtr(x): v.SetPointer(x.Pointer()); return nil
// reflect.Invalid, reflect.Slice, reflect.Array, reflect.Map, reflect.Func, reflect.Chan, reflect.Ptr, reflect.Uintptr, reflect.UnsafePointer:
func (f *StructField) equal(x reflect.Value) int {
	v := f.value
	switch {
	case !f.IsExported():
		return OutOfRange
	case utils.CanStruct(v) && utils.CanStruct(x):
		if reflect.DeepEqual(v.Interface(), x.Interface()) {
			return f.Index()
		}
	case utils.CanDuration(v) && utils.CanDuration(x):
		if utils.Duration(v) == utils.Duration(x) {
			return f.Index()
		}
	case utils.CanTime(v) && utils.CanTime(x):
		if utils.Time(v) == utils.Time(x) {
			return f.Index()
		}
	case utils.CanError(v) && utils.CanError(x):
		errV := utils.Error(v)
		errX := utils.Error(x)
		if errV.Error() == errX.Error() { // if errors.Is(errV, errX) {
			return f.Index()
		}
	case utils.CanString(v) && utils.CanString(x):
		if v.String() == x.String() {
			return f.Index()
		}
	case utils.CanBool(v) && utils.CanBool(x):
		if v.Bool() == x.Bool() {
			return f.Index()
		}
	case utils.CanInt(v) && utils.CanInt(x):
		int64V := v.Int()
		int64X := x.Int()
		if v.OverflowInt(int64X) {
			return OutOfRange
		}
		if int64V == int64X {
			return f.Index()
		}
	case utils.CanUint(v) && utils.CanUint(x):
		uint64V := v.Uint()
		uint64X := x.Uint()
		if v.OverflowUint(uint64X) {
			return OutOfRange
		}
		if uint64V == uint64X {
			return f.Index()
		}
	case utils.CanFloat(v) && utils.CanFloat(x):
		float64V := v.Float()
		float64X := x.Float()
		if v.OverflowFloat(float64X) {
			return OutOfRange
		}
		if float64V == float64X {
			return f.Index()
		}
	case utils.CanComplex(v) && utils.CanComplex(x):
		complex128V := v.Complex()
		complex128X := x.Complex()
		if v.OverflowComplex(complex128X) {
			return OutOfRange
		}
		if complex128V == complex128X {
			return f.Index()
		}
	case utils.CanBytes(v) && utils.CanBytes(x):
		if bytes.Equal(v.Bytes(), x.Bytes()) {
			return f.Index()
		}
	case f.AssignableTo(x), v.CanInterface():
		if reflect.DeepEqual(v.Interface(), x.Interface()) {
			return f.Index()
		}
	}
	return OutOfRange
}

// G e t t e r s
// Getter methods return f's underlying value for a specifc type T. It panics if f's kind is not T.

// TO REVISIT

// Get returns the value of the field as interface.
// reflect.Struct, reflect.Slice, reflect.Array, reflect.Map, reflect.Interface,
// reflect.Func, reflect.Chan, reflect.Uintptr, reflect.UnsafePointer, reflect.Invalid
// Unexported struct fields will be neglected.
func (f *StructField) Get() interface{} {
	v := f.Indirect()
	switch {
	case !f.IsExported():
		return nil
	case utils.CanDuration(v):
		return utils.Duration(v)
	case utils.CanTime(v):
		return utils.Time(v)
	case utils.CanError(v):
		return utils.Error(v).Error()
	case utils.CanString(v):
		return v.String()
	case utils.CanBool(v):
		return v.Bool()
	case utils.CanInt(v):
		return v.Int()
	case utils.CanUint(v):
		return v.Uint()
	case utils.CanFloat(v):
		return v.Float()
	case utils.CanComplex(v):
		return v.Complex()
	case utils.CanBytes(v):
		return v.Bytes()
	default:
		return v.Interface()
	}
}

func (f *StructField) Time() time.Time         { v := f.value; return utils.Time(v) }
func (f *StructField) Duration() time.Duration { v := f.value; return utils.Duration(v) }
func (f *StructField) Error() error            { v := f.value; return utils.Error(v) }
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
	s := IndirectStruct(f.value)
	s.Parent = f.Parent
	f.Parent.Error = s.Err()
	return s
}

// C h e c k e r s
// Checker methods report whether type requested can be used without panicking.

func (f *StructField) CanNil() bool       { return utils.CanNil(f.value) }
func (f *StructField) CanPtr() bool       { return utils.CanPtr(f.value) }
func (f *StructField) CanTime() bool      { return utils.CanTime(f.value) }
func (f *StructField) CanDuration() bool  { return utils.CanDuration(f.value) }
func (f *StructField) CanError() bool     { return utils.CanError(f.value) }
func (f *StructField) CanString() bool    { return utils.CanString(f.value) }
func (f *StructField) CanBool() bool      { return utils.CanBool(f.value) }
func (f *StructField) CanInt() bool       { return utils.CanInt(f.value) }
func (f *StructField) CanUint() bool      { return utils.CanUint(f.value) }
func (f *StructField) CanFloat() bool     { return utils.CanFloat(f.value) }
func (f *StructField) CanComplex() bool   { return utils.CanComplex(f.value) }
func (f *StructField) CanBytes() bool     { return utils.CanBytes(f.value) }
func (f *StructField) CanSlice() bool     { return utils.CanSlice(f.value) }
func (f *StructField) CanMap() bool       { return utils.CanMap(f.value) }
func (f *StructField) CanStruct() bool    { return utils.CanStruct(f.value) }
func (f *StructField) CanInterface() bool { return utils.CanInterface(f.value) }

// S e t t e r s
// Setter methods assigns x to the field f. no assignment is carried out if CanSet
// returns false. As in Go, x's value must be assignable to f's type.

// SetZero sets the field to its zero value.
// Unsettable struct fields will return an error.
func (f *StructField) SetZero() error {
	v, ctx := f.value, fmt.Sprintf("could not set field %s to zero-value", f.FullName())
	if !v.CanSet() {
		return errors.Wrap(ErrNotSettable, ctx)
	}
	v.Set(utils.Zero(v))
	return nil
}

// SetNil sets the field to its zero value.
// Unsettable/Un-nillable struct fields will return an error.
func (f *StructField) SetNil() error {
	v, ctx := f.value, fmt.Sprintf("could not set field %s to nil", f.FullName())
	if !v.CanSet() {
		return errors.Wrap(ErrNotSettable, ctx)
	}
	if !utils.CanNil(v) {
		return errors.Wrap(ErrNotNillable, ctx)
	}
	v.Set(utils.Zero(v))
	return nil
}

// SetTime sets the field to the time.Time value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetTime(x time.Time) {
	v := f.value
	if v.CanSet() {
		utils.PresetIndirect(v).Set(reflect.ValueOf(x))
	}
}

// SetDuration sets the field to the time.Duration value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetDuration(x time.Duration) {
	v := f.value
	if v.CanSet() {
		utils.PresetIndirect(v).Set(reflect.ValueOf(x))
	}
}

// SetError sets the field to the error value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetError(x error) {
	v := f.value
	if v.CanSet() {
		utils.PresetIndirect(v).Set(reflect.ValueOf(x))
	}
}

// SetString sets the field to the string value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetString(x string) {
	v := f.value
	if v.CanSet() {
		utils.PresetIndirect(v).SetString(x)
	}
}

// SetBool sets the field to the bool value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetBool(x bool) {
	v := f.value
	if v.CanSet() {
		utils.PresetIndirect(v).SetBool(x)
	}
}

// SetInt sets the field to the int64 value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetInt(x int64) {
	v := f.value
	if v.CanSet() {
		utils.PresetIndirect(v).SetInt(x)
	}
}

// SetUint sets the field to the uint64 value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetUint(x uint64) {
	v := f.value
	if v.CanSet() {
		utils.PresetIndirect(v).SetUint(x)
	}
}

// SetFloat sets the field to the float64 value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetFloat(x float64) {
	v := f.value
	if v.CanSet() {
		utils.PresetIndirect(v).SetFloat(x)
	}
}

// SetComplex sets the field to the complex128 value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetComplex(x complex128) {
	v := f.value
	if v.CanSet() {
		utils.PresetIndirect(v).SetComplex(x)
	}
}

// SetBytes sets the field to the slice of bytes value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetBytes(x []byte) {
	v := f.value
	if v.CanSet() {
		utils.PresetIndirect(v).SetBytes(x)
	}
}

// SetInterface sets the field to the interface value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetInterface(x interface{}) { // NOTE: Not used!
	v := f.value
	if v.CanSet() {
		utils.PresetIndirect(v).Set(reflect.ValueOf(x))
	}
}

// SetStruct sets the field to the StructValue value x.
// Unsettable struct fields will return an error.
func (f *StructField) SetStruct(x *StructValue) {
	v := f.value
	if v.CanSet() {
		utils.PresetIndirect(v).Set(x.value)
	}
}

// TO REVISIT

// Value returns the underlying value of the field.
// Unexported struct fields will be neglected.
func (f *StructField) Value() reflect.Value {
	v := f.value
	if f.IsExported() {
		return v
	}
	f.Parent.setErrorsf(ErrNotExported, "could not get value of field %s", f.FullName())
	t := v.Type()
	return reflect.New(t).Elem()
}

// Indirect returns the value that StructField f.value points to.
// If f.value is a nil pointer, Indirect returns a zero Value.
// If f.value is not a pointer, Indirect returns f.value.
func (f *StructField) Indirect() reflect.Value {
	return reflect.Indirect(f.value)
}

// IndirectType returns the type that field f points to.
// If f is a pointer, IndirectType returns the type f points to.
// If f is not a pointer, IndirectType returns the type of f.
func (f *StructField) IndirectType() reflect.Type {
	return utils.IndirectType(f.value)
}

// AssignableTo reports whether a field value is assignable to reflect Value x.
//
// NOTE: case utils.CanInterface(v) && utils.CanInterface(x):
//       no need; x came from interface{} dest
func (f *StructField) AssignableTo(x reflect.Value) bool {
	v := f.value
	vT := v.Type()
	xT := x.Type()
	switch {
	case vT.AssignableTo(xT):
		return true
	case vT.String() == xT.String():
		return true
	case utils.CanInterface(v):
		return true
	}
	return false
}

// TO REVISIT

// Set sets the field to a given value dest. It returns an error if the field is not
// settable (not addressable or not exported) or if the given value's type doesn't
// match the fields type.
//
// The pointers are not expected for field ...
//
// It has already been established that x is not nil
// by bailing out on the 'dest == nil' condition above.
//
// - date     <- date
//   date     <- text
// - duration <- duration
//   duration <- text
//   duration <- number
// - error    <- error
//   error    <- text
// - text     <- text
//   text     <- bool
//   text     <- number
//   text     <- []byte
//   text     <- date
// - bool     <- bool
//   bool     <- text
//   bool     <- number
// - number   <- number
//   number   <- bool
//   number   <- float (losing decimal point value)
// - float    <- float
//   float    <- number
// - []byte   <- []byte
//   []byte   <- text
// - complex  <- complex
//
// NOTE: Set might benefit from using reflect.Type.AssignableTo() or ConvertibleTo().
func (f *StructField) Set(dest interface{}) error {
	fullname := f.FullName()

	if !f.CanSet() {
		return errors.Wrapf(ErrNotSettable, "could not set field %s", fullname)
	}

	// Set(nil) <=> SetNil()
	if dest == nil {
		return f.SetNil()
	}

	v := f.value
	x := reflect.ValueOf(dest)

	if utils.CanPtr(v) && !utils.CanPtr(x) {
		v = utils.PresetIndirect(v)
	}

	// Assignables
	switch {
	case f.AssignableTo(x):
		if fullname == "Interface" {
			fmt.Printf("%q is assignable\n", fullname)
		}
		v.Set(x)
		return nil
	case utils.CanTime(v):
		switch {
		case utils.CanTime(x):
			v.Set(x)
			return nil
		case utils.CanString(x):
			txt := x.String()
			t, err := utils.StringToTime(txt)
			if err != nil {
				return err
			}
			x = reflect.ValueOf(t)
			v.Set(x)
			return nil
		}
	case utils.CanDuration(v):
		switch {
		case utils.CanDuration(x):
			v.Set(x)
			return nil
		case utils.CanString(x):
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
		case utils.CanInt(x):
			v.Set(reflect.ValueOf(time.Duration(x.Int())))
			return nil
		case utils.CanUint(x):
			v.Set(reflect.ValueOf(time.Duration(int64(x.Uint()))))
			return nil
		case utils.CanFloat(x):
			v.Set(reflect.ValueOf(time.Duration(int64(x.Float()))))
			return nil
		}
	case utils.CanError(v):
		switch {
		case utils.CanError(x):
			v.Set(x)
			return nil
		case utils.CanString(x):
			v.Set(reflect.ValueOf(errors.New(x.String())))
			return nil
		}
	case utils.CanString(v):
		switch {
		case utils.CanString(x):
			v.SetString(x.String())
			return nil
		case utils.CanBool(x):
			v.SetString(fmt.Sprintf("%t", x.Bool()))
			return nil
		case utils.CanInt(x):
			v.SetString(fmt.Sprintf("%d", x.Int()))
			return nil
		case utils.CanUint(x):
			v.SetString(fmt.Sprintf("%d", x.Uint()))
			return nil
		case utils.CanFloat(x):
			v.SetString(fmt.Sprintf("%f", x.Float()))
			return nil
		case utils.CanBytes(x):
			v.SetString(string(x.Bytes()))
			return nil
		case utils.CanTime(x):
			v.SetString(utils.Time(x).Format(time.RFC3339))
			return nil
		}
	case utils.CanBool(v):
		switch {
		case utils.CanBool(x):
			v.SetBool(x.Bool())
			return nil
		case utils.CanString(x):
			switch strings.ToLower(x.String()) {
			case "true", "yes", "y", "ok", "1":
				v.SetBool(true)
				return nil
			default:
				v.SetBool(false)
				return nil
			}
		case utils.CanInt(x):
			switch x.Int() {
			case 1:
				v.SetBool(true)
				return nil
			case 0:
				v.SetBool(true)
				return nil
			}
		case utils.CanUint(x):
			switch x.Uint() {
			case 1:
				v.SetBool(false)
				return nil
			case 0:
				v.SetBool(false)
				return nil
			}
		case utils.CanFloat(x):
			switch x.Float() {
			case 1.0:
				v.SetBool(false)
				return nil
			case 0.0:
				v.SetBool(false)
				return nil
			}
		}
	case utils.CanInt(v):
		switch {
		case utils.CanInt(x):
			int64X := x.Int()
			if v.OverflowInt(int64X) {
				return errors.Errorf("field %s(%s) could not represent int64", fullname, x.Type())
			}
			v.SetInt(int64X)
			return nil
		case utils.CanBool(x):
			b := x.Bool()
			if b {
				v.SetInt(1) // var i int64 = 1; v.SetInt(i)
				return nil
			} else {
				v.SetInt(0) // var i int64 = 0; v.SetInt(i)
				return nil
			}
		case utils.CanFloat(x):
			v.SetInt(int64(x.Float()))
			return nil
		}
	case utils.CanUint(v):
		switch {
		case utils.CanUint(x):
			uint64X := x.Uint()
			if v.OverflowUint(uint64X) {
				return errors.Errorf("field %s(%s) could not represent uint64", fullname, x.Type())
			}
			v.SetUint(uint64X)
			return nil
		case utils.CanBool(x):
			b := x.Bool()
			if b {
				v.SetUint(1) // var i uint64 = 1; v.SetUint(i)
				return nil
			} else {
				v.SetUint(0) // var i uint64 = 0; v.SetUint(i)
				return nil
			}
		case utils.CanFloat(x):
			v.SetUint(uint64(x.Float()))
			return nil
		}
	case utils.CanFloat(v):
		switch {
		case utils.CanFloat(x):
			float64X := x.Float()
			if v.OverflowFloat(float64X) {
				return errors.Errorf("field %s(%s) could not represent float64", fullname, x.Type())
			}
			v.SetFloat(float64X)
			return nil
		case utils.CanInt(x):
			v.SetFloat(float64(x.Int()))
			return nil
		case utils.CanUint(x):
			v.SetFloat(float64(x.Uint()))
			return nil
		}
	case utils.CanBytes(v):
		switch {
		case utils.CanBytes(x):
			v.SetBytes(x.Bytes())
			return nil
		case utils.CanString(x):
			v.SetBytes([]byte(x.String()))
			return nil
		}
	case utils.CanComplex(v):
		switch {
		case utils.CanComplex(x):
			complex128X := x.Complex()
			if v.OverflowComplex(complex128X) {
				return errors.Errorf("field %s(%s) could not represent complex128", fullname, x.Type())
			}
			v.SetComplex(complex128X)
			return nil
		}
	}

	// Non-Assignables
	return errors.Errorf("wrong kind of value for field %s. got: %q want: %q", fullname, x.Type(), v.Type())
}
