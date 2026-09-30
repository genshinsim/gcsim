package target

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
)

func TestNewIcdKeyDistinct(t *testing.T) {
	seen := make(map[IcdKey][3]int)
	for char := -1; char <= 4; char++ {
		for tag := range 256 {
			for grp := range 256 {
				k := NewIcdKey(char, attacks.ICDTag(tag), attacks.ICDGroup(grp))
				if prev, ok := seen[k]; ok {
					t.Fatalf("NewIcdKey(%v, %v, %v) collides with %v", char, tag, grp, prev)
				}
				seen[k] = [3]int{char, tag, grp}
			}
		}
	}
}
