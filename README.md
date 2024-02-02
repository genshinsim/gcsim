# Overview

gcsim is a Monte Carlo simulation tool used to model Genshin Impact's combat. The user can input a set of characters, targets, options, and actions to perform, and then gcsim executes these actions. It outputs a variety of results, such as mean DPS and the DPS distribution across iterations. The user can also scroll through a sample of 1 iteration, which comprehensively lists every action, damage instance, reactions, buffs, etc. frame by frame.

# Getting Started

Primary usage of gcsim is through the webapp: [https://gcsim.app](https://gcsim.app). You also can download the latest build and run it as a CLI [here](https://github.com/genshinsim/gcsim/releases). Our [docs](https://docs.gcsim.app/guides/building_a_simulation_basic_tutorial) explain how to write and understand configs.

Any issues or questions can be shared on our [Discord](https://discord.gg/m7jvjdxx7q), where experienced users can take a look.

## Project Status

The project is still under development. While the majority of characters, items, and game mechanics have been implemented, there are still improvements that can be made, which you can find in our [issues](https://github.com/genshinsim/gcsim/issues?q=is%3Aopen+is%3Aissue) and [discussions](https://github.com/genshinsim/gcsim/discussions).

## Contributing

Here are a few ways to help improve the quality of gcsim:
- Record exhaustive frame counts of new unit actions, methodology detailed [here](https://docs.gcsim.app/mechanics/frames/).
- Validate gameplay and sim results, ensure the sim can reproduce damage calculations, reactions, and buff uptimes faithfully.
- Build action lists aka "rotations" for any team composition and submit them to our [Config Database](https://simpact.app/) via [Discord](https://discord.gg/m7jvjdxx7q).
- Use gcsim for gear, rotation, and team comparisons, while scrutinizing both expected and unexpected results. This is the best way potential issues can be spotted.

gcsim is always looking for developers. If you would like to contribute code, please look at the [contributing guidelines](CONTRIBUTING.md).
