// Copyright 2021 Roninzo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package structs implements a generic interface for manipulating Go structs.
// The related API is powered and inspired from the Go reflection package.
//
// Introduction
//
// It was initially developed to provide generic field getters and setters
// for any struct. It has grown into an generic abstraction layer to structs
// powered by the reflect package.
//
// While not being particularly performant compared to other structs packages,
// it provides a "natural feel" API, convenience, flexibility and simpler
// alternative compared to using the reflect package directly.
//
// Throughout the documentation, the t variable is assumed to be of type
// struct T. T has two fields: the first is a string field called "Property" and
// the second is a nested struct via a pointer called "Nested". While "Property" is
// set to value "Test", "Number" inside "Nested" is set to 123456.
//   type T struct{
//      Property string
//      Nested   *N
//   }
//
//   type N struct{
//      Number int
//   }
//
//   t := T{
//      Property: "Test",
//      Nested:   &N{Number: 123456},
//   }
// NOTE: the same applies to all other variables subsequently declared.
//
// Support
//
// The following input are supported throughout the package:
//
//   Types   Description                     Example
//
//   T       a struct                        New(t)
//   *T      a pointer to a struct           New(&t)
//   []T     a slice of struct               New([]T{t})
//   *[]T    a pointer to a slice of struct  New(&[]T{t})
//   []*T    a slice of pointers to struct   New([]*T{&t})
//
//
// NOTE: See the findStruct method for details on the above scenarios.
//
// Implementation
//
// There are two main ways to use for this package to manipulate structs.
// Either by using the objects and methods or by using helper functions.
//
// Objects and methods:
// this approach requires calling of the New method first, which gives access
// to the API of the StructValue, StructField and StructFields objects.
//
// Helper functions:
// A set of self-contained functions that hides internal implementations
// powered by the StructValue object, are all declared in the file called
// helpers.go.
//
// The following table summarizes the above statements:
//
//   Usage     Applicable Go Files  Description
//
//   Objects   structs.go           StructValue object
//             field.go             StructField object
//             fields.go            StructFields object
//             rows.go              StructRows object
//
//   Helpers   helpers.go           Wrapper object functions
//
//
// All objects in this package are linked to the main StructValue object.
// The relationships between each one of them are as follow:
//
//    --------------
//   | *StructValue |<----------------------------------------+
//    --+-----------                                          |
//      |                                                     |
//      |                   ---------------                   |
//      +---> Field(x) --->| *StructField  |---> Parent ----->+
//      |                   ---------------                   |
//      |                                                     |
//      |                   ---------------                   |
//      +---> Fields() --->| *StructFields |---> Parent() --->+
//      |                   ---------------                   |
//      |                                                     |
//      |                   ---------------                   |
//      +---> Rows() ----->| *StructRows   |---> Parent ----->+
//                          ---------------
//
// NOTE: For an exhaustive illustration of package capabilities, please refer
// to the following file: https://github.com/roninzo/structs/example_test.go.
//
// StructValue
//
// The StructValue object is the starting point for manipulating structs using
// objects and methods. To initialize the StructValue object, make use of the
// New method followed by handling any potential error encountered in this
// process.
//
// Example:
//   s, err := structs.New(&t)
//   if err != nil {
//      return err
//   }
//
// From there, several methods provides information about the struct.
//
// Example:
//   fmt.Printf("t has %d field(s)\n", s.NumField())
//
//   // Output:
//   // t has 2 field(s)
//
// NOTE: When possible, all object method names were inspired from  the
// reflect package, trying to reduce the learning curve.
//
// Example:
//   if s.CanSet() {
//      fmt.Println("t is modifiable")
//   }
//
//   // Output:
//   // t is modifiable
//
// The StructValue object is also the gateway to the other two objects declared
// in this package: StructField and StructFields.
//
// Examples:
//   f := s.Field("Property") // f is a StructField object
//   fields := s.Fields()     // fields is a StructFields object
//
// StructField
//
// The StructField object represents one field in the struct, and provides getters
// and setters methods.
//
// Before reading data out of struct fields generically, it is recommended to get extra
// information about the struct field. This is useful if the type of the field is not
// known at runtime.
//
// Example:
//   f := s.Field("Property")
//   if f.CanString() {
//      fmt.Printf("Property was equal to '%s'\n", s.String())
//      if f.CanSet() {
//         err := f.SetString("Verified")
//         if err != nil {
//            return err
//         }
//         fmt.Printf("Property is now equal to '%s'\n", s.String())
//      }
//   }
//
//   // Output:
//   // Property was equal to 'Test'
//   // Property is now equal to 'Verified'
//
// However, if nested struct are present inside t, sub-fields are not made available
// directly. This means that nested structs must be loaded explicitly with the Struct
// method.
//
// Example:
//   f := s.Field("Nested")
//   if err := s.Err(); err != nil {
//      return err
//   }
//   if f.CanStruct() {
//      f = f.Struct().Field("Number")
//      if f.CanInt() {
//         fmt.Printf("Number was equal to %d\n", f.Int())
//         if f.CanSet() {
//            err := f.SetInt(654321)
//            if err != nil {
//               return err
//            }
//            fmt.Printf("Number is now equal to %d\n", f.Int())
//         }
//      }
//   }
//
//   // Output:
//   // Number was equal to 123456
//   // Number is now equal to 654321
//
// StructFields
//
// The StructFields object represents all the fields in the struct. Its main purpose
// at the moment, is to loop through StructField objects in a "for range" loop.
//
// Example:
//   for _, f := range s.Fields() {
//      fmt.Printf("struct field name is: %s\n", f.Name())
//   }
//
//   // Output:
//   // struct field name is: Property
//   // struct field name is: Nested
//
// The other purpose, is to return all struct filed names:
//
// Example:
//   names := s.Fields().Names()
//   fmt.Printf("struct field names are: %v\n", names)
//
//   // Output:
//   // struct field names are: [Property Nested]
//
// StructRows
//
// The StructRows object represents the slice of structs and mostly follows the
// database/sql API. Its main purpose is to loop through elements of a submitted
// slice of structs. Each of those elements can then be manipulated in the same
// manner as the StructValue.
//
// Example:
//   s, err := structs.New([]*T{&t})
//   if err != nil {
//      return err
//   }
//   rows, err := s.Rows()
//   if err != nil {
//      return err
//   }
//   defer rows.Close()
//   for rows.Next() {
//      f := rows.Field("Property")
//      f.Set(f.String() + "s")        // adds an "s" to the Property field
//      fmt.Printf("%s: %s.\n", f.Name(), f.String())
//   }
//
//   // Output:
//   // Property: Tests.
//
// Helper functions
//
// The helper methods in the helpers.go provide advanced functions for transforming
// structs or wrap StructValue object functionalities.
//
// Examples:
//   clone, err := structs.Clone(&t)
//   if err != nil {
//      return err
//   }
//   t2, ok := clone.(*T)
//   if ok {
//      same := structs.Compare(t2, &t)
//      if same {
//         fmt.Println("t and t2 are the same")
//      }
//      t2.Property = "Cloned"
//      if err != nil {
//         return err
//      }
//      same = structs.Compare(t2, &t)
//      if !same {
//         fmt.Println("t and t2 are now different")
//      }
//   }
//
//   // Output:
//   // t and t2 are the same
//   // t and t2 are now different
//
package structs

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

/*   S t r u c t   d e f i n i t i o n   */

// StructValue is the representation of a Go struct powered by the Go reflection package.
// Its interface provides field getters, field setters and much more.
type StructValue struct {
	value         reflect.Value           // Go value of struct via Go reflection.
	rows          reflect.Value           // Go slice of struct values via Go reflection.
	kinds         []reflect.Kind          // Lits of types that preceeds/including the struct.
	fieldsByIndex StructFields            // List of struct fields by index (not recursive).
	fieldsByName  map[string]*StructField // Map of struct fields by names (not recursive).
	Parent        *StructValue            // Parent struct, if nested struct.
	Error         error                   // Error added when struct could not be found.
}

/*   C o n s t r u c t o r   */

// New returns a new StructValue initialized to the struct concrete value
// stored in the interface dest. New(nil) returns the StructValue with an error.
//
// BUG(roninzo): the New method behaves unexpectidely when passing in an
// empty slice of pointers to structs.
func New(dest interface{}, parents ...*StructValue) (*StructValue, error) {
	// Start with zero-value of StructValue.
	s := &StructValue{}
	// Load parent StructValue, if any.
	if len(parents) > 0 {
		s.Parent = parents[0]
	}
	// Invalid nil argument
	if dest == nil {
		err := errors.Errorf(
			"invalid concrete value; want: '%s' or '%s' or '%s', got: 'nil'",
			reflect.Struct,
			reflect.Ptr,
			reflect.Slice,
		)
		return nil, err
	}
	// Find the struct or structs in the interface dest ...
	v := reflect.ValueOf(dest)
	t := reflect.TypeOf(dest)
	s.findStruct(v, t)
	//... or return failure.
	if err := s.Err(); err != nil {
		return nil, err
	}
	return s, nil
}

/*   I m p l e m e n t a t i o n   */

// IsValid returns true if the struct was found.
// It is so when StructValue represents a reflect value.
func (s *StructValue) IsValid() bool {
	if n := len(s.kinds); n > 0 {
		if s.kinds[n-1] == reflect.Struct {
			return true
		}
	}
	return false
}

// Name returns the name of the struct.
// When the struct was not found, it returns zero-value string.
func (s *StructValue) Name() string {
	return s.value.Type().Name()
}

// Namespace returns the same as the Name method, unless StructValue is a nested struct.
// When dealing with a nested struct, parent struct names are looked up and concatenated
// to the response recursively all the way to the top level struct.
func (s *StructValue) Namespace() string {
	n := ""
	p := s.Parent
	for {
		if p != nil {
			if n == "" {
				n = p.Name()
			} else {
				n = fmt.Sprintf("%s.%s", p.Name(), n)
			}
			p = p.Parent
		} else {
			break
		}
	}
	if n == "" {
		return s.Name()
	}
	return fmt.Sprintf("%s.%s", n, s.Name())
}

// Kind returns the struct reflect kind, or the last kind identified when the struct could
// not be found.
func (s *StructValue) Kind() reflect.Kind {
	return reflect.Struct
}

// Type returns the type struct name, including the name of package,
// such as "structs.T".
func (s *StructValue) Type() reflect.Type {
	return s.value.Type()
}

// Value returns the reflect value of the struct when the struct was found, else it returns
// zero-value reflect value.
func (s *StructValue) Value() (v reflect.Value) {
	return s.value
}

// Values returns the values of the struct as a slice of interfaces recursively.
// Unexported struct fields will be neglected.
func (s *StructValue) Values() (values []reflect.Value) {
	for _, f := range s.Fields() {
		if f.IsExported() {
			if f.CanStruct() {
				values = append(values, f.Struct().Values()...)
			} else {
				values = append(values, f.Value())
			}
		}
	}
	return values
}

// PtrValues returns the values of the struct as a slice of interfaces recursively.
// Unexported struct fields will be neglected.
func (s *StructValue) PtrValues() (values []reflect.Value) {
	for _, f := range s.Fields() {
		if f.IsExported() {
			if f.CanStruct() {
				values = append(values, f.Struct().PtrValues()...)
			} else {
				values = append(values, f.PtrValue())
			}
		}
	}
	return values
}

// Debug dumps the StructValue object itself as json string.
func (s *StructValue) Debug() string {
	var p string
	if s.Parent != nil {
		p = s.Parent.Debug()
	}
	fields := make(map[int]string)
	for i, f := range s.fieldsByIndex {
		fields[i] = f.Name()
	}
	d := struct {
		// Value interface{} `json:"value"`
		// Rows   interface{}    `json:"rows"`
		// MaxRow int            `json:"max_row"`
		Kinds  string         `json:"kinds"`
		Parent string         `json:"parent"`
		Fields map[int]string `json:"fields"`
		Error  error          `json:"error"`
	}{
		// Value: s.value.Interface(),
		// Rows:   s.rows.Interface(),
		// MaxRow: s.rows.Len(),
		Kinds:  Kinds(s.kinds...),
		Parent: p,
		Fields: fields,
		Error:  s.Error,
	}
	return Sprint(d)
}

// CanSet reports whether the value of StructValue can be changed.
// A StructValue can be changed only if it is addressable and was not obtained by the use of
// unexported struct fields. If CanSet returns false, calling Set or any type-specific setter
// (e.g., SetBool, SetInt) will panic.
func (s *StructValue) CanSet() bool {
	return s.value.CanSet()
}

// Multiple reports whether the value of StructValue is a slice of structs. If Multiple
// returns false, either the struct was not found or it was not part of a slice of them.
func (s *StructValue) Multiple() bool {
	return s.rows.IsValid()
}

// NumField returns the number of fields in the struct.
// This method is not recursive, which means that nested structs must be dealt with explicitly.
func (s *StructValue) NumField() int {
	if s.fieldsByIndex == nil {
		s.getFields()
	}
	return len(s.fieldsByIndex)
}

// FindStruct recursively finds and returns the StructValue object
// matching provided name, i.e.: the name of struct desired.
func (s *StructValue) FindStruct(name string) *StructValue {
	if s.Name() == name {
		return s
	}
	for _, f := range s.Fields() {
		if f.CanStruct() {
			nested := f.Struct()
			if found := nested.FindStruct(name); found != nil {
				return found
			}
		}
	}
	return nil
}

// IsZero returns true if all struct fields are of zero value.
// Unexported struct fields will be neglected.
func (s *StructValue) IsZero() bool {
	for _, f := range s.Fields() {
		if f.IsExported() {
			if f.CanStruct() {
				if !f.Struct().IsZero() {
					return false
				}
			} else if !f.IsZero() {
				return false
			}
		}
	}
	return true
}

// HasZero returns true if one or more struct fields are of zero value.
// Unexported struct fields will be neglected.
func (s *StructValue) HasZero() bool {
	for _, f := range s.Fields() {
		if f.CanStruct() {
			if f.Struct().HasZero() {
				return true
			}
		} else if f.IsZero() {
			return true
		}
	}
	return false
}

// IsNested returns true if struct is a nested struct within the root struct.
// IsNested returns false if StructValue is the top level struct.
func (s *StructValue) IsNested() bool {
	return s.Parent != nil
}

// HasNested returns true if struct has one or more nested struct within the
// root struct. HasNested returns false if none of the fields is sub-struct.
func (s *StructValue) HasNested() bool {
	for _, f := range s.Fields() {
		if f.CanStruct() {
			return true
		}
	}
	return false
}

// Err gets error from StructValue, then resets internal error.
func (s *StructValue) Err() (err error) {
	err, s.Error = s.Error, err // swap variable values
	return err
}

// Path returns a comma separated string of reflect.Kind.String describing where the
// struct was found inside the interface input of the New method.
func (s *StructValue) Path() string {
	return Kinds(s.kinds...)
}

// Sprint returns struct as a string, similar to the Values method, but in a json indented format.
// When the struct was not found, it returns zero-value string.
// Unexported struct fields will be neglected.
func (s *StructValue) Sprint() string {
	return Sprint(s.value.Interface())
}

// Contains returns index field of struct inside interface dest.
// Unexported struct fields will be neglected.
func (s *StructValue) Contains(dest interface{}) int {
	v := reflect.ValueOf(dest)
	return s.contains(v)
}

// // Contains reports whether value is within interface dest.
// func (s *StructValue) Contains(dest, value interface{}) bool {
// 	s, err := New(dest)
// 	if err == nil {
// 		v := reflect.ValueOf(value)
// 		if i := s.contains(v); i != OutOfRange {
// 			return true
// 		}
// 	}
// 	return false
// }

// // FieldNameByValue returns the field's name of the first instance of the
// // value in dest.
// func (s *StructValue) FieldNameByValue(dest, value interface{}) string {
// 	s, err := New(dest)
// 	if err == nil {
// 		v := reflect.ValueOf(value)
// 		if i := s.contains(v); i != OutOfRange {
// 			if f := s.Field(i); f != nil {
// 				return f.Name()
// 			}
// 		}
// 	}
// 	return ""
// }

// HasField returns true if struct dest has a field called the same as
// argument name.
func (s *StructValue) HasField(dest interface{}, arg interface{}) (bool, error) {
	f := s.Field(arg)
	err := s.Err()
	if err != nil {
		return false, err
	}
	if f != nil {
		return true, nil
	}
	return false, nil
}

// // Index returns the index of the first instance of the value in dest.
// func (s *StructValue) Index(dest, value interface{}) int {
// 	s, err := New(dest)
// 	if err == nil {
// 		v := reflect.ValueOf(value)
// 		if i := s.contains(v); i != OutOfRange {
// 			return i
// 		}
// 	}
// 	return OutOfRange
// }

// Import loops through destination fields of struct s and set their values to the
// corresponding fields from c. Usually, s is a trim-down version of c.
// Unsettable struct fields will be neglected.
func (s *StructValue) Import(c *StructValue) error {
	for _, field := range s.Fields() {
		if field.CanSet() {
			v := field.value
			f := c.Field(field.Name())
			if err := c.Err(); err != nil {
				return err
			}
			// format := "%s.Field: %s, Embedded: %t, Exported: %t, Settable: %t\n%+v\n"
			// fmt.Printf(format,
			// 	c.Name(),
			// 	f.Name(),
			// 	f.IsEmbedded(),
			// 	f.IsExported(),
			// 	f.CanSet(),
			// 	f.field)
			//
			// Organization.Name: GormModel, Embedded: true, Exported: true, Settable: true
			// {Name:GormModel PkgPath: Type:models.GormModel Tag: Offset:0 Index:[0] Anonymous:true}
			//
			// fmt.Printf(format,
			// 	s.Name(),
			// 	field.Name(),
			// 	field.IsEmbedded(),
			// 	field.IsExported(),
			// 	field.CanSet(),
			// 	field.field)
			//
			// CreateOrganizationRes.Name: ID, Embedded: false, Exported: true, Settable: true
			// {Name:ID PkgPath: Type:uint Tag:json:"id" example:"12" Offset:0 Index:[0] Anonymous:false}
			//
			x := f.Value()
			v.Set(x)
			//
			// panic: reflect.Set: value of type models.GormModel is not assignable to type uint
			//
		}
	}
	return nil
}

// Forward loops through destination fields of struct s and set their values to the
// corresponding fields from c. Zero-value fields from c will be neglected.
// Unsettable struct fields will be neglected.
func (s *StructValue) Forward(c *StructValue) error {
	for _, field := range s.Fields() {
		if field.CanSet() {
			v := field.value
			f := c.Field(field.Name())
			if err := c.Err(); err != nil {
				return err
			}
			if !f.IsZero() {
				x := f.Value()
				v.Set(x)
			}
		}
	}
	return nil
}

// MapFunc maps struct with func handler.
func (s *StructValue) MapFunc(handler func(reflect.Value) error) (*StructValue, error) {
	for _, f := range s.Fields() {
		if f.IsExported() {
			if f.CanStruct() {
				if _, err := f.Struct().MapFunc(handler); err != nil {
					return s, err
				}
			} else if f.CanSet() {
				if err := handler(f.value); err != nil {
					return s, err
				}
			}
		}
	}
	return s, nil
}

/*   U n e x p o r t e d   */

// findStruct finds where the struct is inside the reflect value and type.
// By being chainable, findStruct uses withStruct to finalize successful completion
// and exit the switch case below, all in a one-liner, otherwise it returns the
// StructValue explictly at the end.
// See the Support paragraph in the documentation for more details.
func (s *StructValue) findStruct(v reflect.Value, t reflect.Type) *StructValue {
	switch t.Kind() {
	case reflect.Struct:
		// Type interface
		//
		// // PtrTo returns the pointer type with element t.
		// // For example, if t represents type Foo, PtrTo(t) represents *Foo.
		// func PtrTo(t Type) Type {
		// 	return t.(*rtype).ptrTo()
		// }
		return s.withStruct(v, t)
	case reflect.Ptr:
		v, t = s.getElem(v, t)
		switch t.Kind() {
		case reflect.Struct:
			return s.withStruct(v, t)
		case reflect.Slice:
			v, t = s.getElem(v, t)
			if t.Kind() == reflect.Struct {
				return s.withStruct(v, t)
			}
		}
	case reflect.Slice:
		v, t = s.getElem(v, t)
		switch t.Kind() {
		case reflect.Struct:
			return s.withStruct(v, t)
		case reflect.Ptr:
			v, t = s.getElem(v, t)
			if t.Kind() == reflect.Struct {
				return s.withStruct(v, t)
			}
		}
	}
	v, t = s.getElem(v, t)
	return s
}

// withStruct returns StructValue object after setting reflect.Value of struct found.
// By being chainable, withStruct can finalize StructValue and return it at the same time.
func (s *StructValue) withStruct(v reflect.Value, t reflect.Type) *StructValue {
	s.appendKind(t)
	s.value = v
	return s
}

// getElem returns the element or the first element of the reflection pointer value and type.
// It saves an errors inside StructValue if the type's Kind is not Array, Ptr, or Slice.
func (s *StructValue) getElem(v reflect.Value, t reflect.Type) (rv reflect.Value, rt reflect.Type) {
	s.appendKind(t)
	if t.Kind() == reflect.Slice {
		s.rows = v
	}
	rv, rt, err := structValueElem(v, t)
	if err != nil {
		s.wrapErr(err)
	}
	return rv, rt
}

// appendKind appends reflect type t to the slice of kinds in StructValue.
// The first, it initialized kinds as a the slice of kinds.
func (s *StructValue) appendKind(t reflect.Type) {
	if s.kinds == nil {
		s.kinds = make([]reflect.Kind, 0)
	}
	s.kinds = append(s.kinds, t.Kind())
}

// getFields loads and saves all the struct fields.
// This method is not recursive, which means that nested structs must be dealt with explicitly.
//
// NOTE(roninzo): getFields indexes will not necessary follow the top level struct
// filed indexes. Indeed, if any anonymous/embedded struct are found in the immediate list of
// fields, those will be expanded and their related fields will explicitely be parsed, potentially
// pushing remaining struct fields further down the order.
func (s *StructValue) getFields() {
	if s.fieldsByIndex == nil {
		m := make(map[int]unembeddeds)
		explodeEmbedded(s.value, s.value.Type(), m, nil, nil, nil)
		n := len(m)
		s.initFields(n)
		for i := 0; i < n; i++ {
			s.loadField(m[i].Index, m[i].Indexes, m[i].Value, m[i].StructField)
		}
	}
}

// getFieldByIndex loads and saves the struct field indentified by index i. If an error occurred finding
// field, getFieldByIndex returns nil and error is saved in StructValue.
func (s *StructValue) getFieldByIndex(i int) *StructField {
	if s.fieldsByIndex == nil {
		s.getFields()
	}
	if OutOfRange < i && i < s.NumField() { // Try cache first
		return s.fieldsByIndex[i]
	}
	// DEPRECATED(roninzo):
	// if f, ok := s.Type().FieldByIndex(i); ok { // Lookup using Go reflection
	// 	return s.loadField(f.Index[0])
	// }
	s.setErrorf("invalid field index %d", i)
	return nil
}

// getFieldByName loads and saves the struct field indentified by name n. If an error occurred finding
// field, getFieldByName returns nil and error is saved in StructValue.
func (s *StructValue) getFieldByName(n string) *StructField {
	if s.fieldsByIndex == nil {
		s.getFields()
	}
	if f, ok := s.fieldsByName[n]; ok { // Try cache first
		return f
	}
	// DEPRECATED(roninzo):
	// if f, ok := s.Type().FieldByName(n); ok { // Lookup using Go reflection
	// 	return s.loadField(f.Index[0])
	// }
	s.setErrorf("invalid field name %s", n)
	return nil
}

// initFields initializes the struct fields attributes of StructValue.
func (s *StructValue) initFields(c ...int) {
	total := 0
	if len(c) > 0 {
		total = c[0]
	}
	s.fieldsByIndex = make(StructFields, total)
	s.fieldsByName = map[string]*StructField{}
}

// loadField loads, saves and returns the i'th struct field.
func (s *StructValue) loadField(i int, x []int, v reflect.Value, sf reflect.StructField) {
	f := &StructField{
		index:   i,
		indexes: x,
		value:   v,
		field:   sf,
		Parent:  s,
	}
	s.fieldsByIndex[i] = f
	s.fieldsByName[f.Name()] = f
}

// contains returns index field of struct inside interface dest.
// This method is not recursive, which means that nested structs must be dealt with explicitly.
// Unexported struct fields will be neglected.
func (s *StructValue) contains(v reflect.Value) int {
	for _, f := range s.Fields() {
		if f.IsExported() {
			if !f.CanStruct() {
				if i := f.equal(v); i != OutOfRange {
					return f.Index()
				}
			}
		}
	}
	return OutOfRange
}

// getRow returns the StructRows object, which is mainly used to loop through elements of the
// slice of structs. If s is not a slice of structs, nothing happens except saving an internal
// error.
func (s *StructValue) getRow(rownum int) error {
	if s.Multiple() {
		if n := s.rows.Len(); n > 0 {
			if OutOfRange < rownum && rownum < n {
				//
				// Update StructValue value
				s.value = s.rows.Index(rownum)
				if s.value.Kind() == reflect.Ptr {
					s.value = s.value.Elem()
				}
				//
				// Update StructField values
				for _, f := range s.fieldsByIndex {
					f.value = s.value.FieldByIndex(f.indexes)
				}
				return s.setErr(nil)
			}
			return s.setErr(ErrNoRow)
		}
		return s.setErr(ErrNoRows)
	}
	return s.setErr(ErrNoStructs)
}

// setErr sets error to StructValue.
func (s *StructValue) setErr(err error) error {
	s.Error = err
	return s.Error
}

// wrapErr appends error to StructValue.
func (s *StructValue) wrapErr(err error) error {
	if err != nil {
		if s.Error != nil {
			s.Error = errors.Wrap(s.Error, err.Error())
		} else {
			s.Error = err
		}
	} else {
		s.Error = nil
	}
	return s.Error
}

// setError sets mesg as error to StructValue.
func (s *StructValue) setError(mesg string) error {
	if mesg != "" {
		s.Error = errors.New(mesg)
	} else {
		s.Error = nil
	}
	return s.Error
}

// setErrors adds mesg as error to StructValue.
func (s *StructValue) setErrors(err error, mesg string) error {
	if err != nil {
		if mesg != "" {
			s.Error = errors.Wrap(err, mesg)
		}
	} else {
		s.Error = nil
	}
	return s.Error
}

// setErrorf formats mesg according to a format specifier and sets the resulting string
// as error to StructValue.
func (s *StructValue) setErrorf(format string, args ...interface{}) error {
	mesg := fmt.Sprintf(format, args...)
	return s.setError(mesg)
}

// setErrorsf formats mesg according to a format specifier and adds the resulting string
// as error to StructValue.
func (s *StructValue) setErrorsf(err error, format string, args ...interface{}) error {
	mesg := fmt.Sprintf(format, args...)
	return s.setErrors(err, mesg)
}

// destroy is the object destructor, applying zero-value to all its fields.
func (s *StructValue) destroy() error {
	s.value = reflect.Value{}
	s.rows = reflect.Value{}
	s.kinds = nil
	s.fieldsByIndex = nil
	s.fieldsByName = nil
	s.Parent = nil
	s.Error = nil
	return nil
}
