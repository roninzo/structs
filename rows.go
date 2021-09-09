// Copyright 2021 Roninzo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package structs

/*   S t r u c t   d e f i n i t i o n   */

// StructRows represents a single row of a struct from a StructValue containing a slice
// of structs. If StructValue does not contain a slice of structs, StructRows cannot be
// initialized by contructor Rows. StructRows encapsulates high level functions around
// the element of slice of structs.
type StructRows struct {
	rownum      int // index of the slice of structs.
	StructValue     // embedded copy and inherits all fields and methods.
}

/*   C o n s t r u c t o r   */

// Rows returns an iterator, for a slice of structs.
// Rows returns nil if there was an error.
func (s *StructValue) Rows() (*StructRows, error) {
	if s.Multiple() {
		if s.rows.Len() > 0 {
			return &StructRows{OutOfRange, *s}, nil
		}
		return nil, ErrNoRows
	}
	return nil, errNoStructs
}

/*   I m p l e m e n t a t i o n   */

// Index returns the index element in the slice of structs pointing to current struct.
// Index returns OutOfRange, i.e. -1, if the rows are closed.
func (r *StructRows) Index() int {
	if !r.isClosed() {
		return r.rownum
	}
	return OutOfRange
}

// Len returns the number elements in the slice of structs.
// Len returns OutOfRange, i.e. -1, if the rows are closed.
func (r *StructRows) Len() int {
	if !r.isClosed() {
		return r.rows.Len()
	}
	return OutOfRange
}

// MaxRow returns the index of the lasr elements in the slice of structs.
// MaxRow returns OutOfRange, i.e. -1, if the rows are closed.
func (r *StructRows) MaxRow() int {
	if !r.isClosed() {
		return r.Len() - 1
	}
	return OutOfRange
}

// Columns returns the current struct field names.
// Columns returns an error if the rows are closed.
func (r *StructRows) Columns() ([]string, error) {
	if !r.isClosed() {
		return r.Fields().Names(), nil
	}
	return nil, errRowsClosed
}

// Next prepares the next result row for reading an element from the slice of struct.
// It returns true on success, or false if there is no next result row or an error
// happened while preparing it. Err should be consulted to distinguish between
// the two cases.
func (r *StructRows) Next() bool {
	if !r.isClosed() {
		if i := r.rownum + 1; i < r.Len() {
			err := r.getRow(i)
			if err == nil {
				r.rownum = i // confirm new row number
				return true
			}
		}
	}
	return false
}

// Err returns the error, if any, that was encountered during iteration.
// Err may be called after an explicit or implicit Close.
func (r *StructRows) Err() (err error) {
	if r.Error != nil {
		err, r.Error = r.Error, err // swap variable values
	}
	return err
}

// Close closes the Rows, preventing further enumeration. If Next is called
// and returns false and there are no further result rows,
// the Rows are closed automatically and it will suffice to check the
// result of Err. Close is idempotent and does not affect the result of Err.
func (r *StructRows) Close() error {
	return r.destroy()
}

/*   U n e x p o r t e d   */

// isClosed returns true if r is not closed and false if it is.
// Closure prevents further enumeration of StructRows.
func (r *StructRows) isClosed() bool {
	return !r.IsValid()
}
