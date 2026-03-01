package svc

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// RunLaunchctl executes a launchctl command on macOS and returns combined output.
func RunLaunchctl(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "launchctl", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg == "" {
			errMsg = err.Error()
		}
		return "", fmt.Errorf("launchctl %s: %s", strings.Join(args, " "), errMsg)
	}
	return strings.TrimSpace(stdout.String()), nil
}

// RunSystemctl executes a systemctl command on Linux and returns combined output.
func RunSystemctl(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "systemctl", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg == "" {
			errMsg = err.Error()
		}
		return "", fmt.Errorf("systemctl %s: %s", strings.Join(args, " "), errMsg)
	}
	return strings.TrimSpace(stdout.String()), nil
}

// RunService runs the appropriate service manager command for the current OS.
// On darwin it uses launchctl, on linux it uses systemctl.
func RunService(ctx context.Context, args ...string) (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return RunLaunchctl(ctx, args...)
	case "linux":
		return RunSystemctl(ctx, args...)
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// IsDarwin returns true if the current OS is macOS.
func IsDarwin() bool {
	return runtime.GOOS == "darwin"
}

// IsLinux returns true if the current OS is Linux.
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// RunCommand executes an arbitrary command and returns its stdout.
func RunCommand(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg == "" {
			errMsg = err.Error()
		}
		return "", fmt.Errorf("%s %s: %s", name, strings.Join(args, " "), errMsg)
	}
	return strings.TrimSpace(stdout.String()), nil
}
