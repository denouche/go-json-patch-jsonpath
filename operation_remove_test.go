package jsonpatch

import (
	"encoding/json"
	"reflect"
	"testing"
)

type applyRemoveTestCase[T any] struct {
	name              string
	patches           []*PatchRequest[MyStruct]
	input             T
	newEmptyInputFunc func() T
	expectError       bool
	expect            T
}

func TestPatchRequest_applyRemove(t *testing.T) {
	testCases := []applyRemoveTestCase[*MyStruct]{
		{
			name: "remove_string",
			patches: []*PatchRequest[MyStruct]{
				{
					Operation: "remove",
					Path:      "$.field_string",
				},
			},
			newEmptyInputFunc: func() *MyStruct {
				return &MyStruct{}
			},
			input: &MyStruct{
				FieldString: "foo",
			},
			expectError: false,
			expect: &MyStruct{
				FieldString: "",
			},
		},

		{
			name: "remove_string_ptr",
			patches: []*PatchRequest[MyStruct]{
				{
					Operation: "remove",
					Path:      "$.field_string_ptr",
				},
			},
			newEmptyInputFunc: func() *MyStruct {
				return &MyStruct{}
			},
			input: &MyStruct{
				FieldStringPtr: getPtr("foo"),
			},
			expectError: false,
			expect: &MyStruct{
				FieldStringPtr: nil,
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			var patched *MyStruct
			patched = tt.input
			for _, pr := range tt.patches {
				patched, err = pr.remove(patched, tt.newEmptyInputFunc())
			}

			if (err != nil) != tt.expectError {
				t.Errorf("applyRemove() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if !reflect.DeepEqual(patched, tt.expect) {
				bPatched, _ := json.Marshal(patched)
				bExpected, _ := json.Marshal(tt.expect)
				t.Errorf("applyRemove() got    = %s", string(bPatched))
				t.Errorf("applyRemove() expect = %s", string(bExpected))
			}
		})
	}
}
