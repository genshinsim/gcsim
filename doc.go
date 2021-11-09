/**
Package gsim provides a monte carlo simulation to Genshin Impact combat system as a way for
users to calculate a team's overall dps over a certain combat duration

The simulator can be run'd as follows in it's simplest form
  result, err := gsim.Run(config, opts)
  if err != nil {
	handleError(err)
  }
  fmt.Print(result.PrettyPrint())

The simulator is made up of various packages.

- pkg/core
- pkg/character
- pkg/monster
- pkg/parse
- pkg/shield


**/

package gcsim
