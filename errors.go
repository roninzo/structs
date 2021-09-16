// Copyright 2021 Roninzo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package structs

import (
	"errors"
)

var (
	ErrNotExported = errors.New("struct field is not exported")
	ErrNotSettable = errors.New("struct field is not settable")
	ErrNotNillable = errors.New("struct field is not nillable")
	ErrNoStruct    = errors.New("struct not found")
	ErrNoStructs   = errors.New("structs not found")
	ErrNoField     = errors.New("struct field not found")
	ErrNoFields    = errors.New("struct fields not found")
	ErrNoRow       = errors.New("struct row not found")
	ErrNoRows      = errors.New("struct rows not found")
	ErrRowsClosed  = errors.New("struct rows are closed")
	ErrNotReplaced = errors.New("struct field old and new value types does not match") // could not replace value in struct
)
