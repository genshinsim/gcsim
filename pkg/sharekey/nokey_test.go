//go:build !sharekey

package sharekey

import "testing"

func TestNoEmbeddedKey(t *testing.T) {
	if Embedded() != nil {
		t.Fatal("Embedded() != nil without the sharekey tag")
	}
}
