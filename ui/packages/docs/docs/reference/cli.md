---
title: CLI
sidebar_position: 1
---

## Usage

```
./gcsim.exe <options>
```

The gcsim CLI accepts the following options:

| Option | Description | Input | Default | 
| --- | --- | --- | --- |
| `-c` | Which config file to run gcsim on. This option is required. | Path to the config file. | config.txt | 
| `-out` | Which file to output the results of the gcsim command to. | Output file path. | disabled |
| `-sample` | Which file to output a sample result to. | Output file path. | disabled |
| `-sampleMinDps` | Similar to `-sample`, except that we're writing out the min-DPS run as the sample. | Output file path. | disabled |
| `-sampleMaxDps` | Similar to `-sample`, except that we're writing out the max-DPS run as the sample. | Output file path. | disabled |
| `-nr` | Whether to disable running the simulation. Useful in combination with `-sample` if a sample result is all that is needed. | - | disabled |
| `-gz` | Whether to zip up the results. Can be used together with `-out` and `-sample`. | - | disabled |
| `-s` | Whether to serve the results to the "local" gcsim viewer page using the default web browser. | - | disabled |
| `-nb` | Whether to open the default web browser when using `-s`. gcsim will wait until the "local" gcsim viewer page has been opened and then output the results onto that site. | - | disabled | 
| `-ks` | Whether to keep serving results to the "local" gcsim viewer page when using `-s`. | - | disabled | 
| `-substatOptim` | Whether to perform substat optimization on the config file. Use the `-out` flag to output the optimized config to a new config file. | - | disabled | 
| `-substatOptimFull` | Similar to `-subtatOptim`, but the optimized config is output to the config file given by `-c` and gcsim is run on that optimized config. | - | disabled |
| `-options` | Additional options to customize the substat optimizer. | Options string. | disabled |
| `-v` | Whether to enable verbose output. This is exclusive to `-substatOptim` and `-substatOptimFull` at the moment. | - | disabled |
| `-cpuprofile` | Create a CPU profile file. Used to analyse the performance of gcsim. The results can be viewed in the browser via `go tool pprof -http=localhost:3000` for example (requires [Graphviz](https://graphviz.org/)). | Output file path. | disabled |
| `-memprofile` | Create a memory profile file. Used to analyse the performance of gcsim. The results can be viewed in the browser via `go tool pprof -http=localhost:3000` for example (requires [Graphviz](https://graphviz.org/)). | Output file path. | disabled | 
| `-update` | Update the current gcsim CLI with the latest version released on [GitHub](https://github.com/genshinsim/gcsim/releases). | - | false |
| `-version` | Print the version of the current gcsim CLI. | - | false |

### Input

Input for options can be provided either via `<option> <value>` or `<option>=<value>`.

#### Example

```
./gcsim.exe -c="test.txt"
```

or

```
./gcsim.exe -c test.txt
```

:::caution
In case of file paths you might need to wrap it in " for it to be interpreted correctly as shown in the example.
:::

### Additional Options For Substat Optimizer: `-options`

The input has to be specified as `-options="<option list>"`. 
The option list has the following format: `<option>=<value>` with `;` as the separator.

| option | description | default |
| --- | --- | --- |
| `total_liquid_substats` | Total liquid substats available to be assigned across all substats. | 20 |
| `indiv_liquid_cap` | Total liquid substats that can be assigned to a single substat. | 10 |  
| `fixed_substats_count` | Amount of fixed substats that are assigned to all substats. | 2 |
| `fine_tune` | 0 to disable and 1 to enable the fine-tune step. This step will compare ER vs DMG substats after having allocated all substats and cover when the initial ER heuristic fails due to not replacing `.<char>.burst.ready` with `.<char>.burst.ready && .<char>.energy == .<char>.energymax` in conditionals (Instead of `<char>`, put the name of the affected character). | 1 |

#### Example

```
./gcsim.exe -c test.txt -s -substatOptimFull -options="total_liquid_substats=10;fixed_substats_count=4;fine_tune=0"
```
