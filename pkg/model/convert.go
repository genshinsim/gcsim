package model

// TODO: SimDuration & TotalDamage should be a OverviewStats instead?
func (result *SimulationResult) ToDBEntry() *DBEntry {
	names := make([]string, len(result.CharacterDetails))
	for i, c := range result.CharacterDetails {
		names[i] = c.Name
	}

	return &DBEntry{
		SimDuration: &DescriptiveStats{
			Min:  result.Statistics.Duration.Min,
			Max:  result.Statistics.Duration.Max,
			Mean: result.Statistics.Duration.Mean,
			SD:   result.Statistics.Duration.SD,
		},
		TotalDamage: &DescriptiveStats{
			Min:  result.Statistics.TotalDamage.Min,
			Max:  result.Statistics.TotalDamage.Max,
			Mean: result.Statistics.TotalDamage.Mean,
			SD:   result.Statistics.TotalDamage.SD,
		},
		TargetCount:      int32(len(result.TargetDetails)),
		Hash:             *result.SimVersion,
		Config:           result.Config,
		MeanDpsPerTarget: *result.Statistics.TotalDamage.Mean / (float64(len(result.TargetDetails)) * *result.Statistics.Duration.Mean),
		Team:             result.CharacterDetails,
		CharNames:        names,
		Mode:             result.Mode,
	}
}
