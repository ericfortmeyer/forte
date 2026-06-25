package deploy

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/ericfortmeyer/forte/internal/ui"
)

const (
	DeploymentTypeConfig = iota
	DeploymentTypeService
	DeploymentTypeAssets
)
const (
	Command                 = "deploy"
	Example                 = "forte deploy [app-name] [user]"
	DeploymentTypeSeparator = "-"
	ConfigSuffix            = "config"
	AssetsSuffix            = "assets"
	defaultFileOwner        = "root"
	defaultDirOwner         = "root"
)

type DeploymentType int

type Deployment struct {
	Src  string
	Type DeploymentType
}

type CopyCfg struct {
	Src       string
	Dst       string
	DirPerms  fs.FileMode
	FilePerms fs.FileMode
	DirOwner  *user.User
	DirGroup  *user.Group
	FileOwner *user.User
	FileGroup *user.Group
	Chown     ChownFunc
}

type DeployConfig struct {
	AppName       string
	Deployment    Deployment
	WebServerUser *user.User
	DirPerms      os.FileMode
	FilePerms     os.FileMode
	DestRoot      string
	ConfigDest    string
	WebSvcDest    string
	SvcAssetDest  string
	Chown         ChownFunc
}

type Owners struct {
	DirOwner  *user.User
	DirGroup  *user.Group
	FileOwner *user.User
	FileGroup *user.Group
}

type CleanupFunc func(Src string) error

type ChownFunc func(filename string, uid, gid int) error

type RootResolver struct {
	rootDir string
}

type PathResolver struct {
	configDir   string
	webSrvDir   string
	srvAssetDir string
}
type DeployInterface interface {
	Deploy(cfg DeployConfig, cleanup CleanupFunc, out io.Writer) error
	ResolveSrc(srcRoot string, appName string) ([]Deployment, error)
}

func (r *RootResolver) ConfigDir(appName string, p *PathResolver) string {
	return filepath.Join(r.rootDir, p.ConfigDir(appName))
}

func (r *RootResolver) WebServiceDir(appName string, p *PathResolver) string {
	return filepath.Join(r.rootDir, p.WebServiceDir(appName))
}

func (r *RootResolver) ServiceAssetDir(appName string, p *PathResolver) string {
	return filepath.Join(r.rootDir, p.ServiceAssetDir(appName))
}

func (p *PathResolver) ConfigDir(appName string) string {
	return filepath.Join(p.configDir, appName)
}

func (p *PathResolver) WebServiceDir(appName string) string {
	return filepath.Join(p.webSrvDir, appName)
}

func (p *PathResolver) ServiceAssetDir(appName string) string {
	return filepath.Join(p.srvAssetDir, appName)
}

func NewRootResolver(rootDir string) *RootResolver {
	return &RootResolver{rootDir: rootDir}
}

func NewPathResolver(cfgDir, webSrvDir, srvAssetDir string) *PathResolver {
	return &PathResolver{configDir: cfgDir, webSrvDir: webSrvDir, srvAssetDir: srvAssetDir}
}

func Deploy(cfg DeployConfig, cleanup CleanupFunc, out io.Writer) error {
	r := NewRootResolver(cfg.DestRoot)
	p := NewPathResolver(cfg.ConfigDest, cfg.WebSvcDest, cfg.SvcAssetDest)

	switch cfg.Deployment.Type {
	case DeploymentTypeConfig:
		start := time.Now()
		_, _ = fmt.Fprintln(out, ui.Working("Installing config..."))
		if cfgErr := installConfig(cfg, r, p); cfgErr != nil {
			return cfgErr
		}
		_, _ = fmt.Fprintln(out, ui.Working("Installed config in"), time.Since(start))
	case DeploymentTypeService:
		start := time.Now()
		_, _ = fmt.Fprintln(out, ui.Working("Installing service data..."))
		if svcErr := installWebService(cfg, r, p); svcErr != nil {
			return svcErr
		}
		_, _ = fmt.Fprintln(out, ui.Working("Installed service data in"), time.Since(start))
	case DeploymentTypeAssets:
		start := time.Now()
		_, _ = fmt.Fprintln(out, ui.Working("Installing service assets..."))
		if svcErr := installServiceAsset(cfg, r, p); svcErr != nil {
			return svcErr
		}
		_, _ = fmt.Fprintln(out, ui.Working("Installed service assets in"), time.Since(start))
	}

	if cleanup != nil {
		if cleanupErr := cleanup(cfg.Deployment.Src); cleanupErr != nil {
			return cleanupErr
		}
	}

	return nil
}

func ResolveSrc(srcRoot, appName string) ([]Deployment, error) {
	var deployments []Deployment

	servicePath := filepath.Join(srcRoot, appName)
	configPath := filepath.Join(srcRoot, appName+DeploymentTypeSeparator+deploymentTypeSuffix(DeploymentTypeConfig))
	serviceAssetPath := filepath.Join(srcRoot, appName+DeploymentTypeSeparator+deploymentTypeSuffix(DeploymentTypeAssets))

	// Check for config deployment
	if _, err := os.Stat(configPath); err == nil {
		deployments = append(deployments, Deployment{
			Src:  configPath,
			Type: DeploymentTypeConfig,
		})
	}

	// Check for service asset deployment
	if _, err := os.Stat(serviceAssetPath); err == nil {
		deployments = append(deployments, Deployment{
			Src:  serviceAssetPath,
			Type: DeploymentTypeAssets,
		})
	}

	// Check for service deployment
	if _, err := os.Stat(servicePath); err == nil {
		deployments = append(deployments, Deployment{
			Src:  servicePath,
			Type: DeploymentTypeService,
		})
	}

	// If nothing was found, return error
	if len(deployments) == 0 {
		return nil, fmt.Errorf("no deployments found for app %q in %s", appName, srcRoot)
	}

	return deployments, nil
}

func installWebService(cfg DeployConfig, r *RootResolver, p *PathResolver) error {
	dst := r.WebServiceDir(cfg.AppName, p)

	if err := checkDeployableDir(dst); err != nil {
		return err
	}

	owners, err := ownersAndGroups(cfg)
	if err != nil {
		return err
	}

	cpCfg := CopyCfg{
		Src:       cfg.Deployment.Src,
		Dst:       dst,
		DirPerms:  cfg.DirPerms,
		FilePerms: cfg.FilePerms,
		DirOwner:  owners.DirOwner,
		DirGroup:  owners.DirGroup,
		FileOwner: owners.FileOwner,
		FileGroup: owners.FileGroup,
		Chown:     cfg.Chown,
	}

	if err := copyRecursive(cpCfg); err != nil {
		return err
	}

	return nil
}

func installServiceAsset(cfg DeployConfig, r *RootResolver, p *PathResolver) error {
	dst := r.ServiceAssetDir(cfg.AppName, p)

	if err := checkDeployableDir(dst); err != nil {
		return err
	}

	owners, err := ownersAndGroups(cfg)
	if err != nil {
		return err
	}

	cpCfg := CopyCfg{
		Src:       cfg.Deployment.Src,
		Dst:       dst,
		DirPerms:  cfg.DirPerms,
		FilePerms: cfg.FilePerms,
		DirOwner:  owners.DirOwner,
		DirGroup:  owners.DirGroup,
		FileOwner: owners.FileOwner,
		FileGroup: owners.FileGroup,
		Chown:     cfg.Chown,
	}

	if err := copyRecursive(cpCfg); err != nil {
		return err
	}

	return nil
}

func installConfig(cfg DeployConfig, r *RootResolver, p *PathResolver) error {
	dst := r.ConfigDir(cfg.AppName, p)

	if err := checkDeployableDir(dst); err != nil {
		return err
	}

	owners, err := ownersAndGroups(cfg)
	if err != nil {
		return err
	}

	cpCfg := CopyCfg{
		Src:       cfg.Deployment.Src,
		Dst:       dst,
		DirPerms:  cfg.DirPerms,
		FilePerms: cfg.FilePerms,
		DirOwner:  owners.DirOwner,
		DirGroup:  owners.DirGroup,
		FileOwner: owners.FileOwner,
		FileGroup: owners.FileGroup,
		Chown:     cfg.Chown,
	}

	if err := copyRecursive(cpCfg); err != nil {
		return err
	}

	return nil
}

func CleanupProduction(Src string) error {
	// TODO: add optional cleanup
	// if err := os.RemoveAll(Src); err != nil {
	// 	return err
	// }
	return nil
}

func ChownProduction(filename string, uid, gid int) error {
	return os.Chown(filename, uid, gid)
}

func copyRecursive(cfg CopyCfg) error {
	return filepath.Walk(cfg.Src, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(cfg.Src, srcPath)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(cfg.Dst, relPath)

		if info.IsDir() {
			if err := os.MkdirAll(dstPath, cfg.DirPerms); err != nil {
				return err
			}

			uid, err := strconv.Atoi(cfg.DirOwner.Uid)
			if err != nil {
				return err
			}

			gid, err := strconv.Atoi(cfg.DirGroup.Gid)
			if err != nil {
				return err
			}

			if err := cfg.Chown(dstPath, uid, gid); err != nil {
				return err
			}
			return os.Chmod(dstPath, cfg.DirPerms)
		}
		// Skip if destination exists and is up-to-date
		dstInfo, err := os.Stat(dstPath)
		if err != nil && !os.IsNotExist(err) {
			return err // Unexpected error
		}
		srcInfo, err := os.Stat(srcPath)
		if err != nil && !os.IsNotExist(err) {
			return err // Unexpected error
		}

		if dstInfo != nil && !srcInfo.ModTime().After(dstInfo.ModTime()) {
			return nil // Destination is up-to-date, skip
		}

		srcFile, err := os.Open(srcPath)
		if err != nil {
			return err
		}
		defer func() {
			if err := srcFile.Close(); err != nil {
				_, _ = fmt.Printf("warning: failed to close file: %v", err)
			}
		}()

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer func() {
			if err := dstFile.Close(); err != nil {
				_, _ = fmt.Printf("warning: failed to close file: %v", err)
			}
		}()

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return err
		}

		uid, err := strconv.Atoi(cfg.FileOwner.Uid)
		if err != nil {
			return err
		}

		gid, err := strconv.Atoi(cfg.FileGroup.Gid)
		if err != nil {
			return err
		}

		if err := cfg.Chown(dstPath, uid, gid); err != nil {
			return err
		}

		return os.Chmod(dstPath, cfg.FilePerms)
	})
}

func ownersAndGroups(cfg DeployConfig) (Owners, error) {
	dirOwner, err := user.Lookup(defaultDirOwner)
	if err != nil {
		return Owners{}, err
	}
	dirGroup, err := user.LookupGroup(cfg.WebServerUser.Username)
	if err != nil {
		return Owners{}, err
	}
	fileOwner, err := user.Lookup(defaultFileOwner)
	if err != nil {
		return Owners{}, err
	}
	fileGroup, err := user.LookupGroup(cfg.WebServerUser.Username)
	if err != nil {
		return Owners{}, err
	}

	return Owners{DirOwner: dirOwner, DirGroup: dirGroup, FileOwner: fileOwner, FileGroup: fileGroup}, nil
}

func deploymentTypeSuffix(dt DeploymentType) string {
	switch dt {
	case DeploymentTypeAssets:
		return AssetsSuffix
	case DeploymentTypeConfig:
		return ConfigSuffix
	}
	return ""
}

func checkOwned(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return err
		}
		return nil
	}

	worldWritable := info.Mode()&0o02 != 0
	groupWritable := info.Mode()&0o020 != 0
	userWritable := info.Mode()&0o200 != 0

	if worldWritable {
		return nil
	}

	stat := info.Sys().(*syscall.Stat_t)

	currentGID := uint32(os.Getgid())
	if stat.Gid == currentGID && groupWritable {
		return nil
	}

	currentUID := uint32(os.Getuid())
	if stat.Uid == currentUID && userWritable {
		return nil
	}
	return fmt.Errorf("you don't own %s (owner: %d, you: %d)", path, stat.Uid, currentUID)
}

func checkDeployableDir(path string) error {
	// if path doesn't exist, check parent
	target := path
	if _, err := os.Stat(path); os.IsNotExist(err) {
		target = filepath.Dir(path)
	}
	return checkOwned(target)
}
