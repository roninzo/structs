// Copyright 2021 Roninzo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package structs

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

/*   F u n c t i o n s   */

// Name returns the structs's type name within its package. It returns an
// empty string for unnamed types. It returns an error if s's kind is
// not struct.
func Name(dest interface{}) (string, error) {
	s, err := New(dest)
	if err != nil {
		return "", err
	}
	return s.Name(), err
}

// Names returns a slice of field names. For more info refer to StructValue
// types Names() method. It returns an error if s's kind is not struct.
func Names(dest interface{}) ([]string, error) {
	s, err := New(dest)
	if err != nil {
		return nil, err
	}
	return s.Fields().Names(), nil
}

// Fields returns a slice of *StructField. For more info refer to StructValue
// types Fields() method. It returns an error if s's kind is not struct.
func Fields(dest interface{}) (StructFields, error) {
	s, err := New(dest)
	if err != nil {
		return nil, err
	}
	return s.Fields(), nil
}

// Copy makes a copy, for instance, of a struct pointer.
// Do the copy "manually", e.g. create a new struct and copy the fields,
// where pointers or slices/maps/channels/etc must be duplicated manually,
// in a recursive manner.
//
// Resources:
// https://stackoverflow.com/questions/50269322/how-to-copy-struct-and-dereference-all-pointers
func Copy(dest, src interface{}) error {
	// ctx will be the context error returned
	// by this func if anything goes wrong
	ctx := "could not copy data between two structs"
	//
	// Both interfaces must be valid structs for this to work
	s1, err := New(src)
	if err != nil {
		return errors.Wrap(err, ctx)
	}
	s2, err := New(dest)
	if err != nil {
		return errors.Wrap(err, ctx)
	}
	//
	// ctx content can now now be improved
	ctx = fmt.Sprintf("could not copy data between two '%s' structs", s1.Name())
	//
	// StructValue names must be the same
	if s1.Name() != s2.Name() {
		return errors.Wrap(errors.Errorf("target struct name is invalid: want: '%s', got: '%s'", s1.Name(), s2.Name()), ctx)
	}
	//
	// Target struct must be editable
	if !s2.CanSet() {
		return errors.Wrap(errors.Errorf("cannot edit struct %s", s2.Name()), ctx)
	}
	//
	// Both interfaces must be singulars of struct, not multiples
	if s1.Multiple() {
		return errors.Wrap(errors.Errorf("source is a slice of struct %s", s1.Name()), ctx)
	}
	if s2.Multiple() {
		return errors.Wrap(errors.Errorf("target is a slice of struct %s", s2.Name()), ctx)
	}
	return s2.Import(s1)
}

// Clone returns a copy from a struct out of nothing.
func Clone(src interface{}) (interface{}, error) {
	// ctx will be the context error returned
	// by this func if anything goes wrong
	ctx := "could not clone struct"
	//
	// Source interface must be valid struct for this to work
	// Target dest will be the recipient for a copy of src
	s1, err := New(src)
	if err != nil {
		return nil, errors.Wrap(err, ctx)
	}
	t := s1.Type()
	dest := reflect.New(t).Interface()
	s2, err := New(dest)
	if err != nil {
		return nil, errors.Wrap(err, ctx)
	}
	//
	// ctx content can now now be improved
	ctx = fmt.Sprintf("could not clone struct '%s'", s1.Name())
	//
	// Target struct must be editable
	if !s2.CanSet() {
		return nil, errors.Wrap(errors.Errorf("cannot edit struct %s", s2.Name()), ctx)
	}
	//
	// Target must be singular of struct, not multiple
	if s2.Multiple() {
		return nil, errors.Wrap(errors.Errorf("source is a slice of struct %s", s2.Name()), ctx)
	}
	err = s2.Import(s1)
	return dest, err
}

// Transpose loops through target fields and set value of its related
// source field.
func Transpose(dest, src interface{}) error {
	// ctx will be the context error returned
	// by this func if anything goes wrong
	ctx := "could not transpose data between two structs"
	//
	// Both interfaces must be valid structs for this to work
	s1, err := New(src)
	if err != nil {
		return errors.Wrap(err, ctx)
	}
	s2, err := New(dest)
	if err != nil {
		return errors.Wrap(err, ctx)
	}
	//
	// ctx content can now now be improved
	ctx = fmt.Sprintf("could not transpose data between '%s' and '%s' structs", s1.Name(), s2.Name())
	//
	// StructValue names must be diffenent (otherwise, just use Copy)
	if s1.Name() == s2.Name() {
		return errors.Wrap(errors.Errorf("target struct name should not be %s", s1.Name()), ctx)
	}
	//
	// Target struct must be editable
	if !s2.CanSet() {
		return errors.Wrap(errors.Errorf("cannot edit struct %s", s2.Name()), ctx)
	}
	//
	// Both interfaces must be singulars of struct, not multiples
	if s1.Multiple() {
		return errors.Wrap(errors.Errorf("source is a slice of struct %s", s1.Name()), ctx)
	}
	if s2.Multiple() {
		return errors.Wrapf(errors.Errorf("target is a slice of struct %s", s2.Name()), ctx)
	}
	return s2.Import(s1)
}

// Forward copies only non-zero values between two structs, i.e. from src to dest interface.
func Forward(dest, src interface{}) error {
	// ctx will be the context error returned
	// by this func if anything goes wrong
	ctx := "could not copy source non-zero values to target struct"
	//
	// Both interfaces must be valid structs for this to work
	s1, err := New(src)
	if err != nil {
		return errors.Wrap(err, ctx)
	}
	s2, err := New(dest)
	if err != nil {
		return errors.Wrap(err, ctx)
	}
	//
	// ctx content can now now be improved
	ctx = fmt.Sprintf("could not copy source non-zero values to target struct '%s'", s2.Name())
	//
	// StructValue names must be the same
	if s1.Name() != s2.Name() {
		return errors.Wrap(errors.Errorf("target struct name is invalid: want: '%s', got: '%s'", s1.Name(), s2.Name()), ctx)
	}
	//
	// Target struct must be editable
	if !s2.CanSet() {
		return errors.Wrap(errors.Errorf("cannot edit struct %s", s2.Name()), ctx)
	}
	//
	// Both interfaces must be singulars of struct, not multiples
	if s1.Multiple() {
		return errors.Wrap(errors.Errorf("source is a slice of struct %s", s1.Name()), ctx)
	}
	if s2.Multiple() {
		return errors.Wrap(errors.Errorf("target is a slice of struct %s", s2.Name()), ctx)
	}
	return s2.Forward(s1)
}

// Compare returns dest boolean comparing two struct.
func Compare(dest, src interface{}) bool {
	return reflect.DeepEqual(dest, src)
}

// Replace returns a copy of the struct dest with the first n non-overlapping
// instance of old replaced by new.
//
// Counts how many replacing to do until n. if n = -1, then replace all.
func Replace(dest, old, new interface{}, n int) (interface{}, error) {
	// ctx will be the context error returned
	// by this func if anything goes wrong
	ctx := errNotReplaced.Error()
	src, err := Clone(dest)
	if err != nil {
		return nil, errors.Wrap(err, ctx)
	}
	s, err := New(src)
	if err != nil {
		return nil, errors.Wrap(err, ctx)
	}
	c := 0
	v := reflect.ValueOf(old)
	for {
		if n != ReplaceAll {
			if c >= n {
				break
			}
		}
		if i := s.contains(v); i != OutOfRange {
			f := s.Field(i)
			if f == nil {
				return src, errors.New(ctx)
			}
			if err := f.Set(new); err != nil {
				return src, errors.Wrap(err, ctx)
			}
		} else {
			break
		}
		c++
	}
	return src, nil
}

// MapFunc returns a copy of the StructValue s with all its fields modified
// according to the mapping function handler.
//
//    If mapping returns a negative value, the character is
//    dropped from the byte slice with no replacement. The characters in s and the
//    output are interpreted as UTF-8-encoded code points.
//
// BUG(roninzo): the MapFunc method argument dest is also changed. should
// that be the case?
func MapFunc(dest interface{}, handler func(reflect.Value) error) (interface{}, error) {
	// ctx will be the context error returned
	// by this func if anything goes wrong
	ctx := "could not map struct with func"
	clone, err := Clone(dest)
	if err != nil {
		return nil, errors.Wrap(err, ctx)
	}
	s, err := New(clone)
	if err != nil {
		return nil, errors.Wrap(err, ctx)
	}
	if _, err := s.MapFunc(handler); err != nil {
		return nil, errors.Wrap(err, ctx)
	}
	return clone, nil
}