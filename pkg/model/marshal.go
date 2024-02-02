package model

import (
	"google.golang.org/protobuf/encoding/protojson"
)

func marshalOptions() protojson.MarshalOptions {
	return protojson.MarshalOptions{
		AllowPartial:    true,
		UseEnumNumbers:  true, // TODO: prob better if we set to false?
		EmitUnpopulated: false,
	}
}
func unmarshalOptions() protojson.UnmarshalOptions {
	return protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}
}

func (r *SimulationResult) MarshalJSON() ([]byte, error) {
	return marshalOptions().Marshal(r)
}

func (r *SimulationResult) UnmarshalJSON(b []byte) error {
	return unmarshalOptions().Unmarshal(b, r)
}

func (r *SimulationStatistics) MarshalJSON() ([]byte, error) {
	return marshalOptions().Marshal(r)
}

func (r *SignedSimulationStatistics) MarshalJSON() ([]byte, error) {
	return marshalOptions().Marshal(r)
}

func (s *Sample) MarshalJSON() ([]byte, error) {
	return marshalOptions().Marshal(s)
}
