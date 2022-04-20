package structs

import (
	"reflect"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestHelperName(t *testing.T) {
	want := "T2"

	type T2 struct {
		A string
		C bool
	}

	got, err := Name(T2{})
	assert.Equal(t, nil, err)
	assert.Equal(t, want, got)
}

func TestHelperNames(t *testing.T) {
	type T2 struct {
		A string
		C bool
	}

	t2 := T2{
		A: "test",
		C: true,
	}

	want := []string([]string{"A", "C"})
	got, err := Names(&t2)
	assert.Equal(t, nil, err)
	assert.Equal(t, want, got)
}

func TestHelperCopy(t *testing.T) { // To Improve
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

	t2 := T1{}
	assert.Equal(t, nil, Copy(&t2, &t1))
	assert.Equal(t, "test", t2.A)
	assert.Equal(t, true, t2.C)

	assert.NotEqual(t, nil, Copy(nil, nil))
	assert.NotEqual(t, nil, Copy(nil, t1))
}

func TestHelperTranspose(t *testing.T) { // To Improve
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

	type T2 struct {
		A string
		C bool
	}

	t2 := T2{}
	assert.Equal(t, nil, Transpose(&t2, &t1))
	assert.Equal(t, "test", t2.A)
	assert.Equal(t, true, t2.C)

	assert.NotEqual(t, nil, Transpose(nil, nil))
	assert.NotEqual(t, nil, Transpose(nil, t1))

	// r := make([]T2, 0)
	// r = append(r, t2)
	// s := make([][]T2, 0)
	// s = append(s, r)
	// fmt.Printf("type of: %+v\n", reflect.TypeOf(s))
	// fmt.Printf("elem of: %+v\n", reflect.TypeOf(s).Elem())
	// fmt.Printf("kind of: %+v\n", reflect.TypeOf(s).Elem().Kind())
	// fmt.Printf("type of pointer: %+v\n", reflect.TypeOf(&s))
	// assert.NotEqual(t, nil, Transpose(&s, t1))

	return
}

func TestHelperClone(t *testing.T) {
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
	intf, err := Clone(&t1)
	assert.Equal(t, nil, err)
	switch clone := intf.(type) {
	case *T1:
		assert.Equal(t, "test", clone.A)
		assert.Equal(t, 5, clone.B)
		assert.Equal(t, true, clone.C)
	default:
		v := reflect.TypeOf(intf)
		t.Errorf("invalid clone type; want: 'T1', got: %q", v.Kind())
	}
}

func TestHelperForward(t *testing.T) {
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
	forward := T1{
		C: true,
	}
	err := Forward(&forward, &t1)
	assert.Equal(t, nil, err)
	assert.Equal(t, "test", forward.A) // transfered from t1
	assert.Equal(t, 5, forward.B)      // transfered from t1
	assert.Equal(t, true, forward.C)   // not transfered from t1
}

func TestHelperForwardInvalidArgs(t *testing.T) {
	s := 5
	err := Forward(nil, &s)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "could not copy source non-zero values to target struct: \"int\" is not a pointer", err.Error())
}

func TestHelperForwardWithNestedStruct(t *testing.T) {
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

	forward := T3{
		OrgID: 6,
		X: &T2{
			A: "testing",
			C: true,
		},
		UpdatedBy: "roninzo",
	}

	t3 := T3{
		ID:      555,
		OrgID:   4,
		Enabled: true,
		X: &T2{
			A: "decom",
		},
		UpdatedBy: "abcdef",
	}

	err := Forward(&t3, &forward)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, t3.Enabled)
	assert.Equal(t, 555, t3.ID)
	assert.Equal(t, 6, t3.OrgID)
	assert.Equal(t, "testing", t3.X.A)
	assert.Equal(t, true, t3.X.C)
}

func TestHelperForwardWithDirectNestedStruct(t *testing.T) {
	type T2 struct {
		A string
		C bool
	}

	type T4 struct {
		ID        int
		X         T2
		OrgID     int
		Enabled   bool
		UpdatedBy string
	}

	var t4, forward T4
	t4 = T4{
		ID:      555,
		Enabled: true,
		X: T2{
			A: "testing",
			C: false,
		},
	}
	err := Forward(&forward, &t4)
	assert.Equal(t, nil, err)
	assert.Equal(t, 555, t4.ID)
	assert.Equal(t, true, t4.Enabled)
	assert.Equal(t, "testing", t4.X.A)
	assert.Equal(t, 555, forward.ID)
	assert.Equal(t, true, forward.Enabled)
	assert.Equal(t, "testing", forward.X.A)
}

func TestHelperCompare(t *testing.T) {
	testStructA := struct {
		TestInt   int
		TestInt8  int8
		TestInt16 int16
	}{
		TestInt:   12,
		TestInt8:  42,
		TestInt16: 55,
	}

	testStructB := struct {
		TestInt   int
		TestInt8  int8
		TestInt16 int16
	}{
		TestInt:   12,
		TestInt8:  42,
		TestInt16: 55,
	}

	assert.Equal(t, true, Compare(testStructA, testStructB))
}

func TestHelperReplace(t *testing.T) {
	type testStruct struct {
		TestInt     int
		TestInt8    int8
		TestInt16   int16
		TestInt32   int32
		TestInt64   int64
		TestString1 string
		TestString2 string
		TestString3 string
		TestString4 string
		TestBool    bool
		TestFloat32 float32
		TestFloat64 float64
		// TestComplex64  complex64
		// TestComplex128 complex128
	}
	ts := testStruct{
		TestInt:     12,
		TestInt8:    42,
		TestInt16:   55,
		TestInt32:   33,
		TestInt64:   78,
		TestString1: "test",
		TestString2: "test",
		TestString3: "test",
		TestString4: "test",
		TestBool:    false,
		TestFloat32: 13.444,
		TestFloat64: 16.444,
		// TestComplex64:  12333,
		// TestComplex128: 123444455,
	}

	// value, err = Replace(&ts, complex64(12333), complex64(12334), ReplaceAll)
	// assert.Equal(t, nil, err)
	// assert.Equal(t, 12334, value.(*testStruct).TestComplex64)

	// value, err = Replace(&ts, complex64(12333), float32(12334), ReplaceAll)
	// assert.Equal(t, ErrNotReplaced, err)

	value, err := Replace(&ts, "test", "new", 2)
	assert.Equal(t, nil, err)
	assert.Equal(t, "new", value.(*testStruct).TestString1)
	assert.Equal(t, "test", value.(*testStruct).TestString3)

	value, err = Replace(&ts, "test", "new", ReplaceAll)
	assert.Equal(t, nil, err)
	assert.Equal(t, "new", value.(*testStruct).TestString1)
	assert.Equal(t, "new", value.(*testStruct).TestString3)

	value, err = Replace(&ts, 12, 42, ReplaceAll)
	assert.Equal(t, nil, err)
	assert.Equal(t, 42, value.(*testStruct).TestInt)

	value, err = Replace(&ts, false, true, ReplaceAll)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, value.(*testStruct).TestBool)

	value, err = Replace(&ts, float32(13.444), float32(42.444), ReplaceAll)
	assert.Equal(t, nil, err)
	assert.Equal(t, float32(42.444), value.(*testStruct).TestFloat32)
}

func TestHelperMapFunc(t *testing.T) {
	type testStruct struct {
		Username string
		Title    string
		Content  string
	}
	ts := testStruct{
		Username: "Roninzo",
		Title:    "Test title",
		Content:  "Test content",
	}

	res, err := MapFunc(&ts, func(v reflect.Value) error {
		if v.Type().Kind() == reflect.String {
			v.SetString(strings.ToLower(v.String()))
		}
		return nil
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, "roninzo", res.(*testStruct).Username)

	var testErr = errors.New("Test")
	res, err = MapFunc(&ts, func(v reflect.Value) error {
		return testErr
	})
	assert.Contains(t, err.Error(), testErr.Error())
}

/*   B e n c h m a r k s   */

func BenchmarkCompareEqual(b *testing.B) {
	testStructA := struct {
		TestInt   int
		TestInt8  int8
		TestInt16 int16
	}{
		TestInt:   12,
		TestInt8:  42,
		TestInt16: 55,
	}

	testStructB := struct {
		TestInt   int
		TestInt8  int8
		TestInt16 int16
	}{
		TestInt:   12,
		TestInt8:  42,
		TestInt16: 55,
	}

	for n := 0; n < b.N; n++ {
		Compare(testStructA, testStructB)
	}
}

func BenchmarkCompareNotEqual(b *testing.B) {
	testStructA := struct {
		TestInt   int
		TestInt8  int8
		TestInt16 int16
	}{
		TestInt:   12,
		TestInt8:  42,
		TestInt16: 56,
	}

	testStructB := struct {
		TestInt   int
		TestInt8  int8
		TestInt16 int16
	}{
		TestInt:   12,
		TestInt8:  42,
		TestInt16: 55,
	}

	for n := 0; n < b.N; n++ {
		Compare(testStructA, testStructB)
	}
}

func BenchmarkReplace(b *testing.B) {
	type testStruct struct {
		TestInt64      int64
		TestString1    string
		TestString2    string
		TestString3    string
		TestString4    string
		TestBool       bool
		TestFloat32    float32
		TestFloat64    float64
		TestComplex64  complex64
		TestComplex128 complex128
	}
	ts := testStruct{
		TestInt64:      78,
		TestString1:    "test",
		TestString2:    "test",
		TestString3:    "test",
		TestString4:    "test",
		TestBool:       false,
		TestFloat32:    13.444,
		TestFloat64:    16.444,
		TestComplex64:  12333,
		TestComplex128: 123444455,
	}
	for n := 0; n < b.N; n++ {
		Replace(&ts, "test", "new", 2)
	}
}

func BenchmarkMapFunc(b *testing.B) {
	type testStruct struct {
		Username string
		Title    string
		Content  string
	}
	ts := testStruct{
		Username: "Roninzo",
		Title:    "Test title",
		Content:  "Test content",
	}
	for n := 0; n < b.N; n++ {
		MapFunc(&ts, func(v reflect.Value) error {
			if v.Type().Kind() == reflect.String {
				v.SetString(strings.ToLower(v.String()))
			}
			return nil
		})
	}
}
