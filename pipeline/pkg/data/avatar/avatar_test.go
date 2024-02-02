package avatar

import (
	"fmt"
	"math"
	"testing"

	"github.com/genshinsim/gcsim/pkg/model"
)

func TestParseCharacter(t *testing.T) {
	a, err := NewDataSource("../../../data/ExcelBinOutput/")
	if err != nil {
		t.Fatal(err)
	}

	// needs to be typed so the comparison works
	var id int32 = 10000002 // ayaka

	d, err := a.parseChar(id, 0)
	if err != nil {
		t.Fatal(err)
	}

	if d == nil {
		t.Fatalf("result from parse cannot be nil")
	}

	expect(t, "id", id, d.Id)
	expect(t, "quality", model.QualityType_QUALITY_ORANGE, d.Rarity)
	expect(t, "body type", model.BodyType_BODY_GIRL, d.Body)
	expect(t, "region", model.ZoneType_ASSOC_TYPE_INAZUMA, d.Region)
	expect(t, "element", model.Element_Ice, d.Element)
	expect(t, "weapon class", model.WeaponClass_WEAPON_SWORD_ONE_HAND, d.WeaponClass)
	expect(t, "icon", string("UI_AvatarIcon_Ayaka"), d.IconName)
	expect(t, "burst id", int32(10019), d.GetSkillDetails().GetBurst()) // ayaka burst is 10019
	expect(t, "attack id", int32(10024), d.GetSkillDetails().GetAttack())
	expect(t, "skill id", int32(10018), d.GetSkillDetails().GetSkill())
	expectTol(t, "burst energy cost", float64(80), d.GetSkillDetails().GetBurstEnergyCost(), 0.000000001)

	// stat block
	expectTol(t, "base attack", 26.6266, d.GetStats().BaseAtk, 0.00001)
	expectTol(t, "base hp", 1000.98602, d.GetStats().BaseHp, 0.00001)
	expectTol(t, "base def", 61.02659, d.GetStats().BaseDef, 0.00001)
	expect(t, "atk curve", model.AvatarCurveType_GROW_CURVE_ATTACK_S5, d.GetStats().AtkCurve)
	expect(t, "hp curve", model.AvatarCurveType_GROW_CURVE_HP_S5, d.GetStats().HpCurve)
	expect(t, "def curve", model.AvatarCurveType_GROW_CURVE_HP_S5, d.GetStats().DefCruve)
	err = expect(t, "promo data length", len(expectedAyakaCurves), len(d.GetStats().PromoData))
	if err != nil {
		t.FailNow()
	}
	for i, v := range expectedAyakaCurves {
		err = expect(t, fmt.Sprintf("length for promo data idx %v", i), len(v.AddProps), len(d.Stats.PromoData[i].AddProps))
		if err != nil {
			t.FailNow()
		}
		for j, x := range v.AddProps {
			got := d.Stats.PromoData[i].AddProps[j]
			expect(t, fmt.Sprintf("promo data idx %v, stat type (idx %v)", i, j), x.PropType, got.PropType)
			expectTol(t, fmt.Sprintf("promo data idx %v, value (idx %v)", i, j), x.Value, got.Value, 0.00001)
		}
	}

	// make sure traveler is picking up correct skills
	id = 10000007 // lumine

	d, err = a.parseChar(id, 707)
	if err != nil {
		t.Fatal(err)
	}

	if d == nil {
		t.Fatalf("result from parse cannot be nil")
	}

	expect(t, "id", id, d.Id)
	expect(t, "quality", model.QualityType_QUALITY_ORANGE, d.Rarity)
	expect(t, "body type", model.BodyType_BODY_GIRL, d.Body)
	expect(t, "region", model.ZoneType_ASSOC_TYPE_MAINACTOR, d.Region)
	expect(t, "element", model.Element_Electric, d.Element)
	expect(t, "weapon class", model.WeaponClass_WEAPON_SWORD_ONE_HAND, d.WeaponClass)
	expect(t, "icon", string("UI_AvatarIcon_PlayerGirl"), d.IconName)
	expect(t, "burst id", int32(10605), d.GetSkillDetails().GetBurst())
	expect(t, "attack id", int32(100556), d.GetSkillDetails().GetAttack())
	expect(t, "skill id", int32(10602), d.GetSkillDetails().GetSkill())
	expectTol(t, "burst energy cost", float64(80), d.GetSkillDetails().GetBurstEnergyCost(), 0.000000001)
}

func expect(t *testing.T, msg string, expect, got any) error {
	if expect != got {
		err := fmt.Errorf("%v expecting %v, got %v", msg, expect, got)
		t.Error(err)
		return err
	}
	return nil
}

func expectTol(t *testing.T, msg string, expect, got, tol float64) error {
	if math.Abs(expect-got) > tol {
		err := fmt.Errorf("%v expecting %v (with tol %v), got %v", msg, expect, tol, got)
		t.Error(err)
		return err
	}
	return nil
}

var expectedAyakaCurves = []*model.PromotionData{
	{
		MaxLevel: 20,
		AddProps: []*model.PromotionAddProp{
			{
				PropType: model.StatType_FIGHT_PROP_BASE_HP,
			},
			{
				PropType: model.StatType_FIGHT_PROP_BASE_DEFENSE,
			},
			{
				PropType: model.StatType_FIGHT_PROP_BASE_ATTACK,
			},
			{
				PropType: model.StatType_FIGHT_PROP_CRITICAL_HURT,
			},
		},
	},
	{
		MaxLevel: 40,
		AddProps: []*model.PromotionAddProp{
			{
				PropType: model.StatType_FIGHT_PROP_BASE_HP,
				Value:    858.2550048828125,
			},
			{
				PropType: model.StatType_FIGHT_PROP_BASE_DEFENSE,
				Value:    52.32600021362305,
			},
			{
				PropType: model.StatType_FIGHT_PROP_BASE_ATTACK,
				Value:    22.82823371887207,
			},
			{
				PropType: model.StatType_FIGHT_PROP_CRITICAL_HURT,
			},
		},
	},
	{
		MaxLevel: 50,
		AddProps: []*model.PromotionAddProp{
			{
				PropType: model.StatType_FIGHT_PROP_BASE_HP,
				Value:    1468.0677490234375,
			},
			{
				PropType: model.StatType_FIGHT_PROP_BASE_DEFENSE,
				Value:    89.50499725341797,
			},
			{
				PropType: model.StatType_FIGHT_PROP_BASE_ATTACK,
				Value:    39.04829406738281,
			},
			{
				PropType: model.StatType_FIGHT_PROP_CRITICAL_HURT,
				Value:    0.09600000083446503,
			},
		},
	},
	{
		MaxLevel: 60,
		AddProps: []*model.PromotionAddProp{
			{
				PropType: model.StatType_FIGHT_PROP_BASE_HP,
				Value:    2281.1513671875,
			},
			{
				PropType: model.StatType_FIGHT_PROP_BASE_DEFENSE,
				Value:    139.07699584960938,
			},
			{
				PropType: model.StatType_FIGHT_PROP_BASE_ATTACK,
				Value:    60.67504119873047,
			},
			{
				PropType: model.StatType_FIGHT_PROP_CRITICAL_HURT,
				Value:    0.19200000166893005,
			},
		},
	},
	{
		MaxLevel: 70,
		AddProps: []*model.PromotionAddProp{
			{
				PropType: model.StatType_FIGHT_PROP_BASE_HP,
				Value:    2890.964111328125,
			},
			{
				PropType: model.StatType_FIGHT_PROP_BASE_DEFENSE,
				Value:    176.25599670410156,
			},
			{
				PropType: model.StatType_FIGHT_PROP_BASE_ATTACK,
				Value:    76.89510345458984,
			},
			{
				PropType: model.StatType_FIGHT_PROP_CRITICAL_HURT,
				Value:    0.19200000166893005,
			},
		},
	},
	{
		MaxLevel: 80,
		AddProps: []*model.PromotionAddProp{
			{
				PropType: model.StatType_FIGHT_PROP_BASE_HP,
				Value:    3500.77685546875,
			},
			{
				PropType: model.StatType_FIGHT_PROP_BASE_DEFENSE,
				Value:    213.43499755859375,
			},
			{
				PropType: model.StatType_FIGHT_PROP_BASE_ATTACK,
				Value:    93.11516571044922,
			},
			{
				PropType: model.StatType_FIGHT_PROP_CRITICAL_HURT,
				Value:    0.2879999876022339,
			},
		},
	},
	{
		MaxLevel: 90,
		AddProps: []*model.PromotionAddProp{
			{
				PropType: model.StatType_FIGHT_PROP_BASE_HP,
				Value:    4110.58984375,
			},
			{
				PropType: model.StatType_FIGHT_PROP_BASE_DEFENSE,
				Value:    250.61399841308594,
			},
			{
				PropType: model.StatType_FIGHT_PROP_BASE_ATTACK,
				Value:    109.3352279663086,
			},
			{
				PropType: model.StatType_FIGHT_PROP_CRITICAL_HURT,
				Value:    0.3840000033378601,
			},
		},
	},
}
