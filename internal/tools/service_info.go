package tools

import (
	"context"
	"fmt"
	"strings"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-services/internal/svc"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// ServiceInfoSchema returns the JSON Schema for the service_info tool.
func ServiceInfoSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{
				"type":        "string",
				"description": "Name of the service to get info for",
			},
		},
		"required": []any{"name"},
	})
	return s
}

// ServiceInfo returns a tool handler that gets service status and info.
func ServiceInfo() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "name"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		name := helpers.GetString(req.Arguments, "name")

		var output string
		var err error

		if svc.IsDarwin() {
			// Use launchctl print to get detailed service info.
			// Try system domain first, then user domain.
			output, err = svc.RunLaunchctl(ctx, "print", fmt.Sprintf("system/%s", name))
			if err != nil {
				output, err = svc.RunLaunchctl(ctx, "print", fmt.Sprintf("gui/%d/%s", getUserID(), name))
			}
			if err != nil {
				// Fallback: use launchctl list and grep for the service.
				listOutput, listErr := svc.RunLaunchctl(ctx, "list")
				if listErr != nil {
					return helpers.ErrorResult("exec_error", listErr.Error()), nil
				}
				var matched []string
				for _, line := range strings.Split(listOutput, "\n") {
					if strings.Contains(line, name) {
						matched = append(matched, line)
					}
				}
				if len(matched) == 0 {
					return helpers.ErrorResult("not_found", fmt.Sprintf("Service %q not found", name)), nil
				}
				output = strings.Join(matched, "\n")
			}
		} else if svc.IsLinux() {
			output, err = svc.RunSystemctl(ctx, "status", name, "--no-pager")
		} else {
			return helpers.ErrorResult("unsupported_os", "Only macOS and Linux are supported"), nil
		}

		if err != nil {
			return helpers.ErrorResult("exec_error", err.Error()), nil
		}

		return helpers.TextResult(output), nil
	}
}

// getUserID returns the current user's UID. Defaults to 501 on macOS.
func getUserID() int {
	return 501
}
