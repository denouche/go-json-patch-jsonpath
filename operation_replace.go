package jsonpatch

import (
	"encoding/json"
	"fmt"

	"github.com/ohler55/ojg/jp"
)

func (pr *PatchRequest[T]) replace(resource *T, emptyResource *T) (*T, error) {
	compiledPath, err := jp.ParseString(pr.Path)
	if err != nil {
		return nil, fmt.Errorf("fail to compile path %s %s", pr.Path, err.Error())
	}

	b, err := json.Marshal(resource)
	if err != nil {
		return nil, fmt.Errorf("match fail to marshal input resource %s", err.Error())
	}

	var resourceAsMap interface{}
	err = json.Unmarshal(b, &resourceAsMap)
	if err != nil {
		return nil, fmt.Errorf("match fail to unmarshal %s", err.Error())
	}

	err = compiledPath.Set(resourceAsMap, pr.Value)
	if err != nil {
		return nil, err
	}

	return pr.remarshal(resourceAsMap, emptyResource)
}
