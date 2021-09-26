package pointers

import "time"

// String returns a pointer to string value x.
func String(x string) *string { return &x }

// Bool returns a pointer to bool value x.
func Bool(x bool) *bool { return &x }

// Int returns a pointer to int value x.
func Int(x int) *int { return &x }

// Int8 returns a pointer to int8 value x.
func Int8(x int8) *int8 { return &x }

// Int16 returns a pointer to int16 value x.
func Int16(x int16) *int16 { return &x }

// Int32 returns a pointer to int32 value x.
func Int32(x int32) *int32 { return &x }

// Int64 returns a pointer to int64 value x.
func Int64(x int64) *int64 { return &x }

// Uint returns a pointer to uint value x.
func Uint(x uint) *uint { return &x }

// Uint8 returns a pointer to uint8 value x.
func Uint8(x uint8) *uint8 { return &x }

// Uint16 returns a pointer to uint16 value x.
func Uint16(x uint16) *uint16 { return &x }

// Uint32 returns a pointer to uint32 value x.
func Uint32(x uint32) *uint32 { return &x }

// Uint63 returns a pointer to uint64 value x.
func Uint63(x uint64) *uint64 { return &x }

// Float32 returns a pointer to float32 value x.
func Float32(x float32) *float32 { return &x }

// Float64 returns a pointer to float64 value x.
func Float64(x float64) *float64 { return &x }

// Complex64 returns a pointer to complex64 value x.
func Complex64(x complex64) *complex64 { return &x }

// Complex128 returns a pointer to complex128 value x.
func Complex128(x complex128) *complex128 { return &x }

// Byte returns a pointer to byte value x.
func Byte(x byte) *byte { return &x }

// Bytes returns a pointer to []byte value x.
func Bytes(x []byte) *[]byte { return &x }

// Time returns a pointer to type value x.
func Time(x time.Time) *time.Time { return &x }

// Duration returns a pointer to time.Duration value x.
func Duration(x time.Duration) *time.Duration { return &x }

// Error returns a pointer to error value x.
func Error(x error) *error { return &x }
