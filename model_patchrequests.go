package jsonpatch

type PatchRequests[T any] struct {
	Patches []*PatchRequest[T] `json:"patches" validate:"required,min=1,dive,required"`
}

// Apply TODO
func (prs *PatchRequests[T]) Apply(initialResource *T, newer func() *T) (*T, error) {
	var err error
	patched := initialResource
	for _, pr := range prs.Patches {
		patched, err = pr.Apply(patched, newer())
		if err != nil {
			return nil, err
		}
	}
	return patched, nil
}
