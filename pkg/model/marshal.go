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

func (r *SimulationResult) MarshalJson() ([]byte, error) {
	return marshalOptions().Marshal(r)
}

func (r *SimulationResult) UnmarshalJson(b []byte) error {
	return unmarshalOptions().Unmarshal(b, r)
}

func (r *SimulationStatistics) MarshalJson() ([]byte, error) {
	return marshalOptions().Marshal(r)
}

func (r *SignedSimulationStatistics) MarshalJson() ([]byte, error) {
	return marshalOptions().Marshal(r)
}

func (r *Sample) MarshalJson() ([]byte, error) {
	return marshalOptions().Marshal(r)
}
