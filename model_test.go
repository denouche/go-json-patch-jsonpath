package jsonpatch

type MyStruct struct {
	FieldString          string            `json:"field_string"`
	FieldStringPtr       *string           `json:"field_string_ptr"`
	FieldMapStringString map[string]string `json:"field_map_string_string"`
	FieldArrayString     []string          `json:"field_array_string"`
	FieldStruct          *MySubStruct      `json:"field_struct"`
	FieldArrayStruct     []*MySubStruct    `json:"field_array_struct"`
}

func NewMyStruct() *MyStruct {
	return &MyStruct{}
}

type MySubStruct struct {
	FieldIntPtr  *int   `json:"field_int_ptr"`
	FieldBool    bool   `json:"field_bool"`
	FieldString1 string `json:"field_string_sub1"`
	FieldString2 string `json:"field_string_sub2"`
}

func getPtr[T any](in T) *T {
	return &in
}
