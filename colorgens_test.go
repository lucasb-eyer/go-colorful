package colorful

import (
	"math/rand"
	"testing"
	"time"
)

// This is really difficult to test, if you've got a good idea, pull request!

// Check if it returns all valid colors.
func TestColorValidity(t *testing.T) {
	// with default seed
	for i := 0; i < 100; i++ {
		if col := WarmColor(); !col.IsValid() {
			t.Errorf("Warm color %v is not valid! Seed was: default", col)
		}

		if col := FastWarmColor(); !col.IsValid() {
			t.Errorf("Fast warm color %v is not valid! Seed was: default", col)
		}

		if col := HappyColor(); !col.IsValid() {
			t.Errorf("Happy color %v is not valid! Seed was: default", col)
		}

		if col := FastHappyColor(); !col.IsValid() {
			t.Errorf("Fast happy color %v is not valid! Seed was: default", col)
		}
	}

	// with custom seed
	seed := time.Now().UTC().UnixNano()
	rand := rand.New(rand.NewSource(seed))

	for i := 0; i < 100; i++ {
		if col := WarmColorWithRand(rand); !col.IsValid() {
			t.Errorf("Warm color %v is not valid! Seed was: %v", col, seed)
		}

		if col := FastWarmColorWithRand(rand); !col.IsValid() {
			t.Errorf("Fast warm color %v is not valid! Seed was: %v", col, seed)
		}

		if col := HappyColorWithRand(rand); !col.IsValid() {
			t.Errorf("Happy color %v is not valid! Seed was: %v", col, seed)
		}

		if col := FastHappyColorWithRand(rand); !col.IsValid() {
			t.Errorf("Fast happy color %v is not valid! Seed was: %v", col, seed)
		}
	}
}
