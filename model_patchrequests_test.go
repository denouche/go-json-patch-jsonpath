package jsonpatch

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestPatchRequests_Apply(t *testing.T) {
	type testCase[T any] struct {
		name              string
		patchRequests     PatchRequests[T]
		input             *T
		newEmptyInputFunc func() *T
		expect            *T
		expectErr         bool
	}
	tests := []testCase[MyStruct]{
		{
			name: "apply_replace",
			patchRequests: PatchRequests[MyStruct]{
				Patches: []*PatchRequest[MyStruct]{
					{
						Operation: "replace",
						Path:      "$.field_string",
						Value:     "bar",
					},
				},
			},
			newEmptyInputFunc: NewMyStruct,
			input: &MyStruct{
				FieldString: "foo",
			},
			expect: &MyStruct{
				FieldString: "bar",
			},
			expectErr: false,
		},

		{
			name: "apply_remove_string",
			patchRequests: PatchRequests[MyStruct]{
				Patches: []*PatchRequest[MyStruct]{
					{
						Operation: "remove",
						Path:      "$.field_string",
					},
				},
			},
			newEmptyInputFunc: NewMyStruct,
			input: &MyStruct{
				FieldString:    "foo",
				FieldStringPtr: getPtr("fooz"),
			},
			expect: &MyStruct{
				FieldString:    "",
				FieldStringPtr: getPtr("fooz"),
			},
			expectErr: false,
		},

		{
			name: "apply_remove_string_ptr",
			patchRequests: PatchRequests[MyStruct]{
				Patches: []*PatchRequest[MyStruct]{
					{
						Operation: "remove",
						Path:      "$.field_string_ptr",
					},
				},
			},
			newEmptyInputFunc: NewMyStruct,
			input: &MyStruct{
				FieldString:    "foo",
				FieldStringPtr: getPtr("fooz"),
			},
			expect: &MyStruct{
				FieldString:    "foo",
				FieldStringPtr: nil,
			},
			expectErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.patchRequests.Apply(tt.input, tt.newEmptyInputFunc)
			if (err != nil) != tt.expectErr {
				t.Errorf("Apply() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expect) {
				bPatched, _ := json.Marshal(got)
				bExpected, _ := json.Marshal(tt.expect)
				t.Errorf("Apply() got    = %s", string(bPatched))
				t.Errorf("Apply() expect = %s", string(bExpected))
			}
		})
	}
}
