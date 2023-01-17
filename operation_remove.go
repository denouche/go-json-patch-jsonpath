package jsonpatch

import (
	"encoding/json"
	"fmt"

	"github.com/ohler55/ojg/jp"
)

func (pr *PatchRequest[T]) remove(initialResource *T, emptyResource *T) (*T, error) {
	compiledPath, err := jp.ParseString(pr.Path)
	if err != nil {
		return nil, fmt.Errorf("fail to compile path %s %s", pr.Path, err.Error())
	}

	b, err := json.Marshal(initialResource)
	if err != nil {
		return nil, fmt.Errorf("match fail to marshal input resource %s", err.Error())
	}

	var resourceAsMap interface{}
	err = json.Unmarshal(b, &resourceAsMap)
	if err != nil {
		return nil, fmt.Errorf("match fail to unmarshal %s", err.Error())
	}

	resourceAsMapModified, err := compiledPath.Remove(resourceAsMap)
	if err != nil {
		return nil, fmt.Errorf("error while deleting nodes %s", err.Error())
	}

	return pr.remarshal(resourceAsMapModified, emptyResource)
}
