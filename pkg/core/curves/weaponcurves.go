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
	keys.AlleyHunter: {
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
	keys.AmenomaKageuchi: {
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
	keys.AmosBow: {
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
	keys.ApprenticesNotes: {
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
	keys.AquaSimulacra: {
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
	keys.AquilaFavonia: {
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
	keys.AThousandFloatingDreams: {
		AtkCurve:         GROW_CURVE_ATTACK_304,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          44.33580017089844,
		BaseSpecialized:  57.599998474121094,
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
	keys.BalladOfTheBoundlessBlue: {
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
	keys.BalladOfTheFjords: {
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
	keys.BeaconOfTheReedSea: {
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
	keys.BeginnersProtector: {
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
	keys.BlackcliffAgate: {
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
	keys.BlackcliffLongsword: {
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
	keys.BlackcliffPole: {
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
	keys.BlackcliffSlasher: {
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
	keys.BlackcliffWarbow: {
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
	keys.BlackTassel: {
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
	keys.BloodtaintedGreatsword: {
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
	keys.CalamityQueller: {
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
	keys.CashflowSupervision: {
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
	keys.CinnabarSpindle: {
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
	keys.CompoundBow: {
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
	keys.CoolSteel: {
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
	keys.CrescentPike: {
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
	keys.DarkIronSword: {
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
	keys.DebateClub: {
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
	keys.DodocoTales: {
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
	keys.DragonsBane: {
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
	keys.DragonspineSpear: {
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
	keys.DullBlade: {
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
	keys.ElegyForTheEnd: {
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
	keys.EmeraldOrb: {
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
	keys.EndOfTheLine: {
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
	keys.EngulfingLightning: {
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
	keys.EverlastingMoonglow: {
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
	keys.EyeOfPerception: {
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
	keys.FadingTwilight: {
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
	keys.FavoniusCodex: {
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
	keys.FavoniusGreatsword: {
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
	keys.FavoniusLance: {
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
	keys.FavoniusSword: {
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
	keys.FavoniusWarbow: {
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
	keys.FerrousShadow: {
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
	keys.FesteringDesire: {
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
	keys.FilletBlade: {
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
	keys.FinaleOfTheDeep: {
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
	keys.FleuveCendreFerryman: {
		AtkCurve:         GROW_CURVE_ATTACK_201,
		SpecializedCurve: GROW_CURVE_CRITICAL_101,
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
	keys.FlowingPurity: {
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
	keys.ForestRegalia: {
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
	keys.FreedomSworn: {
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
	keys.FruitOfFulfillment: {
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
	keys.HakushinRing: {
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
	keys.HaranGeppakuFutsu: {
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
	keys.HarbingerOfDawn: {
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
	keys.HuntersBow: {
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
	keys.HuntersPath: {
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
	keys.IbisPiercer: {
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
	keys.IronPoint: {
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
	keys.IronSting: {
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
	keys.JadefallsSplendor: {
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
	keys.KagotsurubeIsshin: {
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
	keys.KagurasVerity: {
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
	keys.KatsuragikiriNagamasa: {
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
	keys.KeyOfKhajNisut: {
		AtkCurve:         GROW_CURVE_ATTACK_304,
		SpecializedCurve: GROW_CURVE_CRITICAL_301,
		BaseAtk:          44.33580017089844,
		BaseSpecialized:  0.14399999380111694,
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
	keys.KingsSquire: {
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
	keys.KitainCrossSpear: {
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
	keys.LightOfFoliarIncision: {
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
	keys.LionsRoar: {
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
	keys.LithicBlade: {
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
	keys.LithicSpear: {
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
	keys.LostPrayerToTheSacredWinds: {
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
	keys.LuxuriousSeaLord: {
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
	keys.MagicGuide: {
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
	keys.MakhairaAquamarine: {
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
	keys.MappaMare: {
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
	keys.MailedFlower: {
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
	keys.MemoryOfDust: {
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
	keys.MissiveWindspear: {
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
	keys.MistsplitterReforged: {
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
	keys.MitternachtsWaltz: {
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
	keys.Moonpiercer: {
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
	keys.MouunsMoon: {
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
	keys.OathswornEye: {
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
	keys.OldMercsPal: {
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
	keys.OtherworldlyStory: {
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
	keys.PocketGrimoire: {
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
	keys.PolarStar: {
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
	keys.PortablePowerSaw: {
		AtkCurve:         GROW_CURVE_ATTACK_204,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          41.067100524902344,
		BaseSpecialized:  0.11999999731779099,
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
	keys.PrimordialJadeCutter: {
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
	keys.PrimordialJadeWingedSpear: {
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
	keys.ProspectorsDrill: {
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
	keys.PrototypeAmber: {
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
	keys.PrototypeArchaic: {
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
	keys.PrototypeCrescent: {
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
	keys.PrototypeRancour: {
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
	keys.PrototypeStarglitter: {
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
	keys.RangeGauge: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          42.4010009765625,
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
	keys.RavenBow: {
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
	keys.RecurveBow: {
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
	keys.RedhornStonethresher: {
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
	keys.RightfulReward: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.05999999865889549,
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
	keys.RoyalBow: {
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
	keys.RoyalGreatsword: {
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
	keys.RoyalGrimoire: {
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
	keys.RoyalLongsword: {
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
	keys.RoyalSpear: {
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
	keys.SacrificialBow: {
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
	keys.SacrificialFragments: {
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
	keys.SacrificialGreatsword: {
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
	keys.SacrificialJade: {
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
	keys.SacrificialSword: {
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
	keys.SapwoodBlade: {
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
	keys.ScionOfTheBlazingSun: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.03999999910593033,
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
	keys.SeasonedHuntersBow: {
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
	keys.SerpentSpine: {
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
	keys.SharpshootersOath: {
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
	keys.SilverSword: {
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
	keys.SkyriderGreatsword: {
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
	keys.SkyriderSword: {
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
	keys.SkywardAtlas: {
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
	keys.SkywardBlade: {
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
	keys.SkywardHarp: {
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
	keys.SkywardPride: {
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
	keys.SkywardSpine: {
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
	keys.SnowTombedStarsilver: {
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
	keys.SolarPearl: {
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
	keys.SongOfBrokenPines: {
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
	keys.SongOfStillness: {
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
	keys.SplendorOfTranquilWaters: {
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
	keys.StaffOfHoma: {
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
	keys.StaffOfTheScarletSands: {
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
	keys.SummitShaper: {
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
	keys.SwordOfDescension: {
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
	keys.TalkingStick: {
		AtkCurve:         GROW_CURVE_ATTACK_202,
		SpecializedCurve: GROW_CURVE_CRITICAL_201,
		BaseAtk:          43.734901428222656,
		BaseSpecialized:  0.03999999910593033,
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
	keys.TheAlleyFlash: {
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
	keys.TheBell: {
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
	keys.TheBlackSword: {
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
	keys.TheCatch: {
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
	keys.TheDockhandsAssistant: {
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
	keys.TheFirstGreatMagic: {
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
	keys.TheFlute: {
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
	keys.TheStringless: {
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
	keys.TheUnforged: {
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
	keys.TheViridescentHunt: {
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
	keys.TheWidsith: {
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
	keys.ThrillingTalesOfDragonSlayers: {
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
	keys.ThunderingPulse: {
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
	keys.TidalShadow: {
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
	keys.TomeOfTheEternalFlow: {
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
	keys.ToukabouShigure: {
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
	keys.TravelersHandySword: {
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
	keys.TulaytullahsRemembrance: {
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
	keys.TwinNephrite: {
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
	keys.VortexVanquisher: {
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
	keys.WanderingEvenstar: {
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
	keys.WasterGreatsword: {
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
	keys.WavebreakersFin: {
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
	keys.WhiteIronGreatsword: {
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
	keys.WhiteTassel: {
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
	keys.WindblumeOde: {
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
	keys.WineAndSong: {
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
	keys.WolfFang: {
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
	keys.WolfsGravestone: {
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
	keys.XiphosMoonlight: {
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
}
