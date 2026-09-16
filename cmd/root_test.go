package cmd

import (
	"errors"
	"io"
	"testing"

	"github.com/grpc-kit/pkg/vars"
)

func TestRootCommandReturnsCommandErrors(t *testing.T) {
	cmd := newRootCommand(func() error { return nil })
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"does-not-exist"})

	if err := cmd.Execute(); err == nil {
		t.Fatal("Execute() error = nil, want an unknown-command error")
	}
}

func TestProjectCommandsSkipUserConfig(t *testing.T) {
	previous := vars.ReleaseVersion
	vars.ReleaseVersion = "0.3.9-beta.1"
	t.Cleanup(func() { vars.ReleaseVersion = previous })

	configErr := errors.New("user config must not be loaded")
	loaded := false
	cmd := newRootCommand(func() error {
		loaded = true
		return configErr
	})
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"project", "migrate", writeCommandTestProject(t)})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if loaded {
		t.Fatal("project migrate loaded user generation config")
	}
}

func TestGenerationCommandsLoadUserConfig(t *testing.T) {
	configErr := errors.New("config sentinel")
	loaded := false
	cmd := newRootCommand(func() error {
		loaded = true
		return configErr
	})
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)
	cmd.SetArgs([]string{"new"})

	err := cmd.Execute()
	if !errors.Is(err, configErr) {
		t.Fatalf("Execute() error = %v, want %v", err, configErr)
	}
	if !loaded {
		t.Fatal("new command did not load user generation config")
	}
}

func TestRootCommandDoesNotRegisterAmbiguousUpgradeAliases(t *testing.T) {
	cmd := newRootCommand(func() error { return nil })
	for _, child := range cmd.Commands() {
		if child.Name() == "update" || child.Name() == "upgrade" {
			t.Fatalf("unexpected root command %q", child.Name())
		}
	}
}
