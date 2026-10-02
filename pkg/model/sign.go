package model

import "github.com/genshinsim/gcsim/pkg/sharekey"

// SignedJSON marshals the result and signs those exact bytes. With a nil key the result is
// unsigned and the header is empty.
func (r *SimulationResult) SignedJSON(key *sharekey.Key) ([]byte, string, error) {
	if key != nil {
		r.KeyType = key.Class
	}
	data, err := r.MarshalJSON()
	if err != nil {
		return nil, "", err
	}
	if key == nil {
		return data, "", nil
	}
	return data, key.Sign(data), nil
}
