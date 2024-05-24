package structTags

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	targetCustomTag    = "custom"
	ignoreTagWithValue = "-"
)

var (
	inf interface{} = struct {
		StringVar string `json:"stringVar" custom:"string_var"`
	}{
		StringVar: "example",
	}
)

type scalarStruct struct {
	StringVar  string  `json:"stringVar" custom:"string_var"`
	IntVar     int     `json:"intVar" custom:"int_var"`
	Int8Var    int8    `json:"int8Var" custom:"int8_var"`
	Int16Var   int16   `json:"int16Var" custom:"int16_var"`
	Int32Var   int32   `json:"int32Var" custom:"int32_var"`
	Int64Var   int64   `json:"int64Var" custom:"int64_var"`
	UintVar    uint    `json:"uintVar" custom:"uint_var"`
	Uint8Var   uint8   `json:"uint8Var" custom:"uint8_var"`
	Uint16Var  uint16  `json:"uint16Var" custom:"uint16_var"`
	Uint32Var  uint32  `json:"uint32Var" custom:"uint32_var"`
	Uint64Var  uint64  `json:"uint64Var" custom:"uint64_var"`
	Float32Var float32 `json:"float32Var" custom:"float32_var"`
	Float64Var float64 `json:"float64Var" custom:"float64_var"`
	BoolVar    bool    `json:"boolVar" custom:"bool_var"`
	IgnoredVar string  `json:"ignoredVar" custom:"-"`
}

type sliceStruct struct {
	StringSliceVar  []string  `json:"stringSliceVar" custom:"string_slice_var"`
	IntSliceVar     []int     `json:"intSliceVar" custom:"int_slice_var"`
	Int8SliceVar    []int8    `json:"int8SliceVar" custom:"int8_slice_var"`
	Int16SliceVar   []int16   `json:"int16SliceVar" custom:"int16_slice_var"`
	Int32SliceVar   []int32   `json:"int32SliceVar" custom:"int32_slice_var"`
	Int64SliceVar   []int64   `json:"int64SliceVar" custom:"int64_slice_var"`
	UintSliceVar    []uint    `json:"uintSliceVar" custom:"uint_slice_var"`
	Uint8SliceVar   []uint8   `json:"uint8SliceVar" custom:"uint8_slice_var"`
	Uint16SliceVar  []uint16  `json:"uint16SliceVar" custom:"uint16_slice_var"`
	Uint32SliceVar  []uint32  `json:"uint32SliceVar" custom:"uint32_slice_var"`
	Uint64SliceVar  []uint64  `json:"uint64SliceVar" custom:"uint64_slice_var"`
	Float32SliceVar []float32 `json:"float32SliceVar" custom:"float32_slice_var"`
	Float64SliceVar []float64 `json:"float64SliceVar" custom:"float64_slice_var"`
	BoolSliceVar    []bool    `json:"boolSliceVar" custom:"bool_slice_var"`
	IgnoredSliceVar []string  `json:"ignoredSliceVar" custom:"-"`
}

type mapStruct struct {
	StringMapVar  map[string]string  `json:"stringMapVar" custom:"string_map_var"`
	IntMapVar     map[string]int     `json:"intMapVar" custom:"int_map_var"`
	Int8MapVar    map[string]int8    `json:"int8MapVar" custom:"int8_map_var"`
	Int16MapVar   map[string]int16   `json:"int16MapVar" custom:"int16_map_var"`
	Int32MapVar   map[string]int32   `json:"int32MapVar" custom:"int32_map_var"`
	Int64MapVar   map[string]int64   `json:"int64MapVar" custom:"int64_map_var"`
	UintMapVar    map[string]uint    `json:"uintMapVar" custom:"uint_map_var"`
	Uint8MapVar   map[string]uint8   `json:"uint8MapVar" custom:"uint8_map_var"`
	Uint16MapVar  map[string]uint16  `json:"uint16MapVar" custom:"uint16_map_var"`
	Uint32MapVar  map[string]uint32  `json:"uint32MapVar" custom:"uint32_map_var"`
	Uint64MapVar  map[string]uint64  `json:"uint64MapVar" custom:"uint64_map_var"`
	Float32MapVar map[string]float32 `json:"float32MapVar" custom:"float32_map_var"`
	Float64MapVar map[string]float64 `json:"float64MapVar" custom:"float64_map_var"`
	BoolMapVar    map[string]bool    `json:"boolMapVar" custom:"bool_map_var"`
	IgnoredMapVar map[string]string  `json:"ignoredMapVar" custom:"-"`
}

type parentStruct struct {
	ChildStructVar childStruct `json:"childStructVar" custom:"child_struct_var"`
}

type childStruct struct {
	GrandChildStructVar grandChildStruct `json:"grandChildStructVar" custom:"grand_child_struct_var"`
}

type grandChildStruct struct {
	StringVar string `json:"stringVar" custom:"string_var"`
}

func TestNewCustomMarshaller(t *testing.T) {
	testCases := []struct {
		Name           string
		Input          any
		ExpectedError  error
		ExpectedOutput string
	}{
		{
			Name:           "nil",
			Input:          nil,
			ExpectedError:  errors.New(ErrNilObject),
			ExpectedOutput: "",
		},
		{
			Name:           "nil interface",
			Input:          new(interface{}),
			ExpectedError:  errors.New("failed to marshal ptr: failed to marshal interface: object was nil"),
			ExpectedOutput: "",
		},
		{
			Name:          "interface",
			Input:         &inf,
			ExpectedError: nil,
			ExpectedOutput: `{"string_var":"example"}
`,
		},
		{
			Name: "scalar struct",
			Input: scalarStruct{
				StringVar:  "str",
				IntVar:     1,
				Int8Var:    2,
				Int16Var:   3,
				Int32Var:   4,
				Int64Var:   5,
				UintVar:    6,
				Uint8Var:   7,
				Uint16Var:  8,
				Uint32Var:  9,
				Uint64Var:  10,
				Float32Var: 11.1,
				Float64Var: 12.2,
				BoolVar:    true,
				IgnoredVar: "ignored-example",
			},
			ExpectedError: nil,
			ExpectedOutput: `{"string_var":"str","int_var":1,"int8_var":2,"int16_var":3,"int32_var":4,"int64_var":5,"uint_var":6,"uint8_var":7,"uint16_var":8,"uint32_var":9,"uint64_var":10,"float32_var":11.1,"float64_var":12.2,"bool_var":true}
`,
		},
		{
			Name: "scalar struct ptr",
			Input: &scalarStruct{
				StringVar:  "str",
				IntVar:     1,
				Int8Var:    2,
				Int16Var:   3,
				Int32Var:   4,
				Int64Var:   5,
				UintVar:    6,
				Uint8Var:   7,
				Uint16Var:  8,
				Uint32Var:  9,
				Uint64Var:  10,
				Float32Var: 11.1,
				Float64Var: 12.2,
				BoolVar:    false,
				IgnoredVar: "ignored-example",
			},
			ExpectedError: nil,
			ExpectedOutput: `{"string_var":"str","int_var":1,"int8_var":2,"int16_var":3,"int32_var":4,"int64_var":5,"uint_var":6,"uint8_var":7,"uint16_var":8,"uint32_var":9,"uint64_var":10,"float32_var":11.1,"float64_var":12.2,"bool_var":false}
`,
		},
		{
			Name:          "string slice",
			Input:         []string{"str", "ing"},
			ExpectedError: nil,
			ExpectedOutput: `["str","ing"]
`,
		},
		{
			Name:          "int slice",
			Input:         []int{1, 2},
			ExpectedError: nil,
			ExpectedOutput: `[1,2]
`,
		},
		{
			Name:          "int8 slice",
			Input:         []int8{3, 4},
			ExpectedError: nil,
			ExpectedOutput: `[3,4]
`,
		},
		{
			Name:          "int16 slice",
			Input:         []int16{5, 6},
			ExpectedError: nil,
			ExpectedOutput: `[5,6]
`,
		},
		{
			Name:          "int32 slice",
			Input:         []int32{7, 8},
			ExpectedError: nil,
			ExpectedOutput: `[7,8]
`,
		},
		{
			Name:          "int64 slice",
			Input:         []int64{9, 10},
			ExpectedError: nil,
			ExpectedOutput: `[9,10]
`,
		},
		{
			Name:          "uint slice",
			Input:         []uint{11, 12},
			ExpectedError: nil,
			ExpectedOutput: `[11,12]
`,
		},
		{
			Name:          "uint8 slice",
			Input:         []uint8{13, 14},
			ExpectedError: nil,
			ExpectedOutput: `[13,14]
`,
		},
		{
			Name:          "uint16 slice",
			Input:         []uint16{15, 16},
			ExpectedError: nil,
			ExpectedOutput: `[15,16]
`,
		},
		{
			Name:          "uint32 slice",
			Input:         []uint32{17, 18},
			ExpectedError: nil,
			ExpectedOutput: `[17,18]
`,
		},
		{
			Name:          "uint64 slice",
			Input:         []uint64{19, 20},
			ExpectedError: nil,
			ExpectedOutput: `[19,20]
`,
		},
		{
			Name:          "float32 slice",
			Input:         []float32{21.1, 22.2},
			ExpectedError: nil,
			ExpectedOutput: `[21.1,22.2]
`,
		},
		{
			Name:          "float64 slice",
			Input:         []float64{23.3, 24.4},
			ExpectedError: nil,
			ExpectedOutput: `[23.3,24.4]
`,
		},
		{
			Name: "slice struct",
			Input: sliceStruct{
				StringSliceVar:  []string{"str", "ing"},
				IntSliceVar:     []int{1, 2},
				Int8SliceVar:    []int8{3, 4},
				Int16SliceVar:   []int16{5, 6},
				Int32SliceVar:   []int32{7, 8},
				Int64SliceVar:   []int64{9, 10},
				UintSliceVar:    []uint{11, 12},
				Uint8SliceVar:   []uint8{13, 14},
				Uint16SliceVar:  []uint16{15, 16},
				Uint32SliceVar:  []uint32{17, 18},
				Uint64SliceVar:  []uint64{19, 20},
				Float32SliceVar: []float32{21.1, 22.2},
				Float64SliceVar: []float64{23.3, 24.4},
				BoolSliceVar:    []bool{true, false},
				IgnoredSliceVar: []string{"ignored", "example"},
			},
			ExpectedError: nil,
			ExpectedOutput: `{"string_slice_var":["str","ing"],"int_slice_var":[1,2],"int8_slice_var":[3,4],"int16_slice_var":[5,6],"int32_slice_var":[7,8],"int64_slice_var":[9,10],"uint_slice_var":[11,12],"uint8_slice_var":[13,14],"uint16_slice_var":[15,16],"uint32_slice_var":[17,18],"uint64_slice_var":[19,20],"float32_slice_var":[21.1,22.2],"float64_slice_var":[23.3,24.4],"bool_slice_var":[true,false]}
`,
		},
		{
			Name: "slice struct ptr",
			Input: &sliceStruct{
				StringSliceVar:  []string{"str", "ing"},
				IntSliceVar:     []int{1, 2},
				Int8SliceVar:    []int8{3, 4},
				Int16SliceVar:   []int16{5, 6},
				Int32SliceVar:   []int32{7, 8},
				Int64SliceVar:   []int64{9, 10},
				UintSliceVar:    []uint{11, 12},
				Uint8SliceVar:   []uint8{13, 14},
				Uint16SliceVar:  []uint16{15, 16},
				Uint32SliceVar:  []uint32{17, 18},
				Uint64SliceVar:  []uint64{19, 20},
				Float32SliceVar: []float32{21.1, 22.2},
				Float64SliceVar: []float64{23.3, 24.4},
				BoolSliceVar:    []bool{false, true},
				IgnoredSliceVar: []string{"ignored", "example"},
			},
			ExpectedError: nil,
			ExpectedOutput: `{"string_slice_var":["str","ing"],"int_slice_var":[1,2],"int8_slice_var":[3,4],"int16_slice_var":[5,6],"int32_slice_var":[7,8],"int64_slice_var":[9,10],"uint_slice_var":[11,12],"uint8_slice_var":[13,14],"uint16_slice_var":[15,16],"uint32_slice_var":[17,18],"uint64_slice_var":[19,20],"float32_slice_var":[21.1,22.2],"float64_slice_var":[23.3,24.4],"bool_slice_var":[false,true]}
`,
		},
		{
			Name:          "string map",
			Input:         map[string]string{"str": "ing", "ele": "ment"},
			ExpectedError: nil,
			ExpectedOutput: `{"ele":"ment","str":"ing"}
`,
		},
		{
			Name:          "int map",
			Input:         map[string]int{"a": 1, "b": 2},
			ExpectedError: nil,
			ExpectedOutput: `{"a":1,"b":2}
`,
		},
		{
			Name:          "int8 map",
			Input:         map[string]int8{"a": 3, "b": 4},
			ExpectedError: nil,
			ExpectedOutput: `{"a":3,"b":4}
`,
		},
		{
			Name:          "int16 map",
			Input:         map[string]int16{"a": 5, "b": 6},
			ExpectedError: nil,
			ExpectedOutput: `{"a":5,"b":6}
`,
		},
		{
			Name:          "int32 map",
			Input:         map[string]int32{"a": 7, "b": 8},
			ExpectedError: nil,
			ExpectedOutput: `{"a":7,"b":8}
`,
		},
		{
			Name:          "int64 map",
			Input:         map[string]int64{"a": 9, "b": 10},
			ExpectedError: nil,
			ExpectedOutput: `{"a":9,"b":10}
`,
		},
		{
			Name:          "uint map",
			Input:         map[string]uint{"a": 11, "b": 12},
			ExpectedError: nil,
			ExpectedOutput: `{"a":11,"b":12}
`,
		},
		{
			Name:          "uint8 map",
			Input:         map[string]uint8{"a": 13, "b": 14},
			ExpectedError: nil,
			ExpectedOutput: `{"a":13,"b":14}
`,
		},
		{
			Name:          "uint16 map",
			Input:         map[string]uint16{"a": 15, "b": 16},
			ExpectedError: nil,
			ExpectedOutput: `{"a":15,"b":16}
`,
		},
		{
			Name:          "uint32 map",
			Input:         map[string]uint32{"a": 17, "b": 18},
			ExpectedError: nil,
			ExpectedOutput: `{"a":17,"b":18}
`,
		},
		{
			Name:          "uint64 map",
			Input:         map[string]uint64{"a": 19, "b": 20},
			ExpectedError: nil,
			ExpectedOutput: `{"a":19,"b":20}
`,
		},
		{
			Name:          "float32 map",
			Input:         map[string]float32{"a": 21.1, "b": 22.2},
			ExpectedError: nil,
			ExpectedOutput: `{"a":21.1,"b":22.2}
`,
		},
		{
			Name:          "float64 map",
			Input:         map[string]float64{"a": 23.3, "b": 24.4},
			ExpectedError: nil,
			ExpectedOutput: `{"a":23.3,"b":24.4}
`,
		},
		{
			Name: "map struct",
			Input: mapStruct{
				StringMapVar:  map[string]string{"str": "ing", "ele": "ment"},
				IntMapVar:     map[string]int{"a": 1, "b": 2},
				Int8MapVar:    map[string]int8{"a": 3, "b": 4},
				Int16MapVar:   map[string]int16{"a": 5, "b": 6},
				Int32MapVar:   map[string]int32{"a": 7, "b": 8},
				Int64MapVar:   map[string]int64{"a": 9, "b": 10},
				UintMapVar:    map[string]uint{"a": 11, "b": 12},
				Uint8MapVar:   map[string]uint8{"a": 13, "b": 14},
				Uint16MapVar:  map[string]uint16{"a": 15, "b": 16},
				Uint32MapVar:  map[string]uint32{"a": 17, "b": 18},
				Uint64MapVar:  map[string]uint64{"a": 19, "b": 20},
				Float32MapVar: map[string]float32{"a": 21.1, "b": 22.2},
				Float64MapVar: map[string]float64{"a": 23.3, "b": 24.4},
				BoolMapVar:    map[string]bool{"a": true, "b": false},
				IgnoredMapVar: map[string]string{"a": "ignored", "b": "example"},
			},
			ExpectedError: nil,
			ExpectedOutput: `{"string_map_var":{"ele":"ment","str":"ing"},"int_map_var":{"a":1,"b":2},"int8_map_var":{"a":3,"b":4},"int16_map_var":{"a":5,"b":6},"int32_map_var":{"a":7,"b":8},"int64_map_var":{"a":9,"b":10},"uint_map_var":{"a":11,"b":12},"uint8_map_var":{"a":13,"b":14},"uint16_map_var":{"a":15,"b":16},"uint32_map_var":{"a":17,"b":18},"uint64_map_var":{"a":19,"b":20},"float32_map_var":{"a":21.1,"b":22.2},"float64_map_var":{"a":23.3,"b":24.4},"bool_map_var":{"a":true,"b":false}}
`,
		},
		{
			Name: "map struct ptr",
			Input: &mapStruct{
				StringMapVar:  map[string]string{"str": "ing", "ele": "ment"},
				IntMapVar:     map[string]int{"a": 1, "b": 2},
				Int8MapVar:    map[string]int8{"a": 3, "b": 4},
				Int16MapVar:   map[string]int16{"a": 5, "b": 6},
				Int32MapVar:   map[string]int32{"a": 7, "b": 8},
				Int64MapVar:   map[string]int64{"a": 9, "b": 10},
				UintMapVar:    map[string]uint{"a": 11, "b": 12},
				Uint8MapVar:   map[string]uint8{"a": 13, "b": 14},
				Uint16MapVar:  map[string]uint16{"a": 15, "b": 16},
				Uint32MapVar:  map[string]uint32{"a": 17, "b": 18},
				Uint64MapVar:  map[string]uint64{"a": 19, "b": 20},
				Float32MapVar: map[string]float32{"a": 21.1, "b": 22.2},
				Float64MapVar: map[string]float64{"a": 23.3, "b": 24.4},
				BoolMapVar:    map[string]bool{"a": true, "b": false},
				IgnoredMapVar: map[string]string{"a": "ignored", "b": "example"},
			},
			ExpectedError: nil,
			ExpectedOutput: `{"string_map_var":{"ele":"ment","str":"ing"},"int_map_var":{"a":1,"b":2},"int8_map_var":{"a":3,"b":4},"int16_map_var":{"a":5,"b":6},"int32_map_var":{"a":7,"b":8},"int64_map_var":{"a":9,"b":10},"uint_map_var":{"a":11,"b":12},"uint8_map_var":{"a":13,"b":14},"uint16_map_var":{"a":15,"b":16},"uint32_map_var":{"a":17,"b":18},"uint64_map_var":{"a":19,"b":20},"float32_map_var":{"a":21.1,"b":22.2},"float64_map_var":{"a":23.3,"b":24.4},"bool_map_var":{"a":true,"b":false}}
`,
		},
		{
			Name: "nested struct",
			Input: parentStruct{
				ChildStructVar: childStruct{
					GrandChildStructVar: grandChildStruct{
						StringVar: "str",
					},
				},
			},
			ExpectedError: nil,
			ExpectedOutput: `{"child_struct_var":{"grand_child_struct_var":{"string_var":"str"}}}
`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			b, err := NewCustomMarshaller(targetCustomTag, ignoreTagWithValue).Marshal(testCase.Input)
			assert.Equal(t, testCase.ExpectedError, err)
			assert.Equal(t, testCase.ExpectedOutput, string(b))
		})
	}
}
