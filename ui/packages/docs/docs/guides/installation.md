---
title: Installing the CLI
sidebar_position: 5
---

:::info
Please refer to [this guide](/guides/using_cli) for information on how to use the gcsim CLI.
:::

You have two options to obtain the gcsim [CLI](https://en.wikipedia.org/wiki/Command-line_interface) (`gcsim.exe`):
- Downloading The CLI - *recommended*
- Building The CLI From Source

## Assumptions

The following descriptions assume:
- Windows as the operating system
- main working folder where gcsim source code should be stored is called `TC`

The file and folder names are just examples.

## Downloading The CLI

The CLI (`gcsim.exe`) for the latest version of gcsim can be downloaded from the [Releases page](https://github.com/genshinsim/gcsim/releases) of [gcsim's GitHub repository](https://github.com/genshinsim/gcsim).

## Building The CLI From Source

### Setup

#### Install Go

gcsim is mainly written in a programming langugage called [Go](https://go.dev/). 
Download and install the latest version of Go from [here](https://go.dev/doc/install).

#### Install Git

gcsim uses [Git](https://git-scm.com/) for [version control](https://en.wikipedia.org/wiki/Version_control). 
Download and install the latest version of Git for Windows from [here](https://gitforwindows.org/).

#### Install IDE (VS Code)

In order to have a smooth experience while navigating the gcsim source code, it is recommend to install an [IDE](https://en.wikipedia.org/wiki/Integrated_development_environment).
For this guide, it is recommended to download and install VS Code from [here](https://code.visualstudio.com/).

#### Download Source Code

To obtain the gcsim source code, you need to clone [gcsim's GitHub repository](https://github.com/genshinsim/gcsim) into a folder of your choice.

1. Start VS Code.
2. Open the `TC` folder in VS Code.
3. Click on `Terminal` in the top left menu.
4. Click on `New Terminal`. 
This will open Powershell in the `TC` folder.
5. Type `git clone https://github.com/genshinsim/gcsim.git` into the Powershell window.

Your initial folder structure should now look like this:
```
└── TC
    └── gcsim
```

### Build The CLI (gcsim.exe)

1. Start VS Code.
2. Open the topmost `gcsim` folder in VS Code.
3. Click on `Terminal` in the top left menu.
4. Click on `New Terminal`. 
This will open Powershell in the `gcsim` folder.
5. Type `cd cmd/gcsim` into the Powershell window. This will put you in the folder where you can build the CLI.
6. Type `go build`. This will create a `gcsim.exe` in the current folder.

Now your folder structure should look like this:
```
└── TC 
    └── gcsim 🡐 VS Code is here
        └── cmd
            └── gcsim 🡐 VS Code's Powershell window is here
                └── gcsim.exe
```

:::info
If you modified your local gcsim source code and would like to try out the changes in the CLI, then just build the CLI again.
:::

### Apply Updates

In order to keep your CLI up-to-date with changes made by the gcsim dev team, you need to do the following occasionally:

1. Start VS Code.
2. Open the topmost `gcsim` folder in VS Code.
3. Click on `Terminal` in the top left menu.
4. Click on `New Terminal`. 
This will open Powershell in the `gcsim` folder.
5. Type `git pull` into the Powershell window.
6. Rebuild the `gcsim.exe` as described in steps 5 and 6 in [this section](#building-the-cli-gcsimexe).
