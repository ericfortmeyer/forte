package main

import (
	"bytes"
	"errors"
	"os/user"
	"strings"
	"testing"

	"github.com/ericfortmeyer/forte/internal/deploy"
	"github.com/ericfortmeyer/forte/internal/help"
	forteversion "github.com/ericfortmeyer/forte/internal/version"
)

type mockArchiveIsNotSkippableError struct{}

func (a mockArchiveIsNotSkippableError) Extract(tarGzPath string, destDir string) error {
	return errors.New("Is not skippable")
}
func (a mockArchiveIsNotSkippableError) IsSkippable(err error) bool { return false }

type mockDeployNoop struct{}

func (a mockDeployNoop) ResolveSrc(srcRoot string, appName string) ([]deploy.Deployment, error) {
	return []deploy.Deployment{}, nil
}
func (d mockDeployNoop) Deploy(cfg deploy.DeployConfig, cleanup deploy.CleanupFunc) error {
	return nil
}

type mockDeployResolveSrcFailed struct{}

func (a mockDeployResolveSrcFailed) ResolveSrc(srcRoot string, appName string) ([]deploy.Deployment, error) {
	return []deploy.Deployment{}, errors.New("Deployment src resolution failed")
}
func (d mockDeployResolveSrcFailed) Deploy(cfg deploy.DeployConfig, cleanup deploy.CleanupFunc) error {
	return nil
}

type mockDeploymentFailed struct{}

func (a mockDeploymentFailed) ResolveSrc(srcRoot string, appName string) ([]deploy.Deployment, error) {
	return []deploy.Deployment{
		{
			Type: deploy.DeploymentTypeConfig,
			Src:  "/",
		},
	}, nil
}
func (d mockDeploymentFailed) Deploy(cfg deploy.DeployConfig, cleanup deploy.CleanupFunc) error {
	return errors.New("Deployment failed")
}

type mockArchiveNoop struct{}

func (a mockArchiveNoop) Extract(tarGzPath string, destDir string) error {
	return nil
}
func (a mockArchiveNoop) IsSkippable(err error) bool { return true }

func TestOutputsHelpIfNoSubcommandIsGiven(t *testing.T) {
	var exitCode int
	mockExit := func(code int) { exitCode = code }
	output := &bytes.Buffer{}
	noopUserValidator := func(username string) (*user.User, error) {
		return nil, nil
	}

	Run([]string{}, mockArchiveNoop{}, mockDeployNoop{}, noopUserValidator, mockExit, output)

	if exitCode != 1 {
		t.Fatal("Should have exited with exit code 1")
	}

	if output.String() != help.Help() {
		t.Fatalf("Should have printed help but printed: %s", output)
	}
}

func TestOutputsHelpIfInvalidSubcommandIsGiven(t *testing.T) {
	var exitCode int
	mockExit := func(code int) { exitCode = code }
	output := &bytes.Buffer{}
	noopUserValidator := func(username string) (*user.User, error) {
		return nil, nil
	}

	Run([]string{"INVALID_SUBCOMMAND"}, mockArchiveNoop{}, mockDeployNoop{}, noopUserValidator, mockExit, output)

	if exitCode != 1 {
		t.Fatal("Should have exited with exit code 1")
	}

	if output.String() != "Error: Invalid subcommand: INVALID_SUBCOMMAND. Valid subcommands are deploy, version, and help." {
		t.Fatalf("Should have printed 'Invalid subcommand: INVALID_SUBCOMMAND. Valid subcommands are deploy, version, and help. but printed: %s", output)
	}
}

func TestOutputsHelpIfHelpSubcommandIsGiven(t *testing.T) {
	var exitCode int
	mockExit := func(code int) { exitCode = code }
	output := &bytes.Buffer{}
	noopUserValidator := func(username string) (*user.User, error) {
		return nil, nil
	}

	Run([]string{"help"}, mockArchiveNoop{}, mockDeployNoop{}, noopUserValidator, mockExit, output)

	if exitCode != 0 {
		t.Fatal("Should have exited with exit code 0")
	}

	if output.String() != help.Help() {
		t.Fatalf("Should have printed help but printed: %s", output)
	}
}

func TestOutputsVersionIfVersionSubcommandIsGiven(t *testing.T) {
	var exitCode int
	mockExit := func(code int) { exitCode = code }
	output := &bytes.Buffer{}
	noopUserValidator := func(username string) (*user.User, error) {
		return nil, nil
	}

	Run([]string{"version"}, mockArchiveNoop{}, mockDeployNoop{}, noopUserValidator, mockExit, output)

	if exitCode != 0 {
		t.Fatal("Should have exited with exit code 0")
	}

	if output.String() != forteversion.Version() {
		t.Fatalf("Should have printed version but printed: %s", output)
	}
}

func TestOutputsErrorIfDeploySubcommandIsGivenWithoutAppName(t *testing.T) {
	var exitCode int
	mockExit := func(code int) { exitCode = code }
	output := &bytes.Buffer{}
	noopUserValidator := func(username string) (*user.User, error) {
		return nil, nil
	}

	Run([]string{"deploy"}, mockArchiveNoop{}, mockDeployNoop{}, noopUserValidator, mockExit, output)

	if exitCode != 1 {
		t.Fatal("Should have exited with exit code 1")
	}

	if output.String() != "Application name required" {
		t.Fatalf("Should have printed 'Application name requried' but printed: %s", output)
	}
}

func TestOutputsErrorIfDeploySubcommandIsGivenWithoutSvcUser(t *testing.T) {
	var exitCode int
	mockExit := func(code int) { exitCode = code }
	output := &bytes.Buffer{}
	noopUserValidator := func(username string) (*user.User, error) {
		return nil, nil
	}

	Run(
		[]string{"deploy", "myApp"},
		mockArchiveNoop{},
		mockDeployNoop{},
		noopUserValidator,
		mockExit,
		output,
	)

	if exitCode != 1 {
		t.Fatal("Should have exited with exit code 1")
	}

	if output.String() != "Web service user required" {
		t.Fatalf("Should have printed 'Web service user required' but printed: %s", output)
	}
}

func TestOutputsErrorIfDeploySubcommandIsGivenInvalidSvcUser(t *testing.T) {
	var exitCode int
	mockExit := func(code int) { exitCode = code }
	output := &bytes.Buffer{}
	mockUserValidator := func(username string) (*user.User, error) { return nil, errors.New("it") }

	Run(
		[]string{"deploy", "myApp", "www-data"},
		mockArchiveNoop{},
		mockDeployNoop{},
		mockUserValidator,
		mockExit,
		output,
	)

	if exitCode != 1 {
		t.Fatal("Should have exited with exit code 1")
	}

	if !strings.Contains(output.String(), "Error") || !strings.Contains(output.String(), "www-data") {
		t.Fatalf("Should have printed Error message but printed: %s", output)
	}
}

func TestOutputsErrorIfDeploySubcommandIsGivenAndArchiveErrorIsNotSkippable(t *testing.T) {
	var exitCode int
	mockExit := func(code int) { exitCode = code }
	output := &bytes.Buffer{}
	mockUserValidator := func(username string) (*user.User, error) {
		return &user.User{Username: "www-data"}, nil
	}

	Run(
		[]string{"deploy", "myApp", "www-data"},
		mockArchiveIsNotSkippableError{},
		mockDeployNoop{},
		mockUserValidator,
		mockExit,
		output,
	)

	if exitCode != 1 {
		t.Fatal("Should have exited with exit code 1")
	}

	if output.String() != "Error: Is not skippable" {
		t.Fatalf("Should have printed Error message but printed: %s", output)
	}
}

func TestOutputsErrorIfDeploySubcommandIsGivenAndResolveSrcErrorOccurs(t *testing.T) {
	var exitCode int
	mockExit := func(code int) { exitCode = code }
	output := &bytes.Buffer{}
	mockUserValidator := func(username string) (*user.User, error) {
		return &user.User{Username: "www-data"}, nil
	}

	Run(
		[]string{"deploy", "myApp", "www-data"},
		mockArchiveNoop{},
		mockDeployResolveSrcFailed{},
		mockUserValidator,
		mockExit,
		output,
	)

	if exitCode != 1 {
		t.Fatal("Should have exited with exit code 1")
	}

	if output.String() != "Error: Deployment src resolution failed" {
		t.Fatalf("Should have printed Error message but printed: %s", output)
	}
}

func TestOutputsErrorIfDeploySubcommandIsGivenAndDeployErrorOccurs(t *testing.T) {
	var exitCode int
	mockExit := func(code int) { exitCode = code }
	output := &bytes.Buffer{}
	mockUserValidator := func(username string) (*user.User, error) {
		return &user.User{Username: "www-data"}, nil
	}

	Run(
		[]string{"deploy", "myApp", "www-data"},
		mockArchiveNoop{},
		mockDeploymentFailed{},
		mockUserValidator,
		mockExit,
		output,
	)

	if exitCode != 1 {
		t.Fatal("Should have exited with exit code 1")
	}

	if output.String() != "Error: Deployment failed" {
		t.Fatalf("Should have printed Error message but printed: %s", output)
	}
}
