package curves

import (
	"github.com/genshinsim/gcsim/pkg/core"
)

var WeaponBaseMap = map[string]WeaponBase{
	"akuoumaru": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"alleyhunter": {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"amenomakageuchi": {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"amosbow": {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.1080000028014183,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"apprenticesnotes": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          23.2450008392334,
		BaseSpecialized:  0,
		Specialized:      core.NoStat,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      11.699999809265137,
			},
			{
				MaxLevel: 50,
				Atk:      23.299999237060547,
			},
			{
				MaxLevel: 60,
				Atk:      35,
			},
			{
				MaxLevel: 70,
				Atk:      46.70000076293945,
			},
		},
	},
	"aquasimulacra": {
		AtkCurve:         GROW_CURVE_ATTACK_304,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          44.33580017089844,
		BaseSpecialized:  0.19200000166893005,
		Specialized:      core.CD,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"aquilafavonia": {
		AtkCurve:         GROW_CURVE_ATTACK_302,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          47.5369987487793,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      core.PhyP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"beginnersprotector": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          23.2450008392334,
		BaseSpecialized:  0,
		Specialized:      core.NoStat,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      11.699999809265137,
			},
			{
				MaxLevel: 50,
				Atk:      23.299999237060547,
			},
			{
				MaxLevel: 60,
				Atk:      35,
			},
			{
				MaxLevel: 70,
				Atk:      46.70000076293945,
			},
		},
	},
	"blackcliffagate": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      core.CD,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"blackclifflongsword": {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.07999999821186066,
		Specialized:      core.CD,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"blackcliffpole": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      core.CD,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"blackcliffslasher": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      core.CD,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"blackcliffwarbow": {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.07999999821186066,
		Specialized:      core.CD,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"blacktassel": {
		AtkCurve:         GROW_CURVE_ATTACK_104,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          37.60749816894531,
		BaseSpecialized:  0.10213299840688705,
		Specialized:      core.HPP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"bloodtaintedgreatsword": {
		AtkCurve:         GROW_CURVE_ATTACK_104,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          37.60749816894531,
		BaseSpecialized:  40.79999923706055,
		Specialized:      core.EM,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"calamityqueller": {
		AtkCurve:         GROW_CURVE_ATTACK_303,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          49.137699127197266,
		BaseSpecialized:  0.035999998450279236,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"cinnabarspindle": {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.15013299882411957,
		Specialized:      core.DEFP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"compoundbow": {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.15013299882411957,
		Specialized:      core.PhyP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"coolsteel": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.07660000026226044,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"crescentpike": {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.07506699860095978,
		Specialized:      core.PhyP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"darkironsword": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  30.600000381469727,
		Specialized:      core.EM,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"deathmatch": {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.07999999821186066,
		Specialized:      core.CR,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"debateclub": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.07660000026226044,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"dodocotales": {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"dragonsbane": {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  48,
		Specialized:      core.EM,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"dragonspinespear": {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.15013299882411957,
		Specialized:      core.PhyP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"dullblade": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          23.2450008392334,
		BaseSpecialized:  0,
		Specialized:      core.NoStat,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      11.699999809265137,
			},
			{
				MaxLevel: 50,
				Atk:      23.299999237060547,
			},
			{
				MaxLevel: 60,
				Atk:      35,
			},
			{
				MaxLevel: 70,
				Atk:      46.70000076293945,
			},
		},
	},
	"elegyfortheend": {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"emeraldorb": {
		AtkCurve:         GROW_CURVE_ATTACK_102,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          39.875099182128906,
		BaseSpecialized:  20.399999618530273,
		Specialized:      core.EM,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"engulfinglightning": {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"everlastingmoonglow": {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.1080000028014183,
		Specialized:      core.HPP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"eyeofperception": {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"fadingtwilight": {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.06666699796915054,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"favoniuscodex": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.10000000149011612,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"favoniusgreatsword": {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.1333329975605011,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"favoniuslance": {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.06666699796915054,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"favoniussword": {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.1333329975605011,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"favoniuswarbow": {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.1333329975605011,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"ferrousshadow": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.07660000026226044,
		Specialized:      core.HPP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"festeringdesire": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.10000000149011612,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"filletblade": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.07660000026226044,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"freedomsworn": {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  43.20000076293945,
		Specialized:      core.EM,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"frostbearer": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"hakushinring": {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.06666699796915054,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"halberd": {
		AtkCurve:         GROW_CURVE_ATTACK_102,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          39.875099182128906,
		BaseSpecialized:  0.05106699839234352,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"hamayumi": {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"harangeppakufutsu": {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.07199999690055847,
		Specialized:      core.CR,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"harbingerofdawn": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.10199999809265137,
		Specialized:      core.CD,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"huntersbow": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          23.2450008392334,
		BaseSpecialized:  0,
		Specialized:      core.NoStat,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      11.699999809265137,
			},
			{
				MaxLevel: 50,
				Atk:      23.299999237060547,
			},
			{
				MaxLevel: 60,
				Atk:      35,
			},
			{
				MaxLevel: 70,
				Atk:      46.70000076293945,
			},
		},
	},
	"ironpoint": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          32.93000030517578,
		BaseSpecialized:  0,
		Specialized:      core.NoStat,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      11.699999809265137,
			},
			{
				MaxLevel: 50,
				Atk:      23.299999237060547,
			},
			{
				MaxLevel: 60,
				Atk:      35,
			},
			{
				MaxLevel: 70,
				Atk:      46.70000076293945,
			},
		},
	},
	"ironsting": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  36,
		Specialized:      core.EM,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"kagurasverity": {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.14399999380111694,
		Specialized:      core.CD,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"katsuragikirinagamasa": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.10000000149011612,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"kitaincrossspear": {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  24,
		Specialized:      core.EM,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"lionsroar": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"lithicblade": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"lithicspear": {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"lostprayertothesacredwinds": {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.07199999690055847,
		Specialized:      core.CR,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"luxurioussealord": {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"magicguide": {
		AtkCurve:         GROW_CURVE_ATTACK_104,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          37.60749816894531,
		BaseSpecialized:  40.79999923706055,
		Specialized:      core.EM,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"mappamare": {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  24,
		Specialized:      core.EM,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"memoryofdust": {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.1080000028014183,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"messenger": {
		AtkCurve:         GROW_CURVE_ATTACK_102,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          39.875099182128906,
		BaseSpecialized:  0.06800000369548798,
		Specialized:      core.CD,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"mistsplitterreforged": {
		AtkCurve:         GROW_CURVE_ATTACK_302,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          47.5369987487793,
		BaseSpecialized:  0.09600000083446503,
		Specialized:      core.CD,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"mitternachtswaltz": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.11259999871253967,
		Specialized:      core.PhyP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"mouunsmoon": {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"oathsworneye": {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"oldmercspal": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          32.93000030517578,
		BaseSpecialized:  0,
		Specialized:      core.NoStat,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      11.699999809265137,
			},
			{
				MaxLevel: 50,
				Atk:      23.299999237060547,
			},
			{
				MaxLevel: 60,
				Atk:      35,
			},
			{
				MaxLevel: 70,
				Atk:      46.70000076293945,
			},
		},
	},
	"otherworldlystory": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.08500000089406967,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"pocketgrimoire": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          32.93000030517578,
		BaseSpecialized:  0,
		Specialized:      core.NoStat,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      11.699999809265137,
			},
			{
				MaxLevel: 50,
				Atk:      23.299999237060547,
			},
			{
				MaxLevel: 60,
				Atk:      35,
			},
			{
				MaxLevel: 70,
				Atk:      46.70000076293945,
			},
		},
	},
	"polarstar": {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.07199999690055847,
		Specialized:      core.CR,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"predator": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"primordialjadecutter": {
		AtkCurve:         GROW_CURVE_ATTACK_304,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          44.33580017089844,
		BaseSpecialized:  0.09600000083446503,
		Specialized:      core.CR,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"primordialjadewingedspear": {
		AtkCurve:         GROW_CURVE_ATTACK_302,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          47.5369987487793,
		BaseSpecialized:  0.04800000041723251,
		Specialized:      core.CR,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"prototypeamber": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      core.HPP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"prototypearchaic": {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"prototypecrescent": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"prototyperancour": {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.07506699860095978,
		Specialized:      core.PhyP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"prototypestarglitter": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.10000000149011612,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"rainslasher": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  36,
		Specialized:      core.EM,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"ravenbow": {
		AtkCurve:         GROW_CURVE_ATTACK_102,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          39.875099182128906,
		BaseSpecialized:  20.399999618530273,
		Specialized:      core.EM,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"recurvebow": {
		AtkCurve:         GROW_CURVE_ATTACK_104,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          37.60749816894531,
		BaseSpecialized:  0.10213299840688705,
		Specialized:      core.HPP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"redhornstonethresher": {
		AtkCurve:         GROW_CURVE_ATTACK_304,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          44.33580017089844,
		BaseSpecialized:  0.19200000166893005,
		Specialized:      core.CD,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"royalbow": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"royalgreatsword": {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"royalgrimoire": {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"royallongsword": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"royalspear": {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"rust": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"sacrificialbow": {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.06666699796915054,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"sacrificialfragments": {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  48,
		Specialized:      core.EM,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"sacrificialgreatsword": {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.06666699796915054,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"sacrificialsword": {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.1333329975605011,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"seasonedhuntersbow": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          32.93000030517578,
		BaseSpecialized:  0,
		Specialized:      core.NoStat,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      11.699999809265137,
			},
			{
				MaxLevel: 50,
				Atk:      23.299999237060547,
			},
			{
				MaxLevel: 60,
				Atk:      35,
			},
			{
				MaxLevel: 70,
				Atk:      46.70000076293945,
			},
		},
	},
	"serpentspine": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      core.CR,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"sharpshootersoath": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.10199999809265137,
		Specialized:      core.CD,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"silversword": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          32.93000030517578,
		BaseSpecialized:  0,
		Specialized:      core.NoStat,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      11.699999809265137,
			},
			{
				MaxLevel: 50,
				Atk:      23.299999237060547,
			},
			{
				MaxLevel: 60,
				Atk:      35,
			},
			{
				MaxLevel: 70,
				Atk:      46.70000076293945,
			},
		},
	},
	"skyridergreatsword": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.09560000151395798,
		Specialized:      core.PhyP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"skyridersword": {
		AtkCurve:         GROW_CURVE_ATTACK_104,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          37.60749816894531,
		BaseSpecialized:  0.11333300173282623,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"skywardatlas": {
		AtkCurve:         GROW_CURVE_ATTACK_302,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          47.5369987487793,
		BaseSpecialized:  0.07199999690055847,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"skywardblade": {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"skywardharp": {
		AtkCurve:         GROW_CURVE_ATTACK_302,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          47.5369987487793,
		BaseSpecialized:  0.04800000041723251,
		Specialized:      core.CR,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"skywardpride": {
		AtkCurve:         GROW_CURVE_ATTACK_302,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          47.5369987487793,
		BaseSpecialized:  0.07999999821186066,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"skywardspine": {
		AtkCurve:         GROW_CURVE_ATTACK_302,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          47.5369987487793,
		BaseSpecialized:  0.07999999821186066,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"slingshot": {
		AtkCurve:         GROW_CURVE_ATTACK_104,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          37.60749816894531,
		BaseSpecialized:  0.06800000369548798,
		Specialized:      core.CR,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"snowtombedstarsilver": {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.07506699860095978,
		Specialized:      core.PhyP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"solarpearl": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      core.CR,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"songofbrokenpines": {
		AtkCurve:         GROW_CURVE_ATTACK_303,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          49.137699127197266,
		BaseSpecialized:  0.04500000178813934,
		Specialized:      core.PhyP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"staffofhoma": {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.14399999380111694,
		Specialized:      core.CD,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"summitshaper": {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.1080000028014183,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"swordofdescension": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.07660000026226044,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"thealleyflash": {
		AtkCurve:         GROW_CURVE_ATTACK_203,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          45.06869888305664,
		BaseSpecialized:  11.999987602233887,
		Specialized:      core.EM,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"thebell": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      core.HPP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"theblacksword": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      core.CR,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"thecatch": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.10000000149011612,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"theflute": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"thestringless": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  36,
		Specialized:      core.EM,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"theunforged": {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.1080000028014183,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"theviridescenthunt": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      core.CR,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"thewidsith": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      core.CD,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"thrillingtalesofdragonslayers": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.07660000026226044,
		Specialized:      core.HPP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"thunderingpulse": {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.14399999380111694,
		Specialized:      core.CD,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"travelershandysword": {
		AtkCurve:         GROW_CURVE_ATTACK_102,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          39.875099182128906,
		BaseSpecialized:  0.06373299658298492,
		Specialized:      core.DEFP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"twinnephrite": {
		AtkCurve:         GROW_CURVE_ATTACK_102,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          39.875099182128906,
		BaseSpecialized:  0.03400000184774399,
		Specialized:      core.CR,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"vortexvanquisher": {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.1080000028014183,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
	"wastergreatsword": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          23.2450008392334,
		BaseSpecialized:  0,
		Specialized:      core.NoStat,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      11.699999809265137,
			},
			{
				MaxLevel: 50,
				Atk:      23.299999237060547,
			},
			{
				MaxLevel: 60,
				Atk:      35,
			},
			{
				MaxLevel: 70,
				Atk:      46.70000076293945,
			},
		},
	},
	"wavebreakersfin": {
		AtkCurve:         GROW_CURVE_ATTACK_203,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          45.06869888305664,
		BaseSpecialized:  0.029999999329447746,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"whiteblind": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.11259999871253967,
		Specialized:      core.DEFP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"whiteirongreatsword": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.09560000151395798,
		Specialized:      core.DEFP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"whitetassel": {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.050999999046325684,
		Specialized:      core.CR,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      19.5,
			},
			{
				MaxLevel: 50,
				Atk:      38.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      58.400001525878906,
			},
			{
				MaxLevel: 70,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 80,
				Atk:      97.30000305175781,
			},
			{
				MaxLevel: 90,
				Atk:      116.69999694824219,
			},
		},
	},
	"windblumeode": {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  36,
		Specialized:      core.EM,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"wineandsong": {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.06666699796915054,
		Specialized:      core.ER,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      25.899999618530273,
			},
			{
				MaxLevel: 50,
				Atk:      51.900001525878906,
			},
			{
				MaxLevel: 60,
				Atk:      77.80000305175781,
			},
			{
				MaxLevel: 70,
				Atk:      103.69999694824219,
			},
			{
				MaxLevel: 80,
				Atk:      129.6999969482422,
			},
			{
				MaxLevel: 90,
				Atk:      155.60000610351562,
			},
		},
	},
	"wolfsgravestone": {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.1080000028014183,
		Specialized:      core.ATKP,
		PromotionBonus: []PromoData{
			{
				MaxLevel: 20,
				Atk:      0,
			},
			{
				MaxLevel: 40,
				Atk:      31.100000381469727,
			},
			{
				MaxLevel: 50,
				Atk:      62.20000076293945,
			},
			{
				MaxLevel: 60,
				Atk:      93.4000015258789,
			},
			{
				MaxLevel: 70,
				Atk:      124.5,
			},
			{
				MaxLevel: 80,
				Atk:      155.60000610351562,
			},
			{
				MaxLevel: 90,
				Atk:      186.6999969482422,
			},
		},
	},
}
