package jsonpatch

import (
	"encoding/json"
	"reflect"
	"testing"
)

type applyReplaceTestCase[T any] struct {
	name              string
	patches           []*PatchRequest[T]
	input             *T
	newEmptyInputFunc func() *T
	expectError       bool
	expect            *T
}

func TestPatchRequest_applyReplace(t *testing.T) {
	testCases := []applyReplaceTestCase[MyStruct]{
		{
			name: "replace_string",
			patches: []*PatchRequest[MyStruct]{
				{
					Operation: "replace",
					Path:      "$.field_string",
					Value:     "bar",
				},
			},
			newEmptyInputFunc: NewMyStruct,
			input: &MyStruct{
				FieldString: "foo",
			},
			expectError: false,
			expect: &MyStruct{
				FieldString: "bar",
			},
		},

		{
			name: "replace_int",
			patches: []*PatchRequest[MyStruct]{
				{
					Operation: "replace",
					Path:      "$.field_struct.field_int_ptr",
					Value:     42,
				},
			},
			newEmptyInputFunc: NewMyStruct,
			input: &MyStruct{
				FieldStruct: &MySubStruct{
					FieldIntPtr: getPtr(15),
				},
			},
			expectError: false,
			expect: &MyStruct{
				FieldStruct: &MySubStruct{
					FieldIntPtr: getPtr(42),
				},
			},
		},

		{
			name: "replace_struct",
			patches: []*PatchRequest[MyStruct]{
				{
					Operation: "replace",
					Path:      "$.field_struct",
					Value: &MySubStruct{
						FieldBool: true,
					},
				},
			},
			newEmptyInputFunc: NewMyStruct,
			input: &MyStruct{
				FieldStruct: &MySubStruct{
					FieldBool: false,
				},
			},
			expectError: false,
			expect: &MyStruct{
				FieldStruct: &MySubStruct{
					FieldBool: true,
				},
			},
		},

		{
			name: "replace_map_string_string",
			patches: []*PatchRequest[MyStruct]{
				{
					Operation: "replace",
					Path:      "$.field_map_string_string.field1",
					Value:     "newvalue",
				},
			},
			newEmptyInputFunc: NewMyStruct,
			input: &MyStruct{
				FieldMapStringString: map[string]string{
					"field1": "value1",
					"field2": "value2",
				},
			},
			expectError: false,
			expect: &MyStruct{
				FieldMapStringString: map[string]string{
					"field1": "newvalue",
					"field2": "value2",
				},
			},
		},

		{
			name: "replace_map_string_string_brackets",
			patches: []*PatchRequest[MyStruct]{
				{
					Operation: "replace",
					Path:      "$.field_map_string_string['field1']",
					Value:     "newvalue",
				},
			},
			newEmptyInputFunc: NewMyStruct,
			input: &MyStruct{
				FieldMapStringString: map[string]string{
					"field1": "value1",
					"field2": "value2",
				},
			},
			expectError: false,
			expect: &MyStruct{
				FieldMapStringString: map[string]string{
					"field1": "newvalue",
					"field2": "value2",
				},
			},
		},

		{
			name: "replace_map_string_string_brackets_complex",
			patches: []*PatchRequest[MyStruct]{
				{
					Operation: "replace",
					Path:      "$.field_map_string_string['field1/foo.bar']",
					Value:     "newvalue",
				},
			},
			newEmptyInputFunc: NewMyStruct,
			input: &MyStruct{
				FieldMapStringString: map[string]string{
					"field1/foo.bar": "value1",
					"field2":         "value2",
				},
			},
			expectError: false,
			expect: &MyStruct{
				FieldMapStringString: map[string]string{
					"field1/foo.bar": "newvalue",
					"field2":         "value2",
				},
			},
		},

		{
			name: "replace_map_string_string_brackets_complex_double_quotes",
			patches: []*PatchRequest[MyStruct]{
				{
					Operation: "replace",
					Path:      `$.field_map_string_string["field1/foo.bar"]`,
					Value:     "newvalue",
				},
			},
			newEmptyInputFunc: NewMyStruct,
			input: &MyStruct{
				FieldMapStringString: map[string]string{
					"field1/foo.bar": "value1",
					"field2":         "value2",
				},
			},
			expectError: false,
			expect: &MyStruct{
				FieldMapStringString: map[string]string{
					"field1/foo.bar": "newvalue",
					"field2":         "value2",
				},
			},
		},

		{
			name: "replace_array_string_by_index",
			patches: []*PatchRequest[MyStruct]{
				{
					Operation: "replace",
					Path:      `$.field_array_string[0]`,
					Value:     "newvalue",
				},
			},
			newEmptyInputFunc: NewMyStruct,
			input: &MyStruct{
				FieldArrayString: []string{"m1", "m2"},
			},
			expectError: false,
			expect: &MyStruct{
				FieldArrayString: []string{"newvalue", "m2"},
			},
		},

		{
			name: "replace_struct_array_string_by_filter_by_string",
			patches: []*PatchRequest[MyStruct]{
				{
					Operation: "replace",
					Path:      `$.field_array_struct[?(@.field_string_sub1 == 'string1')].field_string_sub2`,
					Value:     "newvalue",
				},
			},
			newEmptyInputFunc: NewMyStruct,
			input: &MyStruct{
				FieldArrayStruct: []*MySubStruct{
					{
						FieldString1: "string1",
						FieldString2: "str1",
					},
					{
						FieldString1: "string2",
						FieldString2: "str2",
					},
					{
						FieldString1: "string1",
						FieldString2: "str3",
					},
				},
			},
			expectError: false,
			expect: &MyStruct{
				FieldArrayStruct: []*MySubStruct{
					{
						FieldString1: "string1",
						FieldString2: "newvalue",
					},
					{
						FieldString1: "string2",
						FieldString2: "str2",
					},
					{
						FieldString1: "string1",
						FieldString2: "newvalue",
					},
				},
			},
		},

		{
			name: "replace_struct_array_string_by_filter_by_string_and",
			patches: []*PatchRequest[MyStruct]{
				{
					Operation: "replace",
					Path:      `$.field_array_struct[?(@.field_string_sub1 == 'string1' && @.field_int_ptr == 42)].field_string_sub2`,
					Value:     "newvalue",
				},
			},
			newEmptyInputFunc: NewMyStruct,
			input: &MyStruct{
				FieldArrayStruct: []*MySubStruct{
					{
						FieldString1: "string1",
						FieldIntPtr:  getPtr(12),
						FieldString2: "str1",
					},
					{
						FieldString1: "string2",
						FieldIntPtr:  getPtr(42),
						FieldString2: "str2",
					},
					{
						FieldString1: "string1",
						FieldIntPtr:  getPtr(42),
						FieldString2: "str3",
					},
				},
			},
			expectError: false,
			expect: &MyStruct{
				FieldArrayStruct: []*MySubStruct{
					{
						FieldString1: "string1",
						FieldIntPtr:  getPtr(12),
						FieldString2: "str1",
					},
					{
						FieldString1: "string2",
						FieldIntPtr:  getPtr(42),
						FieldString2: "str2",
					},
					{
						FieldString1: "string1",
						FieldIntPtr:  getPtr(42),
						FieldString2: "newvalue",
					},
				},
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			var patched *MyStruct
			patched = tt.input
			for _, pr := range tt.patches {
				patched, err = pr.replace(patched, tt.newEmptyInputFunc())
			}

			if (err != nil) != tt.expectError {
				t.Errorf("applyReplace() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if !reflect.DeepEqual(patched, tt.expect) {
				bPatched, _ := json.Marshal(patched)
				bExpected, _ := json.Marshal(tt.expect)
				t.Errorf("applyReplace() got    = %s", string(bPatched))
				t.Errorf("applyReplace() expect = %s", string(bExpected))
			}
		})
	}
}
