# go-json-patch-jsonpath

![Build Status](https://github.com/denouche/go-json-patch-jsonpath/actions/workflows/build.yaml/badge.svg)
[![GoCover](http://gocover.io/_badge/github.com/denouche/go-json-patch-jsonpath)](http://gocover.io/github.com/denouche/go-json-patch-jsonpath)
[![GoDoc](https://godoc.org/github.com/denouche/go-json-patch-jsonpath?status.svg)](https://godoc.org/github.com/denouche/go-json-patch-jsonpath)


A Golang library implementing JSON Patch (rfc6902) with a digression to use path in JSONPath format (https://www.ietf.org/archive/id/draft-ietf-jsonpath-base-08.txt) instead of JSON Pointer (rfc6901).

## Work in progress

This lib is a work in progress, for now only replace and remove operations has been implemented.

## Installation

```
go get github.com/denouche/go-json-patch-jsonpath
```

## Usage

```golang
package main

import (
	"fmt"
	jsonpatch "github.com/denouche/go-json-patch-jsonpath"
)

type People struct {
	Name    string    `json:"name"`
	Age     int       `json:"age"`
	Friends []*People `json:"friends"`
}

func NewPeople() *People {
	return &People{}
}

func main() {

	person := &People{
		Name: "Foo",
		Age:  42,
		Friends: []*People{
			{Name: "friend 1", Age: 19},
			{Name: "friend 2", Age: 25},
			{Name: "friend 3", Age: 42},
		},
	}

	r := jsonpatch.PatchRequests[People]{
		Patches: []*jsonpatch.PatchRequest[People]{
			{
				Operation: "replace",
				Path:      "$.name",
				Value:     "Bar",
			},
			{
				Operation: "remove",
				Path:      "$.friends[?(@.age < 20)]",
			},
		},
	}

	patchedPerson, err := r.Apply(person, NewPeople)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	fmt.Printf("name:%s age:%d friends:%d\n", person.Name, person.Age, len(person.Friends))
	// name:Foo age:42 friends:3

	fmt.Printf("name:%s age:%d friends:%d\n", patchedPerson.Name, patchedPerson.Age, len(patchedPerson.Friends))
	// name:Bar age:42 friends:2
}
```

## References

JSON Patch: https://www.rfc-editor.org/rfc/rfc6902

JSONPath:
- https://datatracker.ietf.org/doc/draft-ietf-jsonpath-base/
- https://www.ietf.org/archive/id/draft-ietf-jsonpath-base-08.txt
- https://jsonpath.com/

