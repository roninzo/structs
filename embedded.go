package embedded

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/roninzo/structs/pointers"
)

/*   S t r u c t   d e f i n i t i o n   */

type Unembeddeds struct {
	Index       int
	Indexes     []int
	Name        string
	Value       reflect.Value
	StructField reflect.StructField
	Type        reflect.Type
}

/*   I m p l e m e n t a t i o n   */

// Explode catalogs all fields info necessary to subsequently create StructField.
// The catalog will of course contain all fields at the base of the top level struct. However, any encountered
// anonymous/embedded fields will be recursively scanned to also include their fields too, in the same collection.
//
// // NOTE(roninzo): Explode avoids scanning potential embedded struct from third party types
// // (such as time.Time, reflect.Value, etc.) by only expanding on structs that are declared locally to the current package.
// // Hence, the use of the namespace in the program.
func Explode(v reflect.Value, t reflect.Type, m map[int]Unembeddeds, namespace *string, c *int, x []int) {
	if t.Kind() != reflect.Struct {
		return
	}
	if namespace == nil {
		namespace = pointers.String(nameSpace(t)) // <=> tmp := nameSpace(t); namespace = &tmp
	}
	if c == nil {
		c = pointers.Int(0) // <=> tmp := 0; c = &tmp
	}
	if x == nil {
		x = make([]int, 0)
	}
	x = append(x, 0)
	n := len(x)
	for i := 0; i < t.NumField(); i++ {
		x[n-1] = i
		sv := v.Field(i)
		sf := t.Field(i)
		if sf.Anonymous { // && nameSpace(sf.Type) == *namespace {
			Explode(sv, sf.Type, m, namespace, c, x)
		} else {
			tmp := make([]int, n)
			copy(tmp, x)
			unembedded := Unembeddeds{
				Index:       *c,
				Indexes:     tmp,
				Name:        sf.Name,
				Value:       sv,
				StructField: sf,
				Type:        sf.Type,
			}
			m[*c] = unembedded
			*c++
		}
	}
	x = x[:n-1]
}

/*   U n e x p o r t e d   */

// nameSpace returns the name space string part of reflect type.
func nameSpace(t reflect.Type) string {
	s := fmt.Sprintf("%v", t)
	return strings.Split(s, ".")[0]
}
