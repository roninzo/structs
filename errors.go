// Copyright 2021 Roninzo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package structs

import (
	"errors"
)

var (
	errNotExported = errors.New("struct field is not exported")
	errNotSettable = errors.New("struct field is not settable")
	errNotNillable = errors.New("struct field is not nillable")
	errNoStruct    = errors.New("struct not found")
	errNoStructs   = errors.New("structs not found")
	errNoField     = errors.New("struct field not found")
	errNoFields    = errors.New("struct fields not found")
	ErrNoRow       = errors.New("struct row not found")
	ErrNoRows      = errors.New("struct rows not found")
	errRowsClosed  = errors.New("struct rows are closed")
	errNotReplaced = errors.New("struct field old and new value types does not match") // could not replace value in struct
)
