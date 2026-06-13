package fhs

const (
	configDest = "/etc"
	webSvcDest = "/srv"
)

func ConfigDest() string {
	return configDest
}

func WebSvcDest() string {
	return webSvcDest
}
