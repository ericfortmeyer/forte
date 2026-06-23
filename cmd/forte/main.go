package main

import (
	"bytes"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"github.com/ericfortmeyer/forte/internal/archive"
	"github.com/ericfortmeyer/forte/internal/deploy"
	"github.com/ericfortmeyer/forte/internal/fhs"
	"github.com/ericfortmeyer/forte/internal/help"
	forteversion "github.com/ericfortmeyer/forte/internal/version"
)

const ()
const (
	srcRoot                 = "/tmp"
	destRoot                = "" // an empty string resolves to /
	helpCmd                 = "help"
	deployCmd               = "deploy"
	versionCmd              = "version"
	dirPerms                = 0750
	filePerms               = 0640
	deploymentTypeSeparator = "-"
	configSuffix            = "config"
	assetsSuffix            = "assets"
	archiveExt              = ".tar.gz"
)

type userValidator func(username string) (*user.User, error)
type archiveInterface interface {
	Extract(tarGzPath string, destDir string) error
	IsSkippable(err error) bool
}
type deployInterface interface {
	Deploy(cfg deploy.DeployConfig, cleanup deploy.CleanupFunc) error
	ResolveSrc(srcRoot string, appName string) ([]deploy.Deployment, error)
}

type archiveProxy struct{}

func (a archiveProxy) Extract(tarGzPath string, destDir string) error {
	return archive.Extract(tarGzPath, destDir)
}
func (a archiveProxy) IsSkippable(err error) bool { return archive.IsSkippable(err) }

type deployProxy struct{}

func (d deployProxy) ResolveSrc(srcRoot string, appName string) ([]deploy.Deployment, error) {
	return deploy.ResolveSrc(srcRoot, appName)
}

func (d deployProxy) Deploy(cfg deploy.DeployConfig, cleanup deploy.CleanupFunc) error {
	return deploy.Deploy(cfg, cleanup)
}

func main() {
	Run(
		os.Args[1:],
		archiveProxy{},
		deployProxy{},
		user.Lookup,
		exit,
		&bytes.Buffer{},
	)
}

func exit(i int) {
	os.Exit(i)
}

func Run(
	args []string,
	a archiveInterface,
	d deployInterface,
	userValidator userValidator,
	exit func(int),
	out io.Writer,
) {
	if len(args) < 1 {
		out.Write([]byte(help.Help()))
		exit(1)
		return
	}

	cmd := args[0]

	switch cmd {
	case deployCmd:
		if len(args) < 2 {
			out.Write([]byte("Application name required"))
			exit(1)
			return
		}
		if len(args) < 3 {
			out.Write([]byte("Web service user required"))
			exit(1)
			return
		}
		appName := args[1]
		webServerUser := args[2]

		validUser, err := userValidator(webServerUser)
		if err != nil {
			out.Write([]byte("Error: user not found " + webServerUser))
			exit(1)
			return
		}

		archiveNames := []string{
			appName,
			appName + deploymentTypeSeparator + configSuffix,
			appName + deploymentTypeSeparator + assetsSuffix,
		}

		for _, name := range archiveNames {
			tarGzPath := filepath.Join(srcRoot, name+archiveExt)
			destDir := filepath.Join(srcRoot, name)
			if err := a.Extract(tarGzPath, destDir); err != nil {
				if !a.IsSkippable(err) {
					out.Write([]byte("Error: " + err.Error()))
					exit(1)
					return
				} // IsSkippable errors are silently ignored
			}
		}

		if deployments, err := d.ResolveSrc(srcRoot, appName); err != nil {
			out.Write([]byte("Error: " + err.Error()))
			exit(1)
			return
		} else {
			for _, deployment := range deployments {
				cfg := deploy.DeployConfig{
					AppName:       appName,
					Deployment:    deployment,
					WebServerUser: validUser,
					DirPerms:      dirPerms,
					FilePerms:     filePerms,
					Chown:         deploy.ChownProduction,
					DestRoot:      destRoot,
					ConfigDest:    fhs.ConfigDest(),   // TODO: support config file / env var override in future version
					WebSvcDest:    fhs.WebSvcDest(),   // TODO: support config file / env var override in future version
					SvcAssetDest:  fhs.SvcAssetDest(), // TODO: support config file / env var override in future version
				}

				if err := d.Deploy(cfg, deploy.CleanupProduction); err != nil {
					out.Write([]byte("Error: " + err.Error()))
					exit(1)
					return
				}
			}
		}

	case helpCmd:
		out.Write([]byte(help.Help()))
		exit(0)
		return
	case versionCmd:
		out.Write([]byte(forteversion.Version()))
		exit(0)
		return
	}
}
