package asteroids

import (
	"fmt"
	"testing"
)

func TestNewAsteroidManager(t *testing.T) {

	manager := newAsteroidManager

	if manager == nil {
		t.Error("Manager not created")
	}
}

func TestInitAsteroidManager(t *testing.T) {

	manager := newAsteroidManager()
	conn := MakeConnection()

	// Go-routine to be able to read confirmation
	go manager.init(conn.FlipConnection())

	fmt.Println("innan response")
	response := <-conn.read

	if manager.input != conn.write {
		t.Error("Channels incorrect")
	}

	// Testing hardcoded lines in init
	if response.action != "a.manager_ready" || response.result != 200 {
		t.Error("Response incorrect")
	}

	if manager.yMax != 100 {
		t.Error("yMax incorrect")
	}

	if manager.xMax != 100 {
		t.Error("xMax incorrect")
	}

	if manager.maxRoids != 20 {
		t.Error("yMax incorrect")
	}
}
