package version

import "fmt"

const name = "forte"

var version = "dev"

func Version() string {
	return fmt.Sprintf("%s version %s", name, version)
}
