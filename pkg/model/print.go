package model

import (
	"fmt"
	"strings"
)

func (r *SimulationResult) PrettyPrint() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(
		"Average duration of %.2f seconds (min: %.2f max: %.2f std: %.2f)\n",
		r.Statistics.Duration.Mean, r.Statistics.Duration.Min,
		r.Statistics.Duration.Max, r.Statistics.Duration.SD))
	sb.WriteString(fmt.Sprintf(
		"Average %.2f damage over %.2f seconds, resulting in %.0f dps (min: %.2f max: %.2f std: %.2f) \n",
		r.Statistics.TotalDamage.Mean, r.Statistics.Duration.Mean,
		r.Statistics.DPS.Mean, r.Statistics.DPS.Min, r.Statistics.DPS.Max, r.Statistics.DPS.SD))
	sb.WriteString(fmt.Sprintf(
		"Simulation completed %v iterations\n", r.Statistics.Iterations))

	return sb.String()
}
