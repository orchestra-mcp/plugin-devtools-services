package tools

import (
	"context"
	"fmt"

	pluginv1 "github.com/orchestra-mcp/gen-go/orchestra/plugin/v1"
	"github.com/orchestra-mcp/plugin-devtools-services/internal/svc"
	"github.com/orchestra-mcp/sdk-go/helpers"
	"google.golang.org/protobuf/types/known/structpb"
)

// StopServiceSchema returns the JSON Schema for the stop_service tool.
func StopServiceSchema() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{
				"type":        "string",
				"description": "Name of the service to stop",
			},
		},
		"required": []any{"name"},
	})
	return s
}

// StopService returns a tool handler that stops an OS service.
func StopService() func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
	return func(ctx context.Context, req *pluginv1.ToolRequest) (*pluginv1.ToolResponse, error) {
		if err := helpers.ValidateRequired(req.Arguments, "name"); err != nil {
			return helpers.ErrorResult("validation_error", err.Error()), nil
		}

		name := helpers.GetString(req.Arguments, "name")

		var output string
		var err error

		if svc.IsDarwin() {
			output, err = svc.RunLaunchctl(ctx, "stop", name)
		} else if svc.IsLinux() {
			output, err = svc.RunSystemctl(ctx, "stop", name)
		} else {
			return helpers.ErrorResult("unsupported_os", "Only macOS (launchctl) and Linux (systemctl) are supported"), nil
		}

		if err != nil {
			return helpers.ErrorResult("exec_error", err.Error()), nil
		}

		if output == "" {
			output = fmt.Sprintf("Service %q stopped successfully", name)
		}
		return helpers.TextResult(output), nil
	}
}
