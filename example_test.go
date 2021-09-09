package structs_test

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/roninzo/structs"
)

/*   S t r u c t V a l u e   */

func ExampleNew() {
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

	// Output:
	// Name               : T.
	// Value of 1st field : Roninzo.
	// Value of Uint      : 123456.
	// Value of Int       : 5.
	// Dump               : {
	//     "String": "Roninzo",
	//     "Uint": 123456,
	//     "Bool": true,
	//     "Int": 5
	//  }.
	// Set[Error]: wrong kind of value for field T.Bool. got: 'int' want: 'bool'.
	// Value of String    : Roninzo.
	// Value of Uint      : 654321.
	// Value of Int       : 6.
	// Dump               : {
	//     "String": "Roninzo",
	//     "Uint": 654321,
	//     "Bool": true,
	//     "Int": 6
	//  }.
	//
	// Verification       :
	// t.String           : Roninzo.
	// t.Uint             : 654321.
	// t.Int              : 6.
}

func ExampleNew_pointerFields() {
	type Server struct {
		Name    *string `json:"name"`
		ID      *uint   `json:"id"`
		Enabled *bool   `json:"enabled"`
		Count   *int32  `json:"count"`
	}

	server := Server{
		Name:    structs.PtrString("Roninzo"),
		ID:      structs.PtrUint(uint(123456)),
		Enabled: structs.PtrBool(true),
		Count:   structs.PtrInt32(int32(5)),
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error          : %v\n", err)
	}

	fmt.Printf("Name           : %v\n", s.Name())
	fmt.Printf("Value of ID    : %v\n", s.Field("ID").PtrValue())
	fmt.Printf("Value of 0     : %v\n", s.Field(0).PtrValue())
	fmt.Printf("Value of Count : %v\n", s.Field("Count").PtrValue())
	fmt.Printf("Dump           : %s.\n", s.Dump())

	err = s.Field("ID").Set(structs.PtrUint(uint(654321)))
	if err != nil {
		fmt.Printf("Error          : %v\n", err)
	}

	err = s.Field("Count").Set(structs.PtrInt32(int32(6)))
	if err != nil {
		fmt.Printf("Error          : %v\n", err)
	}

	err = s.Field("Enabled").Set(structs.PtrInt32(int32(6))) // not compatible with bool
	if err != nil {
		fmt.Printf("Error          : %v\n", err)
	}

	fmt.Printf("Value of Name  : %v\n", s.Field("Name").PtrValue())
	fmt.Printf("Value of ID    : %v\n", s.Field("ID").PtrValue())
	fmt.Printf("Value of Count : %v\n", s.Field("Count").PtrValue())
	fmt.Printf("Dump           : %s.\n", s.Dump())
	fmt.Printf("\nVerification   :\n")
	fmt.Printf("server.Name    : %s\n", *server.Name)
	fmt.Printf("server.ID      : %d\n", *server.ID)
	fmt.Printf("server.Count   : %d\n", *server.Count)

	// Output:
	// Name           : Server
	// Value of ID    : 123456
	// Value of 0     : Roninzo
	// Value of Count : 5
	// Dump           : {
	//     "name": "Roninzo",
	//     "id": 123456,
	//     "enabled": true,
	//     "count": 5
	//  }.
	// Error          : wrong kind of value for field Server.Enabled. got: '*int32' want: '*bool'
	// Value of Name  : Roninzo
	// Value of ID    : 654321
	// Value of Count : 6
	// Dump           : {
	//     "name": "Roninzo",
	//     "id": 654321,
	//     "enabled": true,
	//     "count": 6
	//  }.
	//
	// Verification   :
	// server.Name    : Roninzo
	// server.ID      : 654321
	// server.Count   : 6
}
func ExampleStructValue_IsValid() {
	type Server struct {
		Name    string
		ID      uint
		Enabled bool
		Count   int32
	}

	server := Server{
		Name:    "Roninzo",
		ID:      123456,
		Enabled: true,
		Count:   5,
	}

	s, err := structs.New(nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("IsValid: %v\n", s.IsValid())
	}

	s, err = structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("IsValid: %v\n", s.IsValid())
	}

	// Output:
	// Error: invalid concrete value; want: 'struct' or 'ptr' or 'slice', got: 'nil'
	// IsValid: true
}

func ExampleStructValue_Name() {
	type Server struct {
		Name    string
		ID      uint
		Enabled bool
		Count   int32
	}

	server := Server{
		Name:    "Roninzo",
		ID:      123456,
		Enabled: true,
		Count:   5,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("Name: %v\n", s.Name())

	// Output:
	// Name: Server
}

func ExampleStructValue_Type() {
	type Server struct {
		Name    string
		ID      uint
		Enabled bool
		Count   int32
	}

	server := Server{
		Name:    "Roninzo",
		ID:      123456,
		Enabled: true,
		Count:   5,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("Type: %v\n", s.Type())

	// Output:
	// Type: structs_test.Server
}

func ExampleStructValue_Kind() {
	type Server struct {
		Name    string
		ID      uint
		Enabled bool
		Count   int32
	}

	server := Server{
		Name:    "Roninzo",
		ID:      123456,
		Enabled: true,
		Count:   5,
	}

	s, err := structs.New(nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	if s != nil {
		fmt.Printf("Kind: %v\n", s.Kind())
	}

	s, err = structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("Kind: %v\n", s.Kind())

	// Output:
	// Error: invalid concrete value; want: 'struct' or 'ptr' or 'slice', got: 'nil'
	// Kind: struct
}

func ExampleStructValue_Value() {
	type Server struct {
		Name    string
		ID      uint
		Enabled bool
		Count   int32
	}

	server := Server{
		Name:    "Roninzo",
		ID:      123456,
		Enabled: true,
		Count:   5,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("Value: %v\n", s.Value())

	// Output:
	// Value: {Roninzo 123456 true 5}
}

func ExampleStructValue_Values() {
	type Program struct {
		Name string `json:"name,omitempty"`
	}

	type Server struct {
		Name       string  `json:"name,omitempty"`
		ID         uint    `json:"id,omitempty"`
		Enabled    bool    `json:"enabled,omitempty"`
		Count      int32   `json:"count,omitempty"`
		Password   string  `json:"-"`
		Program    Program `json:"program,omitempty"`
		unexported bool
	}

	program := Program{
		Name: "Apache",
	}

	server := Server{
		Name:     "Roninzo",
		ID:       123456,
		Enabled:  true,
		Count:    0,
		Password: "abcdefg",
		Program:  program,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// json struct tag omits Password and Count
	for i, value := range s.Values() {
		fmt.Printf("Values[%d]: %v\n", i, value)
	}

	// Output:
	// Values[0]: Roninzo
	// Values[1]: 123456
	// Values[2]: true
	// Values[3]: 0
	// Values[4]: abcdefg
	// Values[5]: Apache
}

func ExampleStructValue_PtrValues() {
	type Program struct {
		Name *string `json:"name,omitempty"`
	}

	type Server struct {
		Name       *string  `json:"name,omitempty"`
		ID         *uint    `json:"id,omitempty"`
		Enabled    *bool    `json:"enabled,omitempty"`
		Count      *int     `json:"count,omitempty"`
		Program    *Program `json:"program,omitempty"`
		Password   *string  `json:"-"`
		unexported *bool
	}

	program := Program{
		Name: structs.PtrString("Apache"),
	}

	server := Server{
		Name:       structs.PtrString("Roninzo"),
		ID:         structs.PtrUint(123456),
		Enabled:    structs.PtrBool(true),
		Count:      structs.PtrInt(0),
		Program:    &program,
		Password:   structs.PtrString("abcdefg"),
		unexported: structs.PtrBool(true),
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v.\n", err)
	}

	// json struct tag omits Password and Count
	for i, value := range s.PtrValues() {
		fmt.Printf("PtrValues[%d]: %v.\n", i, value)
	}

	// Output:
	// PtrValues[0]: Roninzo.
	// PtrValues[1]: 123456.
	// PtrValues[2]: true.
	// PtrValues[3]: 0.
	// PtrValues[4]: Apache.
	// PtrValues[5]: abcdefg.
}

func ExampleStructValue_Dump() {
	type Program struct {
		Name string `json:"name,omitempty"`
	}

	type Server struct {
		Name     string  `json:"name,omitempty"`
		ID       uint    `json:"id,omitempty"`
		Enabled  bool    `json:"enabled,omitempty"`
		Count    int32   `json:"count,omitempty"`
		Password string  `json:"-"`
		Program  Program `json:"program,omitempty"`
	}

	program := Program{
		Name: "Apache",
	}

	server := Server{
		Name:     "Roninzo",
		ID:       123456,
		Enabled:  true,
		Count:    0,
		Password: "abcdefg",
		Program:  program,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// json struct tag omits Password and Count
	fmt.Printf("Dump: %s\n", s.Dump())

	// Output:
	// Dump: {
	//     "name": "Roninzo",
	//     "id": 123456,
	//     "enabled": true,
	//     "program": {
	//        "name": "Apache"
	//     }
	//  }
}

func ExampleStructValue_CanSet() {
	type Program struct {
		Name string `json:"name,omitempty"`
	}

	type Server struct {
		Name     string   `json:"name,omitempty"`
		ID       uint     `json:"id,omitempty"`
		Enabled  bool     `json:"enabled,omitempty"`
		Count    int32    `json:"count,omitempty"`
		Password string   `json:"-"`
		Program  *Program `json:"program,omitempty"`
	}

	program := Program{
		Name: "Apache",
	}

	server := Server{
		Name:     "Roninzo",
		ID:       123456,
		Enabled:  true,
		Count:    0,
		Password: "abcdefg",
		Program:  &program,
	}

	s1, _ := structs.New(&server)
	s2, _ := structs.New(server)

	fmt.Printf("CanSet: %v\n", s1.CanSet())
	fmt.Printf("CanSet: %v\n", s2.CanSet())

	// Output:
	// CanSet: true
	// CanSet: false
}

func ExampleStructValue_Multiple() {
	type Server struct {
		Name     string `json:"name,omitempty"`
		ID       uint   `json:"id,omitempty"`
		Enabled  bool   `json:"enabled,omitempty"`
		Count    int32  `json:"count,omitempty"`
		Password string `json:"-"`
	}

	server := Server{
		Name:     "Roninzo",
		ID:       123456,
		Enabled:  true,
		Count:    0,
		Password: "abcdefg",
	}

	s1, _ := structs.New(&server)
	s2, _ := structs.New([]*Server{&server})

	fmt.Printf("Multiple: %v\n", s1.Multiple())
	fmt.Printf("Multiple: %v\n", s2.Multiple())

	// Output:
	// Multiple: false
	// Multiple: true
}

func ExampleStructValue_Fields() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Count      int32  `json:"count,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Count:      0,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("NumField: %d\n", s.NumField())

	fields := s.Fields()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("Fields.Names: %v\n", fields.Names())

	// Output:
	// NumField: 6
	// Fields.Names: [Name ID Enabled Count Password unexported]
}

func ExampleStructValue_Field() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Count      int32  `json:"count,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Count:      0,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Enabled")

	fmt.Printf("Field.Name: %v\n", f1.Name())
	fmt.Printf("Field.Value: %v\n", f1.Value())
	fmt.Printf("Field.Name: %v\n", f2.Name())
	fmt.Printf("Field.Value: %v\n", f2.Value())

	// Output:
	// Field.Name: Name
	// Field.Value: Roninzo
	// Field.Name: Enabled
	// Field.Value: true
}

func ExampleStructValue_NumField() {
	type Server struct {
		Name     string `json:"name,omitempty"`
		ID       uint   `json:"id,omitempty"`
		Enabled  bool   `json:"enabled,omitempty"`
		Count    int32  `json:"count,omitempty"`
		Password string `json:"-"`
	}

	server := Server{
		Name:     "Roninzo",
		ID:       123456,
		Enabled:  true,
		Count:    0,
		Password: "abcdefg",
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("NumField: %v\n", s.NumField())

	// Output:
	// NumField: 5
}

func ExampleStructValue_FindStruct() {
	type Program struct {
		Name string `json:"name,omitempty"`
	}

	type Server struct {
		Name     string   `json:"name,omitempty"`
		ID       uint     `json:"id,omitempty"`
		Enabled  bool     `json:"enabled,omitempty"`
		Count    int32    `json:"count,omitempty"`
		Password string   `json:"-"`
		Program  *Program `json:"program,omitempty"`
	}

	program := Program{
		Name: "Apache",
	}

	server := Server{
		Name:     "Roninzo",
		ID:       123456,
		Enabled:  true,
		Count:    0,
		Password: "abcdefg",
		Program:  &program,
	}

	s1, _ := structs.New(&server)
	s2 := s1.FindStruct("Program")

	fmt.Printf("Program.Name: %v\n", s2.Name())
	fmt.Printf("Program.Value: %v\n", s2.Value())

	// Output:
	// Program.Name: Program
	// Program.Value: {Apache}
}

func ExampleStructValue_IsZero() {
	type Server struct {
		Name     string `json:"name,omitempty"`
		ID       uint   `json:"id,omitempty"`
		Enabled  bool   `json:"enabled,omitempty"`
		Count    int32  `json:"count,omitempty"`
		Password string `json:"-"`
	}

	server := Server{
		Name:     "Roninzo",
		ID:       123456,
		Enabled:  true,
		Count:    0,
		Password: "abcdefg",
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("IsZero: %v\n", s.IsZero())

	server = Server{
		Name:     "",
		ID:       0,
		Enabled:  false,
		Count:    0,
		Password: "",
	}

	s, err = structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("IsZero: %v\n", s.IsZero())

	// Output:
	// IsZero: false
	// IsZero: true
}

func ExampleStructValue_HasZero() {
	type Server struct {
		Name     string `json:"name,omitempty"`
		ID       uint   `json:"id,omitempty"`
		Enabled  bool   `json:"enabled,omitempty"`
		Count    int32  `json:"count,omitempty"`
		Password string `json:"-"`
	}

	server := Server{
		Name:     "Roninzo",
		ID:       0, // zero-value
		Enabled:  true,
		Count:    5,
		Password: "abcdefg",
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("HasZero: %v\n", s.HasZero())

	server = Server{
		Name:     "Roninzo",
		ID:       123456,
		Enabled:  true,
		Count:    5,
		Password: "abcdefg",
	}

	s, err = structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("HasZero: %v\n", s.HasZero())

	// Output:
	// HasZero: true
	// HasZero: false
}

func ExampleStructValue_IsNested() {
	type Program struct {
		Name string `json:"name,omitempty"`
	}

	type Server struct {
		Name     string   `json:"name,omitempty"`
		ID       uint     `json:"id,omitempty"`
		Enabled  bool     `json:"enabled,omitempty"`
		Count    int32    `json:"count,omitempty"`
		Password string   `json:"-"`
		Program  *Program `json:"program,omitempty"`
	}

	program := Program{
		Name: "Apache",
	}

	server := Server{
		Name:     "Roninzo",
		ID:       123456,
		Enabled:  true,
		Count:    0,
		Password: "abcdefg",
		Program:  &program,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	p := s.FindStruct("Program")

	fmt.Printf("Server.IsNested: %v\n", s.IsNested())
	fmt.Printf("Program.IsNested: %v\n", p.IsNested())

	// Output:
	// Server.IsNested: false
	// Program.IsNested: true
}

func ExampleStructValue_HasNested() {
	type Program struct {
		Name string `json:"name,omitempty"`
	}

	type Server struct {
		Name     string   `json:"name,omitempty"`
		ID       uint     `json:"id,omitempty"`
		Enabled  bool     `json:"enabled,omitempty"`
		Count    int32    `json:"count,omitempty"`
		Password string   `json:"-"`
		Program  *Program `json:"program,omitempty"`
	}

	program := Program{
		Name: "Apache",
	}

	server := Server{
		Name:     "Roninzo",
		ID:       123456,
		Enabled:  true,
		Count:    0,
		Password: "abcdefg",
		Program:  &program,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	p, err := structs.New(&program)
	if p.Error != nil {
		fmt.Printf("Error: %s\n", p.Error)
	}

	fmt.Printf("Server.HasNested: %v\n", s.HasNested())
	fmt.Printf("Program.HasNested: %v\n", p.HasNested())

	// Output:
	// Server.HasNested: true
	// Program.HasNested: false
}

func ExampleStructValue_Namespace() {
	type Program struct {
		Name string `json:"name,omitempty"`
	}

	type Server struct {
		Name     string   `json:"name,omitempty"`
		ID       uint     `json:"id,omitempty"`
		Enabled  bool     `json:"enabled,omitempty"`
		Count    int32    `json:"count,omitempty"`
		Password string   `json:"-"`
		Program  *Program `json:"program,omitempty"`
	}

	program := Program{
		Name: "Apache",
	}

	server := Server{
		Name:     "Roninzo",
		ID:       123456,
		Enabled:  true,
		Count:    0,
		Password: "abcdefg",
		Program:  &program,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	p := s.FindStruct("Program")

	fmt.Printf("Namespace: %v\n", s.Namespace())
	fmt.Printf("Namespace: %v\n", p.Namespace())

	// Output:
	// Namespace: Server
	// Namespace: Server.Program
}

func ExampleStructValue_Path() {
	type Server struct {
		Name     string `json:"name,omitempty"`
		ID       uint   `json:"id,omitempty"`
		Enabled  bool   `json:"enabled,omitempty"`
		Count    int32  `json:"count,omitempty"`
		Password string `json:"-"`
	}

	server := Server{
		Name:     "Roninzo",
		ID:       123456,
		Enabled:  true,
		Count:    5,
		Password: "abcdefg",
	}

	s1, _ := structs.New(server)
	s2, _ := structs.New(&server)
	s3, _ := structs.New([]Server{server})
	s4, _ := structs.New(&[]Server{server})
	s5, _ := structs.New([]*Server{&server})

	fmt.Printf("Path: %v\n", s1.Path())
	fmt.Printf("Path: %v\n", s2.Path())
	fmt.Printf("Path: %v\n", s3.Path())
	fmt.Printf("Path: %v\n", s4.Path())
	fmt.Printf("Path: %v\n", s5.Path())

	// Output:
	// Path: struct
	// Path: ptr,struct
	// Path: slice,struct
	// Path: ptr,slice,struct
	// Path: slice,ptr,struct
}

func ExampleStructValue_Rows() {
	type Server struct {
		Count int32 `json:"count,omitempty"`
	}

	servers := []Server{
		{Count: 5},
		{Count: 6},
	}

	s, err := structs.New(&servers)
	if err != nil {
		fmt.Printf("New: %v.\n", err)
	}

	rows, err := s.Rows()
	if err != nil {
		fmt.Printf("Rows[Error]: %v.\n", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		fmt.Printf("Columns[Error]: %v.\n", err)
	}
	fmt.Printf("Row: %s.\n", s.Dump())
	fmt.Printf("StructValue: %s.\n", s.Debug())
	fmt.Printf("Len: %d.\n", rows.Len())
	fmt.Printf("MaxRow: %d.\n", rows.MaxRow())
	fmt.Printf("Columns: %v.\n", cols)

	for rows.Next() {
		f := rows.Field("Count")
		fmt.Printf("[%d] %s: %d.\n", rows.Index(), f.Name(), f.Int())
		c := f.Int()
		err := f.Set(c * 10)
		if err != nil {
			fmt.Printf("Set[Error]: %v.\n", err)
		}
		fmt.Printf("[%d] %s: %d.\n", rows.Index(), f.Name(), f.Int())
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Err[Error]: %v.\n", err)
	}

	fmt.Println("Verification:")
	fmt.Printf("servers[0].Count: %d.\n", servers[0].Count)
	fmt.Printf("servers[1].Count: %d.\n", servers[1].Count)

	// Output:
	// Row: {
	//     "count": 5
	//  }.
	// StructValue: {
	//     "value": {
	//        "count": 5
	//     },
	//     "rows": [
	//        {
	//           "count": 5
	//        },
	//        {
	//           "count": 6
	//        }
	//     ],
	//     "max_row": 2,
	//     "kinds": "ptr,slice,struct",
	//     "parent": "",
	//     "error": null
	//  }.
	// Len: 2.
	// MaxRow: 1.
	// Columns: [Count].
	// [0] Count: 5.
	// [0] Count: 50.
	// [1] Count: 6.
	// [1] Count: 60.
	// Verification:
	// servers[0].Count: 50.
	// servers[1].Count: 60.
}

func ExampleStructValue_Rows_notMultiple() {
	type Server struct {
		Count int32 `json:"count,omitempty"`
	}

	server := Server{
		Count: 5,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("New[Error]: %v.\n", err)
	}

	rows, err := s.Rows()
	if err != nil {
		fmt.Printf("Rows[Error]: %v.\n", err)
	}
	if rows != nil {
		fmt.Printf("Rows: %v.\n", rows)
	}

	// Output:
	// Rows[Error]: structs not found.
}

func ExampleStructValue_Rows_empty() {
	type Server struct {
		Count int32 `json:"count,omitempty"`
	}

	servers := []Server{}

	s, err := structs.New(&servers)
	if err != nil {
		fmt.Printf("New[Error]: %v.\n", err)
	}

	rows, err := s.Rows()
	if err != nil {
		fmt.Printf("Rows[Error]: %v.\n", err)
	}
	if rows != nil {
		fmt.Printf("Rows: %v.\n", rows)
	}

	// Output:
	// Rows[Error]: struct rows not found.
}

/*   S t r u c t F i e l d   */

func ExampleStructField_IsValid() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Count      int32  `json:"count,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Count:      0,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := &structs.StructField{}

	fmt.Printf("IsValid: %v\n", f1.IsValid())
	fmt.Printf("IsValid: %v\n", f2.IsValid())

	// Output:
	// IsValid: true
	// IsValid: false
}

func ExampleStructField_Name() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Count      int32  `json:"count,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Count:      0,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Enabled")

	fmt.Printf("Field: %v\n", f1.Name())
	fmt.Printf("Field: %v\n", f2.Name())

	// Output:
	// Field: Name
	// Field: Enabled
}

func ExampleStructField_Namespace() {
	type Program struct {
		Name string `json:"name,omitempty"`
	}

	type Server struct {
		Name     string   `json:"name,omitempty"`
		ID       uint     `json:"id,omitempty"`
		Enabled  bool     `json:"enabled,omitempty"`
		Count    int32    `json:"count,omitempty"`
		Password string   `json:"-"`
		Program  *Program `json:"program,omitempty"`
	}

	program := Program{
		Name: "Apache",
	}

	server := Server{
		Name:     "Roninzo",
		ID:       123456,
		Enabled:  true,
		Count:    0,
		Password: "abcdefg",
		Program:  &program,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.FindStruct("Program").Field(0)

	fmt.Printf("Namespace: %v\n", f1.Namespace())
	fmt.Printf("Namespace: %v\n", f2.Namespace())

	// Output:
	// Namespace: Server.Name
	// Namespace: Server.Program.Name
}

func ExampleStructField_Value() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Count      int32  `json:"count,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Count:      0,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Enabled")

	fmt.Printf("Value: %v\n", f1.Value())
	fmt.Printf("Value: %v\n", f2.Value())

	// Output:
	// Value: Roninzo
	// Value: true
}

func ExampleStructField_Type() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Count      int32  `json:"count,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Count:      0,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Enabled")

	fmt.Printf("Type: %v\n", s.Type())
	fmt.Printf("Type: %v\n", f1.Type())
	fmt.Printf("Type: %v\n", f2.Type())

	// Output:
	// Type: structs_test.Server
	// Type: string
	// Type: bool
}

func ExampleStructField_Kind() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Count      int32  `json:"count,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Count:      0,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Enabled")

	fmt.Printf("Kind: %v\n", f1.Kind())
	fmt.Printf("Kind: %v\n", f2.Kind())

	// Output:
	// Kind: string
	// Kind: bool
}

func ExampleStructField_Tag() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Count      int32  `json:"count,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Count:      0,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	t1, _ := s.Field(0).Tag("json")
	t2, _ := s.Field("Enabled").Tag("json")

	fmt.Printf("Tag: %v\n", t1)
	fmt.Printf("Tag: %v\n", t2)

	// Output:
	// Tag: name,omitempty
	// Tag: enabled,omitempty
}

func ExampleStructField_CanSet() {
	type Program struct {
		Name string `json:"name,omitempty"`
	}

	type Server struct {
		Name     string  `json:"name,omitempty"`
		ID       uint    `json:"id,omitempty"`
		Enabled  bool    `json:"enabled,omitempty"`
		Count    int32   `json:"count,omitempty"`
		Password string  `json:"-"`
		Program  Program `json:"program,omitempty"`
	}

	program := Program{
		Name: "Apache",
	}

	server := Server{
		Name:     "Roninzo",
		ID:       123456,
		Enabled:  true,
		Count:    0,
		Password: "abcdefg",
		Program:  program,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.FindStruct("Program").Field("Name")

	fmt.Printf("CanSet: %v\n", f1.CanSet())
	fmt.Printf("CanSet: %v\n", f2.CanSet())

	// Output:
	// CanSet: true
	// CanSet: false
}

func ExampleStructField_IsAnonymous() {
	type Server struct {
		Name     string `json:"name,omitempty"`
		ID       uint   `json:"id,omitempty"`
		Enabled  bool   `json:"enabled,omitempty"`
		Count    int32  `json:"count,omitempty"`
		Password string `json:"-"`
		error    `json:"program,omitempty"`
	}

	server := Server{
		"Roninzo",
		123456,
		true,
		0,
		"abcdefg",
		errors.New("does not work"),
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Password")
	f3 := s.Field(5)

	fmt.Printf("IsAnonymous: %v\n", f1.IsAnonymous())
	fmt.Printf("IsAnonymous: %v\n", f2.IsAnonymous())
	fmt.Printf("IsAnonymous: %v\n", f3.IsAnonymous())

	// Output:
	// IsAnonymous: false
	// IsAnonymous: false
	// IsAnonymous: true
}

func ExampleStructField_IsEmbedded() {
	type Server struct {
		Name     string `json:"name,omitempty"`
		ID       uint   `json:"id,omitempty"`
		Enabled  bool   `json:"enabled,omitempty"`
		Count    int32  `json:"count,omitempty"`
		Password string `json:"-"`
		error    `json:"program,omitempty"`
	}

	server := Server{
		"Roninzo",
		123456,
		true,
		0,
		"abcdefg",
		errors.New("does not work"),
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Password")
	f3 := s.Field(5)

	fmt.Printf("IsEmbedded: %v\n", f1.IsEmbedded())
	fmt.Printf("IsEmbedded: %v\n", f2.IsEmbedded())
	fmt.Printf("IsEmbedded: %v\n", f3.IsEmbedded())

	// Output:
	// IsEmbedded: false
	// IsEmbedded: false
	// IsEmbedded: true
}

func ExampleStructField_IsExported() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Count      int32  `json:"count,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Count:      0,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("IsExported: %v\n", s.Field(0).IsExported())
	fmt.Printf("IsExported: %v\n", s.Field("Password").IsExported())
	fmt.Printf("IsExported: %v\n", s.Field("unexported").IsExported())

	// Output:
	// IsExported: true
	// IsExported: true
	// IsExported: false
}

func ExampleStructField_IsHidden() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Count      int32  `json:"count,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Count:      0,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("IsHidden: %v\n", s.Field(0).IsHidden())
	fmt.Printf("IsHidden: %v\n", s.Field("Password").IsHidden())
	fmt.Printf("IsHidden: %v\n", s.Field("unexported").IsHidden())

	// Output:
	// IsHidden: false
	// IsHidden: true
	// IsHidden: true
}

func ExampleStructField_Zero() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Count      int32  `json:"count,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Count:      5,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field("Enabled")
	f2 := s.Field("Count")

	fmt.Printf("Zero: %v\n", f1.Zero())
	fmt.Printf("Zero: %v\n", f2.Zero())

	// Output:
	// Zero: false
	// Zero: 0
}

func ExampleStructField_IsZero() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Count      int32  `json:"count,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Count:      0,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Count")

	fmt.Printf("IsZero: %v\n", f1.IsZero())
	fmt.Printf("IsZero: %v\n", f2.IsZero())

	// Output:
	// IsZero: false
	// IsZero: true
}

func ExampleStructField_IsNil() {
	type Server struct {
		Name       string  `json:"name,omitempty"`
		ID         *uint   `json:"id,omitempty"`
		Enabled    bool    `json:"enabled,omitempty"`
		Count      *int32  `json:"count,omitempty"`
		Password   *string `json:"-"`
		unexported bool
	}

	var id uint = 5
	server := Server{
		Name:       "Roninzo",
		ID:         &id,
		Enabled:    true,
		Count:      nil,
		Password:   nil,
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("ID")
	f3 := s.Field("Count")
	f4 := s.Field("Password")

	fmt.Printf("IsNil: %v\n", f1.IsNil())
	fmt.Printf("IsNil: %v\n", f2.IsNil())
	fmt.Printf("IsNil: %v\n", f3.IsNil())
	fmt.Printf("IsNil: %v\n", f4.IsNil())

	// Output:
	// IsNil: false
	// IsNil: false
	// IsNil: true
	// IsNil: true
}

func ExampleStructField_Index() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Count      int32  `json:"count,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Count:      0,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Count")

	fmt.Printf("Index: %v\n", f1.Index())
	fmt.Printf("Index: %v\n", f2.Index())

	// Output:
	// Index: 0
	// Index: 3
}

func ExampleStructField_Set() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Count      int32  `json:"count,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Count:      0,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Count")
	f3 := s.Field("unexported")

	fmt.Printf("Value: %v.\n", f1.Value())
	fmt.Printf("Value: %v.\n", f2.Value())
	fmt.Printf("Value: %v.\n", f3.Value())

	err1 := f1.Set("Unknown")
	err2 := f2.Set(10)
	err3 := f3.Set(false)

	fmt.Printf("Value: %v.\n", f1.Value())
	fmt.Printf("Value: %v.\n", f2.Value())
	fmt.Printf("Value: %v.\n", f3.Value())
	fmt.Printf("Error: %v.\n", err1)
	fmt.Printf("Error: %v.\n", err2)
	fmt.Printf("Error: %v.\n", err3)

	// Output:
	// Value: Roninzo.
	// Value: 0.
	// Value: false.
	// Value: Unknown.
	// Value: 10.
	// Value: false.
	// Error: <nil>.
	// Error: <nil>.
	// Error: could not set field Server.unexported: struct field is not settable.
}

/*   S t r u c t F i e l d   -   L o c a l   V a r i a b l e s   */

type structTest struct {
	String          string                 `json:"string,omitempty"`
	Bool            bool                   `json:"bool,omitempty"`
	Int             int                    `json:"int,omitempty"`
	Uint            uint                   `json:"uint,omitempty"`
	Float           float32                `json:"float,omitempty"`
	Complex         complex128             `json:"complex,omitempty"` // `json:"-"` // crashes Dump!
	Bytes           []byte                 `json:"bytes,omitempty"`
	Interface       interface{}            `json:"interface,omitempty"`
	Error           error                  `json:"error,omitempty"`
	Time            time.Time              `json:"time,omitempty"`
	Duration        time.Duration          `json:"duration,omitempty"`
	NestedStruct    structNested           `json:"nested_struct,omitempty"`
	PtrString       *string                `json:"pointer_string,omitempty"`
	PtrBool         *bool                  `json:"pointer_bool,omitempty"`
	PtrInt          *int                   `json:"pointer_int,omitempty"`
	PtrUint         *uint                  `json:"pointer_uint,omitempty"`
	PtrFloat        *float32               `json:"pointer_float,omitempty"`
	PtrComplex      *complex128            `json:"pointer_complex,omitempty"` // `json:"-"` // crashes Dump!
	PtrError        *error                 `json:"pointer_error,omitempty"`
	PtrTime         *time.Time             `json:"pointer_time,omitempty"`
	PtrDuration     *time.Duration         `json:"pointer_duration,omitempty"`
	PtrNestedStruct *structNested          `json:"pointer_nested_struct,omitempty"`
	MapString       map[string]string      `json:"map_string,omitempty"`
	MapBool         map[string]bool        `json:"map_bool,omitempty"`
	MapInt          map[string]int         `json:"mapint,omitempty"`
	MapUint         map[string]uint        `json:"map_uint,omitempty"`
	MapFloat        map[string]float32     `json:"map_float,omitempty"`
	MapComplex      map[string]complex128  `json:"map_complex,omitempty"` // `json:"-"` // crashes Dump!
	MapInterface    map[string]interface{} `json:"map_interface,omitempty"`
	SliceString     []string               `json:"slice_string,omitempty"`
	SliceBool       []bool                 `json:"slice_bool,omitempty"`
	SliceInt        []int                  `json:"slice_int,omitempty"`
	SliceUint       []uint                 `json:"slice_uint,omitempty"`
	SliceFloat      []float32              `json:"slice_float,omitempty"`
	SliceComplex    []complex128           `json:"slice_complex,omitempty"` // `json:"-"` // crashes Dump!
	SliceInterface  []interface{}          `json:"slice_interface,omitempty"`
	SlicePtrString  []*string              `json:"slice_pointer_string,omitempty"`
	SlicePtrBool    []*bool                `json:"slice_pointer_bool,omitempty"`
	SlicePtrInt     []*int                 `json:"slice_pointer_int,omitempty"`
	SlicePtrUint    []*uint                `json:"slice_pointer_uint,omitempty"`
	SlicePtrFloat   []*float32             `json:"slice_pointer_float,omitempty"`
	SlicePtrComplex []*complex128          `json:"slice_pointer_complex,omitempty"`
	Hidden          string                 `json:"-"`
	unexported      bool
}

type structNested struct {
	Uint   uint   `json:"uint,omitempty"`
	String string `json:"string,omitempty"`
}

var structFieldNames []string = []string{
	"String", "PtrString", "MapString", "SliceString", "SlicePtrString",
	"Bool", "PtrBool", "MapBool", "SliceBool", "SlicePtrBool",
	"Int", "PtrInt", "MapInt", "SliceInt", "SlicePtrInt",
	"Uint", "PtrUint", "MapUint", "SliceUint", "SlicePtrUint",
	"Float", "PtrFloat", "MapFloat", "SliceFloat", "SlicePtrFloat",
	"Complex", "PtrComplex", "MapComplex", "SliceComplex", "SlicePtrComplex",
	"Interface", "MapInterface", "SliceInterface",
	"Bytes",
	"Error", "PtrError",
	"Time", "PtrTime",
	"Duration", "PtrDuration",
	"NestedStruct", "PtrNestedStruct",
	"Hidden", "unexported",
}

var structV structTest = structTest{
	String:          "Roninzo",
	Bool:            true,
	Int:             8,
	Uint:            uint(123456),
	Float:           1922.50,
	Complex:         complex(22, 50),
	Bytes:           []byte("Hello world"),
	Interface:       "anything",
	Error:           errors.New("rows not found"),
	Time:            time.Date(2021, time.August, 3, 16, 44, 46, 0, time.UTC),
	Duration:        5 * time.Second,
	NestedStruct:    structNested{Uint: 122334, String: "Apache"},
	PtrString:       structs.PtrString("Roninzo"),
	PtrBool:         structs.PtrBool(true),
	PtrInt:          structs.PtrInt(8),
	PtrUint:         structs.PtrUint(uint(123456)),
	PtrFloat:        structs.PtrFloat32(1922.50),
	PtrComplex:      structs.PtrComplex128(complex(22, 50)),
	PtrError:        structs.PtrError(errors.New("rows not found")),
	PtrTime:         structs.PtrTime(time.Date(2021, time.August, 3, 16, 44, 46, 0, time.UTC)),
	PtrDuration:     structs.PtrDuration(5 * time.Second),
	PtrNestedStruct: &structNested{Uint: 122334, String: "Apache"},
	MapString:       map[string]string{"A": "one", "B": "two", "C": "three"},
	MapBool:         map[string]bool{"A": true, "B": false},
	MapInt:          map[string]int{"A": 1, "B": 2, "C": 3},
	MapUint:         map[string]uint{"A": uint(1), "B": uint(2), "C": uint(3)},
	MapFloat:        map[string]float32{"A": 1.1, "B": 1.2, "C": 1.3},
	MapComplex:      map[string]complex128{"A": complex(1, 1), "B": complex(1, 2), "C": complex(1, 3)},
	MapInterface:    map[string]interface{}{"A": 1, "B": "two", "C": 3.0},
	SliceString:     []string{"one", "two", "three"},
	SliceBool:       []bool{true, false},
	SliceInt:        []int{1, 2, 3},
	SliceUint:       []uint{uint(1), uint(2), uint(3)},
	SliceFloat:      []float32{1.1, 1.2, 1.3},
	SliceComplex:    []complex128{complex(1, 1), complex(1, 2), complex(1, 3)},
	SliceInterface:  []interface{}{1, "two", 3.0},
	SlicePtrString:  []*string{structs.PtrString("one"), structs.PtrString("two"), structs.PtrString("three")},
	SlicePtrBool:    []*bool{structs.PtrBool(true), structs.PtrBool(false)},
	SlicePtrInt:     []*int{structs.PtrInt(1), structs.PtrInt(2), structs.PtrInt(3)},
	SlicePtrUint:    []*uint{structs.PtrUint(uint(1)), structs.PtrUint(uint(2)), structs.PtrUint(uint(3))},
	SlicePtrFloat:   []*float32{structs.PtrFloat32(1.1), structs.PtrFloat32(1.2), structs.PtrFloat32(1.3)},
	SlicePtrComplex: []*complex128{structs.PtrComplex128(complex(1, 1)), structs.PtrComplex128(complex(1, 2)), structs.PtrComplex128(complex(1, 3))},
	Hidden:          "abcdefg",
	unexported:      true,
}

var structX structTest = structTest{
	String:          "ozninoR",
	Bool:            false,
	Int:             3,
	Uint:            uint(654321),
	Float:           7622.50,
	Complex:         complex(-67, -42),
	Bytes:           []byte("Bye bye world"),
	Interface:       3.99,
	Error:           errors.New("not compliant"),
	Time:            time.Date(2021, time.August, 31, 14, 11, 11, 0, time.UTC),
	Duration:        30 * time.Second,
	NestedStruct:    structNested{Uint: 443211, String: "Microsoft IIS"},
	PtrString:       structs.PtrString("ozninoR"),
	PtrBool:         structs.PtrBool(false),
	PtrInt:          structs.PtrInt(3),
	PtrUint:         structs.PtrUint(uint(654321)),
	PtrFloat:        structs.PtrFloat32(7622.50),
	PtrComplex:      structs.PtrComplex128(complex(-67, -42)),
	PtrError:        structs.PtrError(errors.New("not compliant")),
	PtrTime:         structs.PtrTime(time.Date(2021, time.August, 31, 14, 11, 11, 0, time.UTC)),
	PtrDuration:     structs.PtrDuration(30 * time.Second),
	PtrNestedStruct: &structNested{Uint: 443211, String: "Microsoft IIS"},
	MapString:       map[string]string{"D": "four", "E": "five", "F": "six"},
	MapBool:         map[string]bool{"D": false, "E": true},
	MapInt:          map[string]int{"D": 4, "E": 5, "F": 6},
	MapUint:         map[string]uint{"D": uint(4), "E": uint(5), "F": uint(6)},
	MapFloat:        map[string]float32{"D": 1.4, "E": 1.5, "F": 1.6},
	MapComplex:      map[string]complex128{"D": complex(1, 4), "E": complex(1, 5), "F": complex(1, 6)},
	MapInterface:    map[string]interface{}{"D": 4, "E": "five", "F": 6.0},
	SliceString:     []string{"four", "five", "six"},
	SliceBool:       []bool{false, true},
	SliceInt:        []int{4, 5, 6},
	SliceUint:       []uint{uint(4), uint(5), uint(6)},
	SliceFloat:      []float32{1.4, 1.5, 1.6},
	SliceComplex:    []complex128{complex(1, 4), complex(1, 5), complex(1, 6)},
	SliceInterface:  []interface{}{4, "five", 6.0},
	SlicePtrString:  []*string{structs.PtrString("four"), structs.PtrString("five"), structs.PtrString("six")},
	SlicePtrBool:    []*bool{structs.PtrBool(false), structs.PtrBool(true)},
	SlicePtrInt:     []*int{structs.PtrInt(4), structs.PtrInt(5), structs.PtrInt(6)},
	SlicePtrUint:    []*uint{structs.PtrUint(uint(4)), structs.PtrUint(uint(5)), structs.PtrUint(uint(6))},
	SlicePtrFloat:   []*float32{structs.PtrFloat32(1.4), structs.PtrFloat32(1.5), structs.PtrFloat32(1.6)},
	SlicePtrComplex: []*complex128{structs.PtrComplex128(complex(1, 4)), structs.PtrComplex128(complex(1, 5)), structs.PtrComplex128(complex(1, 6))},
	Hidden:          "gfedcba",
	unexported:      false,
}

/*   S t r u c t F i e l d   -   L o c a l   F u n c t i o n s   */

func printStruct(s *structTest) {
	// fmt.Printf("%s", structs.Dump(t))
	// if err := json.NewEncoder(os.Stdout).Encode(s); err != nil {
	// 	panic(err)
	// }
	format := "- %15s: %v.\n"
	formatStruct := "- %15s: %+v.\n"
	formatPointer := "- %15s: *%v.\n"
	formatPtrStruct := "- %15s: *%+v.\n"
	fmt.Printf(format, "String", s.String)
	fmt.Printf(format, "Bool", s.Bool)
	fmt.Printf(format, "Int", s.Int)
	fmt.Printf(format, "Uint", s.Uint)
	fmt.Printf(format, "Float", s.Float)
	fmt.Printf(format, "Complex", s.Complex)
	fmt.Printf(format, "Bytes", string(s.Bytes))
	fmt.Printf(format, "Interface", s.Interface)
	fmt.Printf(format, "Error", s.Error)
	fmt.Printf(format, "Time", s.Time)
	fmt.Printf(format, "Duration", s.Duration)
	fmt.Printf(formatStruct, "NestedStruct", s.NestedStruct)
	if s.PtrString != nil {
		fmt.Printf(formatPointer, "PtrString", *s.PtrString)
	} else {
		fmt.Printf(format, "PtrString", s.PtrString)
	}
	if s.PtrBool != nil {
		fmt.Printf(formatPointer, "PtrBool", *s.PtrBool)
	} else {
		fmt.Printf(format, "PtrBool", s.PtrBool)
	}
	if s.PtrInt != nil {
		fmt.Printf(formatPointer, "PtrInt", *s.PtrInt)
	} else {
		fmt.Printf(format, "PtrInt", s.PtrInt)
	}
	if s.PtrUint != nil {
		fmt.Printf(formatPointer, "PtrUint", *s.PtrUint)
	} else {
		fmt.Printf(format, "PtrUint", s.PtrUint)
	}
	if s.PtrFloat != nil {
		fmt.Printf(formatPointer, "PtrFloat", *s.PtrFloat)
	} else {
		fmt.Printf(format, "PtrFloat", s.PtrFloat)
	}
	if s.PtrComplex != nil {
		fmt.Printf(formatPointer, "PtrComplex", *s.PtrComplex)
	} else {
		fmt.Printf(format, "PtrComplex", s.PtrComplex)
	}
	if s.PtrError != nil {
		fmt.Printf(formatPointer, "PtrError", *s.PtrError)
	} else {
		fmt.Printf(format, "PtrError", s.PtrError)
	}
	if s.PtrTime != nil {
		fmt.Printf(formatPointer, "PtrTime", *s.PtrTime)
	} else {
		fmt.Printf(format, "PtrTime", s.PtrTime)
	}
	if s.PtrDuration != nil {
		fmt.Printf(formatPointer, "PtrDuration", *s.PtrDuration)
	} else {
		fmt.Printf(format, "PtrDuration", s.PtrDuration)
	}
	if s.PtrNestedStruct != nil {
		fmt.Printf(formatPtrStruct, "PtrNestedStruct", *s.PtrNestedStruct)
	} else {
		fmt.Printf(formatStruct, "PtrNestedStruct", s.PtrNestedStruct)
	}
	fmt.Printf(format, "MapString", s.MapString)
	fmt.Printf(format, "MapBool", s.MapBool)
	fmt.Printf(format, "MapInt", s.MapInt)
	fmt.Printf(format, "MapUint", s.MapUint)
	fmt.Printf(format, "MapFloat", s.MapFloat)
	fmt.Printf(format, "MapComplex", s.MapComplex)
	fmt.Printf(format, "MapInterface", s.MapInterface)
	fmt.Printf(format, "SliceString", s.SliceString)
	fmt.Printf(format, "SliceBool", s.SliceBool)
	fmt.Printf(format, "SliceInt", s.SliceInt)
	fmt.Printf(format, "SliceUint", s.SliceUint)
	fmt.Printf(format, "SliceFloat", s.SliceFloat)
	fmt.Printf(format, "SliceComplex", s.SliceComplex)
	fmt.Printf(format, "SliceInterface", s.SliceInterface)
	{
		m := make([]interface{}, 0)
		for _, ptr := range s.SlicePtrString {
			if ptr != nil {
				v := fmt.Sprintf("*%v", *ptr)
				m = append(m, v)
				continue
			}
			m = append(m, nil)
		}
		fmt.Printf(format, "SlicePtrString", m)
	}
	{
		m := make([]interface{}, 0)
		for _, ptr := range s.SlicePtrBool {
			if ptr != nil {
				v := fmt.Sprintf("*%v", *ptr)
				m = append(m, v)
				continue
			}
			m = append(m, nil)
		}
		fmt.Printf(format, "SlicePtrBool", m)
	}
	{
		m := make([]interface{}, 0)
		for _, ptr := range s.SlicePtrInt {
			if ptr != nil {
				v := fmt.Sprintf("*%v", *ptr)
				m = append(m, v)
				continue
			}
			m = append(m, nil)
		}
		fmt.Printf(format, "SlicePtrInt", m)
	}
	{
		m := make([]interface{}, 0)
		for _, ptr := range s.SlicePtrUint {
			if ptr != nil {
				v := fmt.Sprintf("*%v", *ptr)
				m = append(m, v)
				continue
			}
			m = append(m, nil)
		}
		fmt.Printf(format, "SlicePtrUint", m)
	}
	{
		m := make([]interface{}, 0)
		for _, ptr := range s.SlicePtrFloat {
			if ptr != nil {
				v := fmt.Sprintf("*%v", *ptr)
				m = append(m, v)
				continue
			}
			m = append(m, nil)
		}
		fmt.Printf(format, "SlicePtrFloat", m)
	}
	{
		m := make([]interface{}, 0)
		for _, ptr := range s.SlicePtrComplex {
			if ptr != nil {
				v := fmt.Sprintf("*%v", *ptr)
				m = append(m, v)
				continue
			}
			m = append(m, nil)
		}
		fmt.Printf(format, "SlicePtrComplex", m)
	}
	fmt.Printf(format, "Hidden", s.Hidden)
	fmt.Printf(format, "unexported", s.unexported)
}

// func (t *structTest) MarshalJSON() ([]byte, error) {
// 	type Alias structTest
// 	return json.Marshal(&struct {
// 		Complex string `json:"complex"`
// 		// PtrComplex   string `json:"ptr_complex"`
// 		// MapComplex   string `json:"map_complex"`
// 		// SliceComplex string `json:"slice_complex"`
// 		// SlicePtrComplex string `json:"slice_ptr_complex"`
// 		Time    int64  `json:"date"`
// 		Bytes   string `json:"bytes"`
// 		*Alias
// 	}{
// 		Complex: structs.FormatComplex128(t.Complex),
// 		Time:    t.Date.Unix(),
// 		Bytes:   string(t.Bytes),
// 		Alias:   (*Alias)(t),
// 	})
// }

// func (t *structTest) UnmarshalJSON(data []byte) error {
// 	type Alias structTest
// 	aux := &struct {
// 		Complex string `json:"complex"`
// 		// PtrComplex   string `json:"ptr_complex"`
// 		// MapComplex   string `json:"map_complex"`
// 		// SliceComplex string `json:"slice_complex"`
// 		// SlicePtrComplex string `json:"slice_ptr_complex"`
// 		Time    int64  `json:"date"`
// 		Bytes   string `json:"bytes"`
// 		*Alias
// 	}{
// 		Alias: (*Alias)(t),
// 	}
// 	if err := json.Unmarshal(data, &aux); err != nil {
// 		return err
// 	}
// 	t.Complex = structs.ParseComplex128(aux.Complex)
// 	t.Time = time.Unix(aux.Time, 0)
// 	t.Bytes = []byte(aux.Bytes)
// 	return nil
// }

func ExampleStructField_Set_zeroToValue() {
	t := structTest{}

	s, err := structs.New(&t)
	if err != nil {
		fmt.Printf("New[Error]: %v.\n", err)
		return
	}

	err = s.Field("String").Set(structV.String)
	err = s.Field("Bool").Set(structV.Bool)
	err = s.Field("Int").Set(structV.Int)
	err = s.Field("Uint").Set(structV.Uint)
	err = s.Field("Float").Set(structV.Float)
	err = s.Field("Complex").Set(structV.Complex)
	err = s.Field("Bytes").Set(structV.Bytes)
	err = s.Field("Interface").Set(structV.Interface)
	err = s.Field("Error").Set(structV.Error)
	err = s.Field("Time").Set(structV.Time)
	err = s.Field("Duration").Set(structV.Duration)
	err = s.Field("NestedStruct").Set(structV.NestedStruct)
	err = s.Field("PtrString").Set(structV.PtrString)
	err = s.Field("PtrBool").Set(structV.PtrBool)
	err = s.Field("PtrInt").Set(structV.PtrInt)
	err = s.Field("PtrUint").Set(structV.PtrUint)
	err = s.Field("PtrFloat").Set(structV.PtrFloat)
	err = s.Field("PtrComplex").Set(structV.PtrComplex)
	err = s.Field("PtrError").Set(structV.PtrError)
	err = s.Field("PtrTime").Set(structV.PtrTime)
	err = s.Field("PtrDuration").Set(structV.PtrDuration)
	err = s.Field("PtrNestedStruct").Set(structV.PtrNestedStruct)
	err = s.Field("MapString").Set(structV.MapString)
	err = s.Field("MapBool").Set(structV.MapBool)
	err = s.Field("MapInt").Set(structV.MapInt)
	err = s.Field("MapUint").Set(structV.MapUint)
	err = s.Field("MapFloat").Set(structV.MapFloat)
	err = s.Field("MapComplex").Set(structV.MapComplex)
	err = s.Field("MapInterface").Set(structV.MapInterface)
	err = s.Field("SliceString").Set(structV.SliceString)
	err = s.Field("SliceBool").Set(structV.SliceBool)
	err = s.Field("SliceInt").Set(structV.SliceInt)
	err = s.Field("SliceUint").Set(structV.SliceUint)
	err = s.Field("SliceFloat").Set(structV.SliceFloat)
	err = s.Field("SliceComplex").Set(structV.SliceComplex)
	err = s.Field("SliceInterface").Set(structV.SliceInterface)
	err = s.Field("SlicePtrString").Set(structV.SlicePtrString)
	err = s.Field("SlicePtrBool").Set(structV.SlicePtrBool)
	err = s.Field("SlicePtrInt").Set(structV.SlicePtrInt)
	err = s.Field("SlicePtrUint").Set(structV.SlicePtrUint)
	err = s.Field("SlicePtrFloat").Set(structV.SlicePtrFloat)
	err = s.Field("SlicePtrComplex").Set(structV.SlicePtrComplex)
	err = s.Field("Hidden").Set(structV.Hidden)
	err = s.Field("unexported").Set(structV.unexported)
	if err != nil {
		fmt.Printf("Set[Error]: %v.\n", err)
	}
	printStruct(&t)

	// Output:
	// Set[Error]: could not set field structTest.unexported: struct field is not settable.
	// -          String: Roninzo.
	// -            Bool: true.
	// -             Int: 8.
	// -            Uint: 123456.
	// -           Float: 1922.5.
	// -         Complex: (22+50i).
	// -           Bytes: Hello world.
	// -       Interface: anything.
	// -           Error: rows not found.
	// -            Time: 2021-08-03 16:44:46 +0000 UTC.
	// -        Duration: 5s.
	// -    NestedStruct: {Uint:122334 String:Apache}.
	// -       PtrString: *Roninzo.
	// -         PtrBool: *true.
	// -          PtrInt: *8.
	// -         PtrUint: *123456.
	// -        PtrFloat: *1922.5.
	// -      PtrComplex: *(22+50i).
	// -        PtrError: *rows not found.
	// -         PtrTime: *2021-08-03 16:44:46 +0000 UTC.
	// -     PtrDuration: *5s.
	// - PtrNestedStruct: *{Uint:122334 String:Apache}.
	// -       MapString: map[A:one B:two C:three].
	// -         MapBool: map[A:true B:false].
	// -          MapInt: map[A:1 B:2 C:3].
	// -         MapUint: map[A:1 B:2 C:3].
	// -        MapFloat: map[A:1.1 B:1.2 C:1.3].
	// -      MapComplex: map[A:(1+1i) B:(1+2i) C:(1+3i)].
	// -    MapInterface: map[A:1 B:two C:3].
	// -     SliceString: [one two three].
	// -       SliceBool: [true false].
	// -        SliceInt: [1 2 3].
	// -       SliceUint: [1 2 3].
	// -      SliceFloat: [1.1 1.2 1.3].
	// -    SliceComplex: [(1+1i) (1+2i) (1+3i)].
	// -  SliceInterface: [1 two 3].
	// -  SlicePtrString: [*one *two *three].
	// -    SlicePtrBool: [*true *false].
	// -     SlicePtrInt: [*1 *2 *3].
	// -    SlicePtrUint: [*1 *2 *3].
	// -   SlicePtrFloat: [*1.1 *1.2 *1.3].
	// - SlicePtrComplex: [*(1+1i) *(1+2i) *(1+3i)].
	// -          Hidden: abcdefg.
	// -      unexported: false.
}

func ExampleStructField_Set_valueToValue() {
	t := structV

	s, err := structs.New(&t)
	if err != nil {
		fmt.Printf("New[Error]: %v.\n", err)
		return
	}

	err = s.Field("String").Set(structX.String)
	err = s.Field("Bool").Set(structX.Bool)
	err = s.Field("Int").Set(structX.Int)
	err = s.Field("Uint").Set(structX.Uint)
	err = s.Field("Float").Set(structX.Float)
	err = s.Field("Complex").Set(structX.Complex)
	err = s.Field("Bytes").Set(structX.Bytes)
	err = s.Field("Interface").Set(structX.Interface)
	err = s.Field("Error").Set(structX.Error)
	err = s.Field("Time").Set(structX.Time)
	err = s.Field("Duration").Set(structX.Duration)
	err = s.Field("NestedStruct").Set(structX.NestedStruct)
	err = s.Field("PtrString").Set(structX.PtrString)
	err = s.Field("PtrBool").Set(structX.PtrBool)
	err = s.Field("PtrInt").Set(structX.PtrInt)
	err = s.Field("PtrUint").Set(structX.PtrUint)
	err = s.Field("PtrFloat").Set(structX.PtrFloat)
	err = s.Field("PtrComplex").Set(structX.PtrComplex)
	err = s.Field("PtrError").Set(structX.PtrError)
	err = s.Field("PtrTime").Set(structX.PtrTime)
	err = s.Field("PtrDuration").Set(structX.PtrDuration)
	err = s.Field("PtrNestedStruct").Set(structX.PtrNestedStruct)
	err = s.Field("MapString").Set(structX.MapString)
	err = s.Field("MapBool").Set(structX.MapBool)
	err = s.Field("MapInt").Set(structX.MapInt)
	err = s.Field("MapUint").Set(structX.MapUint)
	err = s.Field("MapFloat").Set(structX.MapFloat)
	err = s.Field("MapComplex").Set(structX.MapComplex)
	err = s.Field("MapInterface").Set(structX.MapInterface)
	err = s.Field("SliceString").Set(structX.SliceString)
	err = s.Field("SliceBool").Set(structX.SliceBool)
	err = s.Field("SliceInt").Set(structX.SliceInt)
	err = s.Field("SliceUint").Set(structX.SliceUint)
	err = s.Field("SliceFloat").Set(structX.SliceFloat)
	err = s.Field("SliceComplex").Set(structX.SliceComplex)
	err = s.Field("SliceInterface").Set(structX.SliceInterface)
	err = s.Field("SlicePtrString").Set(structX.SlicePtrString)
	err = s.Field("SlicePtrBool").Set(structX.SlicePtrBool)
	err = s.Field("SlicePtrInt").Set(structX.SlicePtrInt)
	err = s.Field("SlicePtrUint").Set(structX.SlicePtrUint)
	err = s.Field("SlicePtrFloat").Set(structX.SlicePtrFloat)
	err = s.Field("SlicePtrComplex").Set(structX.SlicePtrComplex)
	err = s.Field("Hidden").Set(structX.Hidden)
	err = s.Field("unexported").Set(structX.unexported)
	if err != nil {
		fmt.Printf("Set[Error]: %v.\n", err)
	}
	printStruct(&t)

	// Output:
	// Set[Error]: could not set field structTest.unexported: struct field is not settable.
	// -          String: ozninoR.
	// -            Bool: false.
	// -             Int: 3.
	// -            Uint: 654321.
	// -           Float: 7622.5.
	// -         Complex: (-67-42i).
	// -           Bytes: Bye bye world.
	// -       Interface: 3.99.
	// -           Error: not compliant.
	// -            Time: 2021-08-31 14:11:11 +0000 UTC.
	// -        Duration: 30s.
	// -    NestedStruct: {Uint:443211 String:Microsoft IIS}.
	// -       PtrString: *ozninoR.
	// -         PtrBool: *false.
	// -          PtrInt: *3.
	// -         PtrUint: *654321.
	// -        PtrFloat: *7622.5.
	// -      PtrComplex: *(-67-42i).
	// -        PtrError: *not compliant.
	// -         PtrTime: *2021-08-31 14:11:11 +0000 UTC.
	// -     PtrDuration: *30s.
	// - PtrNestedStruct: *{Uint:443211 String:Microsoft IIS}.
	// -       MapString: map[D:four E:five F:six].
	// -         MapBool: map[D:false E:true].
	// -          MapInt: map[D:4 E:5 F:6].
	// -         MapUint: map[D:4 E:5 F:6].
	// -        MapFloat: map[D:1.4 E:1.5 F:1.6].
	// -      MapComplex: map[D:(1+4i) E:(1+5i) F:(1+6i)].
	// -    MapInterface: map[D:4 E:five F:6].
	// -     SliceString: [four five six].
	// -       SliceBool: [false true].
	// -        SliceInt: [4 5 6].
	// -       SliceUint: [4 5 6].
	// -      SliceFloat: [1.4 1.5 1.6].
	// -    SliceComplex: [(1+4i) (1+5i) (1+6i)].
	// -  SliceInterface: [4 five 6].
	// -  SlicePtrString: [*four *five *six].
	// -    SlicePtrBool: [*false *true].
	// -     SlicePtrInt: [*4 *5 *6].
	// -    SlicePtrUint: [*4 *5 *6].
	// -   SlicePtrFloat: [*1.4 *1.5 *1.6].
	// - SlicePtrComplex: [*(1+4i) *(1+5i) *(1+6i)].
	// -          Hidden: gfedcba.
	// -      unexported: true.
}

func ExampleStructField_Set_valueToZero() {
	t := structX

	s, err := structs.New(&t)
	if err != nil {
		fmt.Printf("New[Error]: %v.\n", err)
		return
	}

	for _, name := range structFieldNames {
		err := s.Field(name).SetZero()
		if err != nil {
			fmt.Printf("SetZero[Error]: %v.\n", err)
		}
	}
	printStruct(&t)

	// Output:
	// SetZero[Error]: could not set field structTest.unexported to zero-value: struct field is not settable.
	// -          String: .
	// -            Bool: false.
	// -             Int: 0.
	// -            Uint: 0.
	// -           Float: 0.
	// -         Complex: (0+0i).
	// -           Bytes: .
	// -       Interface: <nil>.
	// -           Error: <nil>.
	// -            Time: 0001-01-01 00:00:00 +0000 UTC.
	// -        Duration: 0s.
	// -    NestedStruct: {Uint:0 String:}.
	// -       PtrString: <nil>.
	// -         PtrBool: <nil>.
	// -          PtrInt: <nil>.
	// -         PtrUint: <nil>.
	// -        PtrFloat: <nil>.
	// -      PtrComplex: <nil>.
	// -        PtrError: <nil>.
	// -         PtrTime: <nil>.
	// -     PtrDuration: <nil>.
	// - PtrNestedStruct: <nil>.
	// -       MapString: map[].
	// -         MapBool: map[].
	// -          MapInt: map[].
	// -         MapUint: map[].
	// -        MapFloat: map[].
	// -      MapComplex: map[].
	// -    MapInterface: map[].
	// -     SliceString: [].
	// -       SliceBool: [].
	// -        SliceInt: [].
	// -       SliceUint: [].
	// -      SliceFloat: [].
	// -    SliceComplex: [].
	// -  SliceInterface: [].
	// -  SlicePtrString: [].
	// -    SlicePtrBool: [].
	// -     SlicePtrInt: [].
	// -    SlicePtrUint: [].
	// -   SlicePtrFloat: [].
	// - SlicePtrComplex: [].
	// -          Hidden: .
	// -      unexported: false.
}

func ExampleStructField_Set_valueToNil() {
	t := structX

	s, err := structs.New(&t)
	if err != nil {
		fmt.Printf("New[Error]: %v.\n", err)
		return
	}

	for _, name := range structFieldNames {
		err := s.Field(name).SetNil()
		if err != nil {
			fmt.Printf("SetNil[Error]: %v.\n", err)
		}
	}
	printStruct(&t)

	// Output:
	// SetNil[Error]: could not set field structTest.String to nil: struct field is not nillable.
	// SetNil[Error]: could not set field structTest.Bool to nil: struct field is not nillable.
	// SetNil[Error]: could not set field structTest.Int to nil: struct field is not nillable.
	// SetNil[Error]: could not set field structTest.Uint to nil: struct field is not nillable.
	// SetNil[Error]: could not set field structTest.Float to nil: struct field is not nillable.
	// SetNil[Error]: could not set field structTest.Complex to nil: struct field is not nillable.
	// SetNil[Error]: could not set field structTest.Time to nil: struct field is not nillable.
	// SetNil[Error]: could not set field structTest.Duration to nil: struct field is not nillable.
	// SetNil[Error]: could not set field structTest.NestedStruct to nil: struct field is not nillable.
	// SetNil[Error]: could not set field structTest.Hidden to nil: struct field is not nillable.
	// SetNil[Error]: could not set field structTest.unexported to nil: struct field is not settable.
	// -          String: ozninoR.
	// -            Bool: false.
	// -             Int: 3.
	// -            Uint: 654321.
	// -           Float: 7622.5.
	// -         Complex: (-67-42i).
	// -           Bytes: .
	// -       Interface: <nil>.
	// -           Error: <nil>.
	// -            Time: 2021-08-31 14:11:11 +0000 UTC.
	// -        Duration: 30s.
	// -    NestedStruct: {Uint:443211 String:Microsoft IIS}.
	// -       PtrString: <nil>.
	// -         PtrBool: <nil>.
	// -          PtrInt: <nil>.
	// -         PtrUint: <nil>.
	// -        PtrFloat: <nil>.
	// -      PtrComplex: <nil>.
	// -        PtrError: <nil>.
	// -         PtrTime: <nil>.
	// -     PtrDuration: <nil>.
	// - PtrNestedStruct: <nil>.
	// -       MapString: map[].
	// -         MapBool: map[].
	// -          MapInt: map[].
	// -         MapUint: map[].
	// -        MapFloat: map[].
	// -      MapComplex: map[].
	// -    MapInterface: map[].
	// -     SliceString: [].
	// -       SliceBool: [].
	// -        SliceInt: [].
	// -       SliceUint: [].
	// -      SliceFloat: [].
	// -    SliceComplex: [].
	// -  SliceInterface: [].
	// -  SlicePtrString: [].
	// -    SlicePtrBool: [].
	// -     SlicePtrInt: [].
	// -    SlicePtrUint: [].
	// -   SlicePtrFloat: [].
	// - SlicePtrComplex: [].
	// -          Hidden: gfedcba.
	// -      unexported: false.
}

func ExampleStructField_Get() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Count      int32  `json:"count,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Count:      0,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Count")
	f3 := s.Field("unexported")
	f4 := s.Field("Password")
	f5 := s.Field("Undeclared")

	err = s.Err()
	if f5 == nil && err != nil {
		fmt.Printf("Error: %v.\n", err)
	}

	v1, err1 := f1.Get()
	v2, err2 := f2.Get()
	v3, err3 := f3.Get()
	v4, err4 := f4.Get()

	fmt.Printf("Value %-11s: %v.\n", f1.Name(), f1.Value())
	fmt.Printf("Value %-11s: %v.\n", f2.Name(), f2.Value())
	fmt.Printf("Value %-11s: %v.\n", f3.Name(), f3.Value())
	fmt.Printf("Value %-11s: %v.\n", f4.Name(), f4.Value())
	fmt.Printf("Get   %-11s: %-11v, err: %v.\n", f1.Name(), v1, err1)
	fmt.Printf("Get   %-11s: %-11v, err: %v.\n", f2.Name(), v2, err2)
	fmt.Printf("Get   %-11s: %-11v, err: %v.\n", f3.Name(), v3, err3)
	fmt.Printf("Get   %-11s: %-11v, err: %v.\n", f4.Name(), v4, err4)

	// Output:
	// Error: invalid field name Undeclared.
	// Value Name       : Roninzo.
	// Value Count      : 0.
	// Value unexported : false.
	// Value Password   : abcdefg.
	// Get   Name       : Roninzo    , err: <nil>.
	// Get   Count      : 0          , err: <nil>.
	// Get   unexported : <nil>      , err: could not get value of field Server.unexported: struct field is not exported.
	// Get   Password   : abcdefg    , err: <nil>.
}

func ExampleStructField_SetZero() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Count      int32  `json:"count,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Count:      0,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Count")
	f3 := s.Field("unexported")

	fmt.Printf("Value: %v.\n", f1.Value())
	fmt.Printf("Value: %v.\n", f2.Value())
	fmt.Printf("Value: %v.\n", f3.Value())

	err1 := f1.SetZero()
	err2 := f2.SetZero()
	err3 := f3.SetZero()

	fmt.Printf("SetZero: %v.\n", f1.Value())
	fmt.Printf("SetZero: %v.\n", f2.Value())
	fmt.Printf("SetZero: %v.\n", f3.Value())
	fmt.Printf("Error: %v.\n", err1)
	fmt.Printf("Error: %v.\n", err2)
	fmt.Printf("Error: %v.\n", err3)

	// Output:
	// Value: Roninzo.
	// Value: 0.
	// Value: false.
	// SetZero: .
	// SetZero: 0.
	// SetZero: false.
	// Error: <nil>.
	// Error: <nil>.
	// Error: could not set field Server.unexported to zero-value: struct field is not settable.
}

func ExampleStructField_CanTime() {
	type Server struct {
		Name       string    `json:"name,omitempty"`
		ID         uint      `json:"id,omitempty"`
		Enabled    bool      `json:"enabled,omitempty"`
		ChangedAt  time.Time `json:"changed_at,omitempty"`
		Password   string    `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		ChangedAt:  time.Date(2021, time.August, 3, 16, 44, 46, 0, time.UTC),
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("ChangedAt")

	fmt.Printf("CanTime: %v\n", f1.CanTime())
	fmt.Printf("CanTime: %v\n", f2.CanTime())

	// Output:
	// CanTime: false
	// CanTime: true
}

func ExampleStructField_CanDuration() {
	type Server struct {
		Name       string        `json:"name,omitempty"`
		ID         uint          `json:"id,omitempty"`
		Enabled    bool          `json:"enabled,omitempty"`
		TimeOut    time.Duration `json:"time_out,omitempty"`
		Password   string        `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		TimeOut:    5 * time.Second,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("TimeOut")

	fmt.Printf("CanDuration: %v\n", f1.CanDuration())
	fmt.Printf("CanDuration: %v\n", f2.CanDuration())

	// Output:
	// CanDuration: false
	// CanDuration: true
}

func ExampleStructField_CanError() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Err        error  `json:"err,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Err:        errors.New("failed"),
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Err")

	fmt.Printf("CanError: %v\n", f1.CanError())
	fmt.Printf("CanError: %v\n", f2.CanError())

	// Output:
	// CanError: false
	// CanError: true
}

func ExampleStructField_CanString() {
	type Server struct {
		Name       string    `json:"name,omitempty"`
		ID         uint      `json:"id,omitempty"`
		Enabled    bool      `json:"enabled,omitempty"`
		ChangedAt  time.Time `json:"changed_at,omitempty"`
		Password   string    `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		ChangedAt:  time.Date(2021, time.August, 3, 16, 44, 46, 0, time.UTC),
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("ChangedAt")

	fmt.Printf("CanString: %v\n", f1.CanString())
	fmt.Printf("CanString: %v\n", f2.CanString())

	// Output:
	// CanString: true
	// CanString: false
}

func ExampleStructField_CanBool() {
	type Server struct {
		Name       string    `json:"name,omitempty"`
		ID         uint      `json:"id,omitempty"`
		Enabled    bool      `json:"enabled,omitempty"`
		ChangedAt  time.Time `json:"changed_at,omitempty"`
		Password   string    `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		ChangedAt:  time.Date(2021, time.August, 3, 16, 44, 46, 0, time.UTC),
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Enabled")

	fmt.Printf("CanBool: %v\n", f1.CanBool())
	fmt.Printf("CanBool: %v\n", f2.CanBool())

	// Output:
	// CanBool: false
	// CanBool: true
}

func ExampleStructField_CanInt() {
	type Server struct {
		Name       string    `json:"name,omitempty"`
		ID         int       `json:"id,omitempty"`
		Enabled    bool      `json:"enabled,omitempty"`
		ChangedAt  time.Time `json:"changed_at,omitempty"`
		Password   string    `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		ChangedAt:  time.Date(2021, time.August, 3, 16, 44, 46, 0, time.UTC),
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("ID")

	fmt.Printf("CanInt: %v\n", f1.CanInt())
	fmt.Printf("CanInt: %v\n", f2.CanInt())

	// Output:
	// CanInt: false
	// CanInt: true
}

func ExampleStructField_CanUint() {
	type Server struct {
		Name       string    `json:"name,omitempty"`
		ID         uint      `json:"id,omitempty"`
		Enabled    bool      `json:"enabled,omitempty"`
		ChangedAt  time.Time `json:"changed_at,omitempty"`
		Password   string    `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		ChangedAt:  time.Date(2021, time.August, 3, 16, 44, 46, 0, time.UTC),
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("ID")

	fmt.Printf("CanUint: %v\n", f1.CanUint())
	fmt.Printf("CanUint: %v\n", f2.CanUint())

	// Output:
	// CanUint: false
	// CanUint: true
}

func ExampleStructField_CanFloat() {
	type Server struct {
		Name       string  `json:"name,omitempty"`
		ID         uint    `json:"id,omitempty"`
		Enabled    bool    `json:"enabled,omitempty"`
		Price      float32 `json:"price,omitempty"`
		Password   string  `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Price:      10.50,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Price")

	fmt.Printf("CanFloat: %v\n", f1.CanFloat())
	fmt.Printf("CanFloat: %v\n", f2.CanFloat())

	// Output:
	// CanFloat: false
	// CanFloat: true
}

func ExampleStructField_CanComplex() {
	type Server struct {
		Name       string     `json:"name,omitempty"`
		ID         uint       `json:"id,omitempty"`
		Enabled    bool       `json:"enabled,omitempty"`
		Complex    complex128 `json:"complex,omitempty"`
		Password   string     `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Complex:    complex(23, 31),
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Complex")

	fmt.Printf("CanComplex: %v\n", f1.CanComplex())
	fmt.Printf("CanComplex: %v\n", f2.CanComplex())

	// Output:
	// CanComplex: false
	// CanComplex: true
}

func ExampleStructField_CanBytes() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Stream     []byte `json:"stream,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Stream:     []byte("Hello world"),
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Stream")

	fmt.Printf("CanBytes: %v\n", f1.CanBytes())
	fmt.Printf("CanBytes: %v\n", f2.CanBytes())

	// Output:
	// CanBytes: false
	// CanBytes: true
}

func ExampleStructField_CanInterface() {
	type Server struct {
		Name       string      `json:"name,omitempty"`
		ID         uint        `json:"id,omitempty"`
		Enabled    bool        `json:"enabled,omitempty"`
		Anything   interface{} `json:"anything,omitempty"`
		Password   string      `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Anything:   654321,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Anything")

	fmt.Printf("CanInterface: %v\n", f1.CanInterface())
	fmt.Printf("CanInterface: %v\n", f2.CanInterface())

	// Output:
	// CanInterface: false
	// CanInterface: true
}

func ExampleStructField_CanStruct() {
	type Program struct {
		Name string `json:"name,omitempty"`
	}

	type Server struct {
		Name     string  `json:"name,omitempty"`
		ID       uint    `json:"id,omitempty"`
		Enabled  bool    `json:"enabled,omitempty"`
		Count    int32   `json:"count,omitempty"`
		Password string  `json:"-"`
		Program  Program `json:"program,omitempty"`
	}

	program := Program{
		Name: "Apache",
	}

	server := Server{
		Name:     "Roninzo",
		ID:       123456,
		Enabled:  true,
		Count:    0,
		Password: "abcdefg",
		Program:  program,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Program")

	fmt.Printf("CanStruct: %v\n", f1.CanStruct())
	fmt.Printf("CanStruct: %v\n", f2.CanStruct())

	// Output:
	// CanStruct: false
	// CanStruct: true
}

func ExampleStructField_Time() {
	type Server struct {
		Name       string    `json:"name,omitempty"`
		ID         uint      `json:"id,omitempty"`
		Enabled    bool      `json:"enabled,omitempty"`
		ChangedAt  time.Time `json:"changed_at,omitempty"`
		Password   string    `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		ChangedAt:  time.Date(2021, time.August, 3, 16, 44, 46, 0, time.UTC),
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("ChangedAt")

	if f1.CanTime() {
		fmt.Printf("Time: %v\n", f1.Time())
	}
	if f2.CanTime() {
		fmt.Printf("Time: %v\n", f2.Time())
	}

	// Output:
	// Time: 2021-08-03 16:44:46 +0000 UTC
}

func ExampleStructField_Duration() {
	type Server struct {
		Name       string        `json:"name,omitempty"`
		ID         uint          `json:"id,omitempty"`
		Enabled    bool          `json:"enabled,omitempty"`
		TimeOut    time.Duration `json:"time_out,omitempty"`
		Password   string        `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		TimeOut:    5 * time.Second,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("TimeOut")

	if f1.CanDuration() {
		fmt.Printf("Duration: %v\n", f1.Duration())
	}
	if f2.CanDuration() {
		fmt.Printf("Duration: %v\n", f2.Duration())
	}

	// Output:
	// Duration: 5s
}

func ExampleStructField_Error() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Err        error  `json:"err,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Err:        errors.New("failed"),
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Err")

	if f1.CanError() {
		fmt.Printf("Error: %v\n", f1.Error())
	}
	if f2.CanError() {
		fmt.Printf("Error: %v\n", f2.Error())
	}

	// Output:
	// Error: failed
}

func ExampleStructField_String() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field("ID")
	f2 := s.Field("Name")

	if f1.CanString() {
		fmt.Printf("String: %v.\n", f1.String())
	}
	if f2.CanString() {
		fmt.Printf("String: %v.\n", f2.String())
	}

	// Output:
	// String: Roninzo.
}

func ExampleStructField_Bool() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Enabled")

	if f1.CanBool() {
		fmt.Printf("Bool: %v\n", f1.Bool())
	}
	if f2.CanBool() {
		fmt.Printf("Bool: %v\n", f2.Bool())
	}

	// Output:
	// Bool: true
}

func ExampleStructField_Int() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         int    `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("ID")

	if f1.CanInt() {
		fmt.Printf("Int: %v\n", f1.Int())
	}
	if f2.CanInt() {
		fmt.Printf("Int: %v\n", f2.Int())
	}

	// Output:
	// Int: 123456
}

func ExampleStructField_Uint() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("ID")

	if f1.CanUint() {
		fmt.Printf("Uint: %v\n", f1.Uint())
	}
	if f2.CanUint() {
		fmt.Printf("Uint: %v\n", f2.Uint())
	}

	// Output:
	// Uint: 123456
}

func ExampleStructField_Float() {
	type Server struct {
		Name       string  `json:"name,omitempty"`
		ID         uint    `json:"id,omitempty"`
		Enabled    bool    `json:"enabled,omitempty"`
		Price      float32 `json:"price,omitempty"`
		Password   string  `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Price:      22.50,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Price")

	if f1.CanFloat() {
		fmt.Printf("Float: %v\n", f1.Float())
	}
	if f2.CanFloat() {
		fmt.Printf("Float: %v\n", f2.Float())
	}

	// Output:
	// Float: 22.5
}

func ExampleStructField_Complex() {
	type Server struct {
		Name       string     `json:"name,omitempty"`
		ID         uint       `json:"id,omitempty"`
		Enabled    bool       `json:"enabled,omitempty"`
		Complex    complex128 `json:"complex,omitempty"`
		Password   string     `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Complex:    complex(22, 50),
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Complex")

	if f1.CanComplex() {
		fmt.Printf("Complex: %v.\n", f1.Complex())
	}
	if f2.CanComplex() {
		fmt.Printf("Complex: %v.\n", f2.Complex())
	}

	// Output:
	// Complex: (22+50i).
}

func ExampleStructField_Bytes() {
	type Server struct {
		Name       string `json:"name,omitempty"`
		ID         uint   `json:"id,omitempty"`
		Enabled    bool   `json:"enabled,omitempty"`
		Stream     []byte `json:"stream,omitempty"`
		Password   string `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Stream:     []byte("Hello world"),
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Stream")

	if f1.CanBytes() {
		fmt.Printf("Bytes: %v.\n", f1.Bytes())
	}
	if f2.CanBytes() {
		b := f2.Bytes()
		fmt.Printf("Bytes: %v.\n", b)
		fmt.Printf("BytesString: %v.\n", string(b))
	}

	// Output:
	// Bytes: [72 101 108 108 111 32 119 111 114 108 100].
	// BytesString: Hello world.
}

func ExampleStructField_Interface() {
	type Server struct {
		Name       string      `json:"name,omitempty"`
		ID         uint        `json:"id,omitempty"`
		Enabled    bool        `json:"enabled,omitempty"`
		Anything   interface{} `json:"anything,omitempty"`
		Password   string      `json:"-"`
		unexported bool
	}

	server := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Anything:   654321,
		Password:   "abcdefg",
		unexported: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f1 := s.Field(0)
	f2 := s.Field("Anything")

	if f1.CanInterface() {
		fmt.Printf("Interface: %v.\n", f1.Interface())
	}
	if f2.CanInterface() {
		fmt.Printf("Interface: %v.\n", f2.Interface())
	}

	// Output:
	// Interface: 654321.
}

func ExampleStructField_Struct() {
	type Program struct {
		Name string `json:"name,omitempty"`
	}

	type Server struct {
		Name     string  `json:"name,omitempty"`
		ID       uint    `json:"id,omitempty"`
		Enabled  bool    `json:"enabled,omitempty"`
		Count    int32   `json:"count,omitempty"`
		Password string  `json:"-"`
		Program  Program `json:"program,omitempty"`
	}

	program := Program{
		Name: "Apache",
	}

	server := Server{
		Name:     "Roninzo",
		ID:       123456,
		Enabled:  true,
		Count:    0,
		Password: "abcdefg",
		Program:  program,
	}

	s1, _ := structs.New(&server)
	s2 := s1.Field("Program").Struct()

	fmt.Printf("Struct: %v\n", s1.Name())
	fmt.Printf("Struct: %v\n", s2.Name())

	// Output:
	// Struct: Server
	// Struct: Program
}

func ExampleStructField_SetTime() {
	type Server struct {
		ChangedAt time.Time `json:"changed_at,omitempty"`
	}

	server := Server{
		ChangedAt: time.Date(2021, time.August, 3, 16, 44, 46, 0, time.UTC),
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f := s.Field("ChangedAt")

	if f.CanTime() {
		fmt.Printf("Time: %v.\n", f.Time())
		t := time.Date(2021, time.August, 31, 12, 30, 11, 0, time.UTC)
		f.SetTime(t)
		fmt.Printf("SetTime: %v.\n", f.Time())
	}

	// Output:
	// Time: 2021-08-03 16:44:46 +0000 UTC.
	// SetTime: 2021-08-31 12:30:11 +0000 UTC.
}

func ExampleStructField_SetDuration() {
	type Server struct {
		TimeOut time.Duration `json:"time_out,omitempty"`
	}

	server := Server{
		TimeOut: 5 * time.Second,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f := s.Field("TimeOut")

	if f.CanDuration() {
		fmt.Printf("Duration: %v.\n", f.Duration())
		f.SetDuration(30 * time.Second)
		fmt.Printf("SetDuration: %v.\n", f.Duration())
	}

	// Output:
	// Duration: 5s.
	// SetDuration: 30s.
}

func ExampleStructField_SetError() {
	type Server struct {
		Err error `json:"err,omitempty"`
	}

	server := Server{
		Err: errors.New("rows not found"),
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f := s.Field("Err")

	if f.CanError() {
		fmt.Printf("Error: %v.\n", f.Error())
		err := errors.Wrap(f.Error(), "empty table")
		f.SetError(err)
		fmt.Printf("SetError: %v.\n", f.Error())
	}

	// Output:
	// Error: rows not found.
	// SetError: empty table: rows not found.
}

func ExampleStructField_SetString() {
	type Server struct {
		Name string `json:"name,omitempty"`
	}

	server := Server{
		Name: "Roninzo",
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f := s.Field("Name")

	if f.CanString() {
		fmt.Printf("String: %v.\n", f.String())
		f.SetString("ozninoR")
		fmt.Printf("SetString: %v.\n", f.String())
	}

	// Output:
	// String: Roninzo.
	// SetString: ozninoR.
}

func ExampleStructField_SetBool() {
	type Server struct {
		Enabled bool `json:"enabled,omitempty"`
	}

	server := Server{
		Enabled: true,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f := s.Field("Enabled")

	if f.CanBool() {
		fmt.Printf("Bool: %v.\n", f.Bool())
		f.SetBool(false)
		fmt.Printf("SetBool: %v.\n", f.Bool())
	}

	// Output:
	// Bool: true.
	// SetBool: false.
}

func ExampleStructField_SetInt() {
	type Server struct {
		Count int `json:"count,omitempty"`
	}

	server := Server{
		Count: 123456,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f := s.Field("Count")

	if f.CanInt() {
		fmt.Printf("Int: %v.\n", f.Int())
		f.SetInt(654321)
		fmt.Printf("SetInt: %v.\n", f.Int())
	}

	// Output:
	// Int: 123456.
	// SetInt: 654321.
}

func ExampleStructField_SetUint() {
	type Server struct {
		ID uint `json:"id,omitempty"`
	}

	server := Server{
		ID: 123456,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f := s.Field("ID")

	if f.CanUint() {
		fmt.Printf("Uint: %v.\n", f.Uint())
		f.SetUint(654321)
		fmt.Printf("SetUint: %v.\n", f.Uint())
	}

	// Output:
	// Uint: 123456.
	// SetUint: 654321.
}

func ExampleStructField_SetFloat() {
	type Server struct {
		Price float32 `json:"price,omitempty"`
	}

	server := Server{
		Price: 22.50,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f := s.Field("Price")

	if f.CanFloat() {
		fmt.Printf("Float: %v.\n", f.Float())
		f.SetFloat(450.50)
		fmt.Printf("SetUint: %v.\n", f.Float())
	}

	// Output:
	// Float: 22.5.
	// SetUint: 450.5.
}

func ExampleStructField_SetComplex() {
	type Server struct {
		Complex complex128 `json:"complex,omitempty"`
	}

	server := Server{
		Complex: complex(22, 50),
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f := s.Field("Complex")

	if f.CanComplex() {
		fmt.Printf("Complex: %v.\n", f.Complex())
		f.SetComplex(complex(77, 2))
		fmt.Printf("SetComplex: %v.\n", f.Complex())
	}

	// Output:
	// Complex: (22+50i).
	// SetComplex: (77+2i).
}

func ExampleStructField_SetBytes() {
	type Server struct {
		Stream []byte `json:"stream,omitempty"`
	}

	server := Server{
		Stream: []byte("Hello world"),
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f := s.Field("Stream")

	if f.CanBytes() {
		b := f.Bytes()
		fmt.Printf("Bytes: %v.\n", b)
		fmt.Printf("BytesString: %v.\n", string(b))
		f.SetBytes([]byte("Bye bye world"))
		b = f.Bytes()
		fmt.Printf("SetBytes: %v.\n", b)
		fmt.Printf("SetBytesString: %v.\n", string(b))
	}

	// Output:
	// Bytes: [72 101 108 108 111 32 119 111 114 108 100].
	// BytesString: Hello world.
	// SetBytes: [66 121 101 32 98 121 101 32 119 111 114 108 100].
	// SetBytesString: Bye bye world.
}

func ExampleStructField_SetInterface() {
	type Server struct {
		Anything interface{} `json:"anything,omitempty"`
	}

	server := Server{
		Anything: 654321,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f := s.Field("Anything")

	if f.CanInterface() {
		fmt.Printf("Interface: %v.\n", f.Interface())
		f.SetInterface(123456)
		fmt.Printf("SetInterface: %v.\n", f.Interface())
	}

	// Output:
	// Interface: 654321.
	// SetInterface: 123456.
}

func ExampleStructField_SetStruct() {
	type Program struct {
		Name string `json:"name,omitempty"`
	}

	type Server struct {
		Program *Program `json:"program,omitempty"`
	}

	program := Program{
		Name: "Apache",
	}

	server := Server{
		Program: &program,
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	f := s.Field("Program")

	if f.CanStruct() {
		fmt.Printf("Struct: %v.\n", server.Program.Name)
		program2 := Program{
			Name: "Microsoft IIS",
		}
		p, err := structs.New(&program2)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		f.SetStruct(p)
		fmt.Printf("SetStruct: %v.\n", server.Program.Name)
	}

	// Output:
	// Struct: Apache.
	// SetStruct: Microsoft IIS.
}

func ExampleStructField_Equal() {
	type Program struct {
		Name string `json:"name,omitempty"`
	}

	program1 := Program{Name: "Apache"}
	program2 := Program{Name: "Apache"}

	type Server struct {
		Name       string        `json:"name,omitempty"`
		ID         uint          `json:"id,omitempty"`
		Enabled    bool          `json:"enabled,omitempty"`
		Err        error         `json:"err,omitempty"`
		ChangedAt  time.Time     `json:"changed_at,omitempty"`
		TimeOut    time.Duration `json:"time_out,omitempty"`
		Count      int           `json:"count,omitempty"`
		Price      float32       `json:"price,omitempty"`
		Complex    complex128    `json:"complex,omitempty"`
		Stream     []byte        `json:"stream,omitempty"`
		Anything   interface{}   `json:"anything,omitempty"`
		Program    *Program      `json:"program,omitempty"`
		Password   string        `json:"-"`
		unexported bool
	}

	server1 := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Err:        errors.New("rows not found"),
		ChangedAt:  time.Date(2021, time.August, 3, 16, 44, 46, 0, time.UTC),
		TimeOut:    5 * time.Second,
		Count:      8,
		Price:      1922.50,
		Complex:    complex(22, 50),
		Stream:     []byte("Hello world"),
		Anything:   654321,
		Program:    &program1,
		Password:   "abcdefg",
		unexported: true,
	}

	server2 := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Err:        errors.New("rows not found"),
		ChangedAt:  time.Date(2021, time.August, 3, 16, 44, 46, 0, time.UTC),
		TimeOut:    5 * time.Second,
		Count:      8,
		Price:      1922.50,
		Complex:    complex(22, 50),
		Stream:     []byte("Hello world"),
		Anything:   654321,
		Program:    &program2,
		Password:   "abcdefg",
		unexported: true,
	}

	s1, err := structs.New(&server1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	s2, err := structs.New(&server2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	for i, f1 := range s1.Fields() {
		f2 := s2.Field(i)
		fmt.Printf("Equal: %v.\n", f1.Equal(f2))
	}

	// Output:
	// Equal: true.
	// Equal: true.
	// Equal: true.
	// Equal: true.
	// Equal: true.
	// Equal: true.
	// Equal: true.
	// Equal: true.
	// Equal: true.
	// Equal: true.
	// Equal: true.
	// Equal: true.
	// Equal: true.
	// Equal: false.
}

func ExampleStructField_Equal_different() {
	type Program struct {
		Name string `json:"name,omitempty"`
	}

	program1 := Program{Name: "Apache"}
	program2 := Program{Name: "Microsoft IIS"}

	type Server struct {
		Name       string        `json:"name,omitempty"`
		ID         uint          `json:"id,omitempty"`
		Enabled    bool          `json:"enabled,omitempty"`
		Err        error         `json:"err,omitempty"`
		ChangedAt  time.Time     `json:"changed_at,omitempty"`
		TimeOut    time.Duration `json:"time_out,omitempty"`
		Count      int           `json:"count,omitempty"`
		Price      float32       `json:"price,omitempty"`
		Complex    complex128    `json:"complex,omitempty"`
		Stream     []byte        `json:"stream,omitempty"`
		Anything   interface{}   `json:"anything,omitempty"`
		Program    *Program      `json:"program,omitempty"`
		Password   string        `json:"-"`
		unexported bool
	}

	server1 := Server{
		Name:       "Roninzo",
		ID:         123456,
		Enabled:    true,
		Err:        errors.New("rows not found"),
		ChangedAt:  time.Date(2021, time.August, 3, 16, 44, 46, 0, time.UTC),
		TimeOut:    5 * time.Second,
		Count:      8,
		Price:      1922.50,
		Complex:    complex(22, 50),
		Stream:     []byte("Hello world"),
		Anything:   654321,
		Program:    &program1,
		Password:   "abcdefg",
		unexported: true,
	}

	server2 := Server{
		Name:       "ozninoR",
		ID:         654321,
		Enabled:    false,
		Err:        errors.New("not compliant"),
		ChangedAt:  time.Date(2021, time.August, 31, 14, 11, 11, 0, time.UTC),
		TimeOut:    30 * time.Second,
		Count:      3,
		Price:      7622.50,
		Complex:    complex(-67, -42),
		Stream:     []byte("Bye bye world"),
		Anything:   123456,
		Program:    &program2,
		Password:   "gfedcba",
		unexported: false,
	}

	s1, err := structs.New(&server1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	s2, err := structs.New(&server2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	for i, f1 := range s1.Fields() {
		f2 := s2.Field(i)
		fmt.Printf("Equal: %v.\n", f1.Equal(f2))
	}

	// Output:
	// Equal: false.
	// Equal: false.
	// Equal: false.
	// Equal: false.
	// Equal: false.
	// Equal: false.
	// Equal: false.
	// Equal: false.
	// Equal: false.
	// Equal: false.
	// Equal: false.
	// Equal: false.
	// Equal: false.
	// Equal: false.
}

func ExampleStructField_Equal_pointerFields() {
	type Program struct {
		Name *string `json:"name,omitempty"`
	}

	type Server struct {
		Name       *string        `json:"name,omitempty"`
		ID         *uint          `json:"id,omitempty"`
		Enabled    *bool          `json:"enabled,omitempty"`
		Err        *error         `json:"err,omitempty"`
		ChangedAt  *time.Time     `json:"changed_at,omitempty"`
		TimeOut    *time.Duration `json:"time_out,omitempty"`
		Count      *int           `json:"count,omitempty"`
		Price      *float32       `json:"price,omitempty"`
		Stream     *[]byte        `json:"stream,omitempty"`
		Complex    *complex128    `json:"complex,omitempty"`
		Anything   interface{}    `json:"anything,omitempty"`
		Program    *Program       `json:"program,omitempty"`
		Password   *string        `json:"-"`
		unexported *bool
	}

	program1 := Program{
		Name: structs.PtrString("Apache"),
	}

	server1 := Server{
		Name:       structs.PtrString("Roninzo"),
		ID:         structs.PtrUint(123456),
		Enabled:    structs.PtrBool(true),
		Err:        structs.PtrError(errors.New("rows not found")),
		ChangedAt:  structs.PtrTime(time.Date(2021, time.August, 3, 16, 44, 46, 0, time.UTC)),
		TimeOut:    structs.PtrDuration(5 * time.Second),
		Count:      structs.PtrInt(8),
		Price:      structs.PtrFloat32(1922.50),
		Stream:     structs.PtrBytes([]byte("Hello world")),
		Complex:    structs.PtrComplex128(complex(22, 50)),
		Anything:   "test",
		Program:    &program1,
		Password:   structs.PtrString("abcdefg"),
		unexported: structs.PtrBool(true),
	}

	program2 := Program{
		Name: structs.PtrString("Microsoft IIS"),
	}

	server2 := Server{
		Name:       structs.PtrString("ozninoR"),
		ID:         structs.PtrUint(654321),
		Enabled:    structs.PtrBool(false),
		Err:        structs.PtrError(errors.New("not compliant")),
		ChangedAt:  structs.PtrTime(time.Date(2021, time.August, 31, 14, 11, 11, 0, time.UTC)),
		TimeOut:    structs.PtrDuration(30 * time.Second),
		Count:      structs.PtrInt(3),
		Price:      structs.PtrFloat32(7022.50),
		Stream:     structs.PtrBytes([]byte("Bye bye world")),
		Complex:    structs.PtrComplex128(complex(99, 12)),
		Anything:   654321,
		Program:    &program2,
		Password:   structs.PtrString("gfedcba"),
		unexported: structs.PtrBool(false),
	}

	// backup <- server1 via cloning
	clone, err := structs.Clone(&server1)
	if err != nil {
		fmt.Printf("Clone[Error]: %v.\n", err)
		return
	}
	backup, ok := clone.(*Server)
	if !ok {
		fmt.Printf("TypeAssertion[Error]: %v.\n", err)
		return
	}

	s1, err := structs.New(&server1)
	if err != nil {
		fmt.Printf("New1[Error]: %v.\n", err)
		return
	}

	s2, err := structs.New(backup)
	if err != nil {
		fmt.Printf("New2[Error]: %v.\n", err)
		return
	}

	fmt.Println("Compare server1 and backup (equal)")
	for i, f1 := range s1.Fields() {
		f2 := s2.Field(i)
		fmt.Printf("[%d]Equal: %v.\n", i, f1.Equal(f2))
	}

	s2, err = structs.New(&server2)
	if err != nil {
		fmt.Printf("New3[Error]: %v.\n", err)
		return
	}

	fmt.Println("Compare server1 and server2 (different)")
	for i, f1 := range s1.Fields() {
		f2 := s2.Field(i)
		fmt.Printf("[%d]Equal: %v.\n", i, f1.Equal(f2))
	}

	// Output:
	// Compare server1 and backup (equal)
	// [0]Equal: true.
	// [1]Equal: true.
	// [2]Equal: true.
	// [3]Equal: true.
	// [4]Equal: true.
	// [5]Equal: true.
	// [6]Equal: true.
	// [7]Equal: true.
	// [8]Equal: true.
	// [9]Equal: true.
	// [10]Equal: true.
	// [11]Equal: true.
	// [12]Equal: true.
	// [13]Equal: false.
	// Compare server1 and server2 (different)
	// [0]Equal: false.
	// [1]Equal: false.
	// [2]Equal: false.
	// [3]Equal: false.
	// [4]Equal: false.
	// [5]Equal: false.
	// [6]Equal: false.
	// [7]Equal: false.
	// [8]Equal: false.
	// [9]Equal: false.
	// [10]Equal: false.
	// [11]Equal: false.
	// [12]Equal: false.
	// [13]Equal: false.
}

/*   S t r u c t F i e l d s   */

func ExampleStructFields_Names() {
	type Module struct {
		Name string `json:"name"`
	}

	type Server struct {
		Name       *string `json:"name,omitempty"`
		ID         uint    `json:"id,omitempty"`
		Enabled    bool    `json:"enabled,omitempty"`
		Count      int32   `json:"count,omitempty"`
		Password   string  `json:"-"`
		unexported bool
		Module     *Module `json:"module"`
	}

	var name string = "Roninzo"
	server := Server{
		Name:       &name,
		ID:         123456,
		Enabled:    true,
		Count:      0,
		Password:   "abcdefg",
		unexported: true,
		Module:     &Module{Name: "Power Supply"},
	}

	s, err := structs.New(&server)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("Fields.Names: %v\n", s.Fields().Names())

	// Output:
	// Fields.Names: [Name ID Enabled Count Password unexported Module]
}

/*   S t r u c t R o w s   */

func ExampleStructRows_Index() {
	type Server struct {
		Count int32 `json:"count,omitempty"`
	}
	servers := []Server{
		{Count: 5},
		{Count: 6},
	}
	s, err := structs.New(&servers)
	if err != nil {
		fmt.Printf("New[Error]: %v.\n", err)
	}
	rows, err := s.Rows()
	if err != nil {
		fmt.Printf("Rows[Error]: %v.\n", err)
	}
	for rows.Next() {
		fmt.Printf("struct row index: %d.\n", rows.Index())
	}
	err = rows.Close()
	if err != nil {
		fmt.Printf("Close[Error]: %v.\n", err)
	}
	fmt.Printf("struct row index: %d.\n", rows.Index())
	if err := rows.Err(); err != nil {
		fmt.Printf("Err[Error]: %v.\n", err)
	}

	// Output:
	// struct row index: 0.
	// struct row index: 1.
	// struct row index: -1.
}

func ExampleStructRows_Len() {
	type Server struct {
		Count int32 `json:"count,omitempty"`
	}
	servers := []Server{
		{Count: 5},
		{Count: 6},
	}
	s, err := structs.New(&servers)
	if err != nil {
		fmt.Printf("New[Error]: %v.\n", err)
	}
	rows, err := s.Rows()
	if err != nil {
		fmt.Printf("Rows[Error]: %v.\n", err)
	}
	fmt.Printf("struct number of rows: %d.\n", rows.Len())
	err = rows.Close()
	if err != nil {
		fmt.Printf("Close[Error]: %v.\n", err)
	}
	fmt.Printf("struct number of rows: %d.\n", rows.Len())
	if err := rows.Err(); err != nil {
		fmt.Printf("Err[Error]: %v.\n", err)
	}

	// Output:
	// struct number of rows: 2.
	// struct number of rows: -1.
}

func ExampleStructRows_MaxRow() {
	type Server struct {
		Count int32 `json:"count,omitempty"`
	}
	servers := []Server{
		{Count: 5},
		{Count: 6},
	}
	s, err := structs.New(&servers)
	if err != nil {
		fmt.Printf("New[Error]: %v.\n", err)
	}
	rows, err := s.Rows()
	if err != nil {
		fmt.Printf("Rows[Error]: %v.\n", err)
	}
	fmt.Printf("struct last row index: %d.\n", rows.MaxRow())
	err = rows.Close()
	if err != nil {
		fmt.Printf("Close[Error]: %v.\n", err)
	}
	fmt.Printf("struct last row index: %d.\n", rows.MaxRow())
	if err := rows.Err(); err != nil {
		fmt.Printf("Err[Error]: %v.\n", err)
	}

	// Output:
	// struct last row index: 1.
	// struct last row index: -1.
}

func ExampleStructRows_Columns() {
	type Server struct {
		Count int32 `json:"count,omitempty"`
	}
	servers := []Server{
		{Count: 5},
		{Count: 6},
	}
	s, err := structs.New(&servers)
	if err != nil {
		fmt.Printf("New[Error]: %v.\n", err)
	}
	rows, err := s.Rows()
	if err != nil {
		fmt.Printf("Rows[Error]: %v.\n", err)
	}
	cols, err := rows.Columns()
	if err != nil {
		fmt.Printf("Columns[Error]: %v.\n", err)
	} else {
		fmt.Printf("struct row column names: %s.\n", cols)
	}
	err = rows.Close()
	if err != nil {
		fmt.Printf("Close[Error]: %v.\n", err)
	}
	cols, err = rows.Columns()
	if err != nil {
		fmt.Printf("Columns[Error]: %v.\n", err)
	} else {
		fmt.Printf("struct row column names: %s.\n", cols)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("Err[Error]: %v.\n", err)
	}

	// Output:
	// struct row column names: [Count].
	// Columns[Error]: struct rows are closed.
}

func ExampleStructRows_Next() {
	type Server struct {
		Count int32 `json:"count,omitempty"`
	}
	servers := []Server{
		{Count: 5},
		{Count: 6},
	}
	s, err := structs.New(&servers)
	if err != nil {
		fmt.Printf("New[Error]: %v.\n", err)
	}
	rows, err := s.Rows()
	if err != nil {
		fmt.Printf("Rows[Error]: %v.\n", err)
	}
	for rows.Next() {
		fmt.Printf("struct row index: %d.\n", rows.Index())
	}
	err = rows.Close()
	if err != nil {
		fmt.Printf("Close[Error]: %v.\n", err)
	}
	fmt.Printf("struct row index: %d.\n", rows.Index())
	if err := rows.Err(); err != nil {
		fmt.Printf("Err[Error]: %v.\n", err)
	}

	// Output:
	// struct row index: 0.
	// struct row index: 1.
	// struct row index: -1.
}

func ExampleStructRows_Close() {
	type Server struct {
		Count int32 `json:"count,omitempty"`
	}
	servers := []Server{
		{Count: 5},
		{Count: 6},
	}
	s, err := structs.New(&servers)
	if err != nil {
		fmt.Printf("New[Error]: %v.\n", err)
	}
	rows, err := s.Rows()
	if err != nil {
		fmt.Printf("Rows[Error]: %v.\n", err)
	}
	for rows.Next() {
		fmt.Printf("struct row index: %d.\n", rows.Index())
	}
	err = rows.Close()
	if err != nil {
		fmt.Printf("Close[Error]: %v.\n", err)
	}
	fmt.Printf("struct row index: %d.\n", rows.Index())
	if err := rows.Err(); err != nil {
		fmt.Printf("Err[Error]: %v.\n", err)
	}

	// Output:
	// struct row index: 0.
	// struct row index: 1.
	// struct row index: -1.
}

func ExampleStructRows_Err() {
	type Server struct {
		Count int32 `json:"count,omitempty"`
	}
	servers := []Server{
		{Count: 5},
		{Count: 6},
	}
	s, err := structs.New(&servers)
	if err != nil {
		fmt.Printf("New[Error]: %v.\n", err)
	}
	rows, err := s.Rows()
	if err != nil {
		fmt.Printf("Err[Error]: %v.\n", err)
	}
	defer rows.Close()
	fmt.Printf("struct row index: %d.\n", rows.Index())
	_ = rows.Field("InvalidName")
	if err := rows.Err(); err != nil {
		fmt.Printf("Err[Error]: %v.\n", err)
	}

	// Output:
	// struct row index: -1.
	// Err[Error]: invalid field name InvalidName.
}
