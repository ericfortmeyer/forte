package main

import (
	"fmt"
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

func main() {
	cmd := helpCmd
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case deployCmd:
		if len(os.Args) < 3 {
			fmt.Println("Application name required")
			os.Exit(1)
		}
		appName := os.Args[2]

		if len(os.Args) < 4 {
			fmt.Println("Web service user required")
			os.Exit(1)
		}
		webServerUser := os.Args[3]

		validUser, err := user.Lookup(webServerUser)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: user %q not found\n", webServerUser)
			os.Exit(1)
		}

		archiveNames := []string{
			appName,
			appName + deploymentTypeSeparator + configSuffix,
			appName + deploymentTypeSeparator + assetsSuffix,
		}

		for _, name := range archiveNames {
			tarGzPath := filepath.Join(srcRoot, name+archiveExt)
			destDir := filepath.Join(srcRoot, name)
			if err := archive.Extract(tarGzPath, destDir); err != nil {
				if !archive.IsSkippable(err) {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
					os.Exit(1)
				} // IsSkippable errors are silently ignored
			}
		}

		deployments, err := deploy.ResolveSrc(srcRoot, appName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}

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

			if err := deploy.Deploy(cfg, deploy.CleanupProduction); err != nil {
				fmt.Fprintf(os.Stderr, "Error: deployment failed: %v\n", err)
				os.Exit(1)
			}
		}

		os.Exit(0)
	case helpCmd:
		fmt.Println(help.Help())
		os.Exit(0)
	case versionCmd:
		fmt.Println(forteversion.Version())
		os.Exit(0)
	default:
		fmt.Println(help.Help())
		os.Exit(1)
	}
}
