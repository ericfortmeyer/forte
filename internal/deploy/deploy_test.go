package deploy

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/ericfortmeyer/forte/internal/fhs"
)

func TestPathResolverConfigDir(t *testing.T) {
	// Given
	fakeCfgDir := "/some/fake/dir"
	fakeAppName := "fake_app_name"
	expected := filepath.Join(fakeCfgDir, fakeAppName)
	sut := NewPathResolver(fakeCfgDir, "", "")

	// When
	actual := sut.ConfigDir(fakeAppName)

	// Then
	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestPathResolverWebServiceDir(t *testing.T) {
	// Given
	fakeSvcDir := "/some/fake/dir"
	fakeAppName := "fake_app_name"
	expected := filepath.Join(fakeSvcDir, fakeAppName)
	sut := NewPathResolver("", fakeSvcDir, "")

	// When
	actual := sut.WebServiceDir(fakeAppName)

	// Then
	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestPathResolverServiceAssetDir(t *testing.T) {
	// Given
	fakeSvcAssetDir := "/some/fake/dir"
	fakeAppName := "fake_app_name"
	expected := filepath.Join(fakeSvcAssetDir, fakeAppName)
	sut := NewPathResolver("", "", fakeSvcAssetDir)

	// When
	actual := sut.ServiceAssetDir(fakeAppName)

	// Then
	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestRootResolverConfigDir(t *testing.T) {
	// Given
	fakeRootDir := "/chroot/here"
	fakeCfgDir := "/some/fake/dir"
	fakeAppName := "fake_app_name"
	expected := fakeRootDir + fakeCfgDir + "/" + fakeAppName
	p := NewPathResolver(fakeCfgDir, "", "")
	sut := NewRootResolver(fakeRootDir)

	// When
	actual := sut.ConfigDir(fakeAppName, p)

	// Then
	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestRootResolverWebServiceDir(t *testing.T) {
	// Given
	fakeRootDir := "/chroot/here"
	fakeSvcDir := "/some/fake/dir"
	fakeAppName := "fake_app_name"
	expected := fakeRootDir + fakeSvcDir + "/" + fakeAppName
	p := NewPathResolver("", fakeSvcDir, "")
	sut := NewRootResolver(fakeRootDir)

	// When
	actual := sut.WebServiceDir(fakeAppName, p)

	// Then
	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestRootResolverServiceAssetDir(t *testing.T) {
	// Given
	fakeRootDir := "/chroot/here"
	fakeSvcAssetDir := "/some/fake/dir"
	fakeAppName := "fake_app_name"
	expected := fakeRootDir + fakeSvcAssetDir + "/" + fakeAppName
	p := NewPathResolver("", "", fakeSvcAssetDir)
	sut := NewRootResolver(fakeRootDir)

	// When
	actual := sut.ServiceAssetDir(fakeAppName, p)

	// Then
	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

func TestDeployInstallsConfig(t *testing.T) {
	rootDir := t.TempDir()
	fakeAppName := "fake_app"
	mockChown := func(filename string, uid, gid int) error { return nil }

	testUser := &user.User{
		Uid:      "33", // www-data on many systems
		Gid:      "33",
		Username: "www-data",
		Name:     "www-data",
		HomeDir:  "/var/www",
	}

	deployment := Deployment{
		Src:  filepath.Join("testdata", fakeAppName+DeploymentTypeSeparator+ConfigSuffix),
		Type: DeploymentTypeConfig,
	}

	cfg := DeployConfig{
		AppName:       fakeAppName,
		Deployment:    deployment,
		WebServerUser: testUser,
		DirPerms:      0750,
		FilePerms:     0640,
		ConfigDest:    fhs.ConfigDest(),
		WebSvcDest:    fhs.WebSvcDest(),
		Chown:         mockChown,
		DestRoot:      rootDir,
	}
	expectedDest := filepath.Join(cfg.DestRoot+cfg.ConfigDest, fakeAppName)

	cfgErr := Deploy(cfg, nil)

	if cfgErr != nil && !os.IsExist(cfgErr) {
		t.Errorf("The expected destination folder %s does not exist", expectedDest)
	}

	_, err := os.ReadFile(filepath.Join(expectedDest, "app_info.php"))
	if err != nil {
		t.Errorf("The configuration files did not copy as expected. Checked %s", expectedDest)
	}

	finfo, err := os.Stat(filepath.Join(expectedDest, "app_info.php"))
	if err != nil {
		if mode := finfo.Mode().Perm(); mode != 0640 {
			t.Errorf("file mode = %o, want 0640", mode)
		}
	}
}

func TestDeployInstallsService(t *testing.T) {
	rootDir := t.TempDir()
	fakeAppName := "fake_app"
	mockChown := func(filename string, uid, gid int) error { return nil }

	deployment := Deployment{
		Src:  filepath.Join("testdata", fakeAppName),
		Type: DeploymentTypeService,
	}

	cfg := DeployConfig{
		AppName:    fakeAppName,
		Deployment: deployment,
		WebServerUser: &user.User{
			Uid:      "33",
			Gid:      "33",
			Username: "www-data",
			HomeDir:  "/var/www",
		},
		DirPerms:   0750,
		FilePerms:  0640,
		ConfigDest: fhs.ConfigDest(),
		WebSvcDest: fhs.WebSvcDest(),
		Chown:      mockChown,
		DestRoot:   rootDir,
	}

	p := NewPathResolver(fhs.ConfigDest(), fhs.WebSvcDest(), fhs.SvcAssetDest())
	r := NewRootResolver(rootDir)
	destDir := r.WebServiceDir(fakeAppName, p)

	// Deploy
	if err := Deploy(cfg, nil); err != nil {
		t.Fatalf("Deploy failed: %v", err)
	}

	// Verify files copied
	requiredFiles := []string{
		"src/ItemService.php",
		"public/index.php",
	}
	for _, file := range requiredFiles {
		if _, err := os.ReadFile(filepath.Join(destDir, file)); err != nil {
			t.Errorf("missing file: %s", file)
		}
	}
}

func TestDeployInstallsServiceAssets(t *testing.T) {
	rootDir := t.TempDir()
	fakeAppName := "fake_app"
	mockChown := func(filename string, uid, gid int) error { return nil }

	deployment := Deployment{
		Src:  filepath.Join("testdata", fakeAppName+DeploymentTypeSeparator+AssetsSuffix),
		Type: DeploymentTypeAssets,
	}

	cfg := DeployConfig{
		AppName:    fakeAppName,
		Deployment: deployment,
		WebServerUser: &user.User{
			Uid:      "33",
			Gid:      "33",
			Username: "www-data",
			HomeDir:  "/var/www",
		},
		DirPerms:     0750,
		FilePerms:    0640,
		ConfigDest:   fhs.ConfigDest(),
		WebSvcDest:   fhs.WebSvcDest(),
		SvcAssetDest: fhs.SvcAssetDest(),
		Chown:        mockChown,
		DestRoot:     rootDir,
	}

	p := NewPathResolver(fhs.ConfigDest(), fhs.WebSvcDest(), fhs.SvcAssetDest())
	r := NewRootResolver(rootDir)
	destDir := r.ServiceAssetDir(fakeAppName, p)

	// Deploy
	if err := Deploy(cfg, nil); err != nil {
		t.Fatalf("Deploy failed: %v", err)
	}

	// Verify files copied
	requiredFiles := []string{
		"app.css",
	}
	for _, file := range requiredFiles {
		if _, err := os.ReadFile(filepath.Join(destDir, file)); err != nil {
			t.Errorf("missing file: %s", file)
		}
	}
}

func TestResolveSrc(t *testing.T) {
	tests := []struct {
		name            string
		setup           func(t *testing.T, srcRoot string)
		appName         string
		wantDeployments []DeploymentType // Expected types in order
		wantErr         bool
	}{
		{
			name: "only config version exists",
			setup: func(t *testing.T, srcRoot string) {
				configPath := filepath.Join(srcRoot, "myapp"+DeploymentTypeSeparator+ConfigSuffix)
				if err := os.MkdirAll(configPath, 0755); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			},
			appName:         "myapp",
			wantDeployments: []DeploymentType{DeploymentTypeConfig},
			wantErr:         false,
		},
		{
			name: "only assets version exists",
			setup: func(t *testing.T, srcRoot string) {
				assetsPath := filepath.Join(srcRoot, "myapp"+DeploymentTypeSeparator+AssetsSuffix)
				if err := os.MkdirAll(assetsPath, 0755); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			},
			appName:         "myapp",
			wantDeployments: []DeploymentType{DeploymentTypeAssets},
			wantErr:         false,
		},
		{
			name: "only service version exists",
			setup: func(t *testing.T, srcRoot string) {
				servicePath := filepath.Join(srcRoot, "myapp")
				if err := os.MkdirAll(servicePath, 0755); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			},
			appName:         "myapp",
			wantDeployments: []DeploymentType{DeploymentTypeService},
			wantErr:         false,
		},
		{
			name: "both config and service exist",
			setup: func(t *testing.T, srcRoot string) {
				servicePath := filepath.Join(srcRoot, "myapp")
				configPath := filepath.Join(srcRoot, "myapp"+DeploymentTypeSeparator+ConfigSuffix)
				if err := os.MkdirAll(servicePath, 0755); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				if err := os.MkdirAll(configPath, 0755); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			},
			appName:         "myapp",
			wantDeployments: []DeploymentType{DeploymentTypeConfig, DeploymentTypeService},
			wantErr:         false,
		},
		{
			name: "both assets and service exist",
			setup: func(t *testing.T, srcRoot string) {
				servicePath := filepath.Join(srcRoot, "myapp")
				assetsPath := filepath.Join(srcRoot, "myapp"+DeploymentTypeSeparator+AssetsSuffix)
				if err := os.MkdirAll(servicePath, 0755); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				if err := os.MkdirAll(assetsPath, 0755); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			},
			appName:         "myapp",
			wantDeployments: []DeploymentType{DeploymentTypeAssets, DeploymentTypeService},
			wantErr:         false,
		},
		{
			name: "both assets and config exist",
			setup: func(t *testing.T, srcRoot string) {
				assetsPath := filepath.Join(srcRoot, "myapp"+DeploymentTypeSeparator+AssetsSuffix)
				configPath := filepath.Join(srcRoot, "myapp"+DeploymentTypeSeparator+ConfigSuffix)
				if err := os.MkdirAll(assetsPath, 0755); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				if err := os.MkdirAll(configPath, 0755); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			},
			appName:         "myapp",
			wantDeployments: []DeploymentType{DeploymentTypeConfig, DeploymentTypeAssets},
			wantErr:         false,
		},
		{
			name: "assets, service, and config exist",
			setup: func(t *testing.T, srcRoot string) {
				servicePath := filepath.Join(srcRoot, "myapp")
				assetsPath := filepath.Join(srcRoot, "myapp"+DeploymentTypeSeparator+AssetsSuffix)
				configPath := filepath.Join(srcRoot, "myapp"+DeploymentTypeSeparator+ConfigSuffix)
				if err := os.MkdirAll(servicePath, 0755); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				if err := os.MkdirAll(assetsPath, 0755); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
				if err := os.MkdirAll(configPath, 0755); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			},
			appName:         "myapp",
			wantDeployments: []DeploymentType{DeploymentTypeConfig, DeploymentTypeAssets, DeploymentTypeService},
			wantErr:         false,
		},
		{
			name: "no version exists",
			setup: func(t *testing.T, srcRoot string) {
				// create nothing
			},
			appName:         "nonexistent",
			wantDeployments: nil,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcRoot := t.TempDir()
			tt.setup(t, srcRoot)

			got, err := ResolveSrc(srcRoot, tt.appName)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("ResolveSrc() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check number of deployments
			if len(got) != len(tt.wantDeployments) {
				t.Errorf("ResolveSrc() returned %d deployments, want %d", len(got), len(tt.wantDeployments))
				return
			}

			// Check each deployment type and path validity
			for i, deployment := range got {
				if deployment.Type != tt.wantDeployments[i] {
					t.Errorf("ResolveSrc()[%d].Type = %v, want %v", i, deployment.Type, tt.wantDeployments[i])
				}

				// Verify Src path exists
				if _, err := os.Stat(deployment.Src); err != nil {
					t.Errorf("ResolveSrc()[%d].Src path does not exist: %v", i, err)
				}

				// Verify Src matches expected pattern
				var expectedSrc string
				switch deployment.Type {
				case DeploymentTypeConfig:
					expectedSrc = filepath.Join(srcRoot, tt.appName+DeploymentTypeSeparator+ConfigSuffix)
				case DeploymentTypeAssets:
					expectedSrc = filepath.Join(srcRoot, tt.appName+DeploymentTypeSeparator+AssetsSuffix)
				default:
					expectedSrc = filepath.Join(srcRoot, tt.appName)
				}
				if deployment.Src != expectedSrc {
					t.Errorf("ResolveSrc()[%d].Src = %q, want %q", i, deployment.Src, expectedSrc)
				}
			}
		})
	}
}

func TestOwnersAndGroupsFailsDirGroupLookup(t *testing.T) {
	cfg := DeployConfig{
		WebServerUser: &user.User{Username: "nonexistent-group-xyz"},
	}

	_, err := ownersAndGroups(cfg)
	if err == nil {
		t.Error("Expected ownersAndGroups to fail when group doesn't exist")
	}
}

func TestOwnersAndGroupsFailsFileGroupLookup(t *testing.T) {
	// This tests the fourth lookup that also fails on nonexistent group
	cfg := DeployConfig{
		WebServerUser: &user.User{Username: "nonexistent-group-123"},
	}

	_, err := ownersAndGroups(cfg)
	if err == nil {
		t.Error("Expected ownersAndGroups to fail")
	}
}

func TestCopyRecursiveHandlesSourceNotExist(t *testing.T) {
	cfg := CopyCfg{
		Src:       "/nonexistent/source/path",
		Dst:       t.TempDir(),
		DirPerms:  0755,
		FilePerms: 0644,
		DirOwner:  &user.User{Uid: "0"},
		DirGroup:  &user.Group{Gid: "0"},
		FileOwner: &user.User{Uid: "0"},
		FileGroup: &user.Group{Gid: "0"},
		Chown:     func(string, int, int) error { return nil },
	}

	err := copyRecursive(cfg)
	if err == nil {
		t.Error("Expected copyRecursive to fail when source doesn't exist")
	}
}

func TestCopyRecursiveHandlesCreateDstFail(t *testing.T) {
	srcDir := t.TempDir()
	srcFile := filepath.Join(srcDir, "test.txt")
	if err := os.WriteFile(srcFile, []byte("test"), 0644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// Read-only destination directory (can't create subdirs)
	dstDir := t.TempDir()
	if err := os.Chmod(dstDir, 0000); err != nil {
		t.Fatalf("chmod failed: %v", err)
	}
	defer func() {
		if err := os.Chmod(dstDir, 0755); err != nil {
			t.Fatalf("chmod failed: %v", err)
		}
	}()

	cfg := CopyCfg{
		Src:       srcDir,
		Dst:       filepath.Join(dstDir, "readonly"),
		DirPerms:  0755,
		FilePerms: 0644,
		DirOwner:  &user.User{Uid: "0"},
		DirGroup:  &user.Group{Gid: "0"},
		FileOwner: &user.User{Uid: "0"},
		FileGroup: &user.Group{Gid: "0"},
		Chown:     func(string, int, int) error { return nil },
	}

	err := copyRecursive(cfg)
	if err == nil {
		t.Error("Expected copyRecursive to fail when destination is not writable")
	}
}

func TestCopyRecursiveHandlesInvalidUid(t *testing.T) {
	srcDir := t.TempDir()
	srcFile := filepath.Join(srcDir, "test.txt")
	if err := os.WriteFile(srcFile, []byte("test"), 0644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	cfg := CopyCfg{
		Src:       srcDir,
		Dst:       t.TempDir(),
		DirPerms:  0755,
		FilePerms: 0644,
		DirOwner:  &user.User{Uid: "not-a-number"}, // Invalid UID
		DirGroup:  &user.Group{Gid: "0"},
		FileOwner: &user.User{Uid: "0"},
		FileGroup: &user.Group{Gid: "0"},
		Chown:     func(string, int, int) error { return nil },
	}

	err := copyRecursive(cfg)
	if err == nil {
		t.Error("Expected copyRecursive to fail on invalid UID conversion")
	}
}

func TestCopyRecursiveHandlesInvalidFileUid(t *testing.T) {
	srcDir := t.TempDir()
	srcFile := filepath.Join(srcDir, "test.txt")
	if err := os.WriteFile(srcFile, []byte("test"), 0644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	cfg := CopyCfg{
		Src:       srcDir,
		Dst:       t.TempDir(),
		DirPerms:  0755,
		FilePerms: 0644,
		DirOwner:  &user.User{Uid: "0"},
		DirGroup:  &user.Group{Gid: "0"},
		FileOwner: &user.User{Uid: "invalid-uid"}, // Invalid UID
		FileGroup: &user.Group{Gid: "0"},
		Chown:     func(string, int, int) error { return nil },
	}

	err := copyRecursive(cfg)
	if err == nil {
		t.Error("Expected copyRecursive to fail on invalid file UID conversion")
	}
}

func TestCopyRecursiveHandlesChownFailure(t *testing.T) {
	srcDir := t.TempDir()
	srcFile := filepath.Join(srcDir, "test.txt")
	if err := os.WriteFile(srcFile, []byte("test"), 0644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	mockChown := func(filename string, uid, gid int) error {
		return fmt.Errorf("permission denied")
	}

	cfg := CopyCfg{
		Src:       srcDir,
		Dst:       t.TempDir(),
		DirPerms:  0755,
		FilePerms: 0644,
		DirOwner:  &user.User{Uid: "0"},
		DirGroup:  &user.Group{Gid: "0"},
		FileOwner: &user.User{Uid: "0"},
		FileGroup: &user.Group{Gid: "0"},
		Chown:     mockChown,
	}

	err := copyRecursive(cfg)
	if err == nil {
		t.Error("Expected copyRecursive to fail when chown fails")
	}
}

func TestCopyRecursiveHandlesFileOpenFailure(t *testing.T) {
	srcDir := t.TempDir()
	srcFile := filepath.Join(srcDir, "test.txt")
	if err := os.WriteFile(srcFile, []byte("test"), 0000); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	defer func() {
		if err := os.Chmod(srcFile, 0644); err != nil {
			t.Fatalf("chmod failed: %v", err)
		}
	}()

	cfg := CopyCfg{
		Src:       srcDir,
		Dst:       t.TempDir(),
		DirPerms:  0755,
		FilePerms: 0644,
		DirOwner:  &user.User{Uid: "0"},
		DirGroup:  &user.Group{Gid: "0"},
		FileOwner: &user.User{Uid: "0"},
		FileGroup: &user.Group{Gid: "0"},
		Chown:     func(string, int, int) error { return nil },
	}

	err := copyRecursive(cfg)
	if err == nil {
		t.Error("Expected copyRecursive to fail when source file can't be opened")
	}
}

func TestDeployCleanupFailure(t *testing.T) {
	rootDir := t.TempDir()
	fakeAppName := "fake_app"
	mockChown := func(filename string, uid, gid int) error { return nil }

	deployment := Deployment{
		Src:  filepath.Join("testdata", fakeAppName),
		Type: DeploymentTypeService,
	}

	cfg := DeployConfig{
		AppName:    fakeAppName,
		Deployment: deployment,
		WebServerUser: &user.User{
			Uid: "0", Gid: "0", Username: "root", HomeDir: "/root",
		},
		DirPerms:   0755,
		FilePerms:  0644,
		ConfigDest: "/etc",
		WebSvcDest: "/srv",
		Chown:      mockChown,
		DestRoot:   rootDir,
	}

	failingCleanup := func(src string) error {
		return fmt.Errorf("cleanup failed")
	}

	// Deploy succeeds but cleanup fails
	err := Deploy(cfg, failingCleanup)
	if err == nil || err.Error() != "cleanup failed" {
		t.Errorf("Expected cleanup error, got: %v", err)
	}
}

func TestInstallConfigPropagatesError(t *testing.T) {
	cfg := DeployConfig{
		AppName: "testapp",
		Deployment: Deployment{
			Src:  "/nonexistent",
			Type: DeploymentTypeConfig,
		},
		WebServerUser: &user.User{Username: "nonexistent"},
		DestRoot:      t.TempDir(),
		ConfigDest:    "/etc",
		WebSvcDest:    "/srv",
		Chown:         func(string, int, int) error { return nil },
	}

	err := Deploy(cfg, nil)
	if err == nil {
		t.Error("Expected Deploy to propagate installConfig error")
	}
}

func TestInstallServicePropagatesError(t *testing.T) {
	cfg := DeployConfig{
		AppName: "testapp",
		Deployment: Deployment{
			Src:  "/nonexistent",
			Type: DeploymentTypeService,
		},
		WebServerUser: &user.User{Username: "nonexistent"},
		DestRoot:      t.TempDir(),
		ConfigDest:    "/etc",
		WebSvcDest:    "/srv",
		Chown:         func(string, int, int) error { return nil },
	}

	err := Deploy(cfg, nil)
	if err == nil {
		t.Error("Expected Deploy to propagate installWebService error")
	}
}

func TestInstallServiceAssetPropagatesError(t *testing.T) {
	cfg := DeployConfig{
		AppName: "testapp",
		Deployment: Deployment{
			Src:  "/nonexistent",
			Type: DeploymentTypeAssets,
		},
		WebServerUser: &user.User{Username: "nonexistent"},
		DestRoot:      t.TempDir(),
		ConfigDest:    "/etc",
		WebSvcDest:    "/srv",
		SvcAssetDest:  "/srv/assets",
		Chown:         func(string, int, int) error { return nil },
	}

	err := Deploy(cfg, nil)
	if err == nil {
		t.Error("Expected Deploy to propagate installServiceAsset error")
	}
}

func TestDeploymentTypeSuffixDefault(t *testing.T) {
	// Tests the default case in deploymentTypeSuffix
	result := deploymentTypeSuffix(999) // Invalid type
	if result != "" {
		t.Errorf("Expected empty string for invalid type, got %q", result)
	}
}

func TestCopyRecursiveInvalidGid(t *testing.T) {
	srcDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(srcDir, "test.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	cfg := CopyCfg{
		Src:       srcDir,
		Dst:       t.TempDir(),
		DirPerms:  0755,
		FilePerms: 0644,
		DirOwner:  &user.User{Uid: "0"},
		DirGroup:  &user.Group{Gid: "invalid-gid"}, // Invalid GID
		FileOwner: &user.User{Uid: "0"},
		FileGroup: &user.Group{Gid: "0"},
		Chown:     func(string, int, int) error { return nil },
	}

	err := copyRecursive(cfg)
	if err == nil {
		t.Error("Expected error on invalid GID")
	}
}
