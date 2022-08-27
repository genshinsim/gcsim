# Overview

gcsim is a Monte Carlo simulation tool used to model team dps. It allows for the formation of arbitrary combination of any characters.

# Getting started

You can download the latest build [here](https://github.com/genshinsim/gcsim/releases). You can also use the webapp: [https://gcsim.app](https://gcsim.app).

## Project status

The project is still under development. While every character up until 2.8 has been implemented with proper frames and hitlag, there are still improvements that can be made, which you can find in our [issues](https://github.com/genshinsim/gcsim/issues?q=is%3Aopen+is%3Aissue) and [discussions](https://github.com/genshinsim/gcsim/discussions).

## Contributing

If you are looking to contribute, the following are some key areas that you can help out with. If you're interested, please don't hesistate to reach out on the sim's [discord](https://discord.gg/m7jvjdxx7q)

- Comprehensive frame counts of new character's abilities, including but not limited to:
  - Hit mark frame (not just animation frame, which is currently in the KQM library)
  - Cooldown start frame
  - Energy drain frame
  - Particle proc frame (relative to hit mark)
- Team comp damage validation between actual game play and sim. This is done by recording videos of actual gameplay footage, and then reproducing the same team/artifact/items/targets in the simulator and comparing results frame by frame. This means both comparing the damage output as well as the reactions. The simulator should be able to reproduce actual gameplay faithfully.
- Building action list for various common team comps and comparing/validating the result vs actual in game gameplay
- Helping with documenting the sim in the wiki
- Further testing of in game reactions, primarily dendro.
- Testing of 5* constellations
- Just in general using the sim for calculations/weapon comparisons/day 1 testing etc...
- If you would like to contribute code please take a look at the [contributing guidelines](CONTRIBUTING.md)
