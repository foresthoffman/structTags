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

type simpleTestStruct struct {
	Field4   string `json:"field4" custom:"field_4"`
	Ignored3 string `json:"ignored3" custom:"-"`
}

type slicedTestStruct struct {
	Field2   string              `json:"field2" custom:"field_2"`
	Field3   string              `json:"field3" custom:"field_3"`
	Ignored2 string              `json:"ignored2" custom:"-"`
	Items    []simpleTestStruct  `json:"simpleItems" custom:"simple_items"`
	ItemsPtr *[]simpleTestStruct `json:"simpleItemsPtr" custom:"simple_items_ptr"`
}

type nestedTestStruct struct {
	Field1    string            `json:"field1" custom:"field_1"`
	Ignored1  string            `json:"ignored1" custom:"-"`
	Sliced    slicedTestStruct  `json:"sliced" custom:"sliced"`
	SlicedPtr *slicedTestStruct `json:"slicedPtr" custom:"sliced_ptr"`
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
			Name: "slice of nested structs",
			Input: []nestedTestStruct{
				{
					Field1:   "testing",
					Ignored1: "SUPERSECRET",
					Sliced: slicedTestStruct{
						Field2:   "critical",
						Field3:   "another",
						Ignored2: "shouldn't see me #1",
						Items: []simpleTestStruct{
							{
								Field4:   "sunny",
								Ignored3: "TOP_SECRET_TOKEN #1",
							},
						},
					},
					SlicedPtr: &slicedTestStruct{
						Field2:   "special",
						Field3:   "different",
						Ignored2: "shouldn't see me #2",
						Items: []simpleTestStruct{
							{
								Field4:   "cloudy",
								Ignored3: "TOP_SECRET_TOKEN #2",
							},
						},
					},
				},
				{
					Field1:   "debugging",
					Ignored1: "PASSWORD123",
					Sliced: slicedTestStruct{
						Field2:   "taco",
						Field3:   "tortilla",
						Ignored2: "shouldn't see me #3",
						Items: []simpleTestStruct{
							{
								Field4:   "stormy",
								Ignored3: "TOP_SECRET_TOKEN #3",
							},
						},
					},
					SlicedPtr: &slicedTestStruct{
						Field2:   "guacamole",
						Field3:   "burrito",
						Ignored2: "shouldn't see me #4",
						Items: []simpleTestStruct{
							{
								Field4:   "tornado",
								Ignored3: "TOP_SECRET_TOKEN #4",
							},
						},
					},
				},
			},
			ExpectedError: nil,
			ExpectedOutput: `{"field_1":"testing","sliced":{"field_2":"critical","field_3":"another","simple_items":[{"field_4":"sunny"}],"simple_items_ptr":{}},"sliced_ptr":{"field_2":"special","field_3":"different","simple_items":[{"field_4":"cloudy"}],"simple_items_ptr":{}}}
{"field_1":"debugging","sliced":{"field_2":"taco","field_3":"tortilla","simple_items":[{"field_4":"stormy"}],"simple_items_ptr":{}},"sliced_ptr":{"field_2":"guacamole","field_3":"burrito","simple_items":[{"field_4":"tornado"}],"simple_items_ptr":{}}}
`,
		},
		{
			Name: "nested struct",
			Input: nestedTestStruct{
				Field1:   "testing",
				Ignored1: "SUPERSECRET",
				Sliced: slicedTestStruct{
					Field2:   "critical",
					Field3:   "another",
					Ignored2: "shouldn't see me #1",
					Items: []simpleTestStruct{
						{
							Field4:   "sunny",
							Ignored3: "TOP_SECRET_TOKEN #1",
						},
					},
				},
				SlicedPtr: &slicedTestStruct{
					Field2:   "special",
					Field3:   "different",
					Ignored2: "shouldn't see me #2",
					Items: []simpleTestStruct{
						{
							Field4:   "cloudy",
							Ignored3: "TOP_SECRET_TOKEN #2",
						},
					},
				},
			},
			ExpectedError: nil,
			ExpectedOutput: `{"field_1":"testing","sliced":{"field_2":"critical","field_3":"another","simple_items":[{"field_4":"sunny"}],"simple_items_ptr":{}},"sliced_ptr":{"field_2":"special","field_3":"different","simple_items":[{"field_4":"cloudy"}],"simple_items_ptr":{}}}
`,
		},
		{
			Name: "simple struct",
			Input: simpleTestStruct{
				Field4:   "cloudy",
				Ignored3: "TOP_SECRET_TOKEN #2",
			},
			ExpectedError: nil,
			ExpectedOutput: `{"field_4":"cloudy"}
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
