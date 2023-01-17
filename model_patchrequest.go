package jsonpatch

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	PatchOperationAdd     = "add"
	PatchOperationRemove  = "remove"
	PatchOperationReplace = "replace"
)

type PatchOperation string

type PatchRequest[T any] struct {
	Operation PatchOperation `json:"op" validate:"required,oneof=remove replace"` // TODO implements add
	Path      string         `json:"path" validate:"required,jsonpath,ne=$"`
	Value     any            `json:"value"`
}

// Apply TODO
func (pr *PatchRequest[T]) Apply(initialResource *T, emptyResource *T) (*T, error) {
	switch pr.Operation {
	case PatchOperationReplace:
		return pr.replace(initialResource, emptyResource)
	case PatchOperationRemove:
		return pr.remove(initialResource, emptyResource)
	case PatchOperationAdd:
		//TODO make the implementation
		fallthrough // fallthrough for now
	default:
		return nil, errors.New("operation not implemented")
	}
}

func (pr *PatchRequest[T]) remarshal(resourceAsMap interface{}, emptyResource *T) (*T, error) {
	newBytes, err := json.Marshal(resourceAsMap)
	if err != nil {
		return nil, fmt.Errorf("match fail to marshal input resource %s", err.Error())
	}

	err = json.Unmarshal(newBytes, &emptyResource)
	if err != nil {
		return nil, fmt.Errorf("match fail to unmarshal %s", err.Error())
	}

	return emptyResource, nil
}
