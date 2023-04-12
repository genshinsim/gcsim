package weapon

import (
	"fmt"
	"log"
	"math"
	"testing"

	"github.com/genshinsim/gcsim/pkg/model"
)

func TestParseWeapon(t *testing.T) {
	d, err := NewDataSource("../../../data/ExcelBinOutput/")
	if err != nil {
		t.Fatal(err)
	}

	id := int32(14405)
	w, err := d.parseWeapon(id)
	if err != nil {
		t.Fatal(err)
	}

	if w == nil {
		log.Fatal("result from parse cannot be nil")
	}

	expect(t, "id", id, w.Id)
	expect(t, "quality", int32(4), w.Rarity)
	expect(t, "weapon class", model.WeaponClass_WEAPON_CATALYST, w.WeaponClass)
	expect(t, "image name", "UI_EquipIcon_Catalyst_Resurrection", w.ImageName)

	//stats check
	if expect(t, "weapon props length", len(expectedSPWeaponProps), len(w.BaseStats.BaseProps)) == nil {
		for i, v := range expectedSPWeaponProps {
			got := w.BaseStats.BaseProps[i]
			expect(t, fmt.Sprintf("weapon prop idx %v, stat type", i), v.PropType, got.PropType)
			expect(t, fmt.Sprintf("weapon prop idx %v, curve type", i), v.Curve, got.Curve)
			expectTol(t, fmt.Sprintf("weapon prop idx %v, initial value", i), v.InitialValue, got.InitialValue, 0.00001)
		}
	}

	if expect(t, "promo data length", len(expectedSPPromoteCurve), len(w.BaseStats.PromoData)) == nil {
		for i, v := range expectedSPPromoteCurve {
			if expect(t, fmt.Sprintf("length for promo data idx %v", i), len(v.AddProps), len(w.BaseStats.PromoData[i].AddProps)) == nil {
				for j, x := range v.AddProps {
					got := w.BaseStats.PromoData[i].AddProps[j]
					expect(t, fmt.Sprintf("promo data idx %v, stat type (idx %v)", i, j), x.PropType, got.PropType)
					expectTol(t, fmt.Sprintf("promo data idx %v, value (idx %v)", i, j), x.Value, got.Value, 0.00001)
				}
			}
		}
	}

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

var expectedSPWeaponProps = []*model.WeaponProp{
	{
		PropType:     model.StatType_FIGHT_PROP_BASE_ATTACK,
		InitialValue: 42.4010009765625,
		Curve:        model.WeaponCurveType_GROW_CURVE_ATTACK_201,
	},
	{
		PropType:     model.StatType_FIGHT_PROP_CRITICAL,
		InitialValue: 0.05999999865889549,
		Curve:        model.WeaponCurveType_GROW_CURVE_CRITICAL_201,
	},
}

var expectedSPPromoteCurve = []*model.PromotionData{
	{
		MaxLevel: 20,
	},
	{
		MaxLevel: 40,
		AddProps: []*model.PromotionAddProp{
			{
				PropType: model.StatType_FIGHT_PROP_BASE_ATTACK,
				Value:    25.899999618530273,
			},
		},
	},
	{
		MaxLevel: 50,
		AddProps: []*model.PromotionAddProp{
			{
				PropType: model.StatType_FIGHT_PROP_BASE_ATTACK,
				Value:    51.900001525878906,
			},
		},
	},
	{
		MaxLevel: 60,
		AddProps: []*model.PromotionAddProp{
			{
				PropType: model.StatType_FIGHT_PROP_BASE_ATTACK,
				Value:    77.80000305175781,
			},
		},
	},
	{
		MaxLevel: 70,
		AddProps: []*model.PromotionAddProp{
			{
				PropType: model.StatType_FIGHT_PROP_BASE_ATTACK,
				Value:    103.69999694824219,
			},
		},
	},
	{
		MaxLevel: 80,
		AddProps: []*model.PromotionAddProp{
			{
				PropType: model.StatType_FIGHT_PROP_BASE_ATTACK,
				Value:    129.6999969482422,
			},
		},
	},
	{
		MaxLevel: 90,
		AddProps: []*model.PromotionAddProp{
			{
				PropType: model.StatType_FIGHT_PROP_BASE_ATTACK,
				Value:    155.60000610351562,
			},
		},
	},
}
