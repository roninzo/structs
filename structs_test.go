package structs

import (
	"reflect"
	"testing"

	"github.com/roninzo/structs/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewStruct(t *testing.T) {
	type T1 struct {
		A string
		B int
		C bool
	}

	t1 := T1{
		A: "test",
		B: 5,
		C: false,
	}
	s, err := New(t1)
	assert.Equal(t, nil, err)
	assert.Equal(t, utils.Kinds(reflect.Struct), s.Path())
	assert.Equal(t, reflect.Struct.String(), s.Kind().String())
	assert.Equal(t, 3, s.NumField())
	assert.Equal(t, false, s.HasNested())
	assert.Equal(t, "T1", s.Name())
	assert.Equal(t, "T1", s.FullName())
	assert.Equal(t, false, s.IsNested())
	assert.Equal(t, false, s.IsZero())
	assert.Equal(t, true, s.HasZero())
	assert.Equal(t, false, s.CanSet())
	assert.Equal(t, false, s.Multiple())
	assert.Equal(t, []string{"A", "B", "C"}, s.Fields().Names())
}

func TestNewNil(t *testing.T) {
	s, err := New(nil)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "invalid concrete value; want: \"struct\" or \"ptr\" or \"slice\", got: <nil>", err.Error())
	if s != nil {
		t.Error("s is not nil")
	}
}

func TestNewPtrToStruct(t *testing.T) {
	type T1 struct {
		A string
		B int
		C bool
	}

	t1 := T1{
		A: "test",
		B: 5,
		C: false,
	}
	s, err := New(&t1)
	assert.Equal(t, nil, err)
	assert.Equal(t, utils.Kinds(reflect.Ptr, reflect.Struct), s.Path())
	assert.Equal(t, "T1", s.Name())
	assert.Equal(t, "T1", s.FullName())
	assert.Equal(t, false, s.IsNested())
	assert.Equal(t, false, s.IsZero())
	assert.Equal(t, true, s.HasZero())
	assert.Equal(t, true, s.CanSet())
	assert.Equal(t, false, s.Multiple())
	assert.Equal(t, []string{"A", "B", "C"}, s.Fields().Names())
}

func TestNewPtrToSliceOfStruct(t *testing.T) {
	type T1 struct {
		A string
		B int
		C bool
	}

	t1 := T1{
		A: "test",
		B: 5,
		C: false,
	}
	second := []T1{t1}
	s, err := New(&second)
	assert.Equal(t, nil, err)
	assert.Equal(t, utils.Kinds(reflect.Ptr, reflect.Slice, reflect.Struct), s.Path())
	assert.Equal(t, "T1", s.Name())
	assert.Equal(t, "T1", s.FullName())
	assert.Equal(t, false, s.IsNested())
	assert.Equal(t, false, s.IsZero())
	assert.Equal(t, true, s.HasZero())
	assert.Equal(t, true, s.CanSet())
	assert.Equal(t, true, s.Multiple())
	assert.Equal(t, []string{"A", "B", "C"}, s.Fields().Names())
}

func TestNewSliceOfStruct(t *testing.T) {
	type T1 struct {
		A string
		B int
		C bool
	}

	t1 := T1{
		A: "test",
		B: 5,
		C: false,
	}
	second := []T1{t1}
	s, err := New(second)
	assert.Equal(t, nil, err)
	assert.Equal(t, utils.Kinds(reflect.Slice, reflect.Struct), s.Path())
	assert.Equal(t, "T1", s.Name())
	assert.Equal(t, "T1", s.FullName())
	assert.Equal(t, false, s.IsNested())
	assert.Equal(t, false, s.IsZero())
	assert.Equal(t, true, s.HasZero())
	assert.Equal(t, true, s.CanSet())
	assert.Equal(t, true, s.Multiple())
	assert.Equal(t, []string{"A", "B", "C"}, s.Fields().Names())
}

func TestNewSliceOfPtrToStruct(t *testing.T) {
	type T1 struct {
		A string
		B int
		C bool
	}

	t1 := T1{
		A: "test",
		B: 5,
		C: false,
	}
	second := []*T1{&t1}
	s, err := New(second)
	assert.Equal(t, nil, err)
	assert.Equal(t, utils.Kinds(reflect.Slice, reflect.Ptr, reflect.Struct), s.Path())
	assert.Equal(t, "T1", s.Name())
	assert.Equal(t, "T1", s.FullName())
	assert.Equal(t, false, s.IsNested())
	assert.Equal(t, false, s.IsZero())
	assert.Equal(t, true, s.HasZero())
	assert.Equal(t, true, s.CanSet())
	assert.Equal(t, true, s.Multiple())
	assert.Equal(t, []string{"A", "B", "C"}, s.Fields().Names())
}

func TestStruct(t *testing.T) {
	type T2 struct {
		A string
		C bool
	}

	type T3 struct {
		ID        int
		X         *T2
		OrgID     int
		Enabled   bool
		UpdatedBy string
	}

	t3 := T3{
		ID:      555,
		Enabled: true,
		X: &T2{
			A: "testing",
			C: false,
		},
	}
	s, err := New(&t3)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(555), s.FindStruct("T3").Field("ID").Int())
	assert.Equal(t, "testing", s.FindStruct("T2").Field("A").String())
}

func TestFieldNil(t *testing.T) {
	type T2 struct {
		A string
		C bool
	}

	t2 := T2{
		A: "testing",
		C: false,
	}

	s, err := New(&t2)
	assert.Equal(t, nil, err)

	f := s.Field(nil)
	err = s.Err()
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "invalid nil argument", err.Error())
	assert.NotEqual(t, nil, f)
}

func TestNewUnsupported(t *testing.T) {
	type T1 struct {
		A string
		B int
		C bool
	}

	m1 := map[string]T1{
		"subtesting": {
			A: "testing",
			B: 44,
			C: false,
		},
	}

	_, err := New(&m1)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "\"map\" is an unsupported pointer to a struct", err.Error())
}

func TestNestedStruct(t *testing.T) {
	type T1 struct {
		A string
		B int
		C bool
	}

	type T2 struct {
		A string
		B *T1
		C bool
	}

	type T3 struct {
		ID        int
		X         *T2
		OrgID     int
		Enabled   bool
		UpdatedBy string
	}

	t1 := T1{
		A: "test1",
		B: 123456,
		C: true,
	}

	t2 := T2{
		A: "test2",
		B: &t1,
		C: true,
	}

	t3 := T3{
		ID:        123456,
		X:         &t2,
		Enabled:   true,
		UpdatedBy: "roninzo",
	}

	s3, err := New(&t3)
	s2 := s3.FindStruct("T2")
	s1 := s3.FindStruct("T1")
	assert.Equal(t, nil, err)
	assert.Equal(t, "T3", s3.FullName())
	assert.Equal(t, "T3.T2", s2.FullName())
	assert.Equal(t, "T3.T2.T1", s1.FullName())
	// T3
	assert.Equal(t, false, s3.IsZero())
	assert.Equal(t, true, s3.HasZero())
	assert.Equal(t, int64(123456), s3.Field("ID").Int())
	assert.Equal(t, true, s3.Field("Enabled").Bool())
	assert.Equal(t, "roninzo", s3.Field("UpdatedBy").String())
	// T2
	assert.Equal(t, "test2", s2.Field("A").String())
	assert.Equal(t, true, s2.Field("C").Bool())
	// T1
	assert.Equal(t, "test1", s1.Field("A").String())
	assert.Equal(t, int64(123456), s1.Field("B").Int())
	assert.Equal(t, true, s1.Field("C").Bool())
}
