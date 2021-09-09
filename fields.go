// Copyright 2021 Roninzo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package structs

/*   T y p e   d e f i n i t i o n   */

// StructFields represents all struct fields that encapsulates high level functions around
// the struct fields.
type StructFields []*StructField

/*   C o n s t r u c t o r   */

// Fields returns all the fields of the struct in a slice.
// This method is not recursive, which means that nested structs must be dealt with explicitly.
func (s *StructValue) Fields() StructFields {
	s.getFields()
	return s.fieldsByIndex
}

/*   I m p l e m e n t a t i o n   */

// Names returns all the field names of the struct.
// This method is not recursive, which means that nested structs must be dealt with explicitly.
func (fields StructFields) Names() []string {
	n := len(fields)
	names := make([]string, n)
	for i, f := range fields {
		names[i] = f.Name()
	}
	return names
}

// Parent returns the related StructValue object (which is a level above StructFields).
func (fields StructFields) Parent() *StructValue {
	return fields[0].Parent
}
