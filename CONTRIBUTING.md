# Contributing Guidelines

If you would like to contribute code to this project, please follow these pull request guidelines:

0. (Optional but encouraged) Find at least one maintainer interested in your PR.
1. Fork the project.
2. Create a branch specifically for the feature you are contributing.
3. (Optional but encouraged) Rebase your branch as needed. Please see quick reference if you are new to git.
4. After you are happy with your work, please make sure to submit a Pull Request from the feature branch. You are heavily discouraged from making pull requests from your main branch, because we may not be able to get to your PR before you make new changes to your PR.

Character PR checklist:

- [ ] New character package
- [ ] Config in character package
- [ ] Run pipeline with added config (generates character curve, talent stats, `.generated.json` files)
- [ ] Character key
- [ ] Shortcuts for character key
- [ ] Update `mode_gcsim.js` with shortcuts for syntax highlighting
- [ ] Add Character package to imports
- [ ] Normal Attack
- [ ] Charge Attack / Aimed Shot
- [ ] Skill
- [ ] Burst
- [ ] A1
- [ ] A4
- [ ] C1
- [ ] C2
- [ ] C3
- [ ] C4
- [ ] C5
- [ ] C6
- [ ] Other necessary talents (custom dash/jump, low/high plunge, ...)
- [ ] Hitlag
- [ ] ICD
- [ ] StrikeType
- [ ] PoiseDMG (blunt attacks only for now)
- [ ] Hitboxes
- [ ] Attack durability
- [ ] Particles
- [ ] Frames
- [ ] Update documentation
- [ ] Xingqiu/Yelan N0 (optional)
- [ ] Xianyun Plunge (optional)

Weapon PR checklist:

- [ ] New weapon package
- [ ] Config in weapon package
- [ ] Run pipeline with added config (generates weapon curve, `.generated.json` files)
- [ ] Weapon key
- [ ] Shortcuts for weapon key
- [ ] Add weapon package to imports
- [ ] Weapon passive (might include an attack)
- [ ] Update documentation

Artifact PR checklist:

- [ ] New artifact package
- [ ] Config in artifacts package
- [ ] Run pipeline with added config (generates `.generated.json` files)
- [ ] Artifact key
- [ ] Shortcuts for artifact key
- [ ] Add artifact package to imports
- [ ] 2pc
- [ ] 4pc
- [ ] Update documentation

Please try to be explicit about what is complete or incomplete.

Items may be omitted when irrelevant.

<details><summary>Click to expand copy-paste friendly version</summary>
  
```
Character PR checklist:

- [ ] New character package
- [ ] Config in character package
- [ ] Run pipeline with added config (generates character curve, talent stats, `.generated.json` files)
- [ ] Character key
- [ ] Shortcuts for character key
- [ ] Update `mode_gcsim.js` with shortcuts for syntax highlighting
- [ ] Add Character package to imports
- [ ] Normal Attack
- [ ] Charge Attack / Aimed Shot
- [ ] Skill
- [ ] Burst
- [ ] A1
- [ ] A4
- [ ] C1
- [ ] C2
- [ ] C3
- [ ] C4
- [ ] C5
- [ ] C6
- [ ] Other necessary talents (custom dash/jump, low/high plunge, ...)
- [ ] Hitlag
- [ ] ICD
- [ ] StrikeType
- [ ] PoiseDMG (blunt attacks only for now)
- [ ] Hitboxes
- [ ] Attack durability
- [ ] Particles
- [ ] Frames
- [ ] Update documentation
- [ ] Xingqiu/Yelan N0 (optional)
- [ ] Xianyun Plunge (optional)

Weapon PR checklist:

- [ ] New weapon package
- [ ] Config in weapon package
- [ ] Run pipeline with added config (generates weapon curve, `.generated.json` files)
- [ ] Weapon key
- [ ] Shortcuts for weapon key
- [ ] Add weapon package to imports
- [ ] Weapon passive (might include an attack)
- [ ] Update documentation

Artifact PR checklist:

- [ ] New artifact package
- [ ] Config in artifacts package
- [ ] Run pipeline with added config (generates `.generated.json` files)
- [ ] Artifact key
- [ ] Shortcuts for artifact key
- [ ] Add artifact package to imports
- [ ] 2pc
- [ ] 4pc
- [ ] Update documentation

````
</details>


# Git/Github Quick Reference Guide
For those who are new to git/github

```git checkout -b newbranchname``` creates a new branch from your current branch
If you have committed code, but upon finishing your feature, the main branch has progressed, you are encouraged to rebase it to ensure it still works.
Please reach out for help if you are not sure how to do this step, the following steps can be dangerous and you can lose your work if not done correctly.

To rebase your branch you will need to run the command
````

git rebase --onto <newparent> <oldparent>
git push -f

````
Where new parent is the commitment hash of the newest commit on the main branch and old parent is the commitment hash of the oldest common commitment between your feature branch and the main branch.


# Dev Environment Setup

Install the following:
- Golang v1.21+: https://go.dev/doc/install
    - Alternative: https://github.com/stefanmaric/g this will keep golang up to date and configure all the necessary paths for you
- Git: https://github.com/git-guides/install-git
- (optional) Protobuf compiler v21+: https://github.com/protocolbuffers/protobuf/releases
- vscode is the “officially supported” IDE for the repo: https://code.visualstudio.com/
    - Install the recommended extensions.
- golangci-lint: If using vscode, should give popup to install
    - Alternative: manual install https://golangci-lint.run/usage/install/#local-installation
    - Linters that use diffs (for example: gofmt) require [additional setup on Windows](https://github.com/golangci/golangci-lint/issues/307#issuecomment-1001301930) (This assumes you already installed Git for Windows!):
        - Open PowerShell as Administrator
        - execute the following command: setx /M PATH "$($env:path);C:\Program Files\git\usr\bin"
        - Restart VS Code if it is currently open for the changes to take effect.

# Local Testing/Building
0. Create a Config file
1. Navigate to ```./gcsim/cmd/gcsim```
2. Run ```go build``` to build the executable and then feed your config file in e.g. ```./gcsim.exe -c config.txt -sample config -gz``` OR run ```go run . --c config.txt -sample config -gz```
3. Upload the generated sample file to the [Sample page](https://gcsim.app/sample/upload) to confirm everything is working accordingly, and optionally share the sample file in discord for debugging help.
````
