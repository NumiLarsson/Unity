package asteroids

import (
	"testing"
)

func TestNewAsteroidManager(t *testing.T) {

	manager := newAsteroidManager

	if manager == nil {
		t.Error("Manager not created")
	}
}
