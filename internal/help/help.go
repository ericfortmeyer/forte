package help

func Help() string {
	return `usage: forte <command> [<args>]

  forte help                          Show this help
  forte version                       Display Forte version
  forte deploy <app-name> <user-name> Deploy an application
`
}
