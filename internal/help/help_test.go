package help

import "testing"

const expected = `usage: forte <command> [<args>]

  forte help                          Show this help
  forte version                       Display Forte version
  forte deploy <app-name> <user-name> Deploy an application`

func TestHelp(t *testing.T) {
	actual := Help()
	if actual != expected {
		t.Error("The help command did not produce the expected output")
	}
}
