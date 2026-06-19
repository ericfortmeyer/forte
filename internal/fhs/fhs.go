package fhs

const (
	configDest    = "/etc"
	webSvcDest    = "/srv"
	svcAssetsDest = "/srv/assets"
)

func ConfigDest() string {
	return configDest
}

func WebSvcDest() string {
	return webSvcDest
}

func SvcAssetDest() string {
	return svcAssetsDest
}
