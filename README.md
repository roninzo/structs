Golang structs
[![license](https://img.shields.io/badge/license-MIT-green "The MIT License (MIT)")](LICENSE)
[![build](https://img.shields.io/badge/build-passing-green "Go build status")](structs.go)
[![coverage](https://img.shields.io/badge/coverage-85%25-yellowgreen?logo=codecov "Unit tests coverage")](structs_test.go) 
=======

Package structs implements a generic interface for manipulating Go structs.
The related API is powered and inspired from the Go reflection package.

## Installation

```
go get github.com/roninzo/structs
```

## Usage

### Example

main.go:
```go
package main

import (
	"fmt"

	"github.com/roninzo/structs"
)

func main() {
	type T struct {
		String string
		Uint   uint
		Bool   bool
		Int    int32
	}

	t := T{
		String: "Roninzo",
		Uint:   123456,
		Bool:   true,
		Int:    5,
	}

	s, err := structs.New(&t)
	if err != nil {
		fmt.Printf("New[Error]: %v.\n", err)
		return
	}

	fmt.Printf("Name               : %v.\n", s.Name())
	fmt.Printf("Value of 1st field : %v.\n", s.Field(0).Value())
	fmt.Printf("Value of Uint      : %v.\n", s.Field("Uint").Value())
	fmt.Printf("Value of Int       : %v.\n", s.Field("Int").Value())
	fmt.Printf("Dump               : %s.\n", s.Dump())

	err = s.Field("Uint").Set(uint(654321))
	if err != nil {
		fmt.Printf("Set[Error]: %v.\n", err)
	}

	err = s.Field("Int").Set(6)
	if err != nil {
		fmt.Printf("Set[Error]: %v.\n", err)
	}

	err = s.Field("Bool").Set(6)
	if err != nil {
		fmt.Printf("Set[Error]: %v.\n", err)
	}

	fmt.Printf("Value of String    : %s.\n", s.Field("String").String()) // syntax for %s verb
	fmt.Printf("Value of Uint      : %d.\n", s.Field("Uint").Uint())     // syntax for %d verb
	fmt.Printf("Value of Int       : %d.\n", s.Field("Int").Int())       // syntax for %d verb
	fmt.Printf("Dump               : %s.\n", s.Dump())
	fmt.Printf("\nVerification       :\n")
	fmt.Printf("t.String           : %s.\n", t.String)
	fmt.Printf("t.Uint             : %d.\n", t.Uint)
	fmt.Printf("t.Int              : %d.\n", t.Int)
}
```

### Execute

```bash
go run main.go
```

### Output

```yaml
Name               : T.
Value of 1st field : Roninzo.
Value of Uint      : 123456.
Value of Int       : 5.
Dump               : {
   "String": "Roninzo",
   "Uint": 123456,
   "Bool": true,
   "Int": 5
 }.
Set[Error]: wrong kind of value for field T.Bool. got: 'int' want: 'bool'.
Value of String    : Roninzo.
Value of Uint      : 654321.
Value of Int       : 6.
Dump               : {
   "String": "Roninzo",
   "Uint": 654321,
   "Bool": true,
   "Int": 6
 }.

Verification       :
t.String           : Roninzo.
t.Uint             : 654321.
t.Int              : 6.
```

## Caveat

*Package API is not final yet*.


## Documentation

- [pkg.go.dev/github.com/roninzo/structs](https://pkg.go.dev/github.com/roninzo/structs)


<!-- 
## Coverage

### Unit Tests

```
ok  	github.com/roninzo/structs	0.336s	coverage: 78.7% of statements
```

### Benchmarks

```
BenchmarkContains-4           	  254612	      4168   ns/op	    1872 B/op	      35 allocs/op
BenchmarkCompareEqual-4       	 3544616	       337.5 ns/op	      32 B/op	       2 allocs/op
BenchmarkCompareNotEqual-4    	 3569827	       336.9 ns/op	      32 B/op	       2 allocs/op
BenchmarkIndex-4              	  331239	      3682   ns/op	    1872 B/op	      35 allocs/op
BenchmarkFieldNameByValue-4   	  296122	      3624   ns/op	    1872 B/op	      35 allocs/op
BenchmarkReplace-4            	  157482	      7606   ns/op	    1920 B/op	      25 allocs/op
BenchmarkMap-4                	 1000000	      1550   ns/op	     672 B/op	      15 allocs/op
PASS
coverage: 27.0% of statements
ok  	github.com/roninzo/structs	10.581s
```
 -->

## Dependencies

- [github.com/pkg/errors](https://github.com/pkg/errors)
- [github.com/stretchr/testify](https://github.com/stretchr/testify)


## Inspired/forked from

- [github.com/fatih/structs](https://github.com/fatih/structs) (No longer maintained)
- [github.com/PumpkinSeed/structs](https://github.com/PumpkinSeed/structs)
<!-- 
- https://go101.org/article/reflection.html
- [github.com/hvoecking/10772475](https://gist.github.com/hvoecking/10772475)
- [github.com/Ompluscator/dynamic-struct](https://github.com/Ompluscator/dynamic-struct)
- [github.com/r3labs/diff](https://github.com/r3labs/diff)
- [github.com/vdobler/ht/blob/populate](https://github.com/vdobler/ht/blob/master/populate/populate.go)
- [github.com/gookit/config](https://github.com/gookit/config)
- https://stackoverflow.com/questions/24348184/get-pointer-to-value-using-reflection
- https://github.com/jmhodges/copyfighter
- https://github.com/geraldywy/go-refcheck
 -->

## To Do
- Extend support for Pointer to struct fields  
- Extend support for Slice of any Type as struct fields  
- Extend support for Map of any Type as struct fields  
- Extend support for Method of struct  
- Extend support for complex64/128 field types
- Implement Diff 
- Implement Map
- Improve performance
- Improve coverage
