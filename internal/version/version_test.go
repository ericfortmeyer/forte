package version

import "testing"

const expected = "forte version dev"

func TestVersion(t *testing.T) {
	actual := Version()

	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}
