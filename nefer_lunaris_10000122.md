# Nefer Reference Dump

- Source page: [Lunaris character page](https://lunaris.moe/character/10000122)
- Primary data endpoint: [Lunaris JSON endpoint](https://api.lunaris.moe/data/latest/en/char/10000122.json)
- Fetched: 2026-03-14

## Overview

- Name: Nefer
- Element: Dendro
- Weapon: Catalyst
- Rarity: 5-star
- Short description: The remarkably resourceful owner of the Curatorium of Secrets.
- Constellation: Ludus Latrunculorum
- Birthday: May 9th
- Level 90 base stats from charlist.json:
  - HP: 12704
  - ATK: 344
  - DEF: 799
  - Ascension stats: CRIT DMG 38.4%, Elemental Mastery 100

## Core Mechanics Summary

- Normal attack name: Striking Serpent.
- Charged Attack enters Slither, drains stamina while moving forward for up to 2.5s, then consumes extra stamina on exit to deal Dendro damage.
- Using Skill or sprinting during Slither does not force exit from Slither.
- Skill name: Senet Strategy: Dance of a Thousand Nights.
- Skill deals AoE Dendro damage, moves Nefer forward, and enters Shadow Dance.
- Shadow Dance increases interruption resistance.
- While in Shadow Dance, if Nefer has at least 1 Verdant Dew, Charged Attack is replaced by Phantasm Performance and no longer consumes stamina.
- Skill has 2 initial charges.
- Burst name: Sacred Vow: True Eye's Phantasm.
- Burst deals AoE Dendro damage in front and consumes all Veils of Falsehood to increase Burst damage.
- Passive P1 converts on-field Dendro Cores and future Lunar-Bloom outputs into Seeds of Deceit for 15s after Skill, and lets Charged Attack / Phantasm Performance absorb Seeds of Deceit to gain Veil of Falsehood stacks.
- At 3 Veil stacks, or when the third stack is refreshed, Nefer gains +100 EM for 8s.
- P2 strengthens Shadow Dance interaction with Verdant Dew after allies trigger Lunar-Bloom.
- P3 converts Bloom into Lunar-Bloom and scales Lunar-Bloom base damage from Nefer's EM.

## Combat Metadata From API

### Attack Metadata Entries

1. ExtraAttack_MoonLight_Gadget | element=Grass | gauge=0U | icd=Independent | source=Tag: MoonOvergrowDamage (!) | poise=LV2/20
2. ExtraAttack_MoonLight_Gadget | element=Grass | gauge=0U | icd=Independent | source=Tag: MoonOvergrowDamage (!) | poise=LV2/20
3. ExtraAttack_MoonLight_Gadget | element=Grass | gauge=0U | icd=Independent | source=Tag: MoonOvergrowDamage (!) | poise=LV3/30
4. ElementalArt | element=Grass | gauge=1U | icd=Independent | source=Tag: ElementalArt (!) | poise=LV3/50
5. SkillObj_ElementalBurst_Attacker | element=Grass | gauge=1U | icd=Independent | source=Tag: Nefer_ElementalBurst (!) | poise=LV3/50
6. SkillObj_ElementalBurst_Attacker | element=Grass | gauge=1U | icd=Independent | source=Tag: Nefer_ElementalBurst (!) | poise=LV5/100
7. FallingAnthem | element=Grass | gauge=0U | icd=Independent | source=Tag: FallingAttack (!) | poise=LV2/5
8. FallingAnthem | element=Grass | gauge=1U | icd=Independent | source=Tag: FallingAttack (!) | poise=LV3/50
9. FallingAnthem | element=Grass | gauge=1U | icd=Independent | source=Tag: FallingAttack (!) | poise=LV4/100
10. Constellation_6_Gadget | element=Grass | gauge=0U | icd=Independent | source=Tag: MoonOvergrowDamage (!) | poise=LV3/30

### Energy Generation Metadata

- ElementalArt_DropBall: 66% chance, 3 particles, internal cd 0.2s
- ElementalArt_DropBall: 33% chance, 2 particles, internal cd 0.2s

## Talents

### Normal Attack: Striking Serpent

Description:

- Normal Attack: Performs up to 4 kicks that deal Dendro DMG with the ferocity and grace of a striking serpent.
- Charged Attack: Nefer enters the Slither state, consuming Stamina to move rapidly forward for up to 2.5s. When the skill button is released, the duration ends, or Stamina runs out, Nefer will exit the Slither state and consume a certain amount of additional Stamina to deal Dendro DMG to opponents. When in the Shadow Dance state, additional Stamina consumption is decreased.
- Additional note: unleashing Senet Strategy: Dance of a Thousand Nights or sprinting while in Slither does not cause exit from Slither.
- Plunging Attack: Nefer plunges toward the ground, damaging opponents in her path and dealing AoE Dendro DMG on impact.

Multipliers (Talent Lv. 1-15):

- 1-Hit DMG: 38.07%, 40.93%, 43.78%, 47.59%, 50.44%, 53.3%, 57.11%, 60.91%, 64.72%, 68.53%, 72.34%, 76.14%, 80.9%, 85.66%, 90.42%
- 2-Hit DMG: 37.56%, 40.38%, 43.2%, 46.96%, 49.77%, 52.59%, 56.35%, 60.1%, 63.86%, 67.62%, 71.37%, 75.13%, 79.82%, 84.52%, 89.21%
- 3-Hit DMG: 25.24%x2, 27.13%x2, 29.03%x2, 31.55%x2, 33.44%x2, 35.34%x2, 37.86%x2, 40.38%x2, 42.91%x2, 45.43%x2, 47.96%x2, 50.48%x2, 53.64%x2, 56.79%x2, 59.94%x2
- 4-Hit DMG: 60.99%, 65.57%, 70.14%, 76.24%, 80.82%, 85.39%, 91.49%, 97.59%, 103.69%, 109.79%, 115.89%, 121.99%, 129.61%, 137.24%, 144.86%
- Charged Attack Charging Stamina Drain: 18.15/s at all levels
- Charged Attack DMG: 130.88%, 140.7%, 150.51%, 163.6%, 173.42%, 183.23%, 196.32%, 209.41%, 222.5%, 235.58%, 248.67%, 261.76%, 278.12%, 294.48%, 310.84%
- Charged Attack Stamina Cost: 50 at all levels
- Shadow Dance Charged Attack Stamina Cost: 25 at all levels
- Plunge DMG: 56.83%, 61.45%, 66.08%, 72.69%, 77.31%, 82.6%, 89.87%, 97.14%, 104.41%, 112.34%, 120.27%, 128.2%, 136.12%, 144.05%, 151.98%
- Low/High Plunge DMG: 113.63%/141.93%, 122.88%/153.49%, 132.13%/165.04%, 145.35%/181.54%, 154.59%/193.1%, 165.16%/206.3%, 179.7%/224.45%, 194.23%/242.61%, 208.77%/260.76%, 224.62%/280.57%, 240.48%/300.37%, 256.34%/320.18%, 272.19%/339.98%, 288.05%/359.79%, 303.9%/379.59%

### Elemental Skill: Senet Strategy: Dance of a Thousand Nights

Description:

- Nefer charges forward, deals AoE Dendro DMG, and enters Shadow Dance.
- While in Shadow Dance, if Nefer has at least 1 Verdant Dew, Charged Attacks are replaced with Phantasm Performance and no longer consume stamina.
- Shadow Dance increases interruption resistance.
- Two initial charges.

Multipliers (Talent Lv. 1-15):

- Skill DMG: 76.38% ATK+152.77% Elemental Mastery, 82.11% ATK+164.23% Elemental Mastery, 87.84% ATK+175.68% Elemental Mastery, 95.48% ATK+190.96% Elemental Mastery, 101.21% ATK+202.42% Elemental Mastery, 106.94% ATK+213.88% Elemental Mastery, 114.58% ATK+229.15% Elemental Mastery, 122.21% ATK+244.43% Elemental Mastery, 129.85% ATK+259.71% Elemental Mastery, 137.49% ATK+274.98% Elemental Mastery, 145.13% ATK+290.26% Elemental Mastery, 152.77% ATK+305.54% Elemental Mastery, 162.32% ATK+324.63% Elemental Mastery, 171.86% ATK+343.73% Elemental Mastery, 181.41% ATK+362.82% Elemental Mastery
- Phantasm Performance 1-Hit DMG (Nefer): 24.64% ATK+49.28% Elemental Mastery, 26.49% ATK+52.98% Elemental Mastery, 28.34% ATK+56.67% Elemental Mastery, 30.8% ATK+61.6% Elemental Mastery, 32.65% ATK+65.3% Elemental Mastery, 34.5% ATK+68.99% Elemental Mastery, 36.96% ATK+73.92% Elemental Mastery, 39.42% ATK+78.85% Elemental Mastery, 41.89% ATK+83.78% Elemental Mastery, 44.35% ATK+88.7% Elemental Mastery, 46.82% ATK+93.63% Elemental Mastery, 49.28% ATK+98.56% Elemental Mastery, 52.36% ATK+104.72% Elemental Mastery, 55.44% ATK+110.88% Elemental Mastery, 58.52% ATK+117.04% Elemental Mastery
- Phantasm Performance 1-Hit DMG (Shades): 96% Elemental Mastery, 103.2% Elemental Mastery, 110.4% Elemental Mastery, 120% Elemental Mastery, 127.2% Elemental Mastery, 134.4% Elemental Mastery, 144% Elemental Mastery, 153.6% Elemental Mastery, 163.2% Elemental Mastery, 172.8% Elemental Mastery, 182.4% Elemental Mastery, 192% Elemental Mastery, 204% Elemental Mastery, 216% Elemental Mastery, 228% Elemental Mastery
- Phantasm Performance 2-Hit DMG (Nefer): 32.03% ATK+64.06% Elemental Mastery, 34.43% ATK+68.87% Elemental Mastery, 36.84% ATK+73.67% Elemental Mastery, 40.04% ATK+80.08% Elemental Mastery, 42.44% ATK+84.88% Elemental Mastery, 44.84% ATK+89.69% Elemental Mastery, 48.05% ATK+96.1% Elemental Mastery, 51.25% ATK+102.5% Elemental Mastery, 54.45% ATK+108.91% Elemental Mastery, 57.66% ATK+115.32% Elemental Mastery, 60.86% ATK+121.72% Elemental Mastery, 64.06% ATK+128.13% Elemental Mastery, 68.07% ATK+136.14% Elemental Mastery, 72.07% ATK+144.14% Elemental Mastery, 76.08% ATK+152.15% Elemental Mastery
- Phantasm Performance 2-Hit DMG (Shades): 96% Elemental Mastery, 103.2% Elemental Mastery, 110.4% Elemental Mastery, 120% Elemental Mastery, 127.2% Elemental Mastery, 134.4% Elemental Mastery, 144% Elemental Mastery, 153.6% Elemental Mastery, 163.2% Elemental Mastery, 172.8% Elemental Mastery, 182.4% Elemental Mastery, 192% Elemental Mastery, 204% Elemental Mastery, 216% Elemental Mastery, 228% Elemental Mastery
- Phantasm Performance 3-Hit DMG (Shades): 128% Elemental Mastery, 137.6% Elemental Mastery, 147.2% Elemental Mastery, 160% Elemental Mastery, 169.6% Elemental Mastery, 179.2% Elemental Mastery, 192% Elemental Mastery, 204.8% Elemental Mastery, 217.6% Elemental Mastery, 230.4% Elemental Mastery, 243.2% Elemental Mastery, 256% Elemental Mastery, 272% Elemental Mastery, 288% Elemental Mastery, 304% Elemental Mastery
- Phantasm Performance Charges: 3 at all levels
- Shadow Dance Duration: 9s at all levels
- CD: 9s at all levels

### Elemental Burst: Sacred Vow: True Eye's Phantasm

Description:

- Deals AoE Dendro DMG to opponents ahead.
- Consumes all Veils of Falsehood on cast to increase current Burst damage.

Multipliers (Talent Lv. 1-15):

- 1-Hit DMG: 224.64% ATK+449.28% Elemental Mastery, 241.49% ATK+482.98% Elemental Mastery, 258.34% ATK+516.67% Elemental Mastery, 280.8% ATK+561.6% Elemental Mastery, 297.65% ATK+595.3% Elemental Mastery, 314.5% ATK+628.99% Elemental Mastery, 336.96% ATK+673.92% Elemental Mastery, 359.42% ATK+718.85% Elemental Mastery, 381.89% ATK+763.78% Elemental Mastery, 404.35% ATK+808.7% Elemental Mastery, 426.82% ATK+853.63% Elemental Mastery, 449.28% ATK+898.56% Elemental Mastery, 477.36% ATK+954.72% Elemental Mastery, 505.44% ATK+1010.88% Elemental Mastery, 533.52% ATK+1067.04% Elemental Mastery
- 2-Hit DMG: 336.96% ATK+673.92% Elemental Mastery, 362.23% ATK+724.46% Elemental Mastery, 387.5% ATK+775.01% Elemental Mastery, 421.2% ATK+842.4% Elemental Mastery, 446.47% ATK+892.94% Elemental Mastery, 471.74% ATK+943.49% Elemental Mastery, 505.44% ATK+1010.88% Elemental Mastery, 539.14% ATK+1078.27% Elemental Mastery, 572.83% ATK+1145.66% Elemental Mastery, 606.53% ATK+1213.06% Elemental Mastery, 640.22% ATK+1280.45% Elemental Mastery, 673.92% ATK+1347.84% Elemental Mastery, 716.04% ATK+1432.08% Elemental Mastery, 758.16% ATK+1516.32% Elemental Mastery, 800.28% ATK+1600.56% Elemental Mastery
- DMG Bonus per Veil of Falsehood stack: 13%, 16%, 19%, 22%, 25%, 28%, 31%, 34%, 37%, 40%, 43%, 46%, 49%, 52%, 55%
- CD: 15s at all levels
- Energy Cost: 60 at all levels

## Passives

### P1: A Wager of Moonlight

- On Skill, existing Dendro Cores become Seeds of Deceit.
- For 15s, Lunar-Bloom results that would create Dendro Cores or Bountiful Cores instead create Seeds of Deceit.
- Seeds of Deceit cannot trigger Hyperbloom or Burgeon and do not burst.
- Charged Attack or Phantasm Performance can absorb Seeds of Deceit in range and gain 1 Veil of Falsehood stack per seed.
- At 3 stacks, or when the third stack is refreshed, Nefer gains +100 EM for 8s.

### P2: Daughter of the Dust and Sand

- While in Shadow Dance, for 5s after a party member triggers Lunar-Bloom, Slither provides additional Verdant Dew.
- Every 100 EM above 500 increases this additional provision by 10%, up to 50%.

### P3: Moonsign Benediction: Dusklit Eaves

- Party-triggered Bloom becomes Lunar-Bloom.
- Every point of Nefer EM increases Lunar-Bloom base damage by 0.0175%, up to 14%.
- When Nefer is in the party, the party's Moonsign increases by 1 level.

### P4: Conspiracy of the Golden Vault

- Gains 25% more rewards on a 20-hour Nod-Krai Expedition.

## Constellations

### C1: Planning Breeds Success

- Lunar-Bloom base DMG caused by Phantasm Performance increases by 60% of Nefer's EM.
- This effect is also boosted by Veil of Falsehood.

### C2: Observation Feeds Strategy

- Extends Veil of Falsehood duration by 5s.
- Increases Veil stack cap to 5.
- Causes Phantasm Performance to deal up to 140% of its original DMG.
- Skill instantly grants 2 Veil stacks.
- At 5 stacks, or when the fifth stack is refreshed, Nefer gains +200 EM for 8s instead.

### C3: Deceit Cloaks the Truth

- Skill +3.

### C4: Delusion Ensnares Reason

- On-field Shadow Dance increases Verdant Dew gain rate by 25%.
- While in Shadow Dance, nearby opponents lose 20% Dendro RES.
- After leaving Shadow Dance or moving too far away, the effect is removed after 4.5s.

### C5: Opportunity Hides in the Margins

- Burst +3.

### C6: Victory Flows from the Turning of Tides

- During Phantasm Performance, Nefer's second-stage self hit is converted into AoE Dendro damage equal to 85% EM.
- After Phantasm Performance ends, an additional AoE Dendro damage instance equal to 120% EM is dealt.
- These damage instances count as Lunar-Bloom damage dealt by Phantasm Performance.
- Under Moonsign: Ascendant Gleam, Nefer's Lunar-Bloom damage is elevated by 15%.

## Ascension Materials

- Mora: 7,049,900
- Moonfall Silver x168
- Radiant Antler x46
- Tattered Warrant x18
- Immaculate Warrant x30
- Frost-Etched Warrant x36
- Nagadus Emerald Sliver x1
- Nagadus Emerald Fragment x9
- Nagadus Emerald Chunk x9
- Nagadus Emerald Gemstone x6
- Wanderer's Advice x2
- Hero's Wit x418

## Talent Materials

- Teachings of Elysium x9
- Guide to Elysium x63
- Philosophies of Elysium x114
- Tattered Warrant x18
- Immaculate Warrant x66
- Frost-Etched Warrant x93
- Ascended Sample: Rook x18
- Crown of Insight x3

## Implementation Notes For GCSim

- P3 explicitly converts Bloom into Lunar-Bloom and adds base damage scaling from Nefer EM: 0.0175% per EM, capped at 14%.
- P1 is not only a stat buff; it replaces Dendro Core creation results with a custom gadget or stateful object for a 15s window after Skill.
- Seeds of Deceit do not explode and cannot trigger Hyperbloom or Burgeon.
- Phantasm Performance is a special Charged Attack variant gated by both Shadow Dance and at least 1 Verdant Dew.
- Burst damage bonus is per Veil stack consumed on cast.
- C1 and C6 both modify Lunar-Bloom specifically from Phantasm Performance, so ownership and tagging for that route matter.
- The repository already supports special damage elevation through AttackInfo.Elevation, so the C6 line about Lunar-Bloom damage being elevated by 15% maps to existing infrastructure.
- C4 applies Dendro RES shred only while on field and in Shadow Dance, with a 4.5s linger after exit or range break.
