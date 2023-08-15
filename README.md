# Overview

gcsim is a Monte Carlo simulation tool used to model Genshin Impact's combat. The user inputs a set of characters, targets, options, and actions to perform, and then gcsim will execute these actions. It outputs a variety of results, such as mean DPS and the DPS distribution across iterations. The user can also scroll through a sample of 1 iteration, which comprehensively lists every action, damage instance, reactions, buffs, etc. frame by frame.

# Getting started

Primary usage of gcsim is through the webapp: [https://gcsim.app](https://gcsim.app). You can download the latest build and run it as a CLI [here](https://github.com/genshinsim/gcsim/releases).

## Project status

The project is still under development. While the vast majority of characters, items, and game mechanics have been implemented, there are still improvements that can be made, which you can find in our [issues](https://github.com/genshinsim/gcsim/issues?q=is%3Aopen+is%3Aissue) and [discussions](https://github.com/genshinsim/gcsim/discussions).

## Contributing

If you are looking to contribute, here are some key areas that you can help out with. If you're interested, please reach out to our [Discord](https://discord.gg/m7jvjdxx7q).

- Comprehensive frame counts of new character's actions, methodology detailed [here](https://docs.gcsim.app/mechanics/frames/).
- Validation between actual game play and sim results. This is done by recording footage, and then reproducing the same team/artifact/items/targets in the sim and comparing results frame by frame. This means both the damage output as well as the reactions. The simulator should be able to reproduce it faithfully.
- Building action lists aka "rotations" for various common team comps and submitting them to our [Config Database](https://simpact.app/) via Discord.
- Helping with documenting the sim in the wiki.
- Further testing of in game reactions, primarily dendro.
- Testing of 5* constellations
- Just in general using the sim for calculations/weapon comparisons/day 1 testing etc...
- If you would like to contribute code please take a look at the [contributing guidelines](CONTRIBUTING.md)
