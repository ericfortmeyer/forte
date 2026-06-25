package main

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/ericfortmeyer/forte/internal/archive"
	"github.com/ericfortmeyer/forte/internal/deploy"
	"github.com/ericfortmeyer/forte/internal/fhs"
	"github.com/ericfortmeyer/forte/internal/help"
	"github.com/ericfortmeyer/forte/internal/ui"
	forteversion "github.com/ericfortmeyer/forte/internal/version"
)

const (
	srcRoot   = "/tmp"
	destRoot  = "" // an empty string resolves to /
	dirPerms  = 0750
	filePerms = 0640
)

type userValidator func(username string) (*user.User, error)

type archiveProxy struct{}

func (a archiveProxy) Extract(tarGzPath, destDir string, out io.Writer) error {
	return archive.Extract(tarGzPath, destDir, out)
}
func (a archiveProxy) IsSkippable(err error) bool { return archive.IsSkippable(err) }

type deployProxy struct{}

func (d deployProxy) ResolveSrc(srcRoot string, appName string) ([]deploy.Deployment, error) {
	return deploy.ResolveSrc(srcRoot, appName)
}

func (d deployProxy) Deploy(cfg deploy.DeployConfig, cleanup deploy.CleanupFunc, out io.Writer) error {
	return deploy.Deploy(cfg, cleanup, out)
}

func main() {
	Run(
		os.Args[1:],
		archiveProxy{},
		deployProxy{},
		user.Lookup,
		exit,
		os.Stderr,
	)
}

func exit(i int) {
	os.Exit(i)
}

func Run(
	args []string,
	a archive.ArchiveInterface,
	d deploy.DeployInterface,
	userValidator userValidator,
	exit func(int),
	out io.Writer,
) {
	if len(args) < 1 {
		_, _ = fmt.Fprintln(out, help.Help())
		exit(1)
		return
	}

	cmd := args[0]

	switch cmd {
	case deploy.Command:
		start := time.Now()

		if len(args) < 2 {
			_, _ = fmt.Fprintln(out, ui.Error("Application name required"))
			_, _ = fmt.Fprintln(out, "") // blank line
			_, _ = fmt.Fprintln(out, deploy.Example)
			exit(1)
			return
		}
		if len(args) < 3 {
			_, _ = fmt.Fprintln(out, ui.Error("Web service user required"))
			_, _ = fmt.Fprintln(out, "") // blank line
			_, _ = fmt.Fprintln(out, deploy.Example)
			exit(1)
			return
		}
		appName := args[1]
		webServerUser := args[2]

		validUser, err := userValidator(webServerUser)
		if err != nil {
			_, _ = fmt.Fprintln(out, ui.Error("user not found "+webServerUser))
			exit(1)
			return
		}

		archiveNames := []string{
			appName,
			appName + deploy.DeploymentTypeSeparator + deploy.ConfigSuffix,
			appName + deploy.DeploymentTypeSeparator + deploy.AssetsSuffix,
		}

		for _, name := range archiveNames {
			tarGzPath := filepath.Join(srcRoot, name+archive.TarballExt)
			destDir := filepath.Join(srcRoot, name)
			if err := a.Extract(tarGzPath, destDir, os.Stderr); err != nil {
				if !a.IsSkippable(err) {
					_, _ = fmt.Fprintln(out, ui.Error(err.Error()))
					exit(1)
					return
				} // IsSkippable errors are silently ignored
			}
		}

		if deployments, err := d.ResolveSrc(srcRoot, appName); err != nil {
			_, _ = fmt.Fprintln(out, ui.Error(err.Error()))
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

				if err := d.Deploy(cfg, deploy.CleanupProduction, out); err != nil {
					_, _ = fmt.Fprintln(out, ui.Error(err.Error()))
					exit(1)
					return
				}
			}

			_, _ = fmt.Fprintln(out, ui.Success("Total:"), time.Since(start))
		}

	case help.Command:
		_, _ = fmt.Fprintln(out, help.Help())
		exit(0)
		return
	case forteversion.Command:
		_, _ = fmt.Fprintln(out, forteversion.Version())
		exit(0)
		return
	default:
		_, _ = fmt.Fprintln(out, ui.Error("forte: unknown subcommand: forte "+cmd))
		_, _ = fmt.Fprintln(out, "") // blank line
		_, _ = fmt.Fprintln(out, "Run 'forte help' for more information")
		exit(1)
	}
}
