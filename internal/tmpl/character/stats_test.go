package character

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core"
)

func BenchmarkAddModHeap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		//try adding mods 50 times in heap then calling amount
		mods := make([]core.CharStatMod, 0, 50)
		m := make([]float64, core.EndStatType)
		for i := 0; i < 50; i++ {
			mods = append(mods, core.CharStatMod{
				Expiry: -1,
				Amount: func() ([]float64, bool) {
					//do some math here
					m[core.DmgP] = 5 * 5
					return m, true
				},
				Key: "test",
			})
		}
		//call amount function twice per mod
		for _, v := range mods {
			for i := 0; i < 100; i++ {
				v.Amount()
			}
		}
	}
}

func BenchmarkAddModStack(b *testing.B) {
	for n := 0; n < b.N; n++ {
		//try adding mods 50 times in heap then calling amount
		mods := make([]core.CharStatMod, 0, 50)
		for i := 0; i < 50; i++ {

			mods = append(mods, core.CharStatMod{
				Expiry: -1,
				Amount: func() ([]float64, bool) {
					m := make([]float64, core.EndStatType)
					//do some math here
					m[core.DmgP] = 5 * 5
					return m, true
				},
				Key: "test",
			})
		}
		//call amount function 100x per mod
		for _, v := range mods {
			for i := 0; i < 100; i++ {
				v.Amount()
			}
		}
	}
}
