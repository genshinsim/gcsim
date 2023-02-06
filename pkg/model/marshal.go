package model

import "google.golang.org/protobuf/encoding/protojson"

func marshalOptions() protojson.MarshalOptions {
	return protojson.MarshalOptions{
		AllowPartial:    true,
		UseEnumNumbers:  true, // TODO: prob better if we set to false?
		EmitUnpopulated: false,
	}
}

func (r *SimulationResult) MarshalJson() ([]byte, error) {
	return marshalOptions().Marshal(r)
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

func (r *DBEntry) MarshalJson() ([]byte, error) {
	return marshalOptions().Marshal(r)
}

func (r *DBEntries) MarshalJson() ([]byte, error) {
	return marshalOptions().Marshal(r)
}

func (r *ComputeWork) MarshalJson() ([]byte, error) {
	return marshalOptions().Marshal(r)
}
