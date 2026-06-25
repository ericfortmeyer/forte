package version

import "fmt"

const (
	Command  = "version"
	toolName = "forte"
)

var version = "dev"

func Version() string {
	return fmt.Sprintf("%s version %s", toolName, version)
}
