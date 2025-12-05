package modifier

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/stretchr/testify/assert"
)

func TestUniqueModifier(t *testing.T) {
	calledOnRemove := false
	mod := info.Modifier{
		Name:       keys.TestingMod,
		Duration:   10,
		Durability: 100,
		Stacking:   info.Unique,
		ModifierListeners: info.ModifierListeners{
			OnRemove: func(_ *info.Modifier) {
				calledOnRemove = true
			},
		},
	}
	h := &Handler{}

	// Add the modifier
	ok, err := h.Add(&mod)
	assert.NoError(t, err, "adding modifier should not error")
	assert.True(t, ok, "modifier should be added as new")

	assert.Equal(t, mod.Durability/info.Durability(mod.Duration), mod.DecayRate, "decay rate should be correctly calculated")

	// make sure can't add same modifier again
	mod2 := mod
	ok, err = h.Add(&mod2)
	assert.NoError(t, err, "adding duplicate unique modifier should not error")
	assert.False(t, ok, "duplicate unique modifier should not be added")

	for range 11 {
		h.Tick()
	}

	// Assert that OnRemove was called
	assert.True(t, calledOnRemove, "OnRemove should have been called after 11 ticks")
	assert.Equal(t, 0, len(h.modifiers), "No modifiers should remain after expiration")
}

func TestRefreshModifier(t *testing.T) {
	h := &Handler{}
	removeSrc := 0
	mod := info.Modifier{
		Name:       keys.TestingMod,
		Duration:   10,
		Durability: 100,
		Stacking:   info.Refresh,
		ModifierListeners: info.ModifierListeners{
			OnRemove: func(_ *info.Modifier) {
				removeSrc = 1
			},
		},
	}
	ok, err := h.Add(&mod)
	assert.NoError(t, err, "adding modifier should not error")
	assert.True(t, ok, "modifier should be added as new")

	assert.Equal(t, mod.Durability/info.Durability(mod.Duration), mod.DecayRate, "decay rate should be correctly calculated")

	// adding same modifier should replace all the properties
	mod2 := mod
	mod2.OnRemove = func(_ *info.Modifier) {
		removeSrc = 2
	}
	ok, err = h.Add(&mod2)
	assert.NoError(t, err, "adding duplicate refresh modifier should not error")
	assert.False(t, ok, "duplicate refresh modifier should not be new")

	for range 11 {
		h.Tick()
	}

	assert.Equal(t, 2, removeSrc, "OnRemove should be from the second modifier")
	assert.Equal(t, 0, len(h.modifiers), "No modifiers should remain after expiration")
}

func TestOverlapModifier(t *testing.T) {
	h := &Handler{}
	removeCallCount := 0

	// Create base modifier
	baseMod := info.Modifier{
		Name:       keys.TestingMod,
		Duration:   10,
		Durability: 100,
		Stacking:   info.Overlap,
		ModifierListeners: info.ModifierListeners{
			OnRemove: func(_ *info.Modifier) {
				removeCallCount++
			},
		},
	}

	// Add first modifier
	mod1 := baseMod
	ok, err := h.Add(&mod1)
	assert.NoError(t, err, "adding first modifier should not error")
	assert.True(t, ok, "first modifier should be added as new")
	assert.Equal(t, mod1.Durability/info.Durability(mod1.Duration), mod1.DecayRate, "decay rate should be correctly calculated")

	// Add second modifier (should stack)
	mod2 := baseMod
	ok, err = h.Add(&mod2)
	assert.NoError(t, err, "adding second overlap modifier should not error")
	assert.True(t, ok, "second overlap modifier should be added as new")

	// Add third modifier (should stack)
	mod3 := baseMod
	ok, err = h.Add(&mod3)
	assert.NoError(t, err, "adding third overlap modifier should not error")
	assert.True(t, ok, "third overlap modifier should be added as new")

	// Confirm all 3 modifiers exist
	assert.Equal(t, 3, len(h.modifiers), "Should have 3 stacked modifiers")

	// Tick down to expire all modifiers
	for range 11 {
		h.Tick()
	}

	// Assert that OnRemove was called for all 3 modifiers
	assert.Equal(t, 3, removeCallCount, "OnRemove should have been called 3 times for all stacked modifiers")
	assert.Equal(t, 0, len(h.modifiers), "No modifiers should remain after expiration")
}

func TestOverlapRefreshDurationModifier(t *testing.T) {
	h := &Handler{}
	removeCallCount := 0

	// Create base modifier
	baseMod := info.Modifier{
		Name:       keys.TestingMod,
		Duration:   10,
		Durability: 100,
		Stacking:   info.OverlapRefreshDuration,
		ModifierListeners: info.ModifierListeners{
			OnRemove: func(_ *info.Modifier) {
				removeCallCount++
			},
		},
	}
	expectedDecay := baseMod.Durability / info.Durability(baseMod.Duration)

	// Add first modifier
	mod1 := baseMod
	ok, err := h.Add(&mod1)
	assert.NoError(t, err, "adding first modifier should not error")
	assert.True(t, ok, "first modifier should be added as new")
	assert.Equal(t, expectedDecay, mod1.DecayRate, "decay rate should be correctly calculated")

	// Tick 3 times - durability should decrease by 30 (10 per tick)
	for range 3 {
		h.Tick()
	}
	assert.Equal(t, 1, len(h.modifiers), "Should have 1 modifier")
	assert.Equal(t, baseMod.Durability-3*expectedDecay, h.modifiers[0].Durability, "Durability should be 70 after 3 ticks")

	// Add second modifier (should refresh durability of existing)
	mod2 := baseMod
	ok, err = h.Add(&mod2)
	assert.NoError(t, err, "adding second overlap refresh modifier should not error")
	assert.True(t, ok, "second overlap refresh modifier should be considered new")
	assert.Equal(t, 2, len(h.modifiers), "Should have 2 modifier")
	for i, m := range h.modifiers {
		assert.Equal(t, baseMod.Durability, m.Durability, "Durability on mod %v should be refreshed to 100", i)
	}

	// Tick 3 times again
	for range 3 {
		h.Tick()
	}
	for i, m := range h.modifiers {
		assert.Equal(t, baseMod.Durability-3*expectedDecay, m.Durability, "Durability on mod %v should be 70 after 3 more ticks", i)
	}

	// Add third modifier (should refresh durability again)
	mod3 := baseMod
	ok, err = h.Add(&mod3)
	assert.NoError(t, err, "adding third overlap refresh modifier should not error")
	assert.True(t, ok, "third overlap refresh modifier should be considered new")
	assert.Equal(t, 3, len(h.modifiers), "Should have 3 modifier")
	for i, m := range h.modifiers {
		assert.Equal(t, baseMod.Durability, m.Durability, "Durability on mod %v should be refreshed to 100", i)
	}

	// Tick down to expire the modifier
	for range 11 {
		h.Tick()
	}

	// Assert that OnRemove was called only once (since there's only one modifier)
	assert.Equal(t, 3, removeCallCount, "OnRemove should have been called 3 times")
	assert.Equal(t, 0, len(h.modifiers), "No modifiers should remain after expiration")
}
