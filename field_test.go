package structs

import (
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestField(t *testing.T) {
	type T1 struct {
		A string
		B int
		C bool
	}

	t1 := T1{
		A: "test",
		B: 5,
		C: true,
	}

	s, err := New(&t1)
	assert.Equal(t, nil, err)

	f := s.Field("A")
	assert.NotEqual(t, nil, f)
	assert.Equal(t, nil, s.Err())
	assert.Equal(t, "A", f.Name())
	assert.Equal(t, "test", f.String())

	f = s.Field("B")
	assert.NotEqual(t, nil, f)
	assert.Equal(t, nil, s.Err())
	assert.Equal(t, "B", f.Name())
	assert.Equal(t, 5, int(f.Int()))

	f = s.Field("C")
	assert.NotEqual(t, nil, f)
	assert.Equal(t, nil, s.Err())
	assert.Equal(t, "C", f.Name())
	assert.Equal(t, true, f.Bool())

	f = s.Field(0)
	assert.NotEqual(t, nil, f)
	assert.Equal(t, nil, s.Err())
	assert.Equal(t, "A", f.Name())
	assert.Equal(t, "test", f.String())

	f = s.Field(1)
	assert.NotEqual(t, nil, f)
	assert.Equal(t, nil, s.Err())
	assert.Equal(t, "B", f.Name())
	assert.Equal(t, 5, int(f.Int()))

	f = s.Field(2)
	assert.NotEqual(t, nil, f)
	assert.Equal(t, nil, s.Err())
	assert.Equal(t, "C", f.Name())
	assert.Equal(t, true, f.Bool())
	assert.Equal(t, 2, f.Index())
	assert.Equal(t, "bool", f.Kind().String())
	assert.Equal(t, true, f.CanSet())
	assert.Equal(t, false, f.IsAnonymous())
	assert.Equal(t, true, f.IsExported())
	assert.Equal(t, "T1.C", f.Namespace())
	assert.Equal(t, false, f.CanDuration())
	assert.Equal(t, false, f.CanTime())
}

func TestFieldInvalid(t *testing.T) {
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

	f := s.Field("D") // does not exist!
	err = s.Err()
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "invalid field name D", err.Error())
	if f != nil {
		t.Error("f is not nil!")
	}
}

func TestFieldInvalidArg(t *testing.T) {
	type T2 struct {
		String string
		PtrInt *int
	}

	type T1 struct {
		A string
		B *int
		C bool
		D time.Duration
		E error
		T time.Time
		V reflect.Value
		T2
	}

	t1 := T1{
		"test",
		nil,
		false,
		5 * time.Minute,
		errors.New("test error"),
		time.Date(2021, time.August, 3, 13, 59, 35, 0, time.UTC),
		reflect.ValueOf(time.Now()),
		T2{"t2 string", nil},
	}

	s, err := New(&t1)
	assert.Equal(t, nil, err)

	f := s.Field("A")
	i, err := f.Get()
	assert.Equal(t, nil, err)
	assert.Equal(t, "test", i.(string))

	err = f.SetZero()
	assert.Equal(t, nil, err)
	assert.Equal(t, "", f.String())

	f = s.Field("B")
	assert.Equal(t, true, f.IsNil())

	f = s.Field("D")
	assert.Equal(t, true, f.CanDuration())
	assert.Equal(t, 5*time.Minute, f.Duration())

	f = s.Field("E")
	assert.Equal(t, true, f.CanError())
	assert.Equal(t, "test error", f.Error().Error())

	f = s.Field("T")
	mysqlFormat := "2006-01-02 15:04:05"
	assert.Equal(t, true, f.CanTime())
	assert.Equal(t, "2021-08-03 13:59:35", f.Time().Format(mysqlFormat))

	f = s.Field("V")
	assert.Equal(t, true, f.IsValid())
	assert.Equal(t, false, f.IsEmbedded())
	assert.Equal(t, reflect.Struct, f.Kind())

	f = s.Field(7) // T2.String
	assert.Equal(t, true, f.IsValid())
	assert.Equal(t, false, f.IsEmbedded())
	assert.Equal(t, "String", f.Name())

	f = s.Field(true)
	err = s.Err()
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "invalid argument type; want: 'string' or 'int', got: 'bool'", err.Error())
}
