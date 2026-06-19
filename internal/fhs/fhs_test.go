package fhs

import "testing"

func TestConfigDir(t *testing.T) {
	expected := "/etc"

	actual := ConfigDest()

	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestWebSvcDir(t *testing.T) {
	expected := "/srv"

	actual := WebSvcDest()

	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestSvcAssetDir(t *testing.T) {
	expected := "/srv/assets"

	actual := SvcAssetDest()

	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}
