# Overview

GSim is a Monte Carlo simulation tool used to model team dps. It allows for the formation of arbitrary combination of any characters.

# Getting started

**NEW:** There is now a web version available at [https://genshinsim.github.io/gsimweb/](https://genshinsim.github.io/gsimweb/). This is basically just the program compiled into WASM and embedded into a webpage. Works the same as the release version but with the downside that it is much slower as it is only single threaded.

Go to the [releases page](https://github.com/genshinsim/gsim/releases) and download the latest development release for your platform. Currently development releases are [automatically](https://github.com/genshinsim/gsim/blob/main/.github/workflows/release.yaml) built whenever a new commit is pushed to the main branch.

The archive can be extracted to any folder. Once extracted, simply run the executable (gsim.exe). Doing so should open up your default browser to http://localhost:8081 (If it did not, you can simply browse to this address on a browser).

**Note that Chrome may complain that the file is not commonly downloaded and may be dangerous. This is expected as the sim is not all that popular. You can always build the project from source if you wish**

For more information, please visit the [starter guide](https://github.com/genshinsim/gsim/wiki/Starter) on the wiki.

## Project status

The project is still currently under heavy development. Not every character/weapon/artifacts have been implemented. There are also many quality of life improvements that can be made. Currently this project is developed by one person (me) and as such development speed is not as fast as it can be. If you wish to contribute, please see below for ways you can help out.

## Contributing

If you are looking to contribute, the following are some key areas that you can help out with. If you're interested, please don't hesistate to reach out on the sim's [discord](https://discord.gg/m7jvjdxx7q)

- Comprehensive frame counts of every character's abilities, including but not limited to:
    - Hit mark frame (not just animation frame, which is currently in the KQM library)
    - Cooldown start frame
    - Energy drain frame
    - Particle proc frame (relative to hit mark)
- Team comp damage validation between actual game play and sim. This is done by recording videos of actual gameplay footage, and then reproducing the same team/artifact/items/targets in the simulator and comparing results frame by frame. This means both comparing the damage output as well as the reactions. The simulator should be able to reproduce actual gameplay faithfully.
- Building action list for various common team comps and comparing/validating the result vs actual in game gameplay
- Helping with documenting the sim in the wiki
- Further testing of in game reactions, primarily EC and chain freeze duration.
- Just in general using the sim for calculations/weapon comparisons/day 1 testing etc...

## Credits

- Most of the % data: https://genshin.honeyhunterworld.com/
- Tons of discussions on KeqingMain (to be added up with ppl's discord tags at some point)
- Most if not all the frame data came from https://library.keqingmains.com
- All the folks at KQM that helped out with testing; special thanks to (in no particular order):
  - Yukarix#6534
  - Aluminum#5462
  - Terrapin#8603
