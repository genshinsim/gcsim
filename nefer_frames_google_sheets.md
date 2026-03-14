# Nefer Frame Data Dump

- Fetched: 2026-03-14

## Sources

- Sheet 1: [Workbook 1](https://docs.google.com/spreadsheets/d/1JMqNGKIXK8KM_rpP3eu8lYFnjYbIuX45f_QqxSBsTYY/edit?gid=0#gid=0)
- Export TSV 1: [Workbook 1 TSV export](https://docs.google.com/spreadsheets/d/1JMqNGKIXK8KM_rpP3eu8lYFnjYbIuX45f_QqxSBsTYY/export?format=tsv&gid=0)
- Sheet 2: [Workbook 2](https://docs.google.com/spreadsheets/d/1VpqRMv7FHPayKm9PBbazJs6gL2pr4Ka2F-EI0VgkpDQ/edit?gid=0#gid=0)
- Export TSV 2: [Workbook 2 TSV export](https://docs.google.com/spreadsheets/d/1VpqRMv7FHPayKm9PBbazJs6gL2pr4Ka2F-EI0VgkpDQ/export?format=tsv&gid=0)

## Workbook Page Inventory

- Workbook 1 pages: NormalCharge Attacks, Skill, DashJump, Burst, N0
- Workbook 2 pages: Yelan, Xingqu, Phantasm CA
- An earlier draft only summarized the first visible page from each workbook.
- This version explicitly includes data discovered from the additional pages.

## Interpretation Rules

- Confirmed: all measured trials agree.
- Mixed: trials disagree, so only the measured range is safe.
- Missing: source sheet does not provide a usable value.
- No mixed value is resolved by guessing in this document.

## Workbook 1 Summary

### Page: NormalCharge Attacks

Source notes:

- Recording by caramielle
- Counted by caramielle
- Video: [caramielle recording](https://youtu.be/QeYtpxgiZeo)

Hitmarks:

- N1 hitmark: 13 / 10 / 10 -> Mixed, measured range 10-13
- N2 hitmark: 11 / 8 / 8 -> Mixed, measured range 8-11
- N3-1 hitmark: 10 / 10 / 10 -> Confirmed 10
- N3-2 hitmark: 16 / 16 / 16 -> Confirmed 16
- N4 hitmark: 22 / 24 / 22 -> Mixed, measured range 22-24
- Normal CA hitmark: 24 / 24 / 24 -> Confirmed 24, does not include CA windup
- Phantasm CA hit 1 (Nefer): 10 / 10 / 10 -> Confirmed 10, does not include CA windup
- Phantasm CA hit 2 (LB): 15 / 15 / 15 -> Confirmed 15
- Phantasm CA hit 3 (LB): 8 / 8 / 8 -> Confirmed 8
- Phantasm CA hit 4 (Nefer): 24 / 24 / 25 -> Mixed, measured range 24-25
- Phantasm CA hit 5 (LB): 0 / 1 / 0 -> Mixed, measured range 0-1, source note says measurement follows the damage number
- Verdant Dew consumed timing: 9 / 9 / 9 -> Confirmed 9, source note says start frame is CA start, not CA windup start

Cancel frames:

- N1 -> N2: 15 / 15 / 15 -> Confirmed 15
- N1 -> Normal CA: 40 / 41 / 41 -> Mixed, measured range 40-41, includes CA windup
- N1 -> Phantasm CA: 40 / 40 / 40 -> Confirmed 40, includes CA windup
- N1 -> Walk: 33 / 32 / 32 -> Mixed, measured range 32-33
- N2 -> N3: 21 / 19 / 18 -> Mixed, measured range 18-21
- N2 -> Normal CA: 38 / 38 / 38 -> Confirmed 38
- N2 -> Phantasm CA: 39 / 38 / 39 -> Mixed, measured range 38-39
- N2 -> Walk: 34 / 34 / 34 -> Confirmed 34
- N3 -> N4: 53 / 48 / 48 -> Mixed, measured range 48-53; source row label corrected from N3 -> N3 to N3 -> N4
- N3 -> Normal CA: 62 / 61 / 62 -> Mixed, measured range 61-62
- N3 -> Phantasm CA: 60 / 60 / 60 -> Confirmed 60
- N3 -> Walk: 50 / 50 / 50 -> Confirmed 50
- N4 -> N1: 50 / 53 / 54 -> Mixed, measured range 50-54
- N4 -> Walk: 62 / 63 / 64 -> Mixed, measured range 62-64
- Normal CA -> N1: 29 / 29 / 29 -> Confirmed 29, does not include CA windup
- Normal CA -> Normal CA: 48 / 48 / 49 -> Mixed, measured range 48-49, includes CA windup for the second CA but not the first
- Normal CA -> Walk: 71 / 71 / 71 -> Confirmed 71
- Normal CA -> Swap: 28 / 28 / 28 -> Confirmed 28
- Phantasm CA -> N1: 69 / 69 / 69 -> Confirmed 69, does not include CA windup
- Phantasm CA -> Normal CA: 65 / 65 / 65 -> Confirmed 65, includes CA windup for the second CA but not the first
- Phantasm CA -> Phantasm CA: 65 / 65 / 65 -> Confirmed 65, includes CA windup for the second CA but not the first
- Phantasm CA -> Walk: 105 / 104 / 104 -> Mixed, measured range 104-105
- Phantasm CA -> Swap: 0 / 0 / 0 -> Confirmed 0
- Any N1/N2/N3/N4/CA transition not listed above to Skill, Burst, Dash, Jump, or Swap remains Missing in the source sheet.

### Page: Skill

- Skill hitmark: 24 / 24 / 24 -> Confirmed 24
- Skill CD start: 22 / 22 / 22 -> Confirmed 22
- Skill -> N1: 26 / 26 / 26 -> Confirmed 26
- Skill -> Normal CA: 47 / 49 / 48 -> Mixed, measured range 47-49, includes CA windup
- Skill -> Phantasm CA: 28 / 29 / 29 -> Mixed, measured range 28-29, does not include CA windup per source note
- Skill -> Skill: 74 / 69 / 71 -> Mixed, measured range 69-74
- Skill -> Burst: 26 / 28 / 26 -> Mixed, measured range 26-28
- Skill -> Dash: 37 / 38 / 39 -> Mixed, measured range 37-39
- Skill -> Jump: 38 / 38 / 38 -> Confirmed 38
- Skill -> Walk: 34 / 34 / 34 -> Confirmed 34
- Skill -> Swap: 24 / 25 / 25 -> Mixed, measured range 24-25

### Page: DashJump

- Regular Jump: 0 / 0 / 0 -> Present in source, but semantics are not clear enough to hard-code from the sheet alone
- Regular Dash: 0 / 0 / 0 -> Present in source, but semantics are not clear enough to hard-code from the sheet alone
- Safe conclusion: this page confirms a separate Dash/Jump study exists, but not an implementation-ready hook value.

### Page: Burst

- Burst hitmark 1: 101 / 101 / 99 -> Mixed, measured range 99-101
- Burst hitmark 2: 44 / 44 / 40 -> Mixed, measured range 40-44
- Burst CD start: 0 / 0 / 0 -> Confirmed 0
- Burst energy drained: 6 / 7 / outlier -> Mixed, not stable enough to hard-code from this row alone
- Burst -> N1: 119 / 118 / 119 -> Mixed, measured range 118-119
- Burst -> Normal CA: 131 / 129 / 130 -> Mixed, measured range 129-131, includes CA windup
- Burst -> Phantasm CA: 124 / 123 / 124 -> Mixed, measured range 123-124, includes CA windup
- Burst -> Skill: 119 / 119 / 119 -> Confirmed 119
- Burst -> Dash: 120 / 120 / 119 -> Mixed, measured range 119-120
- Burst -> Jump: 119 / 119 / 120 -> Mixed, measured range 119-120
- Burst -> Walk: 119 / 118 / 119 -> Mixed, measured range 118-119
- Burst -> Swap: 118 / 117 / 118 -> Mixed, measured range 117-118

### Page: N0

- Yelan: 0 / 0 / 0 -> Present on workbook 1 N0 page, but not implementation-ready by itself
- Xingqiu: 0 / 0 / 0 -> Present on workbook 1 N0 page, but not implementation-ready by itself
- Safe conclusion: the workbook 1 N0 page is not enough on its own to hard-code a gcsim integration value.

## Workbook 2 Summary

Source notes:

- Video and count by nekorul
- Video: [nekorul recording](https://youtu.be/rPlOnZsbiu8)
- Start frame is when Nefer moves abruptly, though her weapon disappears 1 frame earlier.
- End frame is when team icons go gray.

Pages present:

- Yelan
- Xingqu
- Phantasm CA

### Yelan Page

Group A (140ms):

- Trial frames: 8, 10, 10, 10, 10, 10
- Success flags: FALSE, TRUE, TRUE, TRUE, TRUE, TRUE
- Average successful frame: 10
- Safe conclusion: 10f works consistently in this setup; 8f failed in one recorded trial.

Group B (135ms):

- Trial frames: 10, 9, 9, 10, 9, 10
- Success flags: TRUE, TRUE, FALSE, TRUE, FALSE, TRUE
- Average successful frame: 9.75
- Source note: 9f is possible but not consistent
- Safe conclusion: 10f is more reliable than 9f in this setup.

Derived timing summary:

- Yelan: 10
- Xingqu: 11
- The sheet does not explicitly define whether these map to startup delay, usable N0 trigger timing, or another integration frame, so interpretation still requires analysis against the intended gcsim hook point.

### Phantasm CA Page

- Workbook 2 contains a separate Phantasm CA study.
- It appears to measure abrupt movement, team gray state, swap input, and success timing around Phantasm CA.
- The extracted rows are not labeled clearly enough to replace the already summarized confirmed values without another interpretation pass.

## Practical Extraction For Implementation

### Confirmed Values Safe To Use Directly

- N3-1 hitmark 10
- N3-2 hitmark 16
- Normal CA hitmark 24 excluding windup
- Phantasm CA hitmarks 10, 15, 8 for hits 1-3
- Verdant Dew consumed timing 9 from CA start
- Skill hitmark 24
- Skill CD start 22
- N1 -> N2 cancel 15
- N1 -> Phantasm CA cancel 40 including windup
- N2 -> Normal CA cancel 38
- N2 -> Walk cancel 34
- N3 -> Phantasm CA cancel 60
- N3 -> Walk cancel 50
- Normal CA -> N1 cancel 29 excluding windup
- Normal CA -> Walk cancel 71
- Normal CA -> Swap cancel 28
- Phantasm CA -> N1 cancel 69 excluding windup
- Phantasm CA -> Normal CA cancel 65
- Phantasm CA -> Phantasm CA cancel 65
- Phantasm CA -> Swap cancel 0
- Skill -> N1 cancel 26
- Skill -> Jump cancel 38
- Skill -> Walk cancel 34
- Burst CD start 0
- Burst -> Skill cancel 119
- Yelan timing summary 10

### Values Requiring Additional Analysis

- Any hitmark or cancel entry marked Mixed above
- Any entry marked Missing above
- Any row from the workbook 1 DashJump page until its semantics are tied to a concrete gcsim hook
- Xingqiu 11 and Yelan 10 semantics relative to the exact gcsim AnimationStartDelay integration point
- Phantasm CA hit 5 (LB) timing, because the measurement itself is based on the damage number and differs between trials
- The workbook 2 Phantasm CA page, because its extracted rows still need an interpretation pass before they can replace existing timing assumptions
