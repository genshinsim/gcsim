package curves

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

var WeaponBaseMap = map[keys.Weapon]WeaponBase{
	keys.Akuoumaru: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      attributes.ATKP,
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
	keys.Alleyhunter: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      attributes.ATKP,
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
	keys.Amenomakageuchi: {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      attributes.ATKP,
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
	keys.Amosbow: {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.1080000028014183,
		Specialized:      attributes.ATKP,
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
	keys.Apprenticesnotes: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          23.2450008392334,
		BaseSpecialized:  0,
		Specialized:      attributes.NoStat,
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
	keys.Aquilafavonia: {
		AtkCurve:         GROW_CURVE_ATTACK_302,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          47.5369987487793,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      attributes.PhyP,
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
	keys.Beginnersprotector: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          23.2450008392334,
		BaseSpecialized:  0,
		Specialized:      attributes.NoStat,
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
	keys.Blackcliffagate: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      attributes.CD,
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
	keys.Blackclifflongsword: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.07999999821186066,
		Specialized:      attributes.CD,
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
	keys.Blackcliffpole: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      attributes.CD,
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
	keys.Blackcliffslasher: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      attributes.CD,
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
	keys.Blackcliffwarbow: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.07999999821186066,
		Specialized:      attributes.CD,
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
	keys.Blacktassel: {
		AtkCurve:         GROW_CURVE_ATTACK_104,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          37.60749816894531,
		BaseSpecialized:  0.10213299840688705,
		Specialized:      attributes.HPP,
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
	keys.Bloodtaintedgreatsword: {
		AtkCurve:         GROW_CURVE_ATTACK_104,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          37.60749816894531,
		BaseSpecialized:  40.79999923706055,
		Specialized:      attributes.EM,
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
	keys.Calamityqueller: {
		AtkCurve:         GROW_CURVE_ATTACK_303,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          49.137699127197266,
		BaseSpecialized:  0.035999998450279236,
		Specialized:      attributes.ATKP,
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
	keys.Cinnabarspindle: {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.15013299882411957,
		Specialized:      attributes.DEFP,
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
	keys.Compoundbow: {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.15013299882411957,
		Specialized:      attributes.PhyP,
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
	keys.Coolsteel: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.07660000026226044,
		Specialized:      attributes.ATKP,
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
	keys.Crescentpike: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.07506699860095978,
		Specialized:      attributes.PhyP,
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
	keys.Darkironsword: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  30.600000381469727,
		Specialized:      attributes.EM,
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
	keys.Deathmatch: {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.07999999821186066,
		Specialized:      attributes.CR,
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
	keys.Debateclub: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.07660000026226044,
		Specialized:      attributes.ATKP,
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
	keys.Dodocotales: {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      attributes.ATKP,
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
	keys.Dragonsbane: {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  48,
		Specialized:      attributes.EM,
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
	keys.Dragonspinespear: {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.15013299882411957,
		Specialized:      attributes.PhyP,
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
	keys.Dullblade: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          23.2450008392334,
		BaseSpecialized:  0,
		Specialized:      attributes.NoStat,
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
	keys.Elegyfortheend: {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      attributes.ER,
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
	keys.Emeraldorb: {
		AtkCurve:         GROW_CURVE_ATTACK_102,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          39.875099182128906,
		BaseSpecialized:  20.399999618530273,
		Specialized:      attributes.EM,
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
	keys.Engulfinglightning: {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      attributes.ER,
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
	keys.Everlastingmoonglow: {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.1080000028014183,
		Specialized:      attributes.HPP,
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
	keys.Eyeofperception: {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      attributes.ATKP,
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
	keys.Favoniuscodex: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.10000000149011612,
		Specialized:      attributes.ER,
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
	keys.Favoniusgreatsword: {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.1333329975605011,
		Specialized:      attributes.ER,
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
	keys.Favoniuslance: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.06666699796915054,
		Specialized:      attributes.ER,
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
	keys.Favoniussword: {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.1333329975605011,
		Specialized:      attributes.ER,
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
	keys.Favoniuswarbow: {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.1333329975605011,
		Specialized:      attributes.ER,
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
	keys.Ferrousshadow: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.07660000026226044,
		Specialized:      attributes.HPP,
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
	keys.Festeringdesire: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.10000000149011612,
		Specialized:      attributes.ER,
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
	keys.Filletblade: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.07660000026226044,
		Specialized:      attributes.ATKP,
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
	keys.Freedomsworn: {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  43.20000076293945,
		Specialized:      attributes.EM,
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
	keys.Frostbearer: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      attributes.ATKP,
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
	keys.Hakushinring: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.06666699796915054,
		Specialized:      attributes.ER,
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
	keys.Halberd: {
		AtkCurve:         GROW_CURVE_ATTACK_102,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          39.875099182128906,
		BaseSpecialized:  0.05106699839234352,
		Specialized:      attributes.ATKP,
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
	keys.Hamayumi: {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      attributes.ATKP,
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
	keys.Harangeppakufutsu: {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.07199999690055847,
		Specialized:      attributes.CR,
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
	keys.Harbingerofdawn: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.10199999809265137,
		Specialized:      attributes.CD,
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
	keys.Huntersbow: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          23.2450008392334,
		BaseSpecialized:  0,
		Specialized:      attributes.NoStat,
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
	keys.Ironpoint: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          32.93000030517578,
		BaseSpecialized:  0,
		Specialized:      attributes.NoStat,
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
	keys.Ironsting: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  36,
		Specialized:      attributes.EM,
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
	keys.Kagurasverity: {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.14399999380111694,
		Specialized:      attributes.CD,
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
	keys.Katsuragikirinagamasa: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.10000000149011612,
		Specialized:      attributes.ER,
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
	keys.Kitaincrossspear: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  24,
		Specialized:      attributes.EM,
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
	keys.Lionsroar: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      attributes.ATKP,
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
	keys.Lithicblade: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      attributes.ATKP,
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
	keys.Lithicspear: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      attributes.ATKP,
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
	keys.Lostprayertothesacredwinds: {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.07199999690055847,
		Specialized:      attributes.CR,
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
	keys.Luxurioussealord: {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      attributes.ATKP,
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
	keys.Magicguide: {
		AtkCurve:         GROW_CURVE_ATTACK_104,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          37.60749816894531,
		BaseSpecialized:  40.79999923706055,
		Specialized:      attributes.EM,
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
	keys.Mappamare: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  24,
		Specialized:      attributes.EM,
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
	keys.Memoryofdust: {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.1080000028014183,
		Specialized:      attributes.ATKP,
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
	keys.Messenger: {
		AtkCurve:         GROW_CURVE_ATTACK_102,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          39.875099182128906,
		BaseSpecialized:  0.06800000369548798,
		Specialized:      attributes.CD,
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
	keys.Mistsplitterreforged: {
		AtkCurve:         GROW_CURVE_ATTACK_302,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          47.5369987487793,
		BaseSpecialized:  0.09600000083446503,
		Specialized:      attributes.CD,
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
	keys.Mitternachtswaltz: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.11259999871253967,
		Specialized:      attributes.PhyP,
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
	keys.Mouunsmoon: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      attributes.ATKP,
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
	keys.Oathsworneye: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      attributes.ATKP,
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
	keys.Oldmercspal: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          32.93000030517578,
		BaseSpecialized:  0,
		Specialized:      attributes.NoStat,
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
	keys.Otherworldlystory: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.08500000089406967,
		Specialized:      attributes.ER,
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
	keys.Pocketgrimoire: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          32.93000030517578,
		BaseSpecialized:  0,
		Specialized:      attributes.NoStat,
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
	keys.Polarstar: {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.07199999690055847,
		Specialized:      attributes.CR,
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
	keys.Predator: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      attributes.ATKP,
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
	keys.Primordialjadecutter: {
		AtkCurve:         GROW_CURVE_ATTACK_304,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          44.33580017089844,
		BaseSpecialized:  0.09600000083446503,
		Specialized:      attributes.CR,
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
	keys.Primordialjadewingedspear: {
		AtkCurve:         GROW_CURVE_ATTACK_302,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          47.5369987487793,
		BaseSpecialized:  0.04800000041723251,
		Specialized:      attributes.CR,
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
	keys.Prototypeamber: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      attributes.HPP,
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
	keys.Prototypearchaic: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      attributes.ATKP,
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
	keys.Prototypecrescent: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      attributes.ATKP,
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
	keys.Prototyperancour: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.07506699860095978,
		Specialized:      attributes.PhyP,
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
	keys.Prototypestarglitter: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.10000000149011612,
		Specialized:      attributes.ER,
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
	keys.Rainslasher: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  36,
		Specialized:      attributes.EM,
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
	keys.Ravenbow: {
		AtkCurve:         GROW_CURVE_ATTACK_102,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          39.875099182128906,
		BaseSpecialized:  20.399999618530273,
		Specialized:      attributes.EM,
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
	keys.Recurvebow: {
		AtkCurve:         GROW_CURVE_ATTACK_104,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          37.60749816894531,
		BaseSpecialized:  0.10213299840688705,
		Specialized:      attributes.HPP,
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
	keys.Redhornstonethresher: {
		AtkCurve:         GROW_CURVE_ATTACK_304,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          44.33580017089844,
		BaseSpecialized:  0.19200000166893005,
		Specialized:      attributes.CD,
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
	keys.Royalbow: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      attributes.ATKP,
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
	keys.Royalgreatsword: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      attributes.ATKP,
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
	keys.Royalgrimoire: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      attributes.ATKP,
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
	keys.Royallongsword: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      attributes.ATKP,
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
	keys.Royalspear: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      attributes.ATKP,
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
	keys.Rust: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      attributes.ATKP,
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
	keys.Sacrificialbow: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.06666699796915054,
		Specialized:      attributes.ER,
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
	keys.Sacrificialfragments: {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  48,
		Specialized:      attributes.EM,
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
	keys.Sacrificialgreatsword: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.06666699796915054,
		Specialized:      attributes.ER,
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
	keys.Sacrificialsword: {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.1333329975605011,
		Specialized:      attributes.ER,
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
	keys.Seasonedhuntersbow: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          32.93000030517578,
		BaseSpecialized:  0,
		Specialized:      attributes.NoStat,
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
	keys.Serpentspine: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      attributes.CR,
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
	keys.Sharpshootersoath: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.10199999809265137,
		Specialized:      attributes.CD,
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
	keys.Silversword: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          32.93000030517578,
		BaseSpecialized:  0,
		Specialized:      attributes.NoStat,
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
	keys.Skyridergreatsword: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.09560000151395798,
		Specialized:      attributes.PhyP,
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
	keys.Skyridersword: {
		AtkCurve:         GROW_CURVE_ATTACK_104,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          37.60749816894531,
		BaseSpecialized:  0.11333300173282623,
		Specialized:      attributes.ER,
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
	keys.Skywardatlas: {
		AtkCurve:         GROW_CURVE_ATTACK_302,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          47.5369987487793,
		BaseSpecialized:  0.07199999690055847,
		Specialized:      attributes.ATKP,
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
	keys.Skywardblade: {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      attributes.ER,
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
	keys.Skywardharp: {
		AtkCurve:         GROW_CURVE_ATTACK_302,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          47.5369987487793,
		BaseSpecialized:  0.04800000041723251,
		Specialized:      attributes.CR,
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
	keys.Skywardpride: {
		AtkCurve:         GROW_CURVE_ATTACK_302,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          47.5369987487793,
		BaseSpecialized:  0.07999999821186066,
		Specialized:      attributes.ER,
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
	keys.Skywardspine: {
		AtkCurve:         GROW_CURVE_ATTACK_302,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          47.5369987487793,
		BaseSpecialized:  0.07999999821186066,
		Specialized:      attributes.ER,
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
	keys.Slingshot: {
		AtkCurve:         GROW_CURVE_ATTACK_104,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          37.60749816894531,
		BaseSpecialized:  0.06800000369548798,
		Specialized:      attributes.CR,
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
	keys.Snowtombedstarsilver: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.07506699860095978,
		Specialized:      attributes.PhyP,
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
	keys.Solarpearl: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      attributes.CR,
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
	keys.Songofbrokenpines: {
		AtkCurve:         GROW_CURVE_ATTACK_303,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          49.137699127197266,
		BaseSpecialized:  0.04500000178813934,
		Specialized:      attributes.PhyP,
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
	keys.Staffofhoma: {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.14399999380111694,
		Specialized:      attributes.CD,
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
	keys.Summitshaper: {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.1080000028014183,
		Specialized:      attributes.ATKP,
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
	keys.Swordofdescension: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.07660000026226044,
		Specialized:      attributes.ATKP,
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
	keys.Thealleyflash: {
		AtkCurve:         GROW_CURVE_ATTACK_203,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          45.06869888305664,
		BaseSpecialized:  11.999987602233887,
		Specialized:      attributes.EM,
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
	keys.Thebell: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      attributes.HPP,
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
	keys.Theblacksword: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      attributes.CR,
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
	keys.Thecatch: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.10000000149011612,
		Specialized:      attributes.ER,
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
	keys.Theflute: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.09000000357627869,
		Specialized:      attributes.ATKP,
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
	keys.Thestringless: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  36,
		Specialized:      attributes.EM,
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
	keys.Theunforged: {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.1080000028014183,
		Specialized:      attributes.ATKP,
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
	keys.Theviridescenthunt: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.05999999865889549,
		Specialized:      attributes.CR,
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
	keys.Thewidsith: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.11999999731779099,
		Specialized:      attributes.CD,
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
	keys.Thrillingtalesofdragonslayers: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.07660000026226044,
		Specialized:      attributes.HPP,
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
	keys.Thunderingpulse: {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.14399999380111694,
		Specialized:      attributes.CD,
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
	keys.Travelershandysword: {
		AtkCurve:         GROW_CURVE_ATTACK_102,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          39.875099182128906,
		BaseSpecialized:  0.06373299658298492,
		Specialized:      attributes.DEFP,
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
	keys.Twinnephrite: {
		AtkCurve:         GROW_CURVE_ATTACK_102,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          39.875099182128906,
		BaseSpecialized:  0.03400000184774399,
		Specialized:      attributes.CR,
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
	keys.Vortexvanquisher: {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.1080000028014183,
		Specialized:      attributes.ATKP,
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
	keys.Wastergreatsword: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          23.2450008392334,
		BaseSpecialized:  0,
		Specialized:      attributes.NoStat,
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
	keys.Wavebreakersfin: {
		AtkCurve:         GROW_CURVE_ATTACK_203,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          45.06869888305664,
		BaseSpecialized:  0.029999999329447746,
		Specialized:      attributes.ATKP,
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
	keys.Whiteblind: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  0.11259999871253967,
		Specialized:      attributes.DEFP,
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
	keys.Whiteirongreatsword: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.09560000151395798,
		Specialized:      attributes.DEFP,
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
	keys.Whitetassel: {
		AtkCurve:         GROW_CURVE_ATTACK_101,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
		BaseAtk:          38.74129867553711,
		BaseSpecialized:  0.050999999046325684,
		Specialized:      attributes.CR,
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
	keys.Windblumeode: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
		BaseSpecialized:  36,
		Specialized:      attributes.EM,
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
	keys.Wineandsong: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.06666699796915054,
		Specialized:      attributes.ER,
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
	keys.Wolfsgravestone: {
		AtkCurve:         GROW_CURVE_ATTACK_301,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          45.9364013671875,
		BaseSpecialized:  0.1080000028014183,
		Specialized:      attributes.ATKP,
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
